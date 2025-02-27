/*
Copyright 2022.

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

package v1beta1

import (
	cinderv1 "github.com/openstack-k8s-operators/cinder-operator/api/v1beta1"
	glancev1 "github.com/openstack-k8s-operators/glance-operator/api/v1beta1"
	keystonev1 "github.com/openstack-k8s-operators/keystone-operator/api/v1beta1"
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	mariadbv1 "github.com/openstack-k8s-operators/mariadb-operator/api/v1beta1"
	placementv1 "github.com/openstack-k8s-operators/placement-operator/api/v1beta1"
	rabbitmqv1 "github.com/rabbitmq/cluster-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OpenStackControlPlaneSpec defines the desired state of OpenStackControlPlane
type OpenStackControlPlaneSpec struct {

	// +kubebuilder:validation:Required
	// Secret - FIXME: make this optional
	Secret string `json:"secret,omitempty"`

	// +kubebuilder:validation:Required
	// StorageClass -
	StorageClass string `json:"storageClass,omitempty"`

	// +kubebuilder:validation:Optional
	// NodeSelector to target subset of worker nodes running control plane services (currently only applies to KeystoneAPI and PlacementAPI)
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	// KeystoneTemplate - Overrides to use when creating the Keystone API Service
	KeystoneTemplate keystonev1.KeystoneAPISpec `json:"keystoneTemplate,omitempty"`

	// +kubebuilder:validation:Optional
	// PlacementTemplate - Overrides to use when creating the Placement API
	PlacementTemplate placementv1.PlacementAPISpec `json:"placementTemplate,omitempty"`

	// +kubebuilder:validation:Optional
	// GlanceTemplate - Overrides to use when creating the Glance Service
	GlanceTemplate glancev1.GlanceSpec `json:"glanceTemplate,omitempty"`

	// +kubebuilder:validation:Optional
	// CinderTemplate - Overrides to use when creating Cinder Resources
	CinderTemplate cinderv1.CinderSpec `json:"cinderTemplate,omitempty"`

	// +kubebuilder:validation:Optional
	// MariadbTemplate - Overrides to use when creating the MariaDB API Service
	MariadbTemplate mariadbv1.MariaDBSpec `json:"mariadbTemplate,omitempty"`

	// +kubebuilder:validation:Optional
	// RabbitmqTemplate - Overrides to use when creating the Rabbitmq cluster
	RabbitmqTemplate rabbitmqv1.RabbitmqClusterSpec `json:"rabbitmqTemplate,omitempty"`
}

// OpenStackControlPlaneStatus defines the observed state of OpenStackControlPlane
type OpenStackControlPlaneStatus struct {
	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +operator-sdk:csv:customresourcedefinitions:displayName="OpenStack ControlPlane"
// +kubebuilder:resource:shortName=osctlplane;osctlplanes
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].status",description="Status"
//+kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[0].message",description="Message"

// OpenStackControlPlane is the Schema for the openstackcontrolplanes API
type OpenStackControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackControlPlaneSpec   `json:"spec,omitempty"`
	Status OpenStackControlPlaneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenStackControlPlaneList contains a list of OpenStackControlPlane
type OpenStackControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackControlPlane{}, &OpenStackControlPlaneList{})
}

// IsReady - returns true if service is ready to serve requests
func (instance OpenStackControlPlane) IsReady() bool {
	return instance.Status.Conditions.IsTrue(OpenStackControlPlaneRabbitMQReadyCondition) &&
		instance.Status.Conditions.IsTrue(OpenStackControlPlaneMariaDBReadyCondition) &&
		instance.Status.Conditions.IsTrue(OpenStackControlPlaneKeystoneAPIReadyCondition) &&
		instance.Status.Conditions.IsTrue(OpenStackControlPlanePlacementAPIReadyCondition) &&
		instance.Status.Conditions.IsTrue(OpenStackControlPlaneGlanceReadyCondition)
	    // TODO add once rabbitmq transportURL is integrated with Cinder:instance.Status.Conditions.IsTrue(OpenStackControlPlaneCinderReadyCondition)
}
