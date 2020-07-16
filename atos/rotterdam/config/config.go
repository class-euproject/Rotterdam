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

package cfg

import (
	"log"
	"os"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////
// configuration struct: is filled with values from json file and environement values
type ConfigurationCluster struct {
	ID                            string
	Name                          string
	Description                   string
	Type                          string
	SO                            string
	DefaultDock                   string
	OpenshiftOauthToken           string
	KubernetesEndPoint            string
	OpenshiftEndPoint             string
	SLALiteEndPoint               string
	PrometheusPushgatewayEndPoint string
	HostIP                        string
	HostPort                      int
	User                          string
	Password                      string
	KeyFile                       string
}

type Configuration struct {
	RestApiVersion     string
	CaaSVersion        string
	FaaSVersion        string
	RulesEngineVersion string
	IMECVersion        string
	SLALiteVersion     string
	ServerPort         int
	Clusters           []ConfigurationCluster
	SLAs               struct {
		CreationDate   string
		ExpirationDate string
	}
	Tasks struct {
		MaxReplicas int
		MinReplicas int
		MaxAllowed  int
		ScaleFactor float64
		Value       int
		Comparator  string
		Action      string
	}
}

// global configuration
var Config Configuration

///////////////////////////////////////////////////////////////////////////////
// Global TASK_QOS_LIST

type CLASS_QOS_TEMPLATE__ELEM_GUARANTEE struct {
	Name       string `json:"name,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

type CLASS_QOS_TEMPLATE__ELEM struct {
	Type          string                               `json:"type,omitempty"`
	GuaranteeName string                               `json:"guaranteeName,omitempty"`
	MaxAllowed    int                                  `json:"maxAllowed,omitempty"`
	Action        string                               `json:"action,omitempty"`
	ScaleFactor   float64                              `json:"scaleFactor,omitempty"`
	Guarantees    []CLASS_QOS_TEMPLATE__ELEM_GUARANTEE `json:"guarantees,omitempty"`
	MaxReplicas   int                                  `json:"maxReplicas,omitempty"`
	MinReplicas   int                                  `json:"minReplicas,omitempty"`
}

type CLASS_QOS_TEMPLATE_LIST []CLASS_QOS_TEMPLATE__ELEM

// global qos_templates
var QosTemplates CLASS_QOS_TEMPLATE_LIST

//
func setEnvValue(val *string, name string) {
	if os.Getenv(name) != "" {
		log.Println("Rotterdam > CAAS > Config > Updating " + name + " value ... " + os.Getenv(name))
		*val = os.Getenv(name)
	}
}

//
func setEnvValueInt(val *int, name string) {
	if os.Getenv(name) != "" {
		log.Println("Rotterdam > CAAS > Config > Setting " + name + " value ... " + os.Getenv(name))
		*val, _ = strconv.Atoi(os.Getenv(name))
	}
}

// Initialize configuration values
func InitConfig(cfg *Configuration) {
	log.Println("Rotterdam > CAAS > Config > Checking configuration values from ENV ...")

	setEnvValue(&cfg.Clusters[0].Type, "Type")
	setEnvValue(&cfg.Clusters[0].KubernetesEndPoint, "KubernetesEndPoint")
	setEnvValue(&cfg.Clusters[0].OpenshiftEndPoint, "OpenshiftEndPoint")
	setEnvValue(&cfg.Clusters[0].HostIP, "HostIP")
	setEnvValue(&cfg.Clusters[0].OpenshiftOauthToken, "OpenshiftOauthToken")
	setEnvValue(&cfg.Clusters[0].SLALiteEndPoint, "SLALiteEndPoint")
	setEnvValueInt(&cfg.ServerPort, "ServerPort")
}
