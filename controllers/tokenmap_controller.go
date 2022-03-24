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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	api "k8ssandra.io/k8ssandra-operator/api/v1alpha1"
)

// TokenmapReconciler reconciles a Tokenmap object
type TokenmapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8ssandra.io,resources=tokenmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8ssandra.io,resources=tokenmaps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8ssandra.io,resources=tokenmaps/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tokenmap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *TokenmapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var tokens api.Tokenmap
	if err := r.Get(ctx, req.NamespacedName, &tokens); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if tokens.Spec.URL != "" {
		// Fetch the JSON file
		resp, err := http.Get(tokens.Spec.URL)
		if err != nil {
			return ctrl.Result{}, err
		}

		jsonFile, err := io.ReadAll(resp.Body)
		if err != nil {
			return ctrl.Result{}, err
		}

		defer resp.Body.Close()

		logger.Info("Received file", "jsonFile", string(jsonFile))

		// Parse JSON file
		tokenMap := make(map[string]interface{})
		err = json.Unmarshal(jsonFile, &tokenMap)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Node count = len(tokenMap)
		node := 0
		for k := range tokenMap {
			text := "printf 'Replacing %s with %s %s' $ORIGINAL_POD_KEY $(hostname) $(hostname -i) "
			command := []string{"sh", "-c", text}
			pod := corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("cluster0-node-%d", node),
					Namespace: req.Namespace,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "busybee",
							Image:   "busybox:1.33.1",
							Command: command,
							Env: []corev1.EnvVar{
								{
									Name:  "ORIGINAL_POD_KEY",
									Value: k,
								},
							},
						},
					},
				},
			}

			err = r.Create(ctx, &pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			node++
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TokenmapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.Tokenmap{}).
		Complete(r)
}
