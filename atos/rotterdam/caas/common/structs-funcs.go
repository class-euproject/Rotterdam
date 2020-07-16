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

package common

import (
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////
// CLASS TASK (JSON) DEFINITION

/*
StructCheckClassTask Checks if CLASS struct is valid (from json)
*/
func StructCheckClassTask(decoder *json.Decoder) (*structs.CLASS_TASK, error) {
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassTask] Checking json object ...")

	//decoder := json.NewDecoder(req.Body)
	var t structs.CLASS_TASK
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassTask] ERROR (1)", err)
		return nil, err
	}

	if len(t.Containers) == 0 {
		log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassTask] ERROR (2) 'containers' is not defined")
		return nil, errors.New("'containers' is not defined")
	}

	tStr, err := CommClassStructToString(t)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassTask] Parsed object (string): " + tStr)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassTask] Sending parsed object ...")

	return &t, nil
}

/*
StructCheckClassCOMPSsTask Checks if CLASS struct is valid (from json)
*/
func StructCheckClassCOMPSsTask(decoder *json.Decoder) (*structs.CLASS_COMPSS_TASK, error) {
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassCOMPSsTask] Checking json object ...")

	//decoder := json.NewDecoder(req.Body)
	var t structs.CLASS_COMPSS_TASK
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassCOMPSsTask] ERROR (1)", err)
		return nil, err
	}

	tStr, err := CommClassCOMPSsStructToString(t)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassCOMPSsTask] Parsed object (string): " + tStr)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClassCOMPSsTask] Sending parsed object ...")

	return &t, nil
}

///////////////////////////////////////////////////////////////////////////////
// K8S_DEPLOYMENT

/*
StructNewDeploymentTemplate creates a new K8s Deployment json
*/
func StructNewDeploymentTemplate(task structs.CLASS_TASK, replicas int) *structs.K8S_DEPLOYMENT {
	var jsonDeployment *structs.K8S_DEPLOYMENT
	jsonDeployment = new(structs.K8S_DEPLOYMENT)

	jsonDeployment.ApiVersion = "apps/v1"
	jsonDeployment.Kind = "Deployment"
	jsonDeployment.Metadata.Name = task.ID

	if task.Replicas > 1 {
		jsonDeployment.Spec.Replicas = task.Replicas
	} else {
		jsonDeployment.Spec.Replicas = replicas
	}

	jsonDeployment.Spec.RevisionHistoryLimit = 10
	jsonDeployment.Spec.Selector.MatchLabels.App = task.ID
	jsonDeployment.Spec.Template.Metadata.Labels.App = task.ID

	// containers
	totalContainers := len(task.Containers)
	jsonDeployment.Spec.Template.Spec.Containers = make([]structs.K8S_DEPLOYMENT_CONTAINER, totalContainers)
	for i := 0; i < totalContainers; i++ {
		jsonDeployment.Spec.Template.Spec.Containers[i].Image = task.Containers[i].Image
		jsonDeployment.Spec.Template.Spec.Containers[i].ImagePullPolicy = "Always"
		jsonDeployment.Spec.Template.Spec.Containers[i].Name = task.Containers[i].Name

		// ports
		totalPorts := len(task.Containers[i].Ports)
		jsonDeployment.Spec.Template.Spec.Containers[i].Ports = make([]structs.K8S_DEPLOYMENT_CONTAINER_PORTS, totalPorts)
		for j := 0; j < totalPorts; j++ {
			jsonDeployment.Spec.Template.Spec.Containers[i].Ports[j].ContainerPort = task.Containers[i].Ports[j].ContainerPort
		}

		// env
		totalEnvs := len(task.Containers[i].Environment)
		if totalEnvs > 0 {
			jsonDeployment.Spec.Template.Spec.Containers[i].Env = make([]structs.K8S_DEPLOYMENT_CONTAINER_ENV, totalEnvs)
			for j := 0; j < totalEnvs; j++ {
				jsonDeployment.Spec.Template.Spec.Containers[i].Env[j].Name = task.Containers[i].Environment[j].Name
				jsonDeployment.Spec.Template.Spec.Containers[i].Env[j].Value = task.Containers[i].Environment[j].Value
			}
		}

		// command & args
		if len(task.Containers[i].Command) >= 1 {
			jsonDeployment.Spec.Template.Spec.Containers[i].Command = task.Containers[i].Command
		}

		if len(task.Containers[i].Args) >= 1 {
			jsonDeployment.Spec.Template.Spec.Containers[i].Args = task.Containers[i].Args
		}
	}

	return jsonDeployment
}

///////////////////////////////////////////////////////////////////////////////
// K8S_SERVICE

/*
StructNewK8sServiceTemplate creates a new K8s Service json
*/
func StructNewK8sServiceTemplate(task structs.CLASS_TASK, ip string) (*structs.K8S_SERVICE, int, string) {
	jsonService, mainPort, mainPortName := StructNewServiceTemplate(task)

	var ips = []string{ip}
	jsonService.Spec.ExternalIPs = ips

	return jsonService, mainPort, mainPortName
}

/*
StructNewServiceTemplate creates a new K8s Service json
*/
func StructNewServiceTemplate(task structs.CLASS_TASK) (*structs.K8S_SERVICE, int, string) {
	var jsonService *structs.K8S_SERVICE
	jsonService = new(structs.K8S_SERVICE)

	jsonService.ApiVersion = "v1"
	jsonService.Kind = "Service"
	jsonService.Metadata.Name = "serv-" + task.ID
	jsonService.Metadata.Labels.App = task.ID
	jsonService.Spec.Selector.App = task.ID

	mainPort := 0
	mainPortName := ""
	// ports
	for _, contElement := range task.Containers {
		for _, portElement := range contElement.Ports {
			jsonService.Spec.Ports =
				append(jsonService.Spec.Ports, structs.K8S_SERVICE_PORT{
					Name:       strconv.Itoa(portElement.ContainerPort) + "-" + strings.ToLower(portElement.Protocol),
					Port:       portElement.HostPort,
					Protocol:   strings.ToUpper(portElement.Protocol),
					TargetPort: portElement.ContainerPort})
			if mainPort == 0 {
				mainPort = portElement.ContainerPort
				mainPortName = strconv.Itoa(portElement.ContainerPort) + "-" + portElement.Protocol
			}
		}
	}

	return jsonService, mainPort, mainPortName
}

/*
StructNewPodServiceTemplate creates a new K8s Service for pods json
*/
func StructNewPodServiceTemplate(hostIP string, pod structs.DB_TASK_POD) *structs.K8S_SERVICE {
	var jsonService *structs.K8S_SERVICE
	jsonService = new(structs.K8S_SERVICE)

	jsonService.ApiVersion = "v1"
	jsonService.Kind = "Service"
	jsonService.Metadata.Name = "serv-" + pod.Name
	jsonService.Spec.Selector.PodName = pod.Name
	jsonService.Spec.ExternalIPs = append(jsonService.Spec.ExternalIPs, hostIP)

	jsonService.Spec.Ports =
		append(jsonService.Spec.Ports, structs.K8S_SERVICE_PORT{
			Name:       strconv.Itoa(pod.Port) + "-" + strings.ToLower(pod.Protocol),
			Port:       pod.Port, // port exposed in Kubernetes / Openshift
			Protocol:   strings.ToUpper(pod.Protocol),
			TargetPort: pod.TargetPort}) // application port

	return jsonService
}

///////////////////////////////////////////////////////////////////////////////
// K8S_ROUTE

/*
StructNewRouteTemplate creates a new K8s Route json
*/
func StructNewRouteTemplate(appID string, mainPort int, mainPortName string, serverIP string) *structs.K8S_ROUTE {
	var jsonRoute *structs.K8S_ROUTE
	jsonRoute = new(structs.K8S_ROUTE)

	jsonRoute.ApiVersion = "route.openshift.io/v1"
	jsonRoute.Kind = "Route"
	jsonRoute.Metadata.Name = "route-" + appID
	jsonRoute.Metadata.Namespace = "class"

	jsonRoute.Spec.Host = appID + "." + serverIP + ".xip.io"
	jsonRoute.Spec.Port.TargetPort = strings.ToLower(mainPortName) //strconv.Itoa(mainPort)

	jsonRoute.Spec.To.Kind = "Service"
	jsonRoute.Spec.To.Name = "serv-" + appID

	return jsonRoute
}

///////////////////////////////////////////////////////////////////////////////
// CLASS_QOS_TEMPLATE_LIST (JSON) DEFINITION

/*
StructCheckClasQoSTemplateList Checks if CLASS struct is valid (from json)
*/
func StructCheckClasQoSTemplateList(req *http.Request) (*cfg.CLASS_QOS_TEMPLATE_LIST, error) {
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClasQoSTemplateList] Checking json object ...")

	decoder := json.NewDecoder(req.Body)
	var t cfg.CLASS_QOS_TEMPLATE_LIST
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClasQoSTemplateList] ERROR (1)", err)
		return nil, err
	}

	tStr, err := CommClassQoSTemplateListToString(t)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClasQoSTemplateList] Parsed object (string): " + tStr)
	log.Println("Rotterdam > CAAS > structs-funcs [StructCheckClasQoSTemplateList] Sending parsed object ...")

	return &t, nil
}

///////////////////////////////////////////////////////////////////////////////
// POD_PATCH

/*
StructNewPodPatch creates a new K8s Patch json for pods
*/
func StructNewPodPatch(podName string) []structs.K8S_POD_PATCH_LINE {
	var lres []structs.K8S_POD_PATCH_LINE

	var jsonPatch *structs.K8S_POD_PATCH_LINE
	jsonPatch = new(structs.K8S_POD_PATCH_LINE)

	jsonPatch.Op = "add"
	jsonPatch.Path = "/metadata/labels/pod-name"
	jsonPatch.Value = podName

	lres = append(lres, *jsonPatch)

	return lres
}
