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
	urls "atos/rotterdam/caas/adapters"
	adapt_common "atos/rotterdam/caas/adapters/common"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/imec/db"
	"errors"
	"log"
	"strconv"
	"strings"
)

// k8sRoute creates a route
func k8sRoute(namespace string, mainPort int, mainPortName string, task structs.CLASS_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [k8sRoute] Generating 'route' json ...")
	// TODO
	// routePrefix := task.Name + "-" + task.Id
	k8sRoute := common.StructNewRouteTemplate(task.ID, mainPort, mainPortName, cfg.Config.Clusters[0].HostIP) // returns *K8S_ROUTE

	strTxt, _ := common.CommRouteStructToString(*k8sRoute)
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [k8sRoute] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new route
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [k8sRoute] Creating a new route in K8s cluster ...")
	status, _, err := common.HTTPPOST(
		urls.GetPathOpenshiftRoutes(cluster, namespace),
		true,
		k8sRoute)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [k8sRoute] ERROR", err)
		return "", err
	}

	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [k8sRoute] RESPONSE: OK")
	return strconv.Itoa(status), nil
}

/*
DeployTask Deploy a task (k8s: deployment & service & volumes ...)
*/
func DeployTask(task structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] Deploying new task ...")

	clusterInfr, _ := imec_db.GetCluster(task.Cluster)
	clusterID := ""
	clusterHost := ""
	if clusterInfr != nil {
		clusterID = clusterInfr.ID
		clusterHost = clusterInfr.HostIP
	}
	namespace := task.Dock
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] cluster id = " + clusterID + ", dock = " + namespace + ", host = " + clusterHost + "")

	// 1. DEPLOYMENT /////
	status, err := adapt_common.K8sDeployment(namespace, task, clusterInfr, true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] ERROR (1)", err)
		return "", err
	} else if status == "200" || status == "201" {
		// 2. SERVICE /////
		status, mainPort, mainPortName, err := adapt_common.K8sService(namespace, task, clusterInfr, true)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] ERROR (2)", err)
			return "", err
		} else if status == "200" || status == "201" {
			// 3. ROUTE /////
			status, err := k8sRoute(namespace, mainPort, mainPortName, task, clusterInfr)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] ERROR (3)", err)
				return "", err
			} else if status == "200" || status == "201" {
				log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] Task deployed with success")
				// save to DB
				dbtask := &structs.DB_TASK{
					DbId:           structs.DB_TABLE_TASK,
					Id:             task.ID,
					Name:           task.Name,
					NameSpace:      namespace,
					Type:           structs.DB_TASK_TYPE_DEFAULT,
					ClusterId:      clusterID,
					AgreementId:    strings.Replace(task.ID, "-", "_", -1),
					Url:            "http://" + task.ID + "." + clusterHost + ".xip.io",
					Status:         "Deployed",
					Replicas:       task.Replicas,
					TaskDefinition: task}
				common.SetTaskValue(task.ID, *dbtask)

				return "ok", nil
			}
		}
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTask] ERROR (4)", err)
	return "", err
}

/*
DeployTaskCompss Deploy a COMPSs task (k8s: deployment & service & volumes ...)
*/
func DeployTaskCompss(task structs.CLASS_TASK) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] Deploying new task ...")

	clusterInfr, _ := imec_db.GetCluster(task.Cluster)
	clusterID := ""
	if clusterInfr != nil {
		clusterID = clusterInfr.ID
	}
	namespace := task.Dock
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] cluster id = " + clusterID + ", dock = " + namespace)

	// 1. DEPLOYMENT /////
	status, err := adapt_common.K8sDeployment(namespace, task, clusterInfr, true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] ERROR (1)", err)
		return "", err
	} else if status == "200" || status == "201" {
		log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] Task deployed with success")
		// save to DB
		dbtask := &structs.DB_TASK{
			DbId:           structs.DB_TABLE_TASK,
			Id:             task.ID,
			Name:           task.Name,
			NameSpace:      namespace,
			Type:           structs.DB_TASK_TYPE_COMPSS,
			ClusterId:      clusterID,
			AgreementId:    strings.Replace(task.ID, "-", "_", -1),
			Status:         "Deployed",
			Replicas:       task.Replicas,
			TaskDefinition: task}
		common.SetTaskValue(task.ID, *dbtask)

		go func() {
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] Starting background tasks...")
			adapt_common.CompssDeploymentBackgroundTasks(task, clusterInfr)
			log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] Background tasks completed")
		}()

		return "ok", nil
	}

	err = errors.New("Task creation failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Openshift > Deployment [DeployTaskCompss] ERROR (4)", err)
	return "", err
}
