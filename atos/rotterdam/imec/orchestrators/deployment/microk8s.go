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

package deployment

import (
	constants "atos/rotterdam/globals/constants"
	db "atos/rotterdam/imec/db"
	"errors"
	"log"
	"strconv"

	sshclient "github.com/helloyi/go-sshclient"
)

// MicroK8s deployment in UBUNTU 18
func deployClusterUbuntu18(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] Deploying MicroK8s in " + infr.HostIP + " ...")

	// update infr status
	infr.Status = constants.ClusterDEPLOYING
	err := db.SetInfrValue(infr.ID, *infr)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR ", err)
	}

	// CONNECTION
	// example: DialWithPasswd("192.168.1.101:22", "vagrant", "vagrant")
	client, err := sshclient.DialWithPasswd(infr.HostIP+":"+strconv.Itoa(infr.HostPort), infr.User, infr.Password)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (1) ", err)
		return err
	}

	// CHECK CONNECTION
	err = client.Cmd("ls -la").Run()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (2) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] << Connected to " + infr.User + "@" + infr.HostIP + ":" + strconv.Itoa(infr.HostPort) + " >>")

	// INSTALL MICROK8S
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 1. Installing microk8s ('sudo snap install microk8s --classic --channel=1.17/stable') ...")
	out, err := client.Cmd("sudo snap install microk8s --classic --channel=1.17/stable").Output()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (3) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 1. OUT: " + string(out))

	// SUDO USERMOD
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 2. Executing 'sudo usermod -a -G microk8s <USER>' ...")
	out, err = client.Cmd("sudo usermod -a -G microk8s " + infr.User).Output()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (4) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 2. OUT: " + string(out))

	// CLOSE CONNECTION
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 3. Closing connection ...")
	client.Close()

	// NEW CONNECTION to apply changes
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 4. Connecting again: su - <USER>")
	client, err = sshclient.DialWithPasswd(infr.HostIP+":"+strconv.Itoa(infr.HostPort), infr.User, infr.Password)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (5) ", err)
		return err
	}

	// CHECK CONNECTION
	err = client.Cmd("ls -la").Run()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (6) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] << Connected to " + infr.User + "@" + infr.HostIP + ":" + strconv.Itoa(infr.HostPort) + " >>")

	// SCRIPT
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 5. Executing script ...")

	script := `
	  #!/bin/bash
	  microk8s.kubectl proxy --port=8001 --address='` + infr.HostIP + `' --accept-hosts='.*' &>/dev/null &
	`
	out, err = client.Script(script).Output()
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (7) ", err)
		return err
	}
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] 5. OUT: " + string(out))

	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] MicroK8s deployed and listening: http://" + infr.HostIP + ":8001")

	// save infr status
	infr.KubernetesEndPoint = "http://" + infr.HostIP + ":8001"
	infr.Status = constants.ClusterRUNNING
	err = db.SetInfrValue(infr.ID, *infr)
	if err != nil {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR (8) ", err)
		return err
	}

	defer client.Close()

	return nil
}

// MicroK8s deployment in UBUNTU 16
func deployClusterUbuntu16(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	return deployClusterUbuntu18(infr)
}

/*
MicroK8s ...
*/
func MicroK8s(infr *db.DB_INFRASTRUCTURE_CLUSTER) error {
	log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] Deploying MicroK8s in " + infr.HostIP + " ...")

	if infr.SO == "ubuntu16" {
		return deployClusterUbuntu16(infr)
	} else if infr.SO == "ubuntu18" {
		return deployClusterUbuntu18(infr)
	} else {
		log.Println("Rotterdam > IMEC > orchestrators > deployment [MicroK8s] ERROR S.O. not defined: [" + infr.SO + "]")
		return errors.New("S.O. not defined: [" + infr.SO + "]")
	}
}
