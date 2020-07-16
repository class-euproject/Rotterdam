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
	"log"
	"strings"

	"github.com/tidwall/buntdb"
)

/*
RDatabase DATABASE
*/
var RDatabase *buntdb.DB = InitDB()

/*
InitDB  initializes database
*/
func InitDB() *buntdb.DB {
	log.Println("Rotterdam > CAAS > DB [InitDB] Initializating Database ...")
	// db, err := buntdb.Open("data.db")
	RDatabase, err := buntdb.Open(":memory:") // Open a file that does not persist to disk.
	if err != nil {
		log.Println("Rotterdam > CAAS > DB [InitDB] ERROR", err)
		return nil
	}
	return RDatabase
}

/*
CloseDB Closes database
*/
func CloseDB() {
	log.Println("Rotterdam > CAAS > DB [CloseDB] Closing Database ...")
	defer RDatabase.Close()
}

///////////////////////////////////////////////////////////////////////////////
// DB_TASK

/*
SetTaskValue ...
*/
func SetTaskValue(id string, dbtask structs.DB_TASK) error {
	id = strings.Replace(id, structs.DB_TASK_PREFIX, "", 1)

	dbtaskStr, err := CommDbTaskStructToString(dbtask)
	if err == nil {
		err = RDatabase.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(structs.DB_TASK_PREFIX+id, dbtaskStr, nil)
			return err
		})
	}

	return err
}

/*
ReadTaskValue ...
*/
func ReadTaskValue(id string) (*structs.DB_TASK, error) {
	id = strings.Replace(id, structs.DB_TASK_PREFIX, "", 1)

	dbtask := &structs.DB_TASK{}
	err := RDatabase.View(func(tx *buntdb.Tx) error {
		dbtaskStr, err := tx.Get(structs.DB_TASK_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > CAAS > DB [ReadTaskValue] ERROR", err)
			return err
		}

		log.Println("Rotterdam > CAAS > DB [ReadTaskValue] [DB_TASK=" + dbtaskStr + "]")
		dbtask, err = CommStringToDbTaskStruct(dbtaskStr)
		return err
	})

	return dbtask, err
}

/*
DBDeleteTask ...
*/
func DBDeleteTask(id string) (string, error) {
	id = strings.Replace(id, structs.DB_TASK_PREFIX, "", 1)

	err := RDatabase.Update(func(tx *buntdb.Tx) error {
		res, err := tx.Delete(structs.DB_TASK_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > CAAS > DB [DBDeleteTask] ERROR", err)
			return err
		}

		log.Println("Rotterdam > CAAS > DB [DBDeleteTask] Task (" + structs.DB_TASK_PREFIX + ")" + id + " deleted [" + res + "]")
		return err
	})

	return id, err
}

/*
DBReadAllTasks ...
*/
func DBReadAllTasks() ([]structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > DB [DBReadAllTasks] Getting All tasks ...")
	var dbtasks []structs.DB_TASK

	err := RDatabase.View(func(tx *buntdb.Tx) error {
		err2 := tx.Ascend("", func(key, value string) bool {
			log.Println(">>>> key: " + key + ", value: " + value)

			dbtask, err := CommStringToDbTaskStruct(value)
			if err == nil && dbtask.DbId == structs.DB_TABLE_TASK {
				dbtasks = append(dbtasks, *dbtask)
			}

			return true
		})
		return err2
	})

	return dbtasks, err
}

/*
DBReadAllDockTasks ...
*/
func DBReadAllDockTasks(dock string) ([]structs.DB_TASK, error) {
	log.Println("Rotterdam > CAAS > DB [DBReadAllTasks] Getting All tasks from namespace [" + dock + "] ...")
	var dbtasks []structs.DB_TASK

	err := RDatabase.View(func(tx *buntdb.Tx) error {
		_ = tx.Ascend("", func(key, value string) bool {
			log.Println(">>>> key: " + key + ", value: " + value)

			dbtask, err := CommStringToDbTaskStruct(value)
			if err == nil && dbtask.DbId == structs.DB_TABLE_TASK && dbtask.NameSpace == dock {
				dbtasks = append(dbtasks, *dbtask)
			}

			return true
		})
		return nil
	})

	return dbtasks, err
}

///////////////////////////////////////////////////////////////////////////////
// DB_TASK_QOS

/*
SetTaskQoSValue ...
*/
func SetTaskQoSValue(id string, dbtaskqos structs.DB_TASK_QOS) error {
	id = strings.Replace(id, structs.DB_TASK_QOS_PREFIX, "", 1)

	dbtaskqosStr, err := CommDbTaskQoSStructToString(dbtaskqos)
	if err == nil {
		err = RDatabase.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(structs.DB_TASK_QOS_PREFIX+id, dbtaskqosStr, nil)
			return err
		})
	}

	return err
}

/*
ReadTaskQoSValue ...
*/
func ReadTaskQoSValue(id string) (*structs.DB_TASK_QOS, error) {
	id = strings.Replace(id, structs.DB_TASK_QOS_PREFIX, "", 1)

	dbtaskqos := &structs.DB_TASK_QOS{}
	err := RDatabase.View(func(tx *buntdb.Tx) error {
		dbtaskqosStr, err := tx.Get(structs.DB_TASK_QOS_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > CAAS > DB [ReadTaskQoSValue] ERROR", err)
			return err
		}

		log.Println("Rotterdam > CAAS > DB [ReadTaskQoSValue] [DB_TASK_QOS=" + dbtaskqosStr + "]")
		dbtaskqos, err = CommStringToDbTaskQoSStruct(dbtaskqosStr)
		return err
	})

	return dbtaskqos, err
}

/*
DBReadTasksQosByAgreement ...

func DBReadTasksQosByAgreement(idAgreement string) ([]structs.DB_TASK_QOS, error) {
	dbtasksqos := make([]structs.DB_TASK_QOS, 0)
	err := RDatabase.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			dbtaskqos, err := CommStringToDbTaskQoSStruct(value)
			if err != nil {
				log.Println("Rotterdam > CAAS > DB [DBReadTasksQosByAgreement] ERROR", err)
				return false
			}

			dbtasksqos = append(dbtasksqos, *dbtaskqos)
			log.Println("Rotterdam > CAAS > DB [DBReadTasksQosByAgreement] key: " + key + ", value: " + value)
			return true
		})
		return err
	})

	return dbtasksqos, err
}*/

/*
DBDeleteTaskQos ...
*/
func DBDeleteTaskQos(id string) (string, error) {
	id = strings.Replace(id, structs.DB_TASK_QOS_PREFIX, "", 1)

	err := RDatabase.Update(func(tx *buntdb.Tx) error {
		res, err := tx.Delete(structs.DB_TASK_QOS_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > CAAS > DB [DBDeleteTaskQos] ERROR", err)
			return err
		}

		log.Println("Rotterdam > CAAS > DB [DBDeleteTaskQos] TaskQoS (" + structs.DB_TASK_QOS_PREFIX + ")" + id + " deleted [" + res + "]")
		return err
	})

	return id, err
}

/*
DBReadAllTasksQos ...
*/
func DBReadAllTasksQos() ([]structs.DB_TASK_QOS, error) {
	dbtasksqos := make([]structs.DB_TASK_QOS, 0)
	err := RDatabase.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, structs.DB_TASK_QOS_PREFIX) {
				dbtaskqos, err := CommStringToDbTaskQoSStruct(value)
				if err != nil {
					log.Println("Rotterdam > CAAS > DB [DBReadAllTasksQos] ERROR", err)
					return false
				}

				dbtasksqos = append(dbtasksqos, *dbtaskqos)
				log.Println("Rotterdam > CAAS > DB [DBReadAllTasksQos] key: " + key + ", value: " + value)
				return true
			}

			return false
		})
		return err
	})

	return dbtasksqos, err
}
