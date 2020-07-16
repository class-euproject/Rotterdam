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

package impl

import (
	adapt_common "atos/rotterdam/caas/adapters/common"
	structs "atos/rotterdam/caas/common/structs"
	imec_db "atos/rotterdam/imec/db"
	"errors"
	"log"
)

// removeDefaultTask: Deletes a task
func removeDefaultTask(namespace string, name string, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] Deleting task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := adapt_common.DelK8sDeployment(namespace, name, cluster, false)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICE /////
		status, err := adapt_common.DelK8sService(namespace, name, cluster, false)
		if err != nil {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (2)", err)
			return "", err
		} else if status == "200" {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] Task removed with success")
			return "ok", nil
		}
	}

	err = errors.New("Task termination failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeDefaultTask] ERROR (4)", err)
	return "", err
}

// removeCOMPSsTask: Deletes a COMPSs task
func removeCOMPSsTask(namespace string, name string, dbtask structs.DB_TASK, cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] Deleting COMPSs task [" + name + "] from [" + namespace + "] ...")

	// 1. DEPLOYMENT /////
	status, err := adapt_common.DelK8sDeployment(namespace, name, cluster, false)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (1)", err)
		return "", err
	} else if status == "200" {
		// 2. SERVICES /////
		for _, pod := range dbtask.Pods {
			// DB_TASK_POD
			status, err := adapt_common.DelK8sService(namespace, pod.Name, cluster, false)
			if err != nil {
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (2)", err)
				return "", err
			} else if status == "200" {
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] Task removed with success")
				return "ok", nil
			}
		}
	}

	err = errors.New("Task termination failed. status = [" + status + "]")
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [removeCOMPSsTask] ERROR (4)", err)
	return "", err
}

/*
RemoveTask Deletes a task
*/
func RemoveTask(dbTask structs.DB_TASK) (string, string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] Deleting task [" + dbTask.Id + "] ...")

	clusterInfr, _ := imec_db.GetCluster(dbTask.ClusterId)
	task := dbTask.TaskDefinition
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] cluster id = " + dbTask.ClusterId + ", dock = " + task.Dock + "")

	// remove task
	if dbTask.Type == structs.DB_TASK_TYPE_DEFAULT {
		res, err := removeDefaultTask(task.Dock, dbTask.Id, clusterInfr)
		return res, dbTask.AgreementId, err
	} else if dbTask.Type == structs.DB_TASK_TYPE_COMPSS {
		res, err := removeCOMPSsTask(task.Dock, dbTask.Id, dbTask, clusterInfr)
		return res, dbTask.AgreementId, err
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > Termination [RemoveTask] WARNING type of task is not defined: " + dbTask.Type)
	return "", "", nil
}
