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

package impl

import (
	urls "atos/rotterdam/caas/adapters"
	adapt_common "atos/rotterdam/caas/adapters/common"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	imec_db "atos/rotterdam/imec/db"
	"log"
	"strconv"
	"time"
)

/*
Get Pods from current task. Retrieved Pods = Expected Pods
*/
// checkPodsFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func checkPodsFromTask(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, id string, expectedReplicas int) (string, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkPodsFromTask] Getting pods ...")

	// get pods from task
	_, result, err := common.HTTPGETStruct(
		urls.GetPathKubernetesPodsApp(cluster, namespace, id),
		//cfg.Config.Clusters[clusterIndex].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/pods?labelSelector=app="+id,
		true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkPodsFromTask] ERROR", err)
		return "error", nil, err
	}

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkPodsFromTask] Retrieved pods = " + strconv.Itoa(len(items)))
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkPodsFromTask] Expected pods = " + strconv.Itoa(expectedReplicas))

	if len(items) == expectedReplicas {
		return "ready", result, err
	}

	return "not-ready", result, err
}

/*
Get Pods from current task. Retrieved Pods = Expected Pods
*/
// curl -X GET -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9"
// http://192.168.7.28:8001/api/v1/namespaces/class/services/serv-g7bv9fu6d0rtaius7okjnk1580403141125383512-5c465cbb4b-xqdl8
func getServiceFromPod(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, idPod string) (string, adapt_common.SERV_POD, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getServiceFromPod] Getting service from Pod [" + idPod + "] ...")

	// get pods from task
	_, result, err := common.HTTPGETString(
		urls.GetPathKubernetesService(cluster, namespace, idPod),
		//cfg.Config.Clusters[clusterIndex].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/services/serv-"+idPod,
		true)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getServiceFromPod] ERROR", err)
		return "error", adapt_common.SERV_POD{}, err
	}

	// metadata
	sp, err := adapt_common.StringToServPodStruct(result)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getServiceFromPod] ERROR", err)
		return "error", adapt_common.SERV_POD{}, err
	}

	if sp.Metadata.Name == "" {
		return "not-ready", adapt_common.SERV_POD{}, err
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getServiceFromPod] Retrieved metadata.name = " + sp.Metadata.Name)

	return "ready", *sp, nil
}

// getPodsInfoFromTask
// curl -X GET -H "Authorization: Bearer Mo6CxHG2ZjZCqh-moIK8fjSorm6aennoAX8Q3xTEFXQ"
// http://192.168.7.28:8001/api/v1/namespaces/class/pods?labelSelector=app=nginx-app
func getPodsInfoFromTask(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, result map[string]interface{}, appPort int, appProtocol string) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

	// final res
	var lres []structs.DB_TASK_POD

	// iterate json (response)
	for _, item := range items { // "items" element ==> PODS INFO
		/* DB_TASK_POD:
		{
			Name       string `json:"name,omitempty"`
			IP         string `json:"ip,omitempty"`         // IP accessed by external apps
			HostIP     string `json:"hostIp,omitempty"`     // node IP
			PodIP      string `json:"podIp,omitempty"`      // internal IP created by Kubernetes / Openshift
			Status     string `json:"status,omitempty"`     // running, unknown
			Port       int    `json:"port,omitempty"`       // port exposed in Kubernetes / Openshift
			TargetPort int    `json:"targetPort,omitempty"` // application port
		}
		*/
		podData := &structs.DB_TASK_POD{}
		podData.IP = urls.GetHostIP(cluster) // cfg.Config.Clusters[clusterIndex].HostIP
		podData.Port = adapt_common.NewRPort()
		podData.TargetPort = appPort
		podData.Protocol = appProtocol

		// "metadata" element
		for key1, value1 := range item.(map[string]interface{}) {
			if key1 == "metadata" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "name" {
						podData.Name = value2.(string) // Pod name
					}
				}
			} else if key1 == "status" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "podIP" {
						podData.PodIP = value2.(string) // IP address from pod
					} else if key2 == "hostIP" {
						podData.HostIP = urls.GetHostIP(cluster) // value2.(string) // IP address from node
					} else if key2 == "phase" {
						podData.Status = value2.(string) // Status / Phase from pod
					}
				}
			}
		}

		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask] POD " + podData.Name)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask]    - Status: " + podData.Status)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [getPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

// patchPods: k8s: patch pods names
func patchPods(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, podName string) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [patchPods] Generating 'pod patch' json ...")
	lK8sPatchPod := common.StructNewPodPatch(podName) // returns []structs.K8S_POD_PATCH_LINE

	strTxt, _ := common.CommPatchPodsListToString(lK8sPatchPod)
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [patchPods] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new deployment
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [patchPods] Patching pod " + podName + " ...")
	status, _, err := common.HTTPPATCHStruct(
		urls.GetPathKubernetesPod(cluster, namespace, podName),
		//cfg.Config.Clusters[clusterIndex].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/pods/"+podName,
		true,
		lK8sPatchPod)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [patchPods] ERROR", err)
		return "error", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [patchPods] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

// podService: k8s: service
func podService(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, namespace string, pod structs.DB_TASK_POD) (string, error) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [podService] Generating 'service' json ...")
	k8sServ := common.StructNewPodServiceTemplate(urls.GetHostIP(cluster), pod) // returns *K8S_SERVICE

	strTxt, _ := common.CommServiceStructToString(*k8sServ)
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [podService] [" + strTxt + "]")

	// CALL to Kubernetes API to launch a new service
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [podService] Creating a new service in K8s cluster ...")
	status, _, err := common.HTTPPOST(
		urls.GetPathKubernetesCreateService(cluster, namespace),
		//cfg.Config.Clusters[clusterIndex].KubernetesEndPoint+"/api/v1/namespaces/"+namespace+"/services",
		true,
		k8sServ)
	if err != nil {
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [podService] ERROR", err)
		return "", err
	}
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [podService] RESPONSE: OK")

	return strconv.Itoa(status), nil
}

/*
Deployment of a COMPSs task
*/
func compssDeploymentBackgroundTasks(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, task structs.CLASS_TASK) {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Executing background tasks ...")

	time.Sleep(10 * time.Second)
	mainPort, mainProtocol := adapt_common.GetMainPort(task)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Checking status of task pods ...")
		str, result, err := checkPodsFromTask(cluster, task.Dock, task.ID, task.Replicas)

		if err == nil && str == "ready" { // Expected Pods = Retrieved Pods
			// get info from pods
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Getting all pods info (ids, names, IPs, ports) ...")
			lTaskPods := getPodsInfoFromTask(cluster, result, mainPort, mainProtocol)

			dbTask, err := common.ReadTaskValue(task.ID)
			if err == nil {
				// update task with pods info
				dbTask.Pods = lTaskPods
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks]Saving info to db ...")
				err = common.SetTaskValue(task.ID, *dbTask)

				if err == nil {
					// update pods: pod-name (labels)
					log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Updating pods (/metadata/labels/pod-name) / Creating and exposing pods services ...")
					for _, pod := range lTaskPods {
						log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Updating pod " + pod.Name)
						status, err := patchPods(cluster, task.Dock, pod.Name)
						// TODO check errors
						log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] status: " + status)

						if err == nil {
							log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Creating and exposing pod's service ...")

							_, _ = podService(cluster, task.Dock, pod)
						}
					}

					log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] Finishing background process ...")

					break
				}
			}
		} else if err == nil {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] str = " + str)
		} else {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssDeploymentBackgroundTasks] ERROR (1) ", err)
		}

		time.Sleep(30 * time.Second)
	}
}

// checkCurrentPodsInfoFromTask
func checkCurrentPodsInfoFromTask(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, result map[string]interface{}, appPort int, appProtocol string, dbTask structs.DB_TASK) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] Getting info from pods ...")

	// items
	items := result["items"].([]interface{})
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] Total pods = " + strconv.Itoa(len(items)))

	// final res
	var lres []structs.DB_TASK_POD

	// iterate json (response)
	for _, item := range items { // "items" element
		/* DB_TASK_POD:
		{
			Name       string `json:"name,omitempty"`
			IP         string `json:"ip,omitempty"`         // IP accessed by external apps
			HostIP     string `json:"hostIp,omitempty"`     // node IP
			PodIP      string `json:"podIp,omitempty"`      // internal IP created by Kubernetes / Openshift
			Status     string `json:"status,omitempty"`     // running, unknown
			Port       int    `json:"port,omitempty"`       // port exposed in Kubernetes / Openshift
			TargetPort int    `json:"targetPort,omitempty"` // application port
		}
		*/
		podData := &structs.DB_TASK_POD{}
		podData.IP = urls.GetHostIP(cluster) //cfg.Config.Clusters[clusterIndex].HostIP
		podData.TargetPort = appPort
		podData.Protocol = appProtocol

		// "metadata" element
		for key1, value1 := range item.(map[string]interface{}) {
			if key1 == "metadata" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "name" {
						podData.Name = value2.(string) // Pod name
					}
				}
			} else if key1 == "status" {
				for key2, value2 := range value1.(map[string]interface{}) {
					if key2 == "podIP" {
						podData.PodIP = value2.(string) // IP address from pod
					} else if key2 == "hostIP" {
						podData.HostIP = urls.GetHostIP(cluster) // value2.(string) // IP address from node
					} else if key2 == "phase" {
						podData.Status = value2.(string) // Status / Phase from pod
					}
				}
			}
		}

		// GET SERVICE, if exists
		ready, srvPod, err := getServiceFromPod(cluster, dbTask.TaskDefinition.Dock, podData.Name)
		if err == nil && ready == "ready" {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] POD " + podData.Name + " already exists")
			podData.Port = srvPod.Spec.Ports[0].Port
		} else {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] Generating new port for POD " + podData.Name)
			podData.Port = adapt_common.NewRPort()
		}

		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] POD " + podData.Name)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask]    - PodIP: " + podData.PodIP)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask]    - HostIP: " + podData.HostIP)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask]    - Status: " + podData.Status)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask]    - Port: " + strconv.Itoa(podData.Port))
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask]    - TargetPort: " + strconv.Itoa(podData.TargetPort))

		lres = append(lres, *podData)
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [checkCurrentPodsInfoFromTask] Total pods response = " + strconv.Itoa(len(lres)))

	return lres
}

// compssScalingOutBackgroundTasks Expose new services and save task
func compssScalingOutBackgroundTasks(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, dbTask structs.DB_TASK, items map[string]interface{}) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Exposing new services ...")

	// get main port and protocol
	mainPort, mainProtocol := adapt_common.GetMainPort(dbTask.TaskDefinition)

	// get list of pods - all pods
	lTaskPods := checkCurrentPodsInfoFromTask(cluster, items, mainPort, mainProtocol, dbTask)

	// add new services
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Updating new pods (/metadata/labels/pod-name) / Creating and exposing new pods services ...")
	for _, pod := range lTaskPods { // iterate all pods (including new pods)
		notfound := true
		for _, podold := range dbTask.Pods { // check if pod is new or old
			if pod.Name == podold.Name {
				notfound = false
				break
			}
		}

		if notfound == true {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Updating pod " + pod.Name)
			status, err := patchPods(cluster, dbTask.TaskDefinition.Dock, pod.Name)
			// TODO check errors
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] status: " + status)

			if err == nil {
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Creating and exposing pod's service ...")
				_, _ = podService(cluster, dbTask.TaskDefinition.Dock, pod)
			}
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Update finished")
	return lTaskPods
}

// compssScalingInBackgroundTasks Removes unused services and save task
func compssScalingInBackgroundTasks(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, dbTask structs.DB_TASK, items map[string]interface{}) []structs.DB_TASK_POD {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingInBackgroundTasks] Removing unused services ...")

	// get main port and protocol
	mainPort, mainProtocol := adapt_common.GetMainPort(dbTask.TaskDefinition)

	// get list of pods
	lTaskPods := checkCurrentPodsInfoFromTask(cluster, items, mainPort, mainProtocol, dbTask)

	// remove 'old' services
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Removing unused services ...")
	for _, podold := range dbTask.Pods {
		notfound := true
		for _, pod := range lTaskPods {
			if pod.Name == podold.Name {
				notfound = false
				break
			}
		}

		// remove service if old pod was not found in the list of current pods
		if notfound == true {
			log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Removing unused service of old pod [" + podold.Name + "] ...")
			_, _ = adapt_common.DelK8sService(dbTask.TaskDefinition.Dock, podold.Name, cluster, false)
		}
	}

	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [compssScalingOutBackgroundTasks] Services update finished")
	return lTaskPods
}

/*
CompssScalingUpdateServices Updates (remove unused / add new) services used to expose pods
Replicas are scaled in / out before calling this method
*/
func CompssScalingUpdateServices(cluster *imec_db.DB_INFRASTRUCTURE_CLUSTER, dbTask structs.DB_TASK, newReplicasValue int) structs.DB_TASK {
	log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [CompssScalingUpdateServices] Executing background tasks ...")

	time.Sleep(20 * time.Second)

	for i := 0; i < 10; i++ {
		// check if pods are running
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [CompssScalingUpdateServices] Checking status of task pods [new_replicas_value=" + strconv.Itoa(newReplicasValue) + "] ...")
		str, items, err := checkPodsFromTask(cluster, dbTask.TaskDefinition.Dock, dbTask.TaskDefinition.ID, newReplicasValue)
		if err == nil && str == "ready" {
			if dbTask.Replicas < newReplicasValue {
				dbTask.Pods = compssScalingOutBackgroundTasks(cluster, dbTask, items)
				break
			} else if dbTask.Replicas > newReplicasValue {
				dbTask.Pods = compssScalingInBackgroundTasks(cluster, dbTask, items)
				break
			} else {
				// don't update the task
				log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [CompssScalingUpdateServices] WARNING type of task is not defined: " + dbTask.Type)
				return dbTask
			}
		}

		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [CompssScalingUpdateServices] Result: " + str)
		log.Println("Rotterdam > CAAS > Adapters > Kubernetes > COMPSs [CompssScalingUpdateServices] Trying again in 20s ...")
		time.Sleep(20 * time.Second)
	}

	dbTask.TaskDefinition.Replicas = newReplicasValue
	dbTask.Replicas = newReplicasValue
	return dbTask
}
