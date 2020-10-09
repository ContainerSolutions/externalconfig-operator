/*


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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	storev1alpha1 "github.com/containersolutions/externalsecret-operator/apis/store/v1alpha1"
)

// SecretStoreReconciler reconciles a SecretStore object
type SecretStoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=store.externalsecret-operator.container-solutions.com,resources=secretstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=store.externalsecret-operator.container-solutions.com,resources=secretstores/status,verbs=get;update;patch

func (r *SecretStoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("secretstore", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *SecretStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storev1alpha1.SecretStore{}).
		Complete(r)
}
