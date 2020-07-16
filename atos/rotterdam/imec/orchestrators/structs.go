//
// ROTTERDAM application
// CLASS Project: https://class-project.eu/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Created on 28 Feb 2020
// @author: Roi Sucasas - ATOS
//

package orchestrators

/*
ORCHESTRATOR represents an orchestrator deployed in an infrastructure
*/
type ORCHESTRATOR struct {
	ID                            string `json:"id,omitempty"`
	InfrID                        string `json:"infrID,omitempty"`
	Type                          string `json:"type,omitempty"` // EDGE: 'microk8s', 'k3s'; CLUSTER: 'openshift', 'kubernetes'
	OpenshiftOauthToken           string `json:"openshiftOauthToken,omitempty"`
	KubernetesEndPoint            string `json:"kubernetesEndPoint,omitempty"`
	OpenshiftEndPoint             string `json:"openshiftEndPoint,omitempty"`
	SLALiteEndPoint               string `json:"slaliteEndPoint,omitempty"`
	PrometheusPushgatewayEndPoint string `json:"prometheusPushgatewayEndPoint,omitempty"`
	Status                        string `json:"status,omitempty"`
}
