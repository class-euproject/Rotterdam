;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.restapi.kubernetes
  (:require [atos.class.logs.logs :as log]
            [atos.class.config :as config]
            [atos.class.db.dbk :as dbk]
            [clojure.data.json :as json]
            [atos.class.restapi.http :as http]
            [atos.class.restapi.schemas :as schemas]))


;; FUNCTION: handle-exception
(defn- handle-exception ""
  [e]
  (do (log/log-exception e) {:response "err" :error "Exception" :code http/CODE_ERROR :message (.getMessage e)}))


;; TODO fix problem with json
;; FUNCTION: get-status
(defn- get-status "Get status of task"
  [task query]
  (try
    (if (= query "deployment")
      (let [status (get-in task [:depl-content "status"])]
        (if (empty? status)
          {:desc "deploying" :replicas 0 :readyReplicas 0}
          (if (and (= (status "replicas") (status "readyReplicas")) (> (status "replicas") 0))
            {:desc "ready" :replicas (status "replicas") :readyReplicas (status "readyReplicas")}
            {:desc "not-ready" :replicas (status "replicas") :readyReplicas (status "readyReplicas")})))
      (let [status (get-in task [:depl-content :status])]
        (if (empty? status)
          {:desc "not-ready" :replicas 0 :readyReplicas 0}
          (if (and (= (status :replicas) (status :readyReplicas)) (> (status :replicas) 0))
            {:desc "ready" :replicas (status :replicas) :readyReplicas (status :readyReplicas)}
            {:desc "not-ready" :replicas (status :replicas) :readyReplicas (status :readyReplicas)}))))
    (catch Exception e (do (log/log-exception e) {:desc "error"}))))


;; TODO fix problem with json
;; FUNCTION: resp-get-task
;;    task ... (dbk/read-task task-name);
;; EXAMPLE: (resp-get-task (dbk/read-task task-name) "deployment");
(defn- resp-get-task "Generates response content"
  [task query]
  (log/info "task: " task)
  (if (= query "deployment")
    (let [status  (get-status task query)]
      {:name      (get-in task [:depl-content "metadata" "name"])
       :dock      (get-in task [:depl-content "metadata" "namespace"])
       :status    status
       :urls      (task :urls)
       :created   (get-in task [:depl-content "metadata" "creationTimestamp"])})
    (let [status  (get-status task query)]
      {:name      (get-in task [:depl-content :metadata :name])
       :dock      (get-in task [:depl-content :metadata :namespace])
       :status    status
       :urls      (task :urls)
       :created   (get-in task [:depl-content :metadata :creationTimestamp])})))


;; FUNCTION: get-api-version
(defn get-api-version "Get Kubernetes API version" [] (http/GET (str (config/get-kubernetes-endpoint) "/api")))


;; FUNCTION: create-new-task
;;  1. POST /apis/apps/v1/namespaces/{namespace}/deployments
;;  2. POST /api/v1/namespaces/{namespace}/services
;;  View `input.txt` for more information
(defn create-new-task "Creates a New DEPLOYMENT and a SERVICE to publish the ports"
  [dock json-body]
  (try
    (let [namespace       dock
          task-name       (json-body "name")
          json-k8-deploy  (schemas/gen-json-k8s-deploy task-name json-body) ; create deplyment json -KUBERNETES-
          ; deploy task: POST /apis/apps/v1/namespaces/{namespace}/deployments
          res-deployment (http/POST (str (config/get-kubernetes-endpoint) "/apis/apps/v1/namespaces/" namespace "/deployments") json-k8-deploy)]
      (if (= (res-deployment :response) "ok")
        (let [json-k8-serv  (schemas/gen-json-k8s-serv task-name json-body) ; create service json -KUBERNETES-
              ; deploy service: POST /api/v1/namespaces/{namespace}/services
              res-service     (http/POST (str (config/get-kubernetes-endpoint) "/api/v1/namespaces/" namespace "/services") json-k8-serv)]
          (if (= (res-service :response) "ok")
            (let [json-k8-serv-route  (schemas/gen-json-k8s-route task-name json-body) ; create route json -OPENSHIFT-
                  ; deploy service: POST /api/v1/namespaces/{namespace}/services
                  res-route     (http/POST (str (config/get-openshift-endpoint) "/apis/route.openshift.io/v1/namespaces/" namespace "/routes") json-k8-serv-route)]
              (if (= (res-route :response) "ok")
                (do
                  (dbk/save-task  task-name
                                  [(str "serv-" task-name)]
                                  ["volumes"]                           ; volumes TODO!
                                  (get-in json-k8-serv [:spec :ports])  ; ports
                                  (json/read-str (res-deployment :content))
                                  (json/read-str (res-service :content))
                                  (str task-name "." (config/get-server-ip) ".xip.io"))
                                  ;(config/get-kubernetes-ext-ip))
                  {:response "ok" :message (str "task '" task-name "' created ['" dock "']") :content (resp-get-task (dbk/read-task task-name) "deployment")})
                {:response "err"  :code (res-route :code)  :message (str "Error creating route for task '" task-name "' (3)")}))
            {:response "err"  :code (res-service :code)  :message (str "Error creating service for task '" task-name "' (2)")}))
        {:response "err"  :code (res-deployment :code)  :message (str "Error creating deployment for task '" task-name "' (1)")}))
    (catch Exception e (handle-exception e))))


;; FUNCTION: get-task
;;  1. GET /apis/apps/v1/namespaces/{namespace}/deployments/{name}
(defn get-task "Returns task"
  [dock name]
  (try
    (let [res-get (http/GET (str (config/get-kubernetes-endpoint) "/apis/apps/v1/namespaces/" dock "/deployments/" name))]
      (if (= (res-get :response) "ok")
        (do
          (dbk/update-task name :depl-content (res-get :content))
          {:response "ok" :content (resp-get-task (dbk/read-task name) "query")})
        res-get))
    (catch Exception e (handle-exception e))))


;; FUNCTION: delete-task
;;  1. DELETE /apis/apps/v1/namespaces/{namespace}/deployments/{name}
;;  2. DELETE /apis/route.openshift.io/v1/namespaces/{namespace}/routes/{name}
;;  3. DELETE /api/v1/namespaces/{namespace}/services/{name}
(defn delete-task "Deletes a task"
  [dock name]
  (try
    (let [res (http/DELETE (str (config/get-kubernetes-endpoint) "/apis/apps/v1/namespaces/" dock "/deployments/" name) {})]
      (if (= (res :response) "ok")
        (let [res2 (http/DELETE (str (config/get-openshift-endpoint) "/apis/route.openshift.io/v1/namespaces/" dock "/routes/route-" name) {})]
          (if (= (res2 :response) "ok")
            (let [res3 (http/DELETE (str (config/get-kubernetes-endpoint) "/api/v1/namespaces/" dock "/services/serv-" name) {})]
              (if (= (res3 :response) "ok")
                {:response "ok" :message (str "task " name " deleted")
                  :content {:dock dock :name name :res "deleted" :res-db (str (dbk/delete-task name))}}
                {:response "err"  :code (res3 :code)  :message (str "Error deleting service from task '" name "' (3)")
                  :content {:dock dock :name name :res "not deleted"}}))
            {:response "err"  :code (res2 :code)  :message (str "Error deleting route from task '" name "' (2)")
              :content {:dock dock :name name :res "not deleted"}}))
        {:response "err"  :code (res :code)  :message (str "Error deleting task '" name "' (1)")
          :content {:dock dock :name name :res "not deleted"}}))
    (catch Exception e (handle-exception e))))


;; FUNCTION: get-task-pods
;; 1. GET /api/v1/namespaces/{namespace}/pods?labelSelector=app%3Dmysql
(defn get-task-pods "Get pods -containers- from task"
  [dock name]
  (try
    (let [res-get (http/GET (str (config/get-kubernetes-endpoint) "/api/v1/namespaces/" dock "/pods?labelSelector=app%3D" name))]
      (if (= (res-get :response) "ok")
        (let [items (get-in res-get [:content :items])]
          {:response "ok" :content {:items (count items)} :status (for [i items] (i :status))})
        res-get))
    (catch Exception e (handle-exception e))))


;; FUNCTION: get-db-tasks
;; 1. GET /api/v1/db/tasks
(defn get-db-tasks "Get database content"
  []
  (try
    {:response "ok" :content (dbk/get-db-tasks)}
    (catch Exception e (handle-exception e))))
