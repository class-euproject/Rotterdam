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
	cfg "atos/rotterdam/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////

// creates the request's body' from a JSON
func httpJSONBody(bodyJSON interface{}) io.Reader {
	bodyBytes, err := json.Marshal(bodyJSON)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [httpJSONBody] ERROR (1)", err)
		return nil
	}
	return bytes.NewReader(bodyBytes)
}

// creates the request's body' from a string
func httpRawDataBody(bodyRawData string) io.Reader {
	return bytes.NewReader([]byte(bodyRawData))
}

// httpRequest prepares and executes the HTTP request //// bodyJSON interface{}
func httpRequest(httpMethod string, url string, auth bool, body io.Reader) (int, []byte, error) {
	log.Println("Rotterdam > CAAS > http [httpRequest] " + httpMethod + " request [" + url + "], auth [" + strconv.FormatBool(auth) + "] ...")

	// create request with headers and body
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [httpRequest] ERROR (2)", err)
		return 0, nil, err
	}

	// Content-Type: json / json-patch+json
	if httpMethod == "PATCH" {
		log.Println("Rotterdam > CAAS > http [httpRequest] Content-Type = application/json-patch+json")
		req.Header.Set("Content-Type", "application/json-patch+json")
	} else {
		log.Println("Rotterdam > CAAS > http [httpRequest] Content-Type = application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// Authorization header
	if auth == true {
		log.Println("Rotterdam > CAAS > http [httpRequest] Using Authorization Bearer ...")
		req.Header.Set("Authorization", "Bearer "+cfg.Config.Clusters[0].OpenshiftOauthToken)
	}

	// CLIENT
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// execute HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [httpRequest] ERROR (3)", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	// get data from response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [httpRequest] ERROR (4)", err)
		return resp.StatusCode, nil, err
	} else if resp.StatusCode >= 400 { // check errors => StatusCode
		log.Println("Rotterdam > CAAS > http [httpRequest] ERROR (5) StatusCode >= 400")
		return resp.StatusCode, nil, errors.New("Rotterdam > CAAS > http [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode) + "")
	}

	log.Println("Rotterdam > CAAS > http [httpRequest] HTTP STATUS: (" + strconv.Itoa(resp.StatusCode) + ") " + http.StatusText(resp.StatusCode))
	//log.Println("Rotterdam > CAAS > http [HTTPRequest] RESPONSE: " + string(data))

	return resp.StatusCode, data, nil
}

///////////////////////////////////////////////////////////////////////////////
// GET

/*
HTTPGET generic GET request
*/
func HTTPGET(url string, auth bool) (int, []byte, error) {
	return httpRequest("GET", url, auth, nil)
}

/*
HTTPGETStruct GET request that returns a struct of type 'map[string]interface{}'
*/
func HTTPGETStruct(url string, auth bool) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HTTPGETStruct] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPGETStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPGETStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

/*
HTTPGETString GET request that returns a string (response)
*/
func HTTPGETString(url string, auth bool) (int, string, error) {
	log.Println("Rotterdam > CAAS > http [HTTPGETString] GET request [" + url + "] ...")

	status, data, err := HTTPGET(url, auth)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPGETString] ERROR (1)", err)
		return status, "", err
	}

	return status, string(data), nil
}

///////////////////////////////////////////////////////////////////////////////
// POST

/*
HTTPPOSTRawData Generic POST request
*/
func HTTPPOSTRawData(url string, auth bool, bodyRawData string) (int, []byte, error) {
	return httpRequest("POST", url, auth, httpRawDataBody(bodyRawData))
}

/*
HTTPPOST Generic POST request
*/
func HTTPPOST(url string, auth bool, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("POST", url, auth, httpJSONBody(bodyJSON))
}

/*
HTTPPOSTStruct POST request that returns a struct of type 'map[string]interface{}'
*/
func HTTPPOSTStruct(url string, auth bool, bodyJSON interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HTTPPOSTStruct] POST request [" + url + "] ...")

	status, data, err := HTTPPOST(url, auth, bodyJSON)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPOSTStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPOSTStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// DELETE

/*
HTTPDELETE Generic DELETE request
*/
func HTTPDELETE(url string, auth bool, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("DELETE", url, auth, httpJSONBody(bodyJSON))
}

/*
HTTPDELETEStruct DELETE request that returns a struct of type 'map[string]interface{}'
*/
func HTTPDELETEStruct(url string, auth bool) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HTTPDELETEStruct] DELETE request [" + url + "] ...")

	type Body struct {
		Content interface{}
	}

	status, data, err := HTTPDELETE(url, auth, Body{})
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPDELETEStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPDELETEStruct] WARNING (1)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PUT

/*
HTTPPUT Generic PUT request
*/
func HTTPPUT(url string, auth bool, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("PUT", url, auth, httpJSONBody(bodyJSON))
}

/*
HTTPPUTStruct PUT request that returns a struct of type 'map[string]interface{}'
*/
func HTTPPUTStruct(url string, auth bool, bodyJSON interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HTTPPUTStruct] PUT request [" + url + "] ...")

	status, data, err := HTTPPUT(url, auth, bodyJSON)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPUTStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPUTStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}

///////////////////////////////////////////////////////////////////////////////
// PATCH

/*
HTTPPATCH Generic PATCH request
*/
func HTTPPATCH(url string, auth bool, bodyJSON interface{}) (int, []byte, error) {
	return httpRequest("PATCH", url, auth, httpJSONBody(bodyJSON))
}

/*
HTTPPATCHStruct PATCH request that returns a struct of type 'map[string]interface{}'
*/
func HTTPPATCHStruct(url string, auth bool, bodyJSON interface{}) (int, map[string]interface{}, error) {
	log.Println("Rotterdam > CAAS > http [HTTPPATCHStruct] PATCH request [" + url + "] ...")

	status, data, err := HTTPPATCH(url, auth, bodyJSON)
	if err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPATCHStruct] ERROR (1)", err)
		return status, nil, err
	}

	// create json
	var objmap map[string]interface{}
	if err := json.Unmarshal(data, &objmap); err != nil {
		log.Println("Rotterdam > CAAS > http [HTTPPATCHStruct] ERROR (2)", err)
		return status, nil, err
	}

	return status, objmap, nil
}
