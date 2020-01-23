/*
Copyright 2020 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	operatorsv2alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v2alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/operator/decorators"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/lib/controller-runtime/source"
)

// OperatorReconciler reconciles a Operator object.
type OperatorReconciler struct {
	client.Client

	log     logr.Logger
	mu      sync.RWMutex
	factory decorators.OperatorFactory
 
	// operators contains the names of Operators the OperatorReconciler has observed exist.
	operators map[types.NamespacedName]struct{}
	source *source.Dynamic
}

// +kubebuilder:rbac:groups=operators.coreos.com,resources=operators,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups=operators.coreos.com,resources=operators/status,verbs=update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch

// SetupWithManager adds the operator reconciler to the given controller manager.
func (r *OperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Trigger operator events from the events of their compoenents.
	enqueueOperator := &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(r.mapComponentRequests),
	}

	// Add reconciler enqueued by dynamic Source watching all GVKs.
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv2alpha1.Operator{}).
		Watches(r.source, enqueueOperator).
		Complete(r)
}

// NewOperatorReconciler constructs and returns an OperatorReconciler.
// As a side effect, the given scheme has operator discovery types added to it
func NewOperatorReconciler(cli client.Client, log logr.Logger, scheme *runtime.Scheme) (*OperatorReconciler, error) {
	factory, err := decorators.NewSchemedOperatorFactory(scheme)
	if err != nil {
		return nil, err
	}

	return &OperatorReconciler{
		Client: cli,

		log:     log,
		factory: factory,
		operators:  map[types.NamespacedName]struct{}{},
		source:  &source.Dynamic{},
	}, nil
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &OperatorReconciler{}

func (r *OperatorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", req)
	log.V(1).Info("reconciling operator")

	// Fetch the Operator from the cache
	ctx := context.TODO()
	in := &operatorsv2alpha1.Operator{}
	if err := r.Get(ctx, req.NamespacedName, in); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Could not find Operator")
			r.unobserve(req.NamespacedName)
			// TODO(njhale): Recreate operator if we can find any components.
		} else {
			log.Error(err, "Error finding Operator")
		}

		return reconcile.Result{}, nil
	}
	r.observe(req.NamespacedName)

	// Wrap with convenience decorator
	operator, err := r.factory.NewOperator(in)
	if err != nil {
		log.Error(err, "Could not wrap Operator with convenience decorator")
		return reconcile.Result{}, nil
	}

	if err = r.updateComponents(ctx, operator); err != nil {
		log.Error(err, "Could not update components")
		return reconcile.Result{}, nil

	}

	if err := r.Update(ctx, operator.Operator); err != nil {
		log.Error(err, "Could not update Operator status")
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, req.NamespacedName, operator.Operator); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *OperatorReconciler) updateComponents(ctx context.Context, operator *decorators.Operator) error {
	selector, err := operator.ComponentSelector()
	if err != nil {
		return err
	}

	components, err := r.listComponents(ctx, selector)
	if err != nil {
		return err
	}

	return operator.SetComponents(components...)
}

func (r *OperatorReconciler) listComponents(ctx context.Context, selector labels.Selector) ([]runtime.Object, error) {
	var componentLists []runtime.Object
	for _, gvk := range r.source.Active() {
		gvk.Kind = gvk.Kind + "List"
		ul := &unstructured.UnstructuredList{}
		ul.SetGroupVersionKind(gvk)
		componentLists = append(componentLists, ul)
	}

	opt := client.MatchingLabelsSelector{Selector: selector}
	for _, list := range componentLists {
		if err := r.List(ctx, list, opt); err != nil {
			return nil, err
		}
	}

	return componentLists, nil
}

func (r *OperatorReconciler) observed(name types.NamespacedName) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.operators[name]
	return ok
}

func (r *OperatorReconciler) observe(name types.NamespacedName) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.operators[name] = struct{}{}
}

func (r *OperatorReconciler) unobserve(name types.NamespacedName) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.operators, name)
}

func (r *OperatorReconciler) mapComponentRequests(obj handler.MapObject) (requests []reconcile.Request) {
	if obj.Meta == nil {
		return
	}

	for _, name := range decorators.OperatorNames(obj.Meta.GetLabels()) {
		// Only enqueue if we can find the operator in our cache
		if r.observed(name) {
			requests = append(requests, reconcile.Request{NamespacedName: name})
			continue
		}

		// Otherwise, best-effort generate a new operator
		// TODO(njhale): Implement verification that the operator-discovery admission webhook accepted this label (JWT or maybe sign a set of fields?)
		operator := &operatorsv2alpha1.Operator{}
		operator.SetName(name.Name)
		if err := r.Create(context.Background(), operator); err != nil && !apierrors.IsAlreadyExists(err) {
			r.log.Error(err, "couldn't generate operator", "operator", name, "component", obj.Meta.GetSelfLink())
		}
	}

	return
}