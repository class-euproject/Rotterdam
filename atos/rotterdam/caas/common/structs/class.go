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

package structs

import (
	cfg "atos/rotterdam/config"
)

///////////////////////////////////////////////////////////////////////////////
// COMMON STRUCTS

/*
Configuration ...
*/
type Configuration struct {
	Adapter        string
	KubeConfigPath string
}

/*
RespHTTP Http Response content
*/
type RespHTTP struct {
	RespStr string
}

/*
ResponseCaaS JSON CaaS Response
*/
type ResponseCaaS struct {
	Resp               string `json:"resp,omitempty"`
	Method             string `json:"method,omitempty"`
	Message            string `json:"message,omitempty"`
	CaaSVersion        string `json:"caasversion,omitempty"`
	FaaSVersion        string `json:"faasversion,omitempty"`
	RulesEngineVersion string `json:"rulesengineversion,omitempty"`
	RestAPIVersion     string `json:"restapi,omitempty"`
	SLALiteVersion     string `json:"slalite,omitempty"`
	IMECVersion        string `json:"imec,omitempty"`
	Content            string `json:"content,omitempty"`
}

/*
ResponseCaaSTask CaaS Response
*/
type ResponseCaaSTask struct {
	Resp        string  `json:"resp,omitempty"`
	Method      string  `json:"method,omitempty"`
	Message     string  `json:"message,omitempty"`
	CaaSVersion string  `json:"caasversion,omitempty"`
	Content     string  `json:"content,omitempty"`
	URL         string  `json:"url,omitempty"`
	Task        DB_TASK `json:"task,omitempty"`
	ID          string  `json:"id,omitempty"`
}

/*
ResponseFaaSTask FaaS Response
*/
type ResponseFaaSTask struct {
	Resp        string  `json:"resp,omitempty"`
	Method      string  `json:"method,omitempty"`
	Message     string  `json:"message,omitempty"`
	FaaSVersion string  `json:"faasversion,omitempty"`
	Content     string  `json:"content,omitempty"`
	URL         string  `json:"url,omitempty"`
	Task        DB_TASK `json:"task,omitempty"`
	ID          string  `json:"id,omitempty"`
}

/*
ResponseCaaSTasks ...
*/
type ResponseCaaSTasks struct {
	Resp        string    `json:"resp,omitempty"`
	Method      string    `json:"method,omitempty"`
	Message     string    `json:"message,omitempty"`
	CaaSVersion string    `json:"caasversion,omitempty"`
	Content     string    `json:"content,omitempty"`
	Tasks       []DB_TASK `json:"tasks,omitempty"`
}

/*
ResponseFaaSTasks
*/
type ResponseFaaSTasks struct {
	Resp        string    `json:"resp,omitempty"`
	Method      string    `json:"method,omitempty"`
	Message     string    `json:"message,omitempty"`
	FaaSVersion string    `json:"faasversion,omitempty"`
	Content     string    `json:"content,omitempty"`
	Tasks       []DB_TASK `json:"tasks,omitempty"`
}

/*
ResponseCaaSTasksQoS
*/
type ResponseCaaSTasksQoS struct {
	Resp        string        `json:"resp,omitempty"`
	Method      string        `json:"method,omitempty"`
	Message     string        `json:"message,omitempty"`
	CaaSVersion string        `json:"caasversion,omitempty"`
	Content     string        `json:"content,omitempty"`
	TasksQoS    []DB_TASK_QOS `json:"tasksqos,omitempty"`
}

/*
ResponseQoSDefinition ...
*/
type ResponseQoSDefinition struct {
	Resp        string                       `json:"resp,omitempty"`
	Method      string                       `json:"method,omitempty"`
	Message     string                       `json:"message,omitempty"`
	CaaSVersion string                       `json:"caasversion,omitempty"`
	Content     string                       `json:"content,omitempty"`
	QoSDef      cfg.CLASS_QOS_TEMPLATE__ELEM `json:"qosdef,omitempty"`
}

/*
ResponseQoSDefinitions ...
*/
type ResponseQoSDefinitions struct {
	Resp        string                         `json:"resp,omitempty"`
	Method      string                         `json:"method,omitempty"`
	Message     string                         `json:"message,omitempty"`
	CaaSVersion string                         `json:"caasversion,omitempty"`
	Content     string                         `json:"content,omitempty"`
	QoSDefs     []cfg.CLASS_QOS_TEMPLATE__ELEM `json:"qosdefs,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// CLASS TASK (JSON) DEFINITION

/*
CLASS_TASK_CONTAINER_PORTS ...
*/
type CLASS_TASK_CONTAINER_PORTS struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	HostPort      int    `json:"hostPort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

/*
CLASS_TASK_CONTAINER_VOLUMES ...
*/
type CLASS_TASK_CONTAINER_VOLUMES struct {
	Name      string `json:"name,omitempty"`
	MountPath string `json:"mounthPath,omitempty"`
}

/*
CLASS_TASK_CONTAINER_ENV ...
*/
type CLASS_TASK_CONTAINER_ENV struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

/*
CLASS_TASK_CONTAINER ...
*/
type CLASS_TASK_CONTAINER struct {
	Name        string                         `json:"name,omitempty"`
	Image       string                         `json:"image"`
	Ports       []CLASS_TASK_CONTAINER_PORTS   `json:"ports,omitempty"`
	Volumes     []CLASS_TASK_CONTAINER_VOLUMES `json:"volumes,omitempty"`
	Environment []CLASS_TASK_CONTAINER_ENV     `json:"environment,omitempty"`
	Command     []string                       `json:"command,omitempty"`
	Args        []string                       `json:"args,omitempty"`
}

/*
CLASS_TASK_QOS_CUSTOM_GUARANTEES ...
*/
type CLASS_TASK_QOS_CUSTOM_GUARANTEES struct {
	Metric    string `json:"metric,omitempty"`
	Condition string `json:"condition,omitempty"`
	Value     string `json:"value,omitempty"`
}

/*
CLASS_TASK_QOS_CUSTOM ...
*/
type CLASS_TASK_QOS_CUSTOM struct {
	Type        string                             `json:"type,omitempty"` // type: "infr" or "task"
	Name        string                             `json:"name,omitempty"`
	Description string                             `json:"description,omitempty"`
	Guarantees  []CLASS_TASK_QOS_CUSTOM_GUARANTEES `json:"guarantees,omitempty"`
	ScaleFactor int                                `json:"scalefactor,omitempty"`
	Action      string                             `json:"action,omitempty"`
	Max         int                                `json:"max,omitempty"`
	Min         int                                `json:"min,omitempty"`
}

/*
CLASS_TASK_QOS ...
*/
type CLASS_TASK_QOS struct {
	Name        string                  `json:"name,omitempty"`
	Description string                  `json:"description,omitempty"`
	Custom      []CLASS_TASK_QOS_CUSTOM `json:"custom,omitempty"`
}

/*
CLASS_TASK ...
	JSON FORMAT - EXAMPLE:
		{
			"name": "adas-ped-detection",
			"dock": "adas-pro",
			"cliuster": "maincluster",
			"qos": {},
			"replicas": 1,
			"containers": [{
					"name": "adas-ped-distance", <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
					"image": "docker.class.com/adas/adas_task_1:0.1.1",
					"ports": [
						{
							"containerPort": "80",
							"hostPort": "80", <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
							"protocol": "tcp",
						}
					],
					"volumes": [
						{
							"name": "workdir",
							"mounthPath": "/usr/share/class/adas/"
						}
					],
					"environment": [
						{
						"name": "COMPS_VERSION",
						"value":"1.2.1"
						},
						{
						"name": "COMPS_MASTER_URL",
						"value":"<ES-URL>"
						}
					]
				}
			]
		}
*/
type CLASS_TASK struct {
	Name       string                  `json:"name"`
	ID         string                  `json:"id,omitempty"`
	Type       string                  `json:"type,omitempty"` // type: "compss" or "default"
	Dock       string                  `json:"dock"`
	Cluster    string                  `json:"cluster"`
	Qos        CLASS_TASK_QOS          `json:"qos,omitempty"`
	QoSCOMPSs  []CLASS_COMPSS_TASK_QOS `json:"qoscompss,omitempty"`
	Replicas   int                     `json:"replicas"`
	Containers []CLASS_TASK_CONTAINER  `json:"containers"`
	Created    string                  `json:"created,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// CLASS COMPSS TASK v2 (JSON) DEFINITION

/*
CLASS_COMPSS_TASK_QOS ...
	JSON FORMAT - EXAMPLE:
		{
			"name": "app-compss-example",
			"dock": "class",
			"qos": [],
			"image": "docker.class.com/adas/adas_task_1:0.1.1",
			"replicas": 15,
			"ports": [8554],
			"command": ["printenv"],
			"args": ["HOSTNAME", "KUBERNETES_PORT"]
		}
*/
type CLASS_COMPSS_TASK_QOS struct {
	QoSId       string  `json:"qosid,omitempty"`
	Metric      string  `json:"metric,omitempty"`     // metric name / prefix
	Comparator  string  `json:"comparator,omitempty"` // "<", "=", ">"
	Value       int     `json:"value,omitempty"`
	Action      string  `json:"action,omitempty"` // "scale_out", "scale_in"
	MaxReplicas int     `json:"maxreplicas,omitempty"`
	MinReplicas int     `json:"minreplicas,omitempty"`
	ScaleFactor float64 `json:"scalefactor,omitempty"`
	MaxAllowed  int     `json:"maxallowed,omitempty"`
}

/*
CLASS_COMPSS_TASK ...
*/
type CLASS_COMPSS_TASK struct {
	ID       string                  `json:"id,omitempty"`
	Name     string                  `json:"name,omitempty"`
	Type     string                  `json:"type,omitempty"` // type: "compss" or "default"
	Dock     string                  `json:"dock,omitempty"`
	Cluster  string                  `json:"cluster,omitempty"`
	Qos      []CLASS_COMPSS_TASK_QOS `json:"qos,omitempty"`
	Image    string                  `json:"image,omitempty"`
	Replicas int                     `json:"replicas,omitempty"`
	Ports    []int                   `json:"ports,omitempty"`
	Command  []string                `json:"command,omitempty"`
	Args     []string                `json:"args,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////

/*
CLASS_FUNCTION_TASK ...
	JSON FORMAT - EXAMPLE:
		{
			"name": "redis-app",
			"dock": "default",
			"cluster": "microk8s_1",
			"runtime": "python2.7",
			"timeout": "180",
			"handler": "helloget.foo",
			"functionType": "text",
			"function": "def foo(event, context):\n    return \"hello world\"\n"
		}
*/
type CLASS_FUNCTION_TASK struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"` // type: "compss", "default", "function"
	Dock         string `json:"dock,omitempty"`
	Cluster      string `json:"cluster,omitempty"`
	Runtime      string `json:"runtime,omitempty"`
	Timeout      string `json:"timeout,omitempty"`
	Handler      string `json:"handler,omitempty"`
	FunctionType string `json:"functionType,omitempty"`
	Function     string `json:"function,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////

/*
ViolationInfo ...
*/
type ViolationInfo struct {
	Agreement_id string
	Status       string
	Client_id    string
	Client_name  string
	Guarantee    string
}

///////////////////////////////////////////////////////////////////////////////

/*
SLA_AGREEMENT_DETAILS_GUARANTIES ...
*/
type SLA_AGREEMENT_DETAILS_GUARANTIES struct {
	Name       string `json:"name,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

/*
SLA_AGREEMENT_DETAILS_ST ...
*/
type SLA_AGREEMENT_DETAILS_ST struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

/*
SLA_AGREEMENT_DETAILS ...
*/
type SLA_AGREEMENT_DETAILS struct {
	ID         string                             `json:"id,omitempty"`
	Type       string                             `json:"type,omitempty"`
	Name       string                             `json:"name,omitempty"`
	Provider   SLA_AGREEMENT_DETAILS_ST           `json:"provider,omitempty"`
	Client     SLA_AGREEMENT_DETAILS_ST           `json:"client,omitempty"`
	Creation   string                             `json:"creation,omitempty"`
	Expiration string                             `json:"expiration,omitempty"`
	Guarantees []SLA_AGREEMENT_DETAILS_GUARANTIES `json:"guarantees,omitempty"`
}

/*
SLA_AGREEMENT ...
*/
type SLA_AGREEMENT struct {
	ID      string                `json:"id,omitempty"`
	Name    string                `json:"name,omitempty"`
	State   string                `json:"state,omitempty"`
	Details SLA_AGREEMENT_DETAILS `json:"details,omitempty"`
}
