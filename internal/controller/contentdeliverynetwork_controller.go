/*
Copyright 2024.

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

package controller

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cdnv3 "github.com/benauro/kube-cdn/api/v3"
)

// ContentDeliveryNetworkReconciler reconciles a ContentDeliveryNetwork object
type ContentDeliveryNetworkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cdn.benauro.gg,resources=contentdeliverynetworks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cdn.benauro.gg,resources=contentdeliverynetworks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cdn.benauro.gg,resources=contentdeliverynetworks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ContentDeliveryNetwork object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ContentDeliveryNetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var cdn cdnv3.ContentDeliveryNetwork
	if err := r.Get(ctx, req.NamespacedName, &cdn); err != nil {
		logger.Error(err, "Unable to fetch ContentDeliveryNetwork")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Update status
	cdn.Status.State = "Ready"
	cdn.Status.Nodes = []string{"10.0.0.1", "10.0.0.2"}
	cdn.Status.LastUpdated = metav1.Now()

	if err := r.Status().Update(ctx, &cdn); err != nil {
		logger.Error(err, "Unable to update ContentDeliveryNetwork status")
		return ctrl.Result{}, err
	}

	// Reconcile again every 5 minutes
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContentDeliveryNetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cdnv3.ContentDeliveryNetwork{}).
		Complete(r)
}