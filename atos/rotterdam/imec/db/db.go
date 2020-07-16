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
// Created on 06 Feb 2020
// @author: Roi Sucasas - ATOS
//

package db

import (
	cfg "atos/rotterdam/config"
	constants "atos/rotterdam/globals/constants"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/buntdb"
)

const (
	// DB Tables
	DB_TABLE_INFRASTRUCTURE_CLUSTER = "INFRASTRUCTURE_CLUSTER"

	DB_INFRASTRUCTURE_CLUSTER_PREFIX = "INFR"
)

/*
DB_INFRASTRUCTURE_CLUSTER Infrastructure / Cluster definition
*/
type DB_INFRASTRUCTURE_CLUSTER struct {
	ID                            string `json:"id,omitempty"`
	Name                          string `json:"name,omitempty"`
	TableID                       string `json:"tableid,omitempty"`
	Description                   string `json:"description,omitempty"`
	Type                          string `json:"type,omitempty"`        // EDGE: 'microk8s', 'k3s'; CLUSTER: 'openshift', 'kubernetes'
	SO                            string `json:"so,omitempty"`          // S.O.: 'ubuntu18', 'ubuntu16', 'centos'
	DefaultDock                   string `json:"defaultDock,omitempty"` // 'class', 'default'
	HostIP                        string `json:"hostIP,omitempty"`
	HostPort                      int    `json:"hostPort,omitempty"`
	User                          string `json:"user,omitempty"`
	Password                      string `json:"password,omitempty"`
	KeyFile                       string `json:"keyFile,omitempty"`
	OpenshiftOauthToken           string `json:"openshiftOauthToken,omitempty"`
	KubernetesEndPoint            string `json:"kubernetesEndPoint,omitempty"`
	OpenshiftEndPoint             string `json:"openshiftEndPoint,omitempty"`
	SLALiteEndPoint               string `json:"slaliteEndPoint,omitempty"`
	PrometheusPushgatewayEndPoint string `json:"prometheusPushgatewayEndPoint,omitempty"`
	Status                        string `json:"status,omitempty"`
}

/*
IMECDatabase DATABASE
*/
var IMECDatabase *buntdb.DB = InitDB()

/*
InitDB  initializes database
*/
func InitDB() *buntdb.DB {
	log.Println("Rotterdam > IMEC > DB [InitDB] Initializating Database ...")
	// db, err := buntdb.Open("data.db")
	IMECDatabase, err := buntdb.Open(":memory:") // Open a file that does not persist to disk.
	if err != nil {
		log.Println("Rotterdam > IMEC > DB [InitDB] ERROR", err)
		return nil
	}
	return IMECDatabase
}

/*
CloseDB Closes database
*/
func CloseDB() {
	log.Println("Rotterdam > IMEC > DB [CloseDB] Closing Database ...")
	defer IMECDatabase.Close()
}

///////////////////////////////////////////////////////////////////////////////

/*
CommStringToStruct Parses a string to a struct of type DB_TASK
*/
func CommStringToDbStruct(ct string) (*DB_INFRASTRUCTURE_CLUSTER, error) {
	log.Println("Rotterdam > IMEC > DB [CommStringToDbStruct] string tp json object / struct [DB_INFRASTRUCTURE_CLUSTER]  ...")

	data := &DB_INFRASTRUCTURE_CLUSTER{}
	err := json.Unmarshal([]byte(ct), data)
	if err != nil {
		log.Println("Rotterdam > IMEC > DB [CommStringToDbStruct] ERROR", err)
		return data, err
	}

	return data, nil
}

/*
CommDbStructToString Parses a struct to a string
*/
func CommDbStructToString(ct DB_INFRASTRUCTURE_CLUSTER) (string, error) {
	log.Println("Rotterdam > IMEC > DB [CommDbStructToString] json object / struct [DB_INFRASTRUCTURE_CLUSTER] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > IMEC > DB [CommDbStructToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
StructCheckClassInfr
*/
func StructCheckClassInfr(req *http.Request) (*DB_INFRASTRUCTURE_CLUSTER, error) {
	log.Println("Rotterdam > IMEC > DB [StructCheckClassInfr] Checking json object ...")

	decoder := json.NewDecoder(req.Body)
	var t DB_INFRASTRUCTURE_CLUSTER
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > IMEC > DB [StructCheckClassInfr] ERROR (1)", err)
		return nil, err
	}

	tStr, err := CommDbStructToString(t)
	log.Println("Rotterdam > IMEC > DB [StructCheckClassInfr] Parsed object (string): " + tStr)
	log.Println("Rotterdam > IMEC > DB [StructCheckClassInfr] Sending parsed object ...")

	return &t, nil
}

///////////////////////////////////////////////////////////////////////////////
// DB_INFRASTRUCTURE_CLUSTER

/*
SetInfrValue ...
*/
func SetInfrValue(id string, dbInfr DB_INFRASTRUCTURE_CLUSTER) error {
	id = strings.Replace(id, DB_INFRASTRUCTURE_CLUSTER_PREFIX, "", 1)

	dbInfrStr, err := CommDbStructToString(dbInfr)
	if err == nil {
		err = IMECDatabase.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(DB_INFRASTRUCTURE_CLUSTER_PREFIX+id, dbInfrStr, nil)
			return err
		})
	}

	return err
}

/*
ReadInfrValue ...
*/
func ReadInfrValue(id string) (*DB_INFRASTRUCTURE_CLUSTER, error) {
	id = strings.Replace(id, DB_INFRASTRUCTURE_CLUSTER_PREFIX, "", 1)

	dbInfr := &DB_INFRASTRUCTURE_CLUSTER{}
	err := IMECDatabase.View(func(tx *buntdb.Tx) error {
		dbInfrStr, err := tx.Get(DB_INFRASTRUCTURE_CLUSTER_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > IMEC > DB [ReadInfrValue] ERROR", err)
			return err
		}

		log.Println("Rotterdam > IMEC > DB [ReadInfrValue] [DB_INFRASTRUCTURE_CLUSTER=" + dbInfrStr + "]")
		dbInfr, err = CommStringToDbStruct(dbInfrStr)
		return err
	})

	return dbInfr, err
}

/*
DBDeleteInfr ...
*/
func DBDeleteInfr(id string) (string, error) {
	id = strings.Replace(id, DB_INFRASTRUCTURE_CLUSTER_PREFIX, "", 1)

	err := IMECDatabase.Update(func(tx *buntdb.Tx) error {
		res, err := tx.Delete(DB_INFRASTRUCTURE_CLUSTER_PREFIX + id)
		if err != nil {
			log.Println("Rotterdam > IMEC > DB [DBDeleteInfr] ERROR", err)
			return err
		}

		log.Println("Rotterdam > IMEC > DB [DBDeleteInfr] Infrastructure (" + DB_INFRASTRUCTURE_CLUSTER_PREFIX + ")" + id + " deleted [" + res + "]")
		return err
	})

	return id, err
}

/*
DBReadAllInfrs ...
*/
func DBReadAllInfrs() ([]DB_INFRASTRUCTURE_CLUSTER, error) {
	log.Println("Rotterdam > IMEC > DB [DBReadAllInfrs] Getting All insfrastructures ...")
	var dbInfrs []DB_INFRASTRUCTURE_CLUSTER

	err := IMECDatabase.View(func(tx *buntdb.Tx) error {
		err2 := tx.Ascend("", func(key, value string) bool {
			log.Println(">>>> key: " + key + ", value: " + value)

			dbInfr, err := CommStringToDbStruct(value)
			if err == nil && dbInfr.TableID == DB_TABLE_INFRASTRUCTURE_CLUSTER {
				dbInfrs = append(dbInfrs, *dbInfr)
			}

			return true
		})
		return err2
	})

	return dbInfrs, err
}

/*
AddConfigInfrsToDB Adds initial (from config) infrs / clusters to DB
*/
func AddConfigInfrsToDB() {
	log.Println("Rotterdam > IMEC > DB [AddConfigInfrsToDB] Adding intial infrastructures / orchestrators to DB ...")

	for index := range cfg.Config.Clusters {
		log.Println("Rotterdam > IMEC > DB [AddConfigInfrsToDB] Adding " + cfg.Config.Clusters[index].ID + " ...")

		var dbObj *DB_INFRASTRUCTURE_CLUSTER
		dbObj = new(DB_INFRASTRUCTURE_CLUSTER)

		dbObj.ID = cfg.Config.Clusters[index].ID
		dbObj.Name = cfg.Config.Clusters[index].Name
		dbObj.Description = cfg.Config.Clusters[index].Description
		dbObj.HostIP = cfg.Config.Clusters[index].HostIP
		dbObj.HostPort = cfg.Config.Clusters[index].HostPort
		dbObj.Type = cfg.Config.Clusters[index].Type
		dbObj.User = cfg.Config.Clusters[index].User
		dbObj.Password = cfg.Config.Clusters[index].Password
		dbObj.KubernetesEndPoint = cfg.Config.Clusters[index].KubernetesEndPoint
		dbObj.OpenshiftEndPoint = cfg.Config.Clusters[index].OpenshiftEndPoint
		dbObj.SLALiteEndPoint = cfg.Config.Clusters[index].SLALiteEndPoint
		dbObj.PrometheusPushgatewayEndPoint = cfg.Config.Clusters[index].PrometheusPushgatewayEndPoint
		dbObj.OpenshiftOauthToken = cfg.Config.Clusters[index].OpenshiftOauthToken
		dbObj.TableID = DB_TABLE_INFRASTRUCTURE_CLUSTER
		dbObj.Status = constants.ClusterRUNNING

		err := SetInfrValue(dbObj.ID, *dbObj)
		if err != nil {
			log.Println("Rotterdam > IMEC > DB [AddConfigInfrsToDB] ERROR (1) Error adding info to database ", err)
		}
	}
}
