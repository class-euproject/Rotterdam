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
// Created on 15 Jan 2020
// @author: Roi Sucasas - ATOS
//

package imec

import (
	cfg "atos/rotterdam/config"
	constants "atos/rotterdam/globals/constants"
	db "atos/rotterdam/imec/db"
	infrastructures "atos/rotterdam/imec/infrastructures"
	deployment "atos/rotterdam/imec/orchestrators/deployment"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

/*
ResponseIMEC IMEC Response
*/
type ResponseIMEC struct {
	Resp               string `json:"resp,omitempty"`
	Method             string `json:"method,omitempty"`
	Message            string `json:"message,omitempty"`
	CaaSVersion        string `json:"caasversion,omitempty"`
	IMECVersion        string `json:"imecversion,omitempty"`
	RulesEngineVersion string `json:"rulesengineversion,omitempty"`
	RestApiVersion     string `json:"restapi,omitempty"`
	SLALiteVersion     string `json:"slalite,omitempty"`
	Content            string `json:"content,omitempty"`
}

/*
ResponseInfr
*/
type ResponseInfr struct {
	Resp        string                       `json:"resp,omitempty"`
	Method      string                       `json:"method,omitempty"`
	Message     string                       `json:"message,omitempty"`
	IMECVersion string                       `json:"imecversion,omitempty"`
	Content     string                       `json:"content,omitempty"`
	URL         string                       `json:"url,omitempty"`
	Infr        db.DB_INFRASTRUCTURE_CLUSTER `json:"infr,omitempty"`
	Id          string                       `json:"id,omitempty"`
}

/*
ResponseInfrs
*/
type ResponseInfrs struct {
	Resp        string                         `json:"resp,omitempty"`
	Method      string                         `json:"method,omitempty"`
	Message     string                         `json:"message,omitempty"`
	IMECVersion string                         `json:"imecversion,omitempty"`
	Content     string                         `json:"content,omitempty"`
	Infrs       []db.DB_INFRASTRUCTURE_CLUSTER `json:"infrs,omitempty"`
}

// generateID
func generateID() string {
	id := shortuuid.NewWithAlphabet("0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx")
	id = strings.ToLower(id)
	return id
}

/*
GetAllInfrastructures Get all infrastructures
*/
func GetAllInfrastructures(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > IMEC > imec [GetAllInfrastructures] Getting infrastructures ...")

	resp, err := db.GetAllInfrs()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " infrastructures retrieved"
		json.NewEncoder(w).Encode(ResponseInfrs{
			Resp:        "ok",
			Method:      "GetAllInfrs",
			Message:     msg,
			IMECVersion: cfg.Config.IMECVersion,
			Infrs:       resp})
	} else {
		json.NewEncoder(w).Encode(ResponseIMEC{
			Resp:        "error",
			Method:      "GetAllInfrs",
			Message:     err.Error(),
			IMECVersion: cfg.Config.IMECVersion})
	}
}

/*
CreateInfrastructure Create new infrastructure
*/
func CreateInfrastructure(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > IMEC > imec [CreateInfrastructure] Creating a new infrastructure ...")

	classInfr, err := infrastructures.ValidateStructInfrJSON(r)
	if err == nil {
		dbInfr := infrastructures.InfrJSONToDBObj(*classInfr, generateID())
		err = db.SetInfrValue(dbInfr.ID, *dbInfr)
		if err == nil {
			// send response with DBTask
			dbInfr, err = db.ReadInfrValue(dbInfr.ID)
			if err == nil {
				json.NewEncoder(w).Encode(ResponseInfr{
					Resp:        "ok",
					Method:      "CreateInfrastructure",
					Message:     "Infrastructure created",
					Id:          dbInfr.ID,
					IMECVersion: cfg.Config.IMECVersion,
					Infr:        *dbInfr})
			} else {
				json.NewEncoder(w).Encode(ResponseIMEC{
					Resp:        "error",
					Method:      "CreateInfrastructure",
					Message:     err.Error(),
					IMECVersion: cfg.Config.IMECVersion})
			}
		} else {
			json.NewEncoder(w).Encode(ResponseIMEC{
				Resp:        "error",
				Method:      "CreateInfrastructure",
				Message:     err.Error(),
				IMECVersion: cfg.Config.IMECVersion})
		}
	} else {
		json.NewEncoder(w).Encode(ResponseIMEC{
			Resp:        "error",
			Method:      "CreateInfrastructure",
			Message:     err.Error(),
			IMECVersion: cfg.Config.IMECVersion})
	}
}

/*
UpdateInfrastructure Update infrastructure
*/
func UpdateInfrastructure(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > IMEC > imec [UpdateInfrastructure] ...")

	NotImplementedFunc(w, r)
}

/*
GetInfrastructure Get infrastructure
*/
func GetInfrastructure(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	params := mux.Vars(r)
	log.Println("Rotterdam > IMEC > imec [GetInfrastructure] id=" + params["id"])

	log.Println("Rotterdam > IMEC > imec [GetInfrastructure] Getting infrastructure " + params["id"] + " ...")
	resp, err := db.GetInfrByID(params["id"])
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " infrastructure retrieved"
		json.NewEncoder(w).Encode(ResponseInfrs{
			Resp:        "ok",
			Method:      "GetInfrastructure",
			Message:     msg,
			IMECVersion: cfg.Config.IMECVersion,
			Infrs:       resp})
	} else {
		json.NewEncoder(w).Encode(ResponseIMEC{
			Resp:        "error",
			Method:      "GetInfrastructure",
			Message:     err.Error(),
			IMECVersion: cfg.Config.IMECVersion})
	}
}

/*
DeleteInfrastructure Delete infrastructure
*/
func DeleteInfrastructure(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > IMEC > imec [DeleteInfrastructure] ...")

	NotImplementedFunc(w, r)
}

/*
DeployCluster Deploys a K8s cluster in a new server
*/
func DeployCluster(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### New Deployment: cluster")
	log.Println("####################################################################################")

	params := mux.Vars(r)
	log.Println("Rotterdam > IMEC > imec [DeployCluster] id=" + params["id"])

	dbInfr, err := db.ReadInfrValue(params["id"])
	if err == nil {
		if dbInfr.Type == "microk8s" {
			// microk8s
			go func() {
				err = deployment.MicroK8s(dbInfr)
				if err != nil {
					log.Println("Rotterdam > IMEC > imec [DeployCluster] ERROR in SetInfrValue: ", err)
					// update infr status ==> error
					dbInfr.Status = constants.ClusterERROR
					err := db.SetInfrValue(dbInfr.ID, *dbInfr)
					if err == nil {
						log.Println("Rotterdam > IMEC > imec [DeployCluster] ERROR ", err)
					}
				}
			}()
			json.NewEncoder(w).Encode(ResponseInfr{
				Resp:        "ok",
				Method:      "DeployCluster",
				Message:     "Deploying MicroK8s infrastructure in " + dbInfr.HostIP + ", cluster id=" + dbInfr.ID,
				Id:          dbInfr.ID,
				IMECVersion: cfg.Config.IMECVersion,
				Infr:        *dbInfr})
		} else if dbInfr.Type == "kubeless" {
			// kubeless
			go func() {
				err = deployment.Kubeless(dbInfr)
				if err != nil {
					log.Println("Rotterdam > IMEC > imec [DeployCluster] ERROR in SetInfrValue: ", err)
					// update infr status ==> error
					dbInfr.Status = constants.ClusterERROR
					err := db.SetInfrValue(dbInfr.ID, *dbInfr)
					if err == nil {
						log.Println("Rotterdam > IMEC > imec [DeployCluster] ERROR ", err)
					}
				}
			}()
			json.NewEncoder(w).Encode(ResponseInfr{
				Resp:        "ok",
				Method:      "DeployCluster",
				Message:     "Deploying Kubeless infrastructure in " + dbInfr.HostIP + ", cluster id=" + dbInfr.ID,
				Id:          dbInfr.ID,
				IMECVersion: cfg.Config.IMECVersion,
				Infr:        *dbInfr})
		} else {
			// error
			json.NewEncoder(w).Encode(ResponseIMEC{
				Resp:        "error",
				Method:      "DeployCluster",
				Message:     "Error: " + dbInfr.Type + " not defined",
				IMECVersion: cfg.Config.IMECVersion})
		}
	} else {
		json.NewEncoder(w).Encode(ResponseIMEC{
			Resp:        "error",
			Method:      "DeployCluster",
			Message:     err.Error(),
			IMECVersion: cfg.Config.IMECVersion})
	}
}

/*
DeleteCluster Delete Cluster
*/
func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > IMEC > imec [DeleteCluster] ...")

	NotImplementedFunc(w, r)
}

/*
GetCluster Get cluster
*/
func GetCluster(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > IMEC > imec [GetCluster] ...")

	NotImplementedFunc(w, r)
}

/*
NotImplementedFunc Default Function for not implemented calls
*/
func NotImplementedFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > IMEC > imec [NotImplementedFunc] -Not implemented-")

	json.NewEncoder(w).Encode(ResponseIMEC{
		Resp:        "ok",
		Method:      "NotImplementedFunc",
		Message:     "not implemented",
		IMECVersion: cfg.Config.IMECVersion})
}
