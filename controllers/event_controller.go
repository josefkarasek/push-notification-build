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
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	githubv1alpha1 "github.com/josefkarasek/push-notification-build/api/v1alpha1"
)

// EventReconciler reconciles a Event object
type EventReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	EventCh   <-chan event.GenericEvent
	Namespace string

	cache cache.Store
}

//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=github.apps.openshift.io,resources=events/finalizers,verbs=update

func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	obj, exists, err := r.cache.GetByKey(req.String())
	if err != nil {
		log.Error(err, "Failed to read Obj from cache", "obj", obj)
	}
	if !exists {
		log.Error(err, "Requested resource not found in cache")
		return ctrl.Result{}, nil
	}

	ge, _ := obj.(event.GenericEvent)
	event, _ := ge.Object.(*githubv1alpha1.Event)

	d, _ := yaml.Marshal(event)
	fmt.Println(string(d))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.cache = cache.NewTTLStore(func(o interface{}) (string, error) {
		evt, ok := o.(event.GenericEvent)
		if !ok {
			return "", fmt.Errorf("received object is not event.GenericEvent")
		}
		return fmt.Sprintf("%s/%s", evt.Object.GetNamespace(), evt.Object.GetName()), nil
	}, 1*time.Minute)

	return ctrl.NewControllerManagedBy(mgr).
		For(&githubv1alpha1.Event{}).
		Watches(&source.Channel{Source: r.EventCh}, &EnqueueAndCacheRequest{cache: r.cache}).
		Complete(r)
}

type EnqueueAndCacheRequest struct {
	cache cache.Store
}

func (e *EnqueueAndCacheRequest) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {}
func (e *EnqueueAndCacheRequest) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {}
func (e *EnqueueAndCacheRequest) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {}
func (e *EnqueueAndCacheRequest) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		return
	}
	if err := e.cache.Add(evt); err != nil {
		return
	}
	q.Add(ctrl.Request{NamespacedName: types.NamespacedName{
		Name:      evt.Object.GetName(),
		Namespace: evt.Object.GetNamespace(),
	}})
}
