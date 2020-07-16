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
	"log"
)

/*
GetCluster Returns the cluster identified by its name / id value
*/
func GetCluster(idCluster string) (*DB_INFRASTRUCTURE_CLUSTER, error) {
	dbInfr, err := ReadInfrValue(idCluster)
	if err == nil {
		return dbInfr, nil
	}

	log.Println("Rotterdam > IMEC > DB [GetCluster] ERROR ", err)
	return nil, err
}

/*
GetAllInfrs Get all infrastructures
*/
func GetAllInfrs() ([]DB_INFRASTRUCTURE_CLUSTER, error) {
	log.Println("Rotterdam > IMEC > DB [GetAllInfrs] Getting infrastructures ...")

	dbinfrs, err := DBReadAllInfrs()
	if err == nil {
		return dbinfrs, nil
	}

	log.Println("Rotterdam > IMEC > DB [GetAllInfrs] ERROR ", err)
	dbinfrs = make([]DB_INFRASTRUCTURE_CLUSTER, 0)
	return dbinfrs, err
}

/*
GetInfrByID Returns an infrastructure by id
*/
func GetInfrByID(id string) ([]DB_INFRASTRUCTURE_CLUSTER, error) {
	log.Println("Rotterdam > IMEC > DB [GetAllInfGetInfrByIDrs] Getting infrastructure [" + id + "] ...")

	var dbInfrs []DB_INFRASTRUCTURE_CLUSTER

	dbInfr, err := ReadInfrValue(id)
	if err == nil {
		dbInfrs = append(dbInfrs, *dbInfr)
		return dbInfrs, nil
	}

	log.Println("Rotterdam > IMEC > DB [GetInfrByID] ERROR ", err)
	dbInfrs = make([]DB_INFRASTRUCTURE_CLUSTER, 0)
	return dbInfrs, err
}
