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

package structs

///////////////////////////////////////////////////////////////////////////////
// K8S_DEPLOYMENT

type K8S_DEPLOYMENT_CONTAINER_PORTS struct {
	ContainerPort int `json:"containerPort,omitempty"`
}

type K8S_DEPLOYMENT_CONTAINER_ENV struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type K8S_DEPLOYMENT_CONTAINER struct {
	Image           string                           `json:"image,omitempty"`
	Name            string                           `json:"name,omitempty"`
	ImagePullPolicy string                           `json:"imagePullPolicy,omitempty"`
	Ports           []K8S_DEPLOYMENT_CONTAINER_PORTS `json:"ports,omitempty"`
	Env             []K8S_DEPLOYMENT_CONTAINER_ENV   `json:"env,omitempty"`
	Command         []string                         `json:"command,omitempty"`
	Args            []string                         `json:"args,omitempty"`
}

type K8S_DEPLOYMENT struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name string `json:"name,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Replicas             int `json:"replicas,omitempty"`
		RevisionHistoryLimit int `json:"revisionHistoryLimit,omitempty"`
		Selector             struct {
			MatchLabels struct {
				App string `json:"app,omitempty"`
			} `json:"matchLabels,omitempty"`
		} `json:"selector,omitempty"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `json:"app,omitempty"`
				} `json:"labels,omitempty"`
			} `json:"metadata,omitempty"`
			Spec struct {
				Containers []K8S_DEPLOYMENT_CONTAINER `json:"containers,omitempty"`
			} `json:"spec,omitempty"`
		} `json:"template,omitempty"`
	} `json:"spec,omitempty"`
}

type K8S_JOB struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name string `json:"name,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Replicas             int `json:"replicas,omitempty"`
		RevisionHistoryLimit int `json:"revisionHistoryLimit,omitempty"`
		Selector             struct {
			MatchLabels struct {
				App string `json:"app,omitempty"`
			} `json:"matchLabels,omitempty"`
		} `json:"selector,omitempty"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `json:"app,omitempty"`
				} `json:"labels,omitempty"`
			} `json:"metadata,omitempty"`
			Spec struct {
				Containers []K8S_DEPLOYMENT_CONTAINER `json:"containers,omitempty"`
			} `json:"spec,omitempty"`
		} `json:"template,omitempty"`
	} `json:"spec,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// K8S_SERVICE

type K8S_SERVICE_PORT struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	TargetPort int    `json:"targetPort,omitempty"`
}

type K8S_SERVICE struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name   string `json:"name,omitempty"`
		Labels struct {
			App string `json:"app,omitempty"`
		} `json:"labels,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Selector struct {
			App     string `json:"app,omitempty"`
			PodName string `json:"pod-name,omitempty"`
		} `json:"selector,omitempty"`
		ExternalIPs []string           `json:"externalIPs,omitempty"`
		Ports       []K8S_SERVICE_PORT `json:"ports,omitempty"`
	} `json:"spec,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// K8S_ROUTE

type K8S_ROUTE struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Host string `json:"host,omitempty"`
		Port struct {
			TargetPort string `json:"targetPort,omitempty"`
		} `json:"port,omitempty"`
		To struct {
			Kind string `json:"kind,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"to,omitempty"`
	} `json:"spec,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// K8S_SCALE

type K8S_SCALE struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name              string `json:"name,omitempty"`
		Namespace         string `json:"namespace,omitempty"`
		SelfLink          string `json:"selfLink,omitempty"`
		Uid               string `json:"uid,omitempty"`
		ResourceVersion   string `json:"resourceVersion,omitempty"`
		CreationTimestamp string `json:"creationTimestamp,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Replicas int `json:"replicas,omitempty"`
	} `json:"spec,omitempty"`
	Status struct {
		Replicas int    `json:"replicas,omitempty"`
		Selector string `json:"selector,omitempty"`
	} `json:"status,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// POD_PATCH

type K8S_POD_PATCH_LINE struct {
	Op    string `json:"op,omitempty"`
	Path  string `json:"path,omitempty"`
	Value string `json:"value,omitempty"`
}

///////////////////////////////////////////////////////////////////////////////
// FUNCTIONS - SERVERLESS - KUBELESS

/*
KubelessFunctionStruct ...
	JSON EXAMPLE:
		{
			"apiVersion": "kubeless.io/v1beta1",
			"kind": "Function",
			"metadata": {
				"name": "get-python",
				"namespace": "default",
				"label": {
					"created-by": "kubeless",
					"function": "get-python"
				}
			},
			"spec": {
				"runtime": "python2.7",
				"timeout": "180",
				"handler": "helloget.foo",
				"deps": "",
				"checksum": "sha256:d251999dcbfdeccec385606fd0aec385b214cfc74ede8b6c9e47af71728f6e9a",
				"function-content-type": "text",
				"function": "def foo(event, context):\n    return \"hello world\"\n"
			}
		}
*/
type KubelessFunctionStruct struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name      string `json:"name,omitempty"`
		NameSpace string `json:"namespace,omitempty"`
		Label     struct {
			CreatedBy string `json:"created-by,omitempty"`
			Function  string `json:"function,omitempty"`
		} `json:"label,omitempty"`
	} `json:"metadata,omitempty"`
	Spec struct {
		Runtime             string `json:"runtime,omitempty"`
		Timeout             string `json:"timeout,omitempty"`
		Handler             string `json:"handler,omitempty"`
		Deps                string `json:"deps,omitempty"`
		Checksum            string `json:"checksum,omitempty"`
		FunctionContentType string `json:"function-content-type,omitempty"`
		Function            string `json:"function,omitempty"`
	} `json:"spec,omitempty"`
}
