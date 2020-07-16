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

package caas

import (
	adapters "atos/rotterdam/caas/adapters"
	adpcommon "atos/rotterdam/caas/adapters/common"
	k8s_adapter "atos/rotterdam/caas/adapters/kubernetes"
	ops_adapter "atos/rotterdam/caas/adapters/openshift"
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	sla "atos/rotterdam/caas/sla"
	cfg "atos/rotterdam/config"
	constants "atos/rotterdam/globals/constants"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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
	log.Println("Rotterdam > CAAS > adapter [init] Selected cluster index set to '0' ...")

	qosTemplatesPath := "./config/qos_templates.json"
	log.Println("Rotterdam > CAAS > adapter [init] Reading content of qos_templates.json file [" + qosTemplatesPath + "] ...")

	file, err := os.Open(qosTemplatesPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.QosTemplates)
	if err != nil {
		panic(err)
	}

	log.Println("Rotterdam > CAAS > adapter [init] List of available QoS elements ...")
	for i := range cfg.QosTemplates {
		log.Println("Rotterdam > CAAS > adapter [init] > [" + strconv.Itoa(i) + "] " + cfg.QosTemplates[i].GuaranteeName)
	}
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
		CaaSVersion: cfg.Config.CaaSVersion})
}

/*
okTaskMessage
*/
func okTaskMessage(w http.ResponseWriter, message string, method string, dbTask structs.DB_TASK) {
	json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
		Resp:        "ok",
		Method:      method,
		Message:     message,
		ID:          dbTask.TaskDefinition.ID, //classTask.Id,
		CaaSVersion: cfg.Config.CaaSVersion,
		URL:         dbTask.Url,
		Task:        dbTask})
}

/*
validateJSONTask validates input (json) and generates a valid CLASS Task struct
*/
func validateJSONTask(r *http.Request) (*structs.CLASS_TASK, error) {
	log.Println("Rotterdam > CAAS > adapter [validateJsonTask] 'Duplicating' r.Body to use it multiple times with decoders ...")
	rBody := r.Body
	var buf bytes.Buffer
	tee := io.TeeReader(rBody, &buf)

	log.Println("Rotterdam > CAAS > adapter [validateJsonTask] Parsing default / old json definition ...")
	decoder := json.NewDecoder(tee)
	classTask, err := common.StructCheckClassTask(decoder)
	if err == nil {
		setClassStructDataAll(classTask, false)
		log.Println("Rotterdam > CAAS > adapter [validateJsonTask] Task with name " + classTask.Name + " received. Task ID is [" + classTask.ID + "]. Type = " + classTask.Type)
		return classTask, nil
	} else {
		log.Println("Rotterdam > CAAS > adapter [validateJsonTask] Parsing new json definition (COMPSs) ...")

		decoder = json.NewDecoder(&buf)
		classCOMPSsTask, err := common.StructCheckClassCOMPSsTask(decoder)
		if err == nil {
			classCOMPSsTask.ID = generateID() // ID
			classTask = structs.TransfCOMPSSTASKtoTASK(classCOMPSsTask)
			log.Println("Rotterdam > CAAS > adapter [validateJsonTask] Task with name " + classCOMPSsTask.Name + " received. Task ID is [" + classCOMPSsTask.ID + "]. Type = " + classTask.Type + ". Using 'COMPSs' JSON format.")
			setClassStructDataDefaultValues(classTask, true)
			return classTask, nil
		}
	}

	return nil, err
}

/*
deploy
*/
func deploy(w http.ResponseWriter, classTask *structs.CLASS_TASK) {
	adpt, err := GetAdapter(classTask.Cluster)
	if err != nil {
		errorMessage(w, err, "deploy")
	}

	log.Println("Rotterdam > CAAS > adapter [deploy] (2.) Deploying task [" + classTask.Name + "] ...")
	_, err = adpt.DeployTask(*classTask)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [deploy] (3.) Creating and starting SLA(s) ...")
		err = sla.CreateStartSLA(*classTask)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [deploy] ERROR when creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := common.ReadTaskValue(classTask.ID)
		if err == nil {
			okTaskMessage(w, "Task deployed", "deploy", *dbTask)
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "deploy",
				Message:     "Task deployed, but SLA not started",
				ID:          classTask.ID,
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		errorMessage(w, err, "deploy")
	}
}

/*
deployCOMPSs
*/
func deployCOMPSs(w http.ResponseWriter, classTask *structs.CLASS_TASK) {
	adpt, err := GetAdapter(classTask.Cluster)
	if err != nil {
		errorMessage(w, err, "deployCOMPSs")
	}

	log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (2.) Deploying task " + classTask.Name + " ...")
	_, err = adpt.DeployTaskCompss(*classTask)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (3.) Adding new Prometheus metric ...")
		err = sla.AddPromMetric("deadlines_missed_" + classTask.ID)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] ERROR adding a new Prometheus metric ", err)
		}

		log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] (4.) Creating and starting SLA ...")
		err = sla.CreateStartCOMPSsSLA(*classTask)
		if err != nil {
			log.Println("Rotterdam > CAAS > adapter [deployCOMPSs] ERROR creating and starting the SLA: ", err)
		}

		// send response with DBTask
		dbTask, err := common.ReadTaskValue(classTask.ID)
		if err == nil {
			okTaskMessage(w, "Task deployed", "deployCOMPSs", *dbTask)
		} else {
			json.NewEncoder(w).Encode(structs.ResponseCaaSTask{
				Resp:        "ok",
				Method:      "deployCOMPSs",
				Message:     "Task deployed, but SLA not started",
				ID:          classTask.ID,
				CaaSVersion: cfg.Config.CaaSVersion})
		}
	} else {
		errorMessage(w, err, "deployCOMPSs")
	}
}

///////////////////////////////////////////////////////////////////////////////

/*
InitializeAdapter initialization function
*/
func InitializeAdapter() {
	log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Initializing adapter for '" + cfg.Config.Clusters[0].Type + "' ...")
	if cfg.Config.Clusters[0].Type == "Openshift" {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Openshift adapter")
		a = ops_adapter.OpenshiftAdapter{}
	} else if cfg.Config.Clusters[0].Type == "Kubernetes" {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] Using Kubernetes adapter")
		a = k8s_adapter.KubernetesAdapter{}
	} else {
		log.Println("Rotterdam > CAAS > adapter [InitializeAdapter] ERROR No adapter selected")
	}
}

/*
GetAdapter initialization function
*/
func GetAdapter(idCluster string) (adapters.Adapter, error) {
	t := common.GetClusterType(idCluster)

	log.Println("Rotterdam > CAAS > adapter [GetAdapter] Getting adapter for '" + t + "' ...")
	if t == "Openshift" {
		log.Println("Rotterdam > CAAS > adapter [GetAdapter] Using Openshift adapter")
		return ops_adapter.OpenshiftAdapter{}, nil
	} else if t == "Kubernetes" || t == "microk8s" {
		log.Println("Rotterdam > CAAS > adapter [GetAdapter] Using Kubernetes adapter")
		return k8s_adapter.KubernetesAdapter{}, nil
	}

	log.Println("Rotterdam > CAAS > adapter [GetAdapter] ERROR No adapter selected")
	return nil, errors.New("adapter: No adapter selected")
}

///////////////////////////////////////////////////////////////////////////////

/*
DeployRotterdamTask Deploys a task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployRotterdamTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/tasks     <<New Task Deployment (undefined type)>>")

	log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] (1.) Validating json ...")
	classTask, err := validateJSONTask(r)
	if err == nil {
		classTaskStr, _ := common.CommClassStructToString(*classTask)
		log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Parsed CLASS TASK object (string): " + classTaskStr)

		if classTask.Type == constants.TypeTaskCOMPSs {
			log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] COMPSs task ...")
			deployCOMPSs(w, classTask)
		} else {
			log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] Default task ...")
			deploy(w, classTask)
		}
	} else {
		log.Println("Rotterdam > CAAS > adapter [DeployRotterdamTask] ERROR JSON not valid")
		errorMessage(w, err, "DeployRotterdamTask")
	}
}

/*
DeployTask Deploys a task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### New Default Task Deployment")

	log.Println("Rotterdam > CAAS > adapter [DeployTask] (1.) Validating json ...")
	classTask, err := validateJSONTask(r)
	if err != nil {
		errorMessage(w, err, "DeployTask")
	} else {
		deploy(w, classTask)
	}
}

/*
DeployTaskCOMPSs Deploys a COMPSs task (k8s: deployment & service & volumes ...). Input example: view CLASS_TASK in structs.go
*/
func DeployTaskCOMPSs(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### New COMPSs Task Deployment")

	log.Println("Rotterdam > CAAS > adapter [DeployTaskCOMPSs] (1.) Validating json ...")
	classTask, err := validateJSONTask(r)
	if err != nil {
		errorMessage(w, err, "DeployTaskCOMPSs")
	} else {
		deployCOMPSs(w, classTask)
	}
}

/*
RemoveRotterdamTask Deletes a task
*/
func RemoveRotterdamTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### DELETE /api/v1/tasks/{id}    <<Task termination (v2)>>")

	log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (1.) Reading input params ...")
	params := mux.Vars(r)

	log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (2.) Getting task [" + params["id"] + "] from DB ...")
	dbTask, err := common.ReadTaskValue(params["id"])
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (3.) Getting adapter for cluster [" + dbTask.ClusterId + "]  ...")
		adpt, err := GetAdapter(dbTask.ClusterId)
		if err != nil {
			errorMessage(w, err, "RemoveRotterdamTask")
		} else {
			log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (4.) Removing task " + dbTask.Id + " ...")
			resp, agreementID, err := adpt.RemoveTask(*dbTask)
			if err == nil {
				if agreementID == constants.SLANotDefined {
					log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (5.) SLA(s) not defined")
				} else {
					log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (5.) Stopping and terminating SLA(s) ...")
					sla.StopTerminateSLA(dbTask.Id)
				}

				log.Println("Rotterdam > CAAS > adapter [RemoveRotterdamTask] (6.) Deleting DBTask ...")
				common.DBDeleteTask(dbTask.Id)

				json.NewEncoder(w).Encode(structs.ResponseCaaS{
					Resp:        "ok",
					Method:      "RemoveRotterdamTask",
					Message:     "Task removed",
					CaaSVersion: cfg.Config.CaaSVersion,
					Content:     resp})
			} else {
				errorMessage(w, err, "RemoveRotterdamTask")
			}
		}
	} else {
		errorMessage(w, err, "RemoveRotterdamTask")
	}
}

/*
RemoveTask Deletes a task
*/
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### Task termination (DEPRECATED)")

	log.Println("Rotterdam > CAAS > adapter [RemoveTask] (1.) Reading input params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [RemoveTask] dock=" + params["dock"] + ", id=" + params["id"])

	log.Println("Rotterdam > CAAS > adapter [RemoveTask] (2.) Getting task [" + params["id"] + "] from DB ...")
	dbTask, err := common.ReadTaskValue(params["id"])
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [RemoveTask] (2.) Removing task " + params["id"] + " ...")
		resp, agreementID, err := a.RemoveTask(*dbTask)
		if err == nil {
			if agreementID == constants.SLANotDefined {
				log.Println("Rotterdam > CAAS > adapter [RemoveTask] (3.) SLA(s) not defined")
			} else {
				log.Println("Rotterdam > CAAS > adapter [RemoveTask] (3.) Stopping and terminating SLA(s) ...")
				sla.StopTerminateSLA(dbTask.Id)
			}

			log.Println("Rotterdam > CAAS > adapter [RemoveTask] (4.) Deleting DBTask ...")
			common.DBDeleteTask(dbTask.Id)

			json.NewEncoder(w).Encode(structs.ResponseCaaS{
				Resp:        "ok",
				Method:      "RemoveTask",
				Message:     "Task removed",
				CaaSVersion: cfg.Config.CaaSVersion,
				Content:     resp})
		} else {
			errorMessage(w, err, "RemoveTask")
		}
	} else {
		errorMessage(w, err, "RemoveTask")
	}
}

/*
GetAllTasks Gets all tasks
*/
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/tasks")

	log.Println("Rotterdam > CAAS > adapter [GetAllTasks] Getting all tasks ...")
	resp, err := adpcommon.GetAllTasks()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks retrieved"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasks{
			Resp:        "ok",
			Method:      "GetAllTasks",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetAllTasks")
	}
}

/*
GetAllTasksQoS Gets all tasks QoS
*/
func GetAllTasksQoS(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetAllTasksQoS] Getting QoS tasks ...")

	resp, err := adpcommon.GetAllTasksQoS()
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks qos retrieved"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasksQoS{
			Resp:        "ok",
			Method:      "GetAllTasksQoS",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			TasksQoS:    resp})
	} else {
		errorMessage(w, err, "GetAllTasksQoS")
	}
}

/*
LoadQoSDefinitions Loads QoS definitions list
*/
func LoadQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### POST /api/v1/qos/definitions")

	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Loading QoS definitions list ...")
	res, err := common.StructCheckClasQoSTemplateList(r)

	if err == nil {
		cfg.QosTemplates = *res

		msg := strconv.Itoa(len(*res)) + " tasks qos loaded"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasksQoS{
			Resp:        "ok",
			Method:      "LoadQoSDefinitions",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		errorMessage(w, err, "LoadQoSDefinitions")
	}
}

/*
GetQoSDefinitions Gets QoS defintions
*/
func GetQoSDefinitions(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/qos/definitions")

	log.Println("Rotterdam > CAAS > adapter [LoadQoSDefinitions] Retrieving QoS definitions list ...")
	json.NewEncoder(w).Encode(structs.ResponseQoSDefinitions{
		Resp:        "ok",
		Method:      "GetQoSDefinitions",
		Message:     "QoS definition retrieved",
		CaaSVersion: cfg.Config.CaaSVersion,
		QoSDefs:     cfg.QosTemplates})
}

/*
GetQoSDefinition Gets QoS defintion by name
*/
func GetQoSDefinition(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/qos/definitions/{name}")

	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetQoSDefinition] name=" + params["name"])

	qosTemplate, found := sla.GetQoSElem(params["name"])
	if found {
		json.NewEncoder(w).Encode(structs.ResponseQoSDefinition{
			Resp:        "ok",
			Method:      "GetQoSDefinition",
			Message:     "QoS definition retrieved",
			CaaSVersion: cfg.Config.CaaSVersion,
			QoSDef:      qosTemplate})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "GetQoSDefinition",
			Message:     "QoS definition not found. Check logs",
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}

/*
GetDockTasks Gets all tasks from Dock
*/
func GetDockTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] Reading input params ...")

	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] dock=" + params["dock"])

	log.Println("Rotterdam > CAAS > adapter [GetDockTasks] (1.) Getting all tasks from dock (" + params["dock"] + ") ...")
	resp, err := adpcommon.GetDockTasks(params["dock"])
	if err == nil {
		msg := strconv.Itoa(len(resp)) + " tasks retrieved [dock=" + params["dock"] + "]"
		json.NewEncoder(w).Encode(structs.ResponseCaaSTasks{
			Resp:        "ok",
			Method:      "GetDockTasks",
			Message:     msg,
			CaaSVersion: cfg.Config.CaaSVersion,
			Tasks:       resp})
	} else {
		errorMessage(w, err, "GetDockTasks")
	}
}

/*
GetTask Gets a task
*/
func GetTask(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/tasks/{id}")

	log.Println("Rotterdam > CAAS > adapter [GetTask] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetTask] dock=" + params["dock"] + ", id=" + params["id"])

	log.Println("Rotterdam > CAAS > adapter [GetTask] (1.) Getting task " + params["id"] + " ...")
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		okTaskMessage(w, "Task retrieved", "GetTask", dbTask)
	} else {
		errorMessage(w, err, "GetTask")
	}
}

/*
GetTaskAllInfo Gets a task including deployment info
*/
func GetTaskAllInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("### GET /api/v1/tasks/{id}/all")

	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] Reading input params ...")
	params := mux.Vars(r)
	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] id=" + params["id"])

	log.Println("Rotterdam > CAAS > adapter [GetTaskAllInfo] (1.) Getting task " + params["id"] + " ...")
	//dbTask, err := a.GetTaskAllInfo(params["id"])
	dbTask, err := adpcommon.GetTask(params["id"])
	if err == nil {
		adapter, err := GetAdapter(dbTask.ClusterId)
		if err == nil {
			adapter.GetTaskAllInfo(params["id"])
			okTaskMessage(w, "Task retrieved", "GetTaskAllInfo", dbTask)
		}
	} //else {
	errorMessage(w, err, "GetTaskAllInfo")
	//}
}

/*
GetConfig Gets configuration
*/
func GetConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetConfig] Getting configuration ...")
	resp, err := a.GetConfig()
	if err == nil {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "GetConfig",
			Message:     "Configuration retrieved",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     resp})
	} else {
		errorMessage(w, err, "GetConfig")
	}
}

/*
GetVersion Gets Rotterdam version
*/
func GetVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [GetVersion] Getting Rotterdam components versions ...")
	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:               "ok",
		Method:             "GetVersion",
		Message:            "Rotterdam components versions retrieved",
		CaaSVersion:        cfg.Config.CaaSVersion,
		RulesEngineVersion: cfg.Config.RulesEngineVersion,
		RestAPIVersion:     cfg.Config.RestApiVersion,
		SLALiteVersion:     cfg.Config.SLALiteVersion,
		IMECVersion:        cfg.Config.IMECVersion})
}

/*
ScaleUpDown ...
*/
func ScaleUpDown(dbtask structs.DB_TASK, replicas int) {
	log.Println("####################################################################################")
	log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling out / in the task ...")

	resp, err := a.ScaleUpDown(dbtask, replicas)
	if err == nil {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] Scaling up / down finalized: " + resp)
	} else {
		log.Println("Rotterdam > CAAS > adapter [ScaleUpDown] ERROR Scaling up / down task " + dbtask.TaskDefinition.Name + " [" + dbtask.Id + "] ...")
	}
}
