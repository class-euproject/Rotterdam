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
// Created on 06 March 2020
// @author: Roi Sucasas - ATOS
//

package common

import (
	urls "atos/rotterdam/caas/adapters"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	imec_db "atos/rotterdam/imec/db"
	"log"
	"strconv"
)

/*
DelK8sDeployment  k8s: deployment
*/
func DelK8sDeployment(namespace string, name string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, error) {
	// CALL to Kubernetes API to delete a deployment
	log.Println("Rotterdam > CAAS > Adapters > Common [delK8sDeployment] Deleting deployment from K8s cluster ...")

	status, _, err := common.HTTPDELETEStruct(urls.GetPathKubernetesDeleteDeployment(cluster, namespace, name), sec)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Common [delK8sDeployment] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Common [delK8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
DelK8sService  k8s: service
*/
func DelK8sService(namespace string, name string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, error) {
	// CALL to Kubernetes API to delete a service
	log.Println("Rotterdam > CAAS > Adapters > Common [delK8sService] !!!! Deleting service from K8s cluster ...")

	status, _, err := common.HTTPDELETEStruct(urls.GetPathKubernetesService(cluster, namespace, name), sec)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Common [delK8sService] ERROR", err)
		return strconv.Itoa(status), err
	}
	log.Println("Rotterdam > CAAS > Adapters > Common [delK8sService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
K8sDeployment k8s: deployment
*/
func K8sDeployment(namespace string, task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sDeployment] Generating 'deployment' json ...")
	k8sDepl := common.StructNewDeploymentTemplate(task, 1) // returns *K8S_DEPLOYMENT

	strTxt, _ := common.CommDeploymentStructToString(*k8sDepl)
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sDeployment] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sDeployment] Creating a new deployment in K8s cluster ...")
	status, _, err := common.HTTPPOSTStruct(
		urls.GetPathKubernetesCreateDeployment(cluster, namespace),
		sec,
		k8sDepl)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Common [K8sDeployment] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sDeployment] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
K8sService k8s: service
*/
func K8sService(namespace string, task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, sec bool) (string, int, string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sService] Generating 'service' json ...")
	k8sServ, mainPort, mainPortName := common.StructNewK8sServiceTemplate(task, cluster.HostIP) // returns *K8S_SERVICE

	strTxt, _ := common.CommServiceStructToString(*k8sServ)
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sService] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new service
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sService] Creating a new service in K8s cluster ...")
	status, _, err := common.HTTPPOSTStruct(
		urls.GetPathKubernetesCreateService(cluster, namespace),
		sec,
		k8sServ)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Common [K8sService] ERROR", err)
		return "", -1, "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Common [K8sService] RESPONSE: OK")

	return strconv.Itoa(status), mainPort, mainPortName, nil
}
