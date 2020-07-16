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

// DB_TASK

const (
	// DB Tables
	DB_TABLE_TASK = "Task"
	DB_TABLE_QOS  = "QoS"

	// DB_TASK 'Type' values
	DB_TASK_TYPE_DEFAULT = "default"
	DB_TASK_TYPE_COMPSS  = "compss"
)

type DB_TASK_POD struct {
	Name       string `json:"name,omitempty"`
	IP         string `json:"ip,omitempty"`         // IP accessed by external apps
	HostIP     string `json:"hostIp,omitempty"`     // node IP
	PodIP      string `json:"podIp,omitempty"`      // internal IP created by Kubernetes / Openshift
	Status     string `json:"status,omitempty"`     // running, unknown
	Protocol   string `json:"protocol,omitempty"`   // TCP
	Port       int    `json:"port,omitempty"`       // port exposed in Kubernetes / Openshift
	TargetPort int    `json:"targetPort,omitempty"` // application port
}

type DB_TASK struct {
	DbId                 string            `json:"dbid,omitempty"`
	Id                   string            `json:"id,omitempty"`
	Name                 string            `json:"name,omitempty"`
	NameSpace            string            `json:"nameSpace,omitempty"`
	Type                 string            `json:"type,omitempty"` // compss, default, function
	Url                  string            `json:"url,omitempty"`
	Status               string            `json:"status,omitempty"`
	AgreementId          string            `json:"agreementId,omitempty"`
	Replicas             int               `json:"replicas,omitempty"`
	TaskDefinition       CLASS_TASK        `json:"taskDefinition,omitempty"`
	COMPSsTaskDefinition CLASS_COMPSS_TASK `json:"compssTaskDefinition,omitempty"`
	Deployment           string            `json:"deployment,omitempty"`
	ClusterId            string            `json:"clusterid,omitempty"`
	Pods                 []DB_TASK_POD     `json:"pods,omitempty"`
}

// DB_TASK_QOS
type DB_TASK_QOS struct {
	DbId            string  `json:"dbid,omitempty"`
	Id              string  `json:"id,omitempty"`
	Type            string  `json:"type,omitempty"`
	IdTask          string  `json:"idtask,omitempty"`
	Guarantee       string  `json:"guarantee,omitempty"`
	TotalViolations int     `json:"totalviolations,omitempty"`
	MaxAllowed      int     `json:"maxallowed,omitempty"`
	Action          string  `json:"action,omitempty"`
	ScaleFactor     float64 `json:"scalefactor,omitempty"`
	MaxReplicas     int     `json:"maxreplicas,omitempty"`
	MinReplicas     int     `json:"minreplicas,omitempty"`
}

var DB_TASK_PREFIX string = "task_"

var DB_TASK_QOS_PREFIX string = "qos_"
