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
// Created on 28 May 2019
// @author: Roi Sucasas - ATOS
//

package impl

import (
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"log"
)

// getK8sDeployment gets the deployment info
func getK8sDeployment(clusterIndex int, namespace string, name string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] Getting deployment from K8s cluster ...")

	// HttpGET_String returns (map[string]interface{}, error)
	_, data, err := common.HTTPGETString(
		cfg.Config.Clusters[clusterIndex].KubernetesEndPoint+"/apis/apps/v1/namespaces/"+namespace+"/deployments/"+name,
		true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [getK8sDeployment] RESPONSE: " + data)
	return data, nil
}

/*
GetTaskAllInfo gets a deployment by id
*/
func GetTaskAllInfo(idTask string) (structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetTaskAllInfo] Getting deployment with id=" + idTask + " ...")

	// get task
	dbTask, err := common.ReadTaskValue(idTask)
	if err == nil {
		task := dbTask.TaskDefinition
		clusterIndex := common.GetClusterIndex(task.Cluster)
		namespace := task.Dock

		// get deployment information
		deploymentData, err := getK8sDeployment(clusterIndex, namespace, idTask)
		if err == nil {
			dbTask.Deployment = deploymentData
			return *dbTask, nil
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetTaskAllInfo] ERROR", err)
	return *dbTask, err
}

/*
GetConfig gets the Kubernetes service configuration
*/
func GetConfig() (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] Getting configuration from K8s cluster ...")

	_, data, err := common.HTTPGETString(cfg.Config.Clusters[0].KubernetesEndPoint+"/api", true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Operations [GetConfig] RESPONSE: " + data)
	return data, err
}
