package main

import (
	"fmt"
	"path/filepath"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/operator"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/features"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(k8sscheme.AddToScheme(scheme))
}

func Manager() (ctrl.Manager, error) {
	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	// Setup a Manager
	setupLog.Info("configuring manager")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0", // TODO(njhale): Enable metrics on non-conflicting port (not 8080)
	})
	if err != nil {
		return nil, err
	}

	// Setup a new controller to reconcile Operators
	setupLog.Info("configuring controller")

	if features.Gate.Enabled(features.OperatorLifecycleManagerV2) {
		setupLog.Info(fmt.Sprintf("feature enabled: %v", features.OperatorLifecycleManagerV2))

		_, err := envtest.InstallCRDs(mgr.GetConfig(), envtest.CRDInstallOptions{
			Paths:              []string{filepath.Join("..", "..", "config", "crd", "bases")},
			ErrorIfPathMissing: true,
		})
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return nil, err
		}

		setupLog.Info("v2alpha1 CRDs installed")

		reconciler, err := operator.NewOperatorReconciler(
			mgr.GetClient(),
			ctrl.Log.WithName("controllers").WithName("Operator"),
			mgr.GetScheme(),
		)
		if err != nil {
			setupLog.Error(err, "unable to create reconciler", "controller", "Operator")
			return nil, err
		}

		if err = reconciler.SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "Operator")
			return nil, err
		}
	}

	setupLog.Info("manager configured")

	return mgr, nil
}
