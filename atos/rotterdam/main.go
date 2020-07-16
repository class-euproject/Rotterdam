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

package main

import (
	urls "atos/rotterdam/caas/adapters"
	cfg "atos/rotterdam/config"
	imec_db "atos/rotterdam/imec/db"
	rest_api "atos/rotterdam/rest-api"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

// fileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/*
 * Main: launch server
 */
func main() {
	log.Printf("................................................................................")

	// set global properties from configuration file
	cfgPath := "./config/config.json"

	// arguments: [1] ... alternative configuration file
	if len(os.Args) == 2 {
		log.Println("Rotterdam > Arguments[1] = " + os.Args[1])
		if fileExists(os.Args[1]) {
			log.Println("Rotterdam > Using configuration from " + os.Args[1])
			cfgPath = os.Args[1]
		}
	} else {
		log.Println("Rotterdam > Using default configuration file")
	}
	log.Println("Rotterdam > Rest API > Getting configuration values from file [" + cfgPath + "] ...")

	cfg.Config = cfg.Configuration{}

	file, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg.Config)
	if err != nil {
		panic(err)
	}

	// check ENV values and modify global properties if specified
	cfg.InitConfig(&cfg.Config)

	// show global configuration
	log.Printf("Rotterdam > Rest API > Rest Api Version ...... %s", cfg.Config.RestApiVersion)
	log.Printf("Rotterdam > Rest API > Rules Engine Version .. %s", cfg.Config.RulesEngineVersion)
	log.Printf("Rotterdam > Rest API > CaaS Version .......... %s", cfg.Config.CaaSVersion)
	log.Printf("Rotterdam > Rest API > SLALite Version ....... %s", cfg.Config.SLALiteVersion)
	log.Printf("Rotterdam > Rest API > IMEC Version .......... %s", cfg.Config.IMECVersion)
	log.Printf("Rotterdam > Rest API > Server Port ........... %d", cfg.Config.ServerPort)
	log.Printf("Rotterdam > Rest API > Clusters Configurations:")
	for index := range cfg.Config.Clusters {
		log.Printf("Rotterdam > Rest API >  => %s", cfg.Config.Clusters[index].Name)
		log.Printf("Rotterdam > Rest API >     ID: %s", cfg.Config.Clusters[index].ID)
		log.Printf("Rotterdam > Rest API >     Description: %s", cfg.Config.Clusters[index].Description)
		log.Printf("Rotterdam > Rest API >     Type: %s", cfg.Config.Clusters[index].Type)
		log.Printf("Rotterdam > Rest API >     Default Dock / Namespace: %s", cfg.Config.Clusters[index].DefaultDock)
		log.Printf("Rotterdam > Rest API >     S.O.: %s", cfg.Config.Clusters[index].SO)
		log.Printf("Rotterdam > Rest API >     Host IP: %s", cfg.Config.Clusters[index].HostIP)
		log.Printf("Rotterdam > Rest API >     Kubernetes Endpoint: %s", cfg.Config.Clusters[index].KubernetesEndPoint)
		log.Printf("Rotterdam > Rest API >     Openshift Endpoint: %s", cfg.Config.Clusters[index].OpenshiftEndPoint)
		log.Printf("Rotterdam > Rest API >     SLALite Endpoint: %s", cfg.Config.Clusters[index].SLALiteEndPoint)
		log.Printf("Rotterdam > Rest API >     Prometheus Pushgateway Endpoint: %s", cfg.Config.Clusters[index].PrometheusPushgatewayEndPoint)
	}
	log.Printf("Rotterdam > Rest API > Swagger UI ............ http://localhost:8333/swaggerui/index.html")
	log.Println("Rotterdam > Rest API > Starting server [port " + strconv.Itoa(cfg.Config.ServerPort) + "] ...")

	// initialize URLs
	urls.Initialize()

	// adding clusters to DB
	imec_db.AddConfigInfrsToDB()

	// initialize REST API
	rest_api.InitializeRESTAPI()
}
