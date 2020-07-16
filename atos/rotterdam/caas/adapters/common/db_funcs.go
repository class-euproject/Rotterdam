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
// Created on 13 June 2019
// @author: Roi Sucasas - ATOS
//

package common

import (
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	"log"
)

/*
GetAllTasks Get all tasks
*/
func GetAllTasks() ([]structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [GetAllTasks] Getting tasks ...")

	dbtasks, err := common.DBReadAllTasks()
	if err == nil {
		return dbtasks, nil
	}

	log.Println("Rotterdam > CAAS > Adapters > Common [GetAllTasks] ERROR ", err)
	dbtasks = make([]structs.DB_TASK, 0)
	return dbtasks, err
}

/*
GetDockTasks get all tasks from dock
*/
func GetDockTasks(dock string) ([]structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [GetDockTasks] Getting tasks from dock [" + dock + "] ...")

	dbtasks, err := common.DBReadAllDockTasks(dock)
	if err == nil {
		return dbtasks, nil
	}

	log.Println("Rotterdam > CAAS > Adapters > Common [GetDockTasks] ERROR ", err)
	dbtasks = make([]structs.DB_TASK, 0)
	return dbtasks, err
}

/*
GetAllTasksQoS Get all tasks qos
*/
func GetAllTasksQoS() ([]structs.DB_TASK_QOS, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [GetAllTasksQoS] Getting tasks QoS ...")

	dbtasks, err := common.DBReadAllTasksQos()
	if err == nil {
		return dbtasks, nil
	}

	log.Println("Rotterdam > CAAS > Adapters > Common [GetAllTasksQoS] ERROR ", err)
	dbtasks = make([]structs.DB_TASK_QOS, 0)
	return dbtasks, err
}

/*
GetTask Returns a task
*/
func GetTask(idTask string) (structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > Adapters > Common [GetTask] Getting Task with id=" + idTask + " ...")

	// get task
	dbTask, err := common.ReadTaskValue(idTask)
	if err == nil {
		return *dbTask, nil
	}

	log.Println("Rotterdam > CAAS > Adapters > Common [GetTask] ERROR", err)
	return *dbTask, err
}
