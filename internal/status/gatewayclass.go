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

package status

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayapi_v1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

// SetGatewayClassAccepted inserts or updates the Accepted condition
// for the provided GatewayClass.
func SetGatewayClassAccepted(ctx context.Context, cli client.Client, gc *gatewayapi_v1beta1.GatewayClass, accepted bool) *gatewayapi_v1beta1.GatewayClass {
	gc.Status.Conditions = mergeConditions(gc.Status.Conditions, computeGatewayClassAcceptedCondition(gc, accepted))
	return gc
}
