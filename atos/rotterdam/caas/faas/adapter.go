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

package faas

import (
	adapters "atos/rotterdam/caas/adapters"
	adpcommon "atos/rotterdam/caas/adapters/common"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	constants "atos/rotterdam/globals/constants"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
)

// adapter
var a adapters.Adapter

///////////////////////////////////////////////////////////////////////////////

// init
func init() {
	log.Println("Rotterdam > CAAS > FAAS > adapter [init] Selected cluster index set to '0' ...")

}

///////////////////////////////////////////////////////////////////////////////

/*
generateID generates an Identifier based on a random string and the current time in nanoseconds
*/
func generateID() string {
	// 1. ID generation
	id := shortuuid.NewWithAlphabet("0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwx") //id := txt + shortuuid.New()
	id = strings.ToLower(id)                                                                     // Kubernetes doesnt allow uppercase names for deployments

	// 2. time in nanoseconds is appended to ID
	now := time.Now()      // current local time
	nsec := now.UnixNano() // number of nanoseconds since January 1, 1970 UTC
	id = id + strconv.FormatInt(nsec, 10)

	return id
}

/*
setClassStructDataDefaultValues
*/
func setClassStructDataDefaultValues(task *structs.CLASS_TASK, isCOMPs bool) { // *structs.CLASS_TASK {
	if len(task.Cluster) == 0 {
		task.Cluster = constants.MainClusterID
	}
	if len(task.Dock) == 0 {
		task.Dock = common.GetClusterDefaultDock(task.Cluster)
		//	task.Dock = constants.DefaultDock
	}

	if isCOMPs == true {
		task.Type = constants.TypeTaskCOMPSs
	} else if len(task.Type) == 0 || task.Type != constants.TypeTaskCOMPSs {
		task.Type = constants.TypeTaskDefault
	}
}

/*
setClassStructDataAll
*/
func setClassStructDataAll(task *structs.CLASS_TASK, isCOMPs bool) { // *structs.CLASS_TASK {
	task.ID = generateID()                         // ID
	task.Created = time.Now().Format(time.RFC3339) // current local time
	setClassStructDataDefaultValues(task, isCOMPs)
}

/*
errorMessage
*/
func errorMessage(w http.ResponseWriter, err error, method string) {
	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "error",
		Method:      method,
		Message:     err.Error(),
		FaaSVersion: cfg.Config.FaaSVersion})
}

/*
okTaskMessage
*/
func okTaskMessage(w http.ResponseWriter, message string, method string, dbTask structs.DB_TASK) {
	json.NewEncoder(w).Encode(structs.ResponseFaaSTask{
		Resp:        "ok",
		Method:      method,
		Message:     message,
		ID:          dbTask.TaskDefinition.ID, //classTask.Id,
		FaaSVersion: cfg.Config.FaaSVersion,
		URL:         dbTask.Url,
		Task:        dbTask})
}

/*
validateJSONTask validates input (json) and generates a valid CLASS Task struct
*/
func validateJSONTask(r *http.Request) (*structs.CLASS_TASK, error) {
	log.Println("Rotterdam > CAAS > FAAS > adapter [validateJsonTask] 'Duplicating' r.Body to use it multiple times with decoders ...")
	rBody := r.Body
	var buf bytes.Buffer
	tee := io.TeeReader(rBody, &buf)

	log.Println("Rotterdam > CAAS > FAAS > adapter [validateJsonTask] Parsing default / old json definition ...")
	decoder := json.NewDecoder(tee)
	classTask, err := common.StructCheckClassTask(decoder)
	if err == nil {
		setClassStructDataAll(classTask, false)
		log.Println("Rotterdam > CAAS > FAAS > adapter [validateJsonTask] Task with name " + classTask.Name + " received. Task id is [" + classTask.ID + "]. Type = " + classTask.Type)
		return classTask, nil
	} else {
		log.Println("Rotterdam > CAAS > FAAS > adapter [validateJsonTask] Parsing new json definition (COMPSs) ...")

		decoder = json.NewDecoder(&buf)
		classCOMPSsTask, err := common.StructCheckClassCOMPSsTask(decoder)
		if err == nil {
			classCOMPSsTask.ID = generateID() // ID
			classTask = structs.TransfCOMPSSTASKtoTASK(classCOMPSsTask)
			log.Println("Rotterdam > CAAS > FAAS > adapter [validateJsonTask] Task with name " + classCOMPSsTask.Name + " received. Task id is [" + classCOMPSsTask.ID + "]. Type = " + classTask.Type + ". Using 'COMPSs' JSON format.")
			setClassStructDataDefaultValues(classTask, true)
			return classTask, nil
		}
	}

	return nil, err
}

///////////////////////////////////////////////////////////////////////////////

/*
DeployFunction Deploys a function
*/
func DeployFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/functions    <<Function Deployment>>")

	log.Println("Rotterdam > CAAS > FAAS > adapter [DeployFunction] (1.) Validating json ...")
	_, err := validateJSONTask(r)
	if err != nil {
		errorMessage(w, err, "DeployFunction")
	} else {
		//deploy(w, classTask)
	}
}

/*
CallFunction Calls a function
*/
func CallFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/functions/{id}    <<Call to a function>>")

	log.Println("Rotterdam > CAAS > FAAS > adapter [CallFunction] (1.) Validating json ...")
	_, err := validateJSONTask(r)
	if err != nil {
		errorMessage(w, err, "CallFunction")
	} else {
		//deploy(w, classTask)
	}
}

/*
RemoveFunction Deletes a function
*/
func RemoveFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### DELETE /api/v1/functions/{id}    <<Function termination>>")

	log.Println("Rotterdam > CAAS > adapter [RemoveFunction] (1.) Reading input params ...")
	params := mux.Vars(r)

	log.Println("Rotterdam > CAAS > FAAS > adapter [RemoveFunction] (2.) Getting task [" + params["id"] + "] from DB ...")
	dbTask, err := common.ReadTaskValue(params["id"])
	if err == nil {
		log.Println("Rotterdam > CAAS > FAAS > adapter [RemoveFunction] (3.) Getting adapter for cluster [" + dbTask.ClusterId + "]  ...")

	} else {
		errorMessage(w, err, "RemoveFunction")
	}
}

/*
GetAllFunctions Gets all functions
*/
func GetAllFunctions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/functions")

	log.Println("Rotterdam > CAAS > FAAS > adapter [GetAllFunctions] (1.) Getting all functions ...")
	// TODO
	resp, err := adpcommon.GetAllTasks()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " functions retrieved"
		json.NewEncoder(w).Encode(structs.ResponseFaaSTasks{
			Resp:        "ok",
			Method:      "GetAllFunctions",
			Message:     msg,
			FaaSVersion: cfg.Config.FaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetAllFunctions")
	}
}

/*
GetFunction Gets a function
*/
func GetFunction(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/functions/{id}")

	log.Println("Rotterdam > CAAS > FAAS > adapter [GetFunction] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > FAAS > adapter [GetFunction] id=" + params["id"])

	log.Println("Rotterdam > CAAS > FAAS > adapter [GetFunction] (1.) Getting function " + params["id"] + " ...")
	// TODO
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		okTaskMessage(w, "Task retrieved", "GetFunction", dbTask)
	} else {
		errorMessage(w, err, "GetFunction")
	}
}
