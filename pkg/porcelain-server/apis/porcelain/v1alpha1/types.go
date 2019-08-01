package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/lib/version"
)

// TODO: Add godoc comments.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstalledOperatorList is a list of InstalledOperator objects.
type InstalledOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []InstalledOperator `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstalledOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// ClusterServiceVersionRef references the CSV which attempted to install the operator.
	ClusterServiceVersionRef *corev1.ObjectReference `json:"clusterServiceVersionRef" protobuf:"bytes,2,opt,name=clusterServiceVersionRef"`
	// SubscriptionRef references the Subscription that installed the referenced CSV, if the CSV was installed via Subscription.
	SubscriptionRef *corev1.ObjectReference `json:"subscriptionRef,omitempty" protobuf:"bytes,3,opt,name=clusterServiceVersionRef"`

	// Fields projected from the referenced CSV

	// CustomResourceDefinitions is the set of CustomResourceDefinitions provided and required by the referenced CSV.
	CustomResourceDefinitions operatorsv1alpha1.CustomResourceDefinitions `json:"customResourceDefinitions,omitempty"`
	// APIServiceDefinitions is the set of APIServices provided and required by the referenced CSV.
	APIServiceDefinitions operatorsv1alpha1.APIServiceDefinitions `json:"apiServiceDefinitions,omitempty"`
	// MinKubeVersion is the minimum kubernetes version the operator is compatible with.
	MinKubeVersion string `json:"minKubeVersion,omitempty"`
	// Version is the semantic version of the operator.
	Version version.OperatorVersion `json:"version,omitempty"`
	// Maturity is a rating of how mature the operator is.
	Maturity string `json:"maturity,omitempty"`
	// DisplayName is the human-readable name used to represent the operator in client displays.
	DisplayName string `json:"displayName"`
	// Description is a brief description of the operator's purpose.
	Description string `json:"description,omitempty"`
	// Keywords defines a set of keywords associated with the operator.
	Keywords []string `json:"keywords,omitempty"`
	// Maintainers defines a set of people and/or organizations responsible for maintaining the operator.
	Maintainers []operatorsv1alpha1.Maintainer `json:"maintainers,omitempty"`
	// Provider is a link to the site of the operator's provider.
	Provider operatorsv1alpha1.AppLink `json:"provider,omitempty"`
	// Links is a set of associated links.
	Links []operatorsv1alpha1.AppLink `json:"links,omitempty"`
	// Icon is the operator's base64 encoded icon.
	Icon []operatorsv1alpha1.Icon `json:"icon,omitempty"`
	// InstallModes specify supported installation types.
	InstallModes []operatorsv1alpha1.InstallMode `json:"installModes,omitempty"`
	// The name of a CSV this one replaces. Should match the `metadata.Name` field of the old CSV.
	Replaces string `json:"replaces,omitempty"`
	// Current condition of the ClusterServiceVersion
	Phase operatorsv1alpha1.ClusterServiceVersionPhase `json:"phase,omitempty"`
	// A human readable message indicating details about why the ClusterServiceVersion is in this condition.
	Message string `json:"message,omitempty"`
	// A brief CamelCase message indicating details about why the ClusterServiceVersion is in this state.
	// e.g. 'RequirementsNotMet'
	Reason operatorsv1alpha1.ConditionReason `json:"reason,omitempty"`

	// Fields projected from the referenced Subscription

	// CatalogSource is the name of the catalog the referenced CSV was resolved from (if any).
	CatalogSourceName string `json:"catalogSourceName,omitempty"`
	// CatalogSourceNamespace is the namespace of the catalog the referenced CSV was resolved from (if any).
	CatalogSourceNamespace string `json:"catalogSourceNamespace,omitempty"`
	// Package is the package the referenced CSV belongs to.
	Package string `json:"package,omitempty"`
	// Channel is the channel the referenced CSV was resolved from.
	Channel string `json:"channel,omitempty"`
}
