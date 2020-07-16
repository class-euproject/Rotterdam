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

package sla

import (
	structs "atos/rotterdam/caas/common/structs"
	cfg "atos/rotterdam/config"
	"log"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

// create a DB Task QoS element
func genDbTaskQos(taskID string, elem cfg.CLASS_QOS_TEMPLATE__ELEM) *structs.DB_TASK_QOS {
	dbtaskqos := &structs.DB_TASK_QOS{
		DbId:            structs.DB_TABLE_QOS,
		Id:              taskID,
		IdTask:          taskID,
		Type:            elem.Type,
		Action:          elem.Action,
		ScaleFactor:     elem.ScaleFactor,
		MaxReplicas:     30,
		MinReplicas:     1,
		Guarantee:       elem.GuaranteeName,
		TotalViolations: 0,
		MaxAllowed:      elem.MaxAllowed}

	return dbtaskqos
}

// create Guarantees element
func genGuarantees(taskName string, elem cfg.CLASS_QOS_TEMPLATE__ELEM) []structs.SLA_AGREEMENT_DETAILS_GUARANTIES {
	guarantees := make([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, 0)

	for _, v := range elem.Guarantees {
		g := structs.SLA_AGREEMENT_DETAILS_GUARANTIES{}
		g.Name = elem.GuaranteeName // i.e. "kubelet_running_pod_count"
		g.Constraint = v.Constraint // i.e. "kubelet_running_pod_count < 50"

		guarantees = append(guarantees, g)
	}

	return guarantees
}

///////////////////////////////////////////////////////////////////////////////

/*
CreateGuarantees Creates Guarantees for default tasks
*/
func CreateGuarantees(task structs.CLASS_TASK) ([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, *structs.DB_TASK_QOS, error) {
	log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Generating Guarantees from task qos [" + task.Qos.Name + "] ...")

	var guarantees []structs.SLA_AGREEMENT_DETAILS_GUARANTIES

	if task.Qos.Custom != nil && len(task.Qos.Custom) > 0 {
		log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Translating custom task QoS to SLA guarantees ...")
		// TODO support for more than one guarantee
		// for index, _ := range task.Qos.Custom {
		if task.Qos.Custom[0].Guarantees != nil {
			guarantees = make([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, 0)

			guarantee := structs.SLA_AGREEMENT_DETAILS_GUARANTIES{}
			guarantee.Name = task.Qos.Custom[0].Guarantees[0].Metric
			guarantee.Constraint = task.Qos.Custom[0].Guarantees[0].Metric + " " +
				task.Qos.Custom[0].Guarantees[0].Condition + " " + task.Qos.Custom[0].Guarantees[0].Value

			guarantees = append(guarantees, guarantee)

			return guarantees, nil, nil
		}
		//}
	}

	log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] Translating catalog task QoS to SLA guarantees ...")

	qosTemplate, found := GetQoSElem(task.Qos.Name)
	if found {
		guarantees = genGuarantees(task.Name, qosTemplate)
		dbtaskqos := genDbTaskQos(task.Name, qosTemplate)
		return guarantees, dbtaskqos, nil
	}

	log.Println("Rotterdam > CAAS > SLA [CreateGuarantees] ERROR QoS element not found: [" + task.Qos.Name + "]")
	return guarantees, nil, nil
}

// getVarStr
func getVarStr(vVal string, vConfig string) string {
	if len(vVal) == 0 {
		return vConfig
	}
	return vVal
}

// getVarInt
func getVarInt(vVal int, vConfig int) int {
	if vVal == 0 {
		return vConfig
	}
	return vVal
}

// getVarFloat
func getVarFloat(vVal float64, vConfig float64) float64 {
	if vVal == 0.0 {
		return vConfig
	}
	return vVal
}

/*
CreateCOMPSsGuarantees Creates Guarantees for COMPSs based tasks
*/
func CreateCOMPSsGuarantees(task structs.CLASS_TASK) ([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, *structs.DB_TASK_QOS, error) {
	log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] Generating guarantees from task qos ...")

	var guarantees []structs.SLA_AGREEMENT_DETAILS_GUARANTIES

	//
	// [Path A] ... custom guarantees defined in task
	//
	if task.Qos.Custom != nil && len(task.Qos.Custom) > 0 {
		log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] A. Translating custom task QoS to SLA guarantees ...")
		// TODO add support for more than one guarantee
		for index, qosCustomElem := range task.Qos.Custom {
			if task.Qos.Custom[0].Guarantees != nil {
				// Not implemented
				log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] [WARNING] index=" + strconv.Itoa(index) + ", name=" + qosCustomElem.Name)
			}
		}
		log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] [WARNING] Not implemented")

		return nil, nil, nil
	}

	//
	// [Paths B and C]
	//
	agreementID := task.ID                                           // agreement ID
	guarantees = make([]structs.SLA_AGREEMENT_DETAILS_GUARANTIES, 0) // guarantees list
	qosTemplate := cfg.CLASS_QOS_TEMPLATE__ELEM{}                    // qosTemplate
	g := structs.SLA_AGREEMENT_DETAILS_GUARANTIES{}                  // guarantee

	if len(task.QoSCOMPSs) > 0 {
		//
		// [Path B] ... []CLASS_COMPSS_TASK_QOS defined in task
		//
		log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] B. Generating guarantees from []CLASS_COMPSS_TASK_QOS defined in task ...")

		found := false
		if len(task.QoSCOMPSs[0].QoSId) > 0 {
			qosTemplate, found = GetQoSElem(task.QoSCOMPSs[0].QoSId)
		}

		if found {
			//
			// B.1. Generate from QoS template
			//
			log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] B.1. Using existing QoS template ...")
			for _, v := range qosTemplate.Guarantees {
				g.Name = v.Name // i.e. "kubelet_running_pod_count"

				if qosTemplate.Type == "app-compss" {
					g.Constraint = strings.Replace(v.Constraint, "deadlines_missed", "deadlines_missed_"+agreementID, 1)
				} else {
					g.Constraint = v.Constraint // i.e. "kubelet_running_pod_count < 50"
				}

				guarantees = append(guarantees, g)
			}
		} else {
			//
			// B.2. generate default COMPSs guarantees
			//
			log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] B.2. Creating QoS template from JSON values ...")
			// qosTemplate
			qosTemplate.Type = "app-compss"
			if len(task.QoSCOMPSs[0].QoSId) > 0 {
				qosTemplate.GuaranteeName = task.QoSCOMPSs[0].QoSId + agreementID
			} else {
				qosTemplate.GuaranteeName = "COMPSs_" + agreementID
			}
			qosTemplate.MaxAllowed = getVarInt(task.QoSCOMPSs[0].MaxAllowed, cfg.Config.Tasks.MaxAllowed)
			qosTemplate.MinReplicas = getVarInt(task.QoSCOMPSs[0].MinReplicas, cfg.Config.Tasks.MinReplicas)
			qosTemplate.MaxReplicas = getVarInt(task.QoSCOMPSs[0].MaxReplicas, cfg.Config.Tasks.MaxReplicas)
			qosTemplate.ScaleFactor = getVarFloat(task.QoSCOMPSs[0].ScaleFactor, cfg.Config.Tasks.ScaleFactor)
			comparator := getVarStr(task.QoSCOMPSs[0].Comparator, cfg.Config.Tasks.Comparator)
			qosTemplate.Action = getVarStr(task.QoSCOMPSs[0].Action, cfg.Config.Tasks.Action)
			value := getVarInt(task.QoSCOMPSs[0].Value, cfg.Config.Tasks.Value)

			// guarantee
			g.Name = task.QoSCOMPSs[0].Metric + "_" + agreementID
			g.Constraint = task.QoSCOMPSs[0].Metric + "_" + agreementID + " " + comparator + " " + strconv.Itoa(value)
		}
		guarantees = append(guarantees, g)

	} else {
		//
		// [Path C] ... no custom guarantees defined in task
		//
		log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] C. Translating catalog task QoS to SLA guarantees ...")

		qosTemplate, found := GetQoSElem(task.Qos.Name)
		if found {
			//
			// C.1. Generate from QoS template
			//
			log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] C.1. Using existing QoS template ...")

			// g.Name = elem.GuaranteeName // i.e. "kubelet_running_pod_count"
			// g.Constraint = v.Constraint // i.e. "kubelet_running_pod_count < 50"
			for _, v := range qosTemplate.Guarantees {
				g.Name = v.Name // i.e. "kubelet_running_pod_count"

				if qosTemplate.Type == "app-compss" {
					g.Constraint = strings.Replace(v.Constraint, "deadlines_missed", "deadlines_missed_"+agreementID, 1)
				} else {
					g.Constraint = v.Constraint // i.e. "kubelet_running_pod_count < 50"
				}

				guarantees = append(guarantees, g)
			}
		} else {
			//
			// C.2. generate default COMPSs guarantees
			//
			log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] WARNING QoS element not found: [" + task.Qos.Name + "]")
			log.Println("Rotterdam > CAAS > SLA [CreateCOMPSsGuarantees] C.2. Using default QoS template ...")

			// qosTemplate
			qosTemplate.Type = "app-compss"
			qosTemplate.GuaranteeName = "DeadlinesMissed_1"
			qosTemplate.MaxAllowed = 0
			qosTemplate.ScaleFactor = 2.0
			qosTemplate.Action = "scale_out"

			// qosTemplate
			qosTemplate.Type = "app-compss"
			qosTemplate.GuaranteeName = "DeadlinesMissed_1"
			qosTemplate.MaxAllowed = cfg.Config.Tasks.MaxAllowed
			qosTemplate.MinReplicas = cfg.Config.Tasks.MinReplicas
			qosTemplate.MaxReplicas = cfg.Config.Tasks.MaxReplicas
			qosTemplate.ScaleFactor = cfg.Config.Tasks.ScaleFactor
			comparator := cfg.Config.Tasks.Comparator
			qosTemplate.Action = cfg.Config.Tasks.Action
			value := cfg.Config.Tasks.Value

			// guarantee
			g.Name = "deadlines_missed_" + agreementID
			g.Constraint = "deadlines_missed_" + agreementID + " " + comparator + " " + strconv.Itoa(value)

			guarantees = append(guarantees, g)
		}
	}

	dbtaskqos := genDbTaskQos(task.ID, qosTemplate)

	return guarantees, dbtaskqos, nil
}

/*
GetQoSElem gets QoS template from global list
*/
func GetQoSElem(name string) (cfg.CLASS_QOS_TEMPLATE__ELEM, bool) {
	for _, qoselem := range cfg.QosTemplates {
		if qoselem.GuaranteeName == name {
			return qoselem, true
		}
	}

	return cfg.CLASS_QOS_TEMPLATE__ELEM{}, false
}
