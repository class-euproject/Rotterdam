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
	common "atos/rotterdam/caas/common"
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"encoding/json"
	"log"
	"net/http"
)

///////////////////////////////////////////////////////////////////////////////
// DEFAULT AND TEST FUNCTIONS

/*
NotImplementedFunc Default Function for not implemented calls
*/
func NotImplementedFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [NotImplementedFunc] -Not implemented-")

	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "ok",
		Method:      "NotImplementedFunc",
		Message:     "not implemented",
		CaaSVersion: cfg.Config.CaaSVersion})
}

// Default Function for not implemented calls
func HomePath(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [HomePath] Returning SwaggerUI path ...")

	json.NewEncoder(w).Encode(structs.ResponseCaaS{
		Resp:        "ok",
		Method:      "HomePath",
		Message:     "UI URL: /swaggerui/",
		CaaSVersion: cfg.Config.CaaSVersion})
}

/*
TestGetRequest Test GET function
*/
func TestGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [TestGetRequest] Testing GET request")

	_, objmap, err := common.HTTPGETStruct("http://postman-echo.com/get?foo1=bar1&foo2=bar2", true)
	if err != nil {
		log.Println("Rotterdam > CAAS > adapter [TestGetRequest] ERROR", err)
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "TestGetRequest",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	} else {
		log.Println("Rotterdam > CAAS > adapter [TestGetRequest] url: " + objmap["url"].(string))
		log.Println("Rotterdam > CAAS > adapter [TestGetRequest] args/foo1: " + objmap["args"].(map[string]interface{})["foo1"].(string))

		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "TestGetRequest",
			Message:     "Test",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     objmap["url"].(string)})
	}
}

type test_struct struct {
	DeviceID  string `json:"deviceID,omitempty"`
	TestID    string `json:"testID,omitempty"`
	ComplexID struct {
		SubDeviceID string `json:"deviceID,omitempty"`
		SubTestID   string `json:"testID,omitempty"`
	} `json:"complexID,omitempty"`
	Complexes []string `json:"complexes,omitempty"`
}

/*
TestPostRequest Test POST function
*/
func TestPostRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Rotterdam > CAAS > adapter [TestPostRequest] Testing POST request")

	decoder := json.NewDecoder(r.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "TestPostRequest",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
	if t.DeviceID != "" {
		log.Println("Rotterdam > CAAS > adapter [TestPostRequest] deviceID: " + t.DeviceID)
	}
	if t.TestID != "" {
		log.Println("Rotterdam > CAAS > adapter [TestPostRequest] testID: " + t.TestID)
	} else {
		log.Println("Rotterdam > CAAS > adapter [TestPostRequest] testID is nil")
	}
	if t.ComplexID.SubDeviceID == "" {
		log.Println("Rotterdam > CAAS > adapter [TestPostRequest] ComplexID is nil ")
	}

	type Payload struct {
		Type     string      `json:"type"`
		Name     string      `json:"name"`
		Data     string      `json:"data"`
		Priority interface{} `json:"priority"`
		Port     interface{} `json:"port"`
		Weight   interface{} `json:"weight"`
	}

	data := Payload{
		Type:     "tipo1",
		Name:     "name1",
		Data:     "datadata",
		Priority: 1,
		Port:     8080,
		Weight:   300}

	_, respBytes, err := common.HTTPPOST("http://postman-echo.com/post", true, data)
	if err == nil {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "ok",
			Method:      "TestPostRequest",
			Message:     "Test",
			CaaSVersion: cfg.Config.CaaSVersion,
			Content:     string(respBytes)})
	} else {
		json.NewEncoder(w).Encode(structs.ResponseCaaS{
			Resp:        "error",
			Method:      "TestPostRequest",
			Message:     err.Error(),
			CaaSVersion: cfg.Config.CaaSVersion})
	}
}
