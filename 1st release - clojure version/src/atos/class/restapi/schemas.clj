;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.restapi.schemas
  (:require [atos.class.config :as config]
            [clojure.data.json :as json]
            [atos.class.logs.logs :as log]
            [atos.class.restapi.http :as http]))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

;; FUNCTION: parse-port-int
(defn- parse-port-int ""
  [port]
  (try
    (if (int? port)
      port
      (Integer/parseInt port))
    (catch Exception e
      (do (log/log-exception e) 666))))


;; FUNCTION: gen-ports-container-data
(defn- gen-ports-container-data ""
  [json-container-ports]
  (for [x json-container-ports]
    {:containerPort (parse-port-int (x "containerPort"))}))
    ;(if (int? (x "containerPort"))
    ;  {:containerPort (x "containerPort")}
    ;  {:containerPort (Integer/parseInt (x "containerPort"))})))


;; FUNCTION: gen-volume-mounts-data
(defn- gen-volume-mounts-data ""
  [json-volume-mounts]
  (for [x json-volume-mounts]
    {:name      (x "name")
     :mountPath (x "mounthPath")}))


;; FUNCTION: gen-env-data
(defn- gen-env-data ""
  [json-envs]
  (for [x json-envs]
    {:name    (x "name")
     :value   (x "value")}))


;; FUNCTION: gen-containers-data
(defn- gen-containers-data ""
  [json-body]
  (let [containers (json-body "containers")]
    (for [x containers]
      {:image           (x "image")
       :name            (x "name")
       :imagePullPolicy "Always"
       :ports           (gen-ports-container-data (x "ports"))
       ;:volumeMounts    (gen-volume-mounts-data (x "volumes"))
       :env             (gen-env-data (x "environment"))})))


;; FUNCTION: gen-ports-service
;; {:name "http" :port 8080 :protocol "TCP" :targetPort 80}
(defn- gen-ports-service ""
  [json-container-ports]
  (for [x json-container-ports]
    {:name "http" :port (parse-port-int (x "hostPort")) :protocol "TCP" :targetPort (parse-port-int (x "containerPort"))}))


;; FUNCTION: gen-service-ports-data
(defn- gen-service-ports-data ""
  [json-body]
  (into []
    (flatten
      (let [containers (json-body "containers")]
        (for [x containers]
          (gen-ports-service (x "ports")))))))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

;; FUNCTION: gen-json-k8s-deploy
(defn gen-json-k8s-deploy ""
  [task-name json-body]
  { :apiVersion "apps/v1"
    :kind "Deployment"
    :metadata {:name task-name}
    :spec {
      :replicas 1
      :revisionHistoryLimit 10
      :selector {:matchLabels {:app task-name}}
      :template {
        :metadata {:labels {:app task-name}}
        :spec {:containers (gen-containers-data json-body)}}}})


;; FUNCTION: gen-json-k8s-serv
(defn gen-json-k8s-serv ""
  [task-name json-body]
  { :kind "Service"
    :apiVersion "v1"
    :metadata {:name (str "serv-" task-name)}
    :spec {
      ; port .......... accessible port in defined external IP
      ; targetPort .... application port
      :ports (gen-service-ports-data json-body)         ;; => [{:name "http" :port 8080 :protocol "TCP" :targetPort 80}]
      :selector {:app task-name}}}
      ;:externalIPs [(config/get-kubernetes-ext-ip)]}}   ;; => "192.168.7.24"
)


;; FUNCTION: gen-json-k8s-route
(defn gen-json-k8s-route ""
  [task-name json-body]
  {
    :apiVersion "route.openshift.io/v1"
    :kind "Route"
    :metadata {
      :name (str "route-" task-name)
      :namespace "class"
    }
    :spec {
      :host (str task-name "." (config/get-server-ip) ".xip.io") ;(str task-name ".192.168.7.28.xip.io")
      :port {
        :targetPort "http"
      }
      :to {
        :kind "Service"
        :name (str "serv-" task-name)
      }
    }
  }
)
