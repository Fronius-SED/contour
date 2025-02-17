// Copyright Project Contour Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objects

import (
	"context"

	"github.com/projectcontour/contour/internal/provisioner/labels"
	"github.com/projectcontour/contour/internal/provisioner/model"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// XDSPort is the network port number of Contour's xDS service.
	XDSPort = int32(8001)
	// EnvoyInsecureContainerPort is the network port number of Envoy's insecure listener.
	EnvoyInsecureContainerPort = int32(8080)
	// EnvoySecureContainerPort is the network port number of Envoy's secure listener.
	EnvoySecureContainerPort = int32(8443)
)

// NewUnprivilegedPodSecurity makes a a non-root PodSecurityContext object
// using 65534 as the user and group ID.
func NewUnprivilegedPodSecurity() *corev1.PodSecurityContext {
	user := int64(65534)
	group := int64(65534)
	nonRoot := true
	return &corev1.PodSecurityContext{
		RunAsUser:    &user,
		RunAsGroup:   &group,
		RunAsNonRoot: &nonRoot,
	}
}

// objectGetter gets an object given a namespace and name.
type objectGetter func(ctx context.Context, cli client.Client, namespace, name string) (client.Object, error)

// EnsureObjectDeleted ensures that an object with the given namespace and name is deleted
// if Contour owner labels exist.
func EnsureObjectDeleted(ctx context.Context, cli client.Client, contour *model.Contour, name string, getter objectGetter) error {
	obj, err := getter(ctx, cli, contour.Namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if !labels.Exist(obj, model.OwnerLabels(contour)) {
		return nil
	}
	if err = cli.Delete(ctx, obj); err == nil || errors.IsNotFound(err) {
		return nil
	}
	return err
}
