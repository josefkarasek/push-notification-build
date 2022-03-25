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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"

	githubv1alpha1 "github.com/josefkarasek/push-notification-build/api/v1alpha1"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	EventCh chan event.GenericEvent
}

//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=apps/finalizers,verbs=update

// Reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("Beginning reconciliation")
	// step 1: fetch App
	// step 2: fetch referenced secret
	// step 3: adopt secret through OwnerReference; requeue (if not already done)
	// step 4: update Key in Secret Manager
	// step 4: update Status

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&githubv1alpha1.App{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
