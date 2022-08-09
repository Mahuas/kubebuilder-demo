## kubebuilder



```shell
kubebuilder init --domain kubebuilder.rym

kubebuilder create api --group autoscaler --version v1 --kind VerticalPodAutoscaler

kubebuilder create webhook --group autoscaler --version v1 --kind VerticalPodAutoscaler --defaulting --programmatic-validation

make install

docker login

make docker-build docker-push IMG=renyuanming/controller:latest

make deploy IMG=renyuanming/controller:latest
```



Controller.go

```golang
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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	autoscalerv1 "operator/client-go-demo/kubebuilder-demo/api/v1"
)

// VerticalPodAutoscalerReconciler reconciles a VerticalPodAutoscaler object
type VerticalPodAutoscalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=autoscaler.kubebuilder.rym,resources=verticalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=autoscaler.kubebuilder.rym,resources=verticalpodautoscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=autoscaler.kubebuilder.rym,resources=verticalpodautoscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VerticalPodAutoscaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *VerticalPodAutoscalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var vpa autoscalerv1.VerticalPodAutoscaler
	if err := r.Get(ctx, req.NamespacedName, &vpa); err != nil {
		log.Log.Error(err, "unable to fetch myapp")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("vpa is ", "vpa", vpa.Name)

	var vpas autoscalerv1.VerticalPodAutoscalerList
	if err := r.List(ctx, &vpas, client.InNamespace(req.Namespace)); err != nil {
		log.Log.Error(err, "unable to fetch myappList")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	for _, app := range vpas.Items {
		log.Log.Info("vpa is", "vpas", app.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VerticalPodAutoscalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autoscalerv1.VerticalPodAutoscaler{}).
		Complete(r)
}

```





types.go

```golang
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VerticalPodAutoscalerSpec defines the desired state of VerticalPodAutoscaler
type VerticalPodAutoscalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	DeploymentName string `json:"deploymentName"`

	//+kubebuilder:validation:Minimum=-1
	RequestCPU *int32 `json:"requestCPU"`

	//+kubebuilder:validation:Minimum=0
	RequestMemory *int32 `json:"requestMemory"`
}

// VerticalPodAutoscalerStatus defines the observed state of VerticalPodAutoscaler
type VerticalPodAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VerticalPodAutoscaler is the Schema for the verticalpodautoscalers API
type VerticalPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VerticalPodAutoscalerSpec   `json:"spec,omitempty"`
	Status VerticalPodAutoscalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VerticalPodAutoscalerList contains a list of VerticalPodAutoscaler
type VerticalPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VerticalPodAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VerticalPodAutoscaler{}, &VerticalPodAutoscalerList{})
}

```





Web hook

```golang
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

```

