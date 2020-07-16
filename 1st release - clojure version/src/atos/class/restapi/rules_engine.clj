;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.restapi.rules-engine
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


(def TMP_COUNTER (atom 0))


;; FUNCTION: process-event
;; 1. curl -X GET -H "Authorization: Bearer ..." http://localhost:8001/apis/apps/v1/namespaces/class/deployments/my-nginx/scale
;; 2. curl -X PUT -d@scale.json -H 'Content-Type: application/json' -H "Authorization: Bearer ..." http://localhost:8001/apis/apps/v1/namespaces/class/deployments/my-nginx/scale
(defn process-event "Creates a New DEPLOYMENT and a SERVICE to publish the ports"
  [dock task-name json-body]
  (swap! TMP_COUNTER inc)
  (if (= @TMP_COUNTER 2)
    (do
      (log/info "rules-engine > Processing EVENT [dock: " dock "] [task: " task-name "] [content: " json-body "]...")
      (try
        ; 1. curl -X GET -H "Authorization: Bearer ..." http://localhost:8001/apis/apps/v1/namespaces/class/deployments/my-nginx/scale
        (let [res-scale-info  (http/GET (str (config/get-kubernetes-endpoint) "/apis/apps/v1/namespaces/" dock "/deployments/" task-name "/scale"))
              log             (log/info "rules-engine > res-scale-info > " res-scale-info)]
          (if (= (res-scale-info :response) "ok")
            ; {
            ;   "kind": "Scale",
            ;   "apiVersion": "autoscaling/v1",
            ;   "metadata": {
            ;     "name": "my-nginx",
            ;     "namespace": "class",
            ;     "selfLink": "/apis/apps/v1/namespaces/class/deployments/my-nginx/scale",
            ;     "uid": "f5bb01e5-fecb-11e8-8947-005056986059",
            ;     "resourceVersion": "4925456",
            ;     "creationTimestamp": "2018-12-13T11:40:58Z"
            ;   },
            ;   "spec": {
            ;     "replicas": 3
            ;   },
            ;   "status": {
            ;     "replicas": 1,
            ;     "selector": "app=my-nginx"
            ;   }
            ; }
            (let [scale-info              (res-scale-info :content)
                  log                     (log/info "rules-engine > scale-info > " scale-info)
                  updated-scale-info      (update-in scale-info ["spec" "replicas"] (constantly 3))
                  log                     (log/info "rules-engine > updated-scale-info > " updated-scale-info)
                  ; 2. curl -X PUT -d@scale.json -H 'Content-Type: application/json' -H "Authorization: Bearer ..." http://localhost:8001/apis/apps/v1/namespaces/class/deployments/my-nginx/scale
                  res-updated-scale-info  (http/PUT (str (config/get-kubernetes-endpoint) "/apis/apps/v1/namespaces/" dock "/deployments/" task-name "/scale") updated-scale-info)]
              (if (= (res-updated-scale-info :response) "ok")
                {:response "ok" :message (str "task '" task-name "' scaled to replicas: [3]") :content (resp-get-task (dbk/read-task task-name) "deployment")}
                {:response "err"  :code (res-updated-scale-info :code)  :message (str "Error scaling task '" task-name "' (2)")}))
            {:response "err"  :code (res-scale-info :code)  :message (str "Error scaling task '" task-name "' (1)")}))
        (catch Exception e (handle-exception e))))
    (log/info "rules-engine > TMP_COUNTER= " @TMP_COUNTER " ==> EVENT [dock: " dock "] [task: " task-name "] [content: " json-body "] not processed...")))
