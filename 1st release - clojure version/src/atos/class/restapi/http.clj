;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.restapi.http
  (:require [atos.class.logs.logs :as log]
            [atos.class.config :as config]
            [cheshire.core :refer :all]
            [clojure.data.json :as json]
            [clj-http.client :as http-client]))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(def CODE_ERROR "-1")
;; common HTTP status codes:
;; OK : The request was fulfilled.
(def HTRESP_OK 200)
;; CREATED : Following a POST command, this indicates success, but the textual part of the response line indicates
;;           the URI by which the newly created document should be known.
(def HTRESP_CREATED 201)
;; ACCEPTED : The request has been accepted for processing, but the processing has not been completed.
(def HTRESP_ACCEPTED 202)
;; Errors...
(def HTRESP_UNAUTHORIZED 401)

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

;; FUNCTION: POST
(defn POST "generic POST operation"
  [endpoint body]
  (try
    (log/info "POST call to: " endpoint)
    (log/info "POST's body: " (json/write-str body))
    (let [resp      (http-client/post endpoint {;:basic-auth ["rsucasas" "redhat"]
                                                :oauth-token (config/get-openshift-oauth-token) :insecure? true
                                                :content-type :json
                                                :accept :json
                                                :body (json/write-str body)})
          resp-log  (log/info "POST response: " (resp :body))]
      (if (or (= HTRESP_ACCEPTED (resp :status)) (= HTRESP_CREATED (resp :status)) (= HTRESP_OK (resp :status)))
        {:response "ok"  :code (resp :status) :content (resp :body)}
        {:response "err" :code (resp :status) :internal-code CODE_ERROR}))
    (catch Exception e (do (log/log-exception e) {:response "err" :error "Exception" :code CODE_ERROR :message (.getMessage e)}))))


;; FUNCTION: PUT
(defn PUT "generic PUT operation"
  [endpoint body]
  (try
    (log/info "PUT call to: " endpoint)
    (log/info "POST's body: " (json/write-str body))
    (let [resp      (http-client/put endpoint {;:basic-auth ["rsucasas" "redhat"]
                                               :oauth-token (config/get-openshift-oauth-token) :insecure? true
                                               :content-type :json
                                               :accept :json
                                               :body (json/write-str body)})
          resp-log  (log/info "PUT response: " (resp :body))]
      (if (or (= HTRESP_ACCEPTED (resp :status)) (= HTRESP_CREATED (resp :status)) (= HTRESP_OK (resp :status)))
        {:response "ok"  :code (resp :status) :content (resp :body)}
        {:response "err" :code (resp :status) :internal-code CODE_ERROR}))
    (catch Exception e (do (log/log-exception e) {:response "err" :error "Exception" :code CODE_ERROR :message (.getMessage e)}))))


;; FUNCTION: DELETE
(defn DELETE "generic DELETE operation"
  [endpoint body]
  (try
    (log/info "DELETE call to: " endpoint)
    (log/info "POST's body: " (json/write-str body))
    (let [resp      (http-client/delete endpoint {;:basic-auth ["rsucasas" "redhat"]
                                                  :oauth-token (config/get-openshift-oauth-token) :insecure? true
                                                  :content-type :json
                                                  :accept :json
                                                  :body (json/write-str body)})
          resp-log  (log/info "DELETE response: " (resp :body))]
      (if (= HTRESP_OK (resp :status))
        {:response "ok"  :code (resp :status) :content (resp :body)}
        {:response "err" :code (resp :status) :internal-code CODE_ERROR}))
    (catch Exception e (do (log/log-exception e) {:response "err" :error "Exception" :code CODE_ERROR :message (.getMessage e)}))))


;; FUNCTION: GET
(defn GET "generic GET operation"
  [endpoint]
  (try
    (log/info "GET call to: " endpoint)
    (let [resp      (http-client/get endpoint {;:basic-auth ["rsucasas" "redhat"]
                                               :oauth-token (config/get-openshift-oauth-token) :insecure? true
                                               :as :json})
          resp-log  (log/info "GET response: " (resp :body))]
      (if (= HTRESP_OK (resp :status))
        {:response "ok"  :code HTRESP_OK :content (resp :body)}
        {:response "err" :code (resp :status) :internal-code CODE_ERROR}))
    (catch Exception e (do (log/log-exception e) {:response "err" :error "Exception" :code CODE_ERROR :message (.getMessage e)}))))
