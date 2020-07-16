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
// Created on 28 Feb 2020
// @author: Roi Sucasas - ATOS
//

package infrastructures

import (
	db "atos/rotterdam/imec/db"
	"encoding/json"
	"log"
	"net/http"
)

/*
INFRASTRUCTURE represents a host / edge device / cluster access point
*/
type INFRASTRUCTURE struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`        // EDGE: 'microk8s', 'k3s'; CLUSTER: 'openshift', 'kubernetes'
	SO          string `json:"so,omitempty"`          // S.O.: 'ubuntu18', 'ubuntu16', 'centos'
	DefaultDock string `json:"defaultDock,omitempty"` // 'class', 'default'
	HostIP      string `json:"hostIP,omitempty"`
	HostPort    int    `json:"hostPort,omitempty"`
	User        string `json:"user,omitempty"`
	Password    string `json:"password,omitempty"`
	KeyFile     string `json:"keyFile,omitempty"`
}

/*
infrJSONToString Parses a struct to a string
*/
func infrJSONToString(ct INFRASTRUCTURE) (string, error) {
	log.Println("Rotterdam > IMEC > infrastructures [infrJSONToString] json object / struct [INFRASTRUCTURE] to string ...")

	out, err := json.Marshal(ct)
	if err != nil {
		log.Println("Rotterdam > IMEC > infrastructures [infrJSONToString] ERROR", err)
		return "", err
	}

	return string(out), nil
}

/*
ValidateStructInfrJSON Validate structure INFRASTRUCTURE
*/
func ValidateStructInfrJSON(req *http.Request) (*INFRASTRUCTURE, error) {
	log.Println("Rotterdam > IMEC > infrastructures [ValidateStructInfrJSON] Checking json object ...")

	decoder := json.NewDecoder(req.Body)
	var t INFRASTRUCTURE
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Rotterdam > IMEC > infrastructures [ValidateStructInfrJSON] ERROR (1)", err)
		return nil, err
	}

	tStr, err := infrJSONToString(t)
	log.Println("Rotterdam > IMEC > infrastructures [ValidateStructInfrJSON] Parsed object (string): " + tStr)
	log.Println("Rotterdam > IMEC > infrastructures [ValidateStructInfrJSON] Sending parsed object ...")

	return &t, nil
}

/*
InfrJSONToDBObj Transforms a INFRASTRUCTURE into a DB_INFRASTRUCTURE_CLUSTER object
*/
func InfrJSONToDBObj(i INFRASTRUCTURE, newID string) *db.DB_INFRASTRUCTURE_CLUSTER {
	log.Println("Rotterdam > IMEC > infrastructures [InfrJSONToDBObj] Transforming INFRASTRUCTURE object ...")

	var dbObj *db.DB_INFRASTRUCTURE_CLUSTER
	dbObj = new(db.DB_INFRASTRUCTURE_CLUSTER)

	dbObj.ID = newID
	dbObj.Name = i.Name
	dbObj.Description = i.Description
	dbObj.HostIP = i.HostIP
	dbObj.HostPort = i.HostPort
	dbObj.Type = i.Type
	dbObj.SO = i.SO
	dbObj.User = i.User
	dbObj.Password = i.Password
	dbObj.TableID = db.DB_TABLE_INFRASTRUCTURE_CLUSTER
	if len(i.DefaultDock) == 0 {
		dbObj.DefaultDock = "default"
	} else {
		dbObj.DefaultDock = i.DefaultDock
	}

	log.Println("Rotterdam > IMEC > infrastructures [InfrJSONToDBObj] Returning DB_INFRASTRUCTURE_CLUSTER object ...")
	return dbObj
}
