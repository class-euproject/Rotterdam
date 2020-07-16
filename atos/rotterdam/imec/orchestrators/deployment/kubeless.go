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
// Created on 14 April 2020
// @author: Roi Sucasas - ATOS
//

package deployment

import (
	constants "atos/rotterdam/globals/constants"
	db "atos/rotterdam/imec/db"
	"errors"
	"log"
	"strconv"

	sshclient "github.com/helloyi/go-sshclient"
)

// Kubeless deployment in UBUNTU 18
func deployKubelessUbuntu18(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] Deploying Kubeless in " + infr.HostIP + " ...")

	// update infr status
	infr.Status = constants.ClusterDEPLOYING
	err := db.SetInfrValue(infr.ID, *infr)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR ", err)
	}

	// CONNECTION
	// example: DialWithPasswd("192.168.1.101:22", "vagrant", "vagrant")
	client, err := sshclient.DialWithPasswd(infr.HostIP+":"+strconv.Itoa(infr.HostPort), infr.User, infr.Password)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR (1) ", err)
		return err
	}

	// CHECK CONNECTION
	err = client.Cmd("ls -la").Run()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR (2) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] << Connected to " + infr.User + "@" + infr.HostIP + ":" + strconv.Itoa(infr.HostPort) + " >>")

	// SCRIPT
	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] 1. Executing script ...")

	script := `
	  #!/bin/bash
	  export RELEASE=$(curl -s https://api.github.com/repos/kubeless/kubeless/releases/latest | grep tag_name | cut -d '"' -f 4)
	  microk8s.kubectl create ns kubeless
	  microk8s.kubectl create -f https://github.com/kubeless/kubeless/releases/download/$RELEASE/kubeless-$RELEASE.yaml
	  sudo apt install unzip

	  export OS=$(uname -s| tr '[:upper:]' '[:lower:]')
	  curl -OL https://github.com/kubeless/kubeless/releases/download/$RELEASE/kubeless_$OS-amd64.zip
	  unzip kubeless_$OS-amd64.zip
	  sudo mv bundles/kubeless_$OS-amd64/kubeless /usr/local/bin/
	  
	  microk8s.kubectl config view --flatten --minify > $HOME/.kube/config
	`
	out, err := client.Script(script).Output()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR (3) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] 2. OUT: " + string(out))

	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] Kubeless deployed and listening: http://" + infr.HostIP + ":8001")

	// save infr status
	// infr.KubernetesEndPoint = "http://" + infr.HostIP + ":8001"
	infr.Status = constants.ClusterRUNNING
	err = db.SetInfrValue(infr.ID, *infr)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR (4) ", err)
		return err
	}

	defer client.Close()

	return nil
}

// Kubeless deployment in UBUNTU 16
func deployClusterKubelessUbuntu16(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	err := deployClusterUbuntu16(infr)
	if err == nil {
		return deployKubelessUbuntu18(infr)
	}
	return err
}

// Kubeless deployment in UBUNTU 18
func deployClusterKubelessUbuntu18(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	err := deployClusterUbuntu18(infr)
	if err == nil {
		return deployKubelessUbuntu18(infr)
	}
	return err
}

/*
Kubeless ...
*/
func Kubeless(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] Deploying Kubeless in " + infr.HostIP + " ...")

	if infr.SO == "ubuntu16" {
		return deployClusterKubelessUbuntu16(infr)
	} else if infr.SO == "ubuntu18" {
		return deployClusterKubelessUbuntu18(infr)
	} else {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [Kubeless] ERROR S.O. not defined: [" + infr.SO + "]")
		return errors.New("S.O. not defined: [" + infr.SO + "]")
	}
}
