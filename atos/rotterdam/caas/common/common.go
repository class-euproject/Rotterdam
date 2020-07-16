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
	constants "atos/rotterdam/globals/constants"
	imec_db "atos/rotterdam/imec/db"
	"encoding/json"
	"log"
	"strconv"
)

/*
CommClassQoSTemplateListToString Parses a struct to a string
*/
func CommClassQoSTemplateListToString(ct cfg.CLASS_QOS_TEMPLATE_LIST) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [CommClassQoSTemplateListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommClassStructToString Parses a struct to a string
*/
func CommClassStructToString(ct structs.CLASS_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [CommClassStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommClassCOMPSsStructToString Parses a struct to a string
*/
func CommClassCOMPSsStructToString(ct structs.CLASS_COMPSS_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > structs-funcs [CommClassCOMPSsStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDeploymentStructToString Parses a struct to a string
*/
func CommDeploymentStructToString(ct structs.K8S_DEPLOYMENT) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDeploymentStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommServiceStructToString Parses a struct to a string
*/
func CommServiceStructToString(ct structs.K8S_SERVICE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommServiceStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommRouteStructToString Parses a struct to a string
*/
func CommRouteStructToString(ct structs.K8S_ROUTE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommRouteStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommPatchPodsListToString Parses a struct to a string
*/
func CommPatchPodsListToString(ct []structs.K8S_POD_PATCH_LINE) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommPatchPodsListToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDbTaskStructToString Parses a struct to a string
*/
func CommDbTaskStructToString(ct structs.DB_TASK) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDbTaskStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommDbTaskQoSStructToString Parses a struct to a string
*/
func CommDbTaskQoSStructToString(ct structs.DB_TASK_QOS) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommDbTaskQoSStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommSLAStructToString Parses a struct to a string
*/
func CommSLAStructToString(ct structs.SLA_AGREEMENT) (string, error) {
	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommSLAStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
CommStringToDbTaskStruct Parses a string to a struct of type DB_TASK
*/
func CommStringToDbTaskStruct(ct string) (*structs.DB_TASK, error) {
	data := &structs.DB_TASK{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToDbTaskStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
CommStringToDbTaskQoSStruct Parses a string to a struct of type DB_TASK_QOS
*/
func CommStringToDbTaskQoSStruct(ct string) (*structs.DB_TASK_QOS, error) {
	data := &structs.DB_TASK_QOS{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToDbTaskQoSStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
CommStringToK8S_SCALE Parses a string to a struct of type K8S_SCALE
*/
func CommStringToK8S_SCALE(ct string) (*structs.K8S_SCALE, error) {
	data := &structs.K8S_SCALE{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > CAAS > common [CommStringToK8S_SCALE] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
GetClusterIndex Returns the index of a cluster identified by its name / id value
*/
func GetClusterIndex(idCluster string) int {
	for index := range cfg.Config.Clusters {
		if idCluster == cfg.Config.Clusters[index].ID {
			log.Println("Rotterdam > CAAS > common [GetClusterIndex] Returning index (" + strconv.Itoa(index) + ") of cluster [" + idCluster + "] ...")
			return index
		}
	}

	log.Println("Rotterdam > CAAS > common [GetClusterIndex] WARNING Returning default cluster index (0) ...")
	return 0
}

/*
GetClusterType Returns the type of a cluster identified by its name / id value
*/
func GetClusterType(idCluster string) string {
	resp, err := imec_db.GetInfrByID(idCluster)
	if err == nil && len(resp) > 0 {
		return resp[0].Type
	}

	log.Println("Rotterdam > CAAS > common [GetClusterType] WARNING Returning default cluster type ...")
	return "Openshift"
}

/*
GetClusterDefaultDock Returns the type of a cluster identified by its name / id value
*/
func GetClusterDefaultDock(idCluster string) string {
	resp, err := imec_db.GetInfrByID(idCluster)
	if err == nil && len(resp) > 0 {
		return resp[0].DefaultDock
	}

	log.Println("Rotterdam > CAAS > common [GetClusterDefaultDock] WARNING Returning default cluster dock ...")
	return constants.DefaultDock
}
