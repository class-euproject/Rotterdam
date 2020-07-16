;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.restapi.handler
  (:use [compojure.core]
        [ring.util.response])
  (:require [compojure.handler :as handler]
            [compojure.route :as route]
            [ring.middleware.json :as middleware]
            [atos.class.logs.logs :as log]
            [atos.class.config :as config]
            [atos.class.restapi.kubernetes :as k8]
            [atos.class.restapi.rules-engine :as rules-engine]
            [atos.class.db.db :as db]
            [cheshire.core :refer [parse-string]]
            [clojure.data.json :as json]
            [clj-http.client :as http-client]
            [ring.middleware.cors :refer [wrap-cors]]))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; APP ROUTES DEFINITION:
(defroutes app-routes
  ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
  ;; WEB - UI:
  ; index page - default
  (GET "/"                  []    (response {:response "ok" :message "Rotterdam-CaaS REST API" :version (str config/get-app-version)
                                             :status "running" :path "/api/v1/"}))
  ; swagger ui
  (GET "/swaggerui/"        []    (resource-response "./swaggerui/index.html" {:root "public"}))
  ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
  ;;  REST-API routes:
  ;;  DEFAULT route: test service
  (GET "/api/"              []    (response {:response "ok" :message "Rotterdam-CaaS REST API" :version (str config/get-app-version)
                                             :status "running" :path "/api/v1/"}))
  (GET "/api/v1/"           []    (response {:response "ok" :message "Rotterdam-CaaS REST API" :version (str config/get-app-version)
                                             :status "running" :path "/api/v1/"}))
  (GET "/api/v1/config"     []    (response (k8/get-api-version)))
  (GET "/api/v1/k8s-proxy"  []    (response {:response "ok" :url (config/get-kubernetes-endpoint)}))
  (GET "/api/v1/version"    []    (response {:response "ok" :version (str config/get-app-version)}))
  (GET "/api/v1/status"     []    (response {:response "ok" :status "running"}))
  ;;  REST-API-URL/test : test services and POST / GET methods
  (context "/test" []
    (defroutes all-users
      ;;
      (GET  "/"           {headers :headers}
                          (do
                            (log/info "> GET /test/")
                            (response {:response "ok" :message "CaaS - CLASS Kubernetes Api running [TESTS]..."})))
      ;;
      (POST "/json"       {body :body}
                          (do
                            (log/info "> POST /test/json [body=" body "]")
                            (response {:response "ok" :content body})))
      ;;
      (POST "/login"      {body :body}
                          (do
                            (log/info "> POST /test/login")
                            (response {:response "ok" :logged true, :token "1"})))
      ;;
      (POST "/deployment" {body :body}
                          (do
                            (log/info "> POST /test/deployment [body=" body "]")
                            (response {:response "ok" :content body})))))
                            ;(response {:response "ok" :content (json/read-str (body "json"))})))))
  ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
  ;; REST-API-URL/api/v1 : authorized services; http://www.class-eu.com/rotterdam/
  (context "/api/v1" []
    (defroutes authorized-users
      ;; GET: Get all tasks from database
      (GET "/docks/tasks"                         {headers :headers}
                                                    (do
                                                      (log/info "> GET /api/v1/db/tasks")
                                                      (response (k8/get-db-tasks))))
      ;; Create a new DEPLOYMENT and SERVICE -manage ports
      (POST "/docks/:dock/tasks"                  {{dock :dock} :params, headers :headers, body :body}
                                                    (do
                                                      ;(log/info "> POST /api/v1/docks/" dock "/tasks [body=" (json/read-str (body "json")) "]")
                                                      (log/info "> POST /api/v1/docks/" dock "/tasks [body=" body "]")
                                                      (response (k8/create-new-task dock body))))
                                                      ;(response (k8/create-new-task dock (json/read-str (body "json"))))))
      ;; GET: Get the description and status of a given task.
      ;; DEPLOYMENT and SERVICE info
      (GET "/docks/:dock/tasks/:name"             {{dock :dock, name :name} :params, headers :headers}
                                                    (do
                                                      (log/info "> GET /api/v1/docks/" dock "/tasks/" name "")
                                                      (response (k8/get-task dock name))))
      ; PUT: Change the definition and re-run a running task.
      (PUT "/docks/:dock/tasks/:name"             {{dock :dock, name :name} :params, headers :headers, body :body}
                                                    (do
                                                      (log/info "> PUT /api/v1/docks/" dock "/tasks/" name " [body=" body "]")
                                                      (response {:response "ok" :message "-not implemented-"})))

      ; DELETE: Stop a running task.
      (DELETE "/docks/:dock/tasks/:name"          {{dock :dock, name :name} :params, headers :headers, body :body}
                                                    (do
                                                      (log/info "> DELETE /api/v1/docks/" dock "/tasks/" name " [body=" body "]")
                                                      (response (k8/delete-task dock name))))

      ; GET: Get the list of containers of a given task.
      (GET "/docks/:dock/tasks/:name/containers"  {{dock :dock, name :name} :params, headers :headers}
                                                    (do
                                                      (log/info "> GET /api/v1/docks/" dock "/tasks/" name "/containers")
                                                      (response (k8/get-task-pods dock name))))

      ; POST: Receives SLA violation
      (POST "/rules-engine/docks/:dock/tasks/:name/sla"  {{dock :dock, name :name} :params, headers :headers, body :body}
                                                          (do
                                                            (log/info "> POST /api/v1/rules-engine/docks/" dock "/tasks/" name "/sla [body=" body "]")
                                                            ;(response {:response "ok" :message "-not implemented-" :message2 "RULES ENGINE INVOKED!"})))))
                                                            (response (rules-engine/process-event "class" "my-nginx" body))))))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
  ; OTHER routes:
  (route/resources "/")
  (route/not-found (response {:message "-Not Found-"})))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; APP - START / RUN SERVER
(def app
  (->
    (handler/api (wrap-cors app-routes
                   :access-control-allow-origin [#".*"] ;[#"localhost"] ;[#".*"]
                   :access-control-allow-methods [:get :put :post :delete]))
    (middleware/wrap-json-response)
    (middleware/wrap-json-body)))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; EXECUTE INIT:
(try
  (do
    (db/open-db config/get-db-path)

    (log/info "###################################################################")
    (log/info "CLASS - Rotterdam-CaaS [version " config/get-app-version "]")
    (log/info "    Version:            " config/get-app-version)
    (log/info "    Configuration file: "  config/get-conf-file)
    (log/info "    Connected to Openshift/Kubernetes")
    (log/info "       + REST API endpoint:            " (config/get-kubernetes-endpoint))
    (log/info "       + Openshift REST API endpoint:  " (config/get-openshift-endpoint))
    (log/info "       + Openshift token:              '..." (subs (config/get-openshift-oauth-token) 20) "'")
    (log/info "       + Serverl IP:                   " (config/get-server-ip))
    (log/info "    Database path: " config/get-db-path)
    (log/info "    Listening in port 18083 ...")
    (log/info "###################################################################"))
  (catch Exception e
    (do
      (log/error "Error initializing Rotterdam-CaaS application: caught exception: " (.getMessage e)
                  "\n       stackTrace: " (clojure.string/join "\n " (map str (.getStackTrace e))))
      (log/info "###################################################################"))))
