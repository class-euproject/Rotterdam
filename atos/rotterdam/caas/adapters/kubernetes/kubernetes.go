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
// Created on 11 June 2019
// @author: Roi Sucasas - ATOS
//

package kubernetes

import (
	"atos/rotterdam/caas/adapters/kubernetes/impl"
	structs "atos/rotterdam/caas/common/structs"
)

/*
Kubernetes Adapter
*/
type KubernetesAdapter struct{}

/*
DeployTask Deploy a task (k8s: deployment & service & volumes ...)
*/
func (a KubernetesAdapter) DeployTask(task structs.CLASS_TASK) (string, error) {
	return impl.DeployTask(task)
}

/*
DeployTaskCompss Deploy a COMPSs task (k8s: deployment & service & volumes ...)
*/
func (a KubernetesAdapter) DeployTaskCompss(task structs.CLASS_TASK) (string, error) {
	return impl.DeployTaskCompss(task)
}

/*
GetTaskAllInfo Gets a task with all the deployment information
*/
func (a KubernetesAdapter) GetTaskAllInfo(idTask string) (structs.DB_TASK, error) {
	return impl.GetTaskAllInfo(idTask)
}

/*
GetConfig k8s configuration
*/
func (a KubernetesAdapter) GetConfig() (string, error) {
	return impl.GetConfig()
}

/*
ScaleUpDown Scale up task
*/
func (a KubernetesAdapter) ScaleUpDown(dbtask structs.DB_TASK, replicas int) (string, error) {
	return impl.ScaleUpDown(dbtask, replicas)
}

/*
RemoveTask Deletes a task
*/
func (a KubernetesAdapter) RemoveTask(dbtask structs.DB_TASK) (string, string, error) {
	return impl.RemoveTask(dbtask)
}
