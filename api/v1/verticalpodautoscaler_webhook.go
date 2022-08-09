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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var verticalpodautoscalerlog = logf.Log.WithName("verticalpodautoscaler-resource")

func (r *VerticalPodAutoscaler) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-autoscaler-kubebuilder-rym-v1-verticalpodautoscaler,mutating=true,failurePolicy=fail,sideEffects=None,groups=autoscaler.kubebuilder.rym,resources=verticalpodautoscalers,verbs=create;update,versions=v1,name=mverticalpodautoscaler.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &VerticalPodAutoscaler{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *VerticalPodAutoscaler) Default() {
	verticalpodautoscalerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.RequestMemory == nil {
		r.Spec.RequestMemory = new(int32)
		*r.Spec.RequestMemory = -1
	}

	if r.Spec.RequestCPU == nil {
		r.Spec.RequestCPU = new(int32)
		*r.Spec.RequestCPU = -1
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-autoscaler-kubebuilder-rym-v1-verticalpodautoscaler,mutating=false,failurePolicy=fail,sideEffects=None,groups=autoscaler.kubebuilder.rym,resources=verticalpodautoscalers,verbs=create;update,versions=v1,name=vverticalpodautoscaler.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &VerticalPodAutoscaler{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *VerticalPodAutoscaler) ValidateCreate() error {
	verticalpodautoscalerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateVerticalPodAutoscaler()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *VerticalPodAutoscaler) ValidateUpdate(old runtime.Object) error {
	verticalpodautoscalerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateVerticalPodAutoscaler()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *VerticalPodAutoscaler) ValidateDelete() error {
	verticalpodautoscalerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return r.validateVerticalPodAutoscaler()
}

func (r *VerticalPodAutoscaler) validateVerticalPodAutoscaler() error {
	var allErrs field.ErrorList
	if err := r.validateVerticalPodAutoscalerSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "autoscaler.kubebuilder.rym", Kind: "VerticalPodAutoscaler"},
		r.Name, allErrs)
}

func (r *VerticalPodAutoscaler) validateVerticalPodAutoscalerSpec() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	return validateDeploymentName(
		r.Spec.DeploymentName,
		field.NewPath("spec").Child("deploymentName"))
}

func validateDeploymentName(deploymentName string, fldPath *field.Path) *field.Error {
	verticalpodautoscalerlog.Info("Debug", "deploymentName", deploymentName)
	if deploymentName != "bilibili" {
		return field.Invalid(fldPath, deploymentName, "deploymentName is nil")
	}
	return nil
}
