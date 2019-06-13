// +build !ignore_autogenerated

/*
Copyright 2019 Red Hat, Inc.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	porcelain "github.com/operator-framework/operator-lifecycle-manager/pkg/porcelain-server/apis/porcelain"
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*InstalledOperator)(nil), (*porcelain.InstalledOperator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstalledOperator_To_porcelain_InstalledOperator(a.(*InstalledOperator), b.(*porcelain.InstalledOperator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*porcelain.InstalledOperator)(nil), (*InstalledOperator)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_porcelain_InstalledOperator_To_v1alpha1_InstalledOperator(a.(*porcelain.InstalledOperator), b.(*InstalledOperator), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InstalledOperatorList)(nil), (*porcelain.InstalledOperatorList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InstalledOperatorList_To_porcelain_InstalledOperatorList(a.(*InstalledOperatorList), b.(*porcelain.InstalledOperatorList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*porcelain.InstalledOperatorList)(nil), (*InstalledOperatorList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_porcelain_InstalledOperatorList_To_v1alpha1_InstalledOperatorList(a.(*porcelain.InstalledOperatorList), b.(*InstalledOperatorList), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_InstalledOperator_To_porcelain_InstalledOperator(in *InstalledOperator, out *porcelain.InstalledOperator, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.ClusterServiceVersionRef = (*v1.ObjectReference)(unsafe.Pointer(in.ClusterServiceVersionRef))
	out.SubscriptionRef = (*v1.ObjectReference)(unsafe.Pointer(in.SubscriptionRef))
	out.CustomResourceDefinitions = in.CustomResourceDefinitions
	out.APIServiceDefinitions = in.APIServiceDefinitions
	out.MinKubeVersion = in.MinKubeVersion
	out.Version = in.Version
	out.Maturity = in.Maturity
	out.DisplayName = in.DisplayName
	out.Description = in.Description
	out.Keywords = *(*[]string)(unsafe.Pointer(&in.Keywords))
	out.Maintainers = *(*[]operatorsv1alpha1.Maintainer)(unsafe.Pointer(&in.Maintainers))
	out.Provider = in.Provider
	out.Links = *(*[]operatorsv1alpha1.AppLink)(unsafe.Pointer(&in.Links))
	out.Icon = *(*[]operatorsv1alpha1.Icon)(unsafe.Pointer(&in.Icon))
	out.InstallModes = *(*[]operatorsv1alpha1.InstallMode)(unsafe.Pointer(&in.InstallModes))
	out.Replaces = in.Replaces
	out.Phase = operatorsv1alpha1.ClusterServiceVersionPhase(in.Phase)
	out.Message = in.Message
	out.Reason = operatorsv1alpha1.ConditionReason(in.Reason)
	out.CatalogSourceName = in.CatalogSourceName
	out.CatalogSourceNamespace = in.CatalogSourceNamespace
	out.Package = in.Package
	out.Channel = in.Channel
	return nil
}

// Convert_v1alpha1_InstalledOperator_To_porcelain_InstalledOperator is an autogenerated conversion function.
func Convert_v1alpha1_InstalledOperator_To_porcelain_InstalledOperator(in *InstalledOperator, out *porcelain.InstalledOperator, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstalledOperator_To_porcelain_InstalledOperator(in, out, s)
}

func autoConvert_porcelain_InstalledOperator_To_v1alpha1_InstalledOperator(in *porcelain.InstalledOperator, out *InstalledOperator, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.ClusterServiceVersionRef = (*v1.ObjectReference)(unsafe.Pointer(in.ClusterServiceVersionRef))
	out.SubscriptionRef = (*v1.ObjectReference)(unsafe.Pointer(in.SubscriptionRef))
	out.CustomResourceDefinitions = in.CustomResourceDefinitions
	out.APIServiceDefinitions = in.APIServiceDefinitions
	out.MinKubeVersion = in.MinKubeVersion
	out.Version = in.Version
	out.Maturity = in.Maturity
	out.DisplayName = in.DisplayName
	out.Description = in.Description
	out.Keywords = *(*[]string)(unsafe.Pointer(&in.Keywords))
	out.Maintainers = *(*[]operatorsv1alpha1.Maintainer)(unsafe.Pointer(&in.Maintainers))
	out.Provider = in.Provider
	out.Links = *(*[]operatorsv1alpha1.AppLink)(unsafe.Pointer(&in.Links))
	out.Icon = *(*[]operatorsv1alpha1.Icon)(unsafe.Pointer(&in.Icon))
	out.InstallModes = *(*[]operatorsv1alpha1.InstallMode)(unsafe.Pointer(&in.InstallModes))
	out.Replaces = in.Replaces
	out.Phase = operatorsv1alpha1.ClusterServiceVersionPhase(in.Phase)
	out.Message = in.Message
	out.Reason = operatorsv1alpha1.ConditionReason(in.Reason)
	out.CatalogSourceName = in.CatalogSourceName
	out.CatalogSourceNamespace = in.CatalogSourceNamespace
	out.Package = in.Package
	out.Channel = in.Channel
	return nil
}

// Convert_porcelain_InstalledOperator_To_v1alpha1_InstalledOperator is an autogenerated conversion function.
func Convert_porcelain_InstalledOperator_To_v1alpha1_InstalledOperator(in *porcelain.InstalledOperator, out *InstalledOperator, s conversion.Scope) error {
	return autoConvert_porcelain_InstalledOperator_To_v1alpha1_InstalledOperator(in, out, s)
}

func autoConvert_v1alpha1_InstalledOperatorList_To_porcelain_InstalledOperatorList(in *InstalledOperatorList, out *porcelain.InstalledOperatorList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]porcelain.InstalledOperator)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_InstalledOperatorList_To_porcelain_InstalledOperatorList is an autogenerated conversion function.
func Convert_v1alpha1_InstalledOperatorList_To_porcelain_InstalledOperatorList(in *InstalledOperatorList, out *porcelain.InstalledOperatorList, s conversion.Scope) error {
	return autoConvert_v1alpha1_InstalledOperatorList_To_porcelain_InstalledOperatorList(in, out, s)
}

func autoConvert_porcelain_InstalledOperatorList_To_v1alpha1_InstalledOperatorList(in *porcelain.InstalledOperatorList, out *InstalledOperatorList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]InstalledOperator)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_porcelain_InstalledOperatorList_To_v1alpha1_InstalledOperatorList is an autogenerated conversion function.
func Convert_porcelain_InstalledOperatorList_To_v1alpha1_InstalledOperatorList(in *porcelain.InstalledOperatorList, out *InstalledOperatorList, s conversion.Scope) error {
	return autoConvert_porcelain_InstalledOperatorList_To_v1alpha1_InstalledOperatorList(in, out, s)
}
