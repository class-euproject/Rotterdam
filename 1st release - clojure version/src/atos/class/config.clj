;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.config)

;; default configuration
(def default-config
  {
    :app {
      ;; REST-API properties
      :name "class-k8s-rest-api"
      :version "0.0.9"
      :conf-file "default"
      ;; USED in logs
      :log-message "Rotterdam-CaaS REST API [0.0.9] "
      ;; database path
      :db-path "/tmp/store"
    }
    :apis {
      ;; server
      :server-ip "192.168.7.28"
      ;; oauth-token
      :openshift-oauth-token "9LQDZOcE1MuYtoMoLshRCco61kJMhwiXI8rWhc7z06w"
      ;; KUBERNETES connection data
      :kubernetes {
        :url "http://192.168.7.28:8001"
      }
      :openshift {
        :url "https://192.168.7.28:8443"
    }}
  })

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; FUNCTION: get-resource
(defn- get-resource "get file or nil if error"
  [path]
  (when path (-> (Thread/currentThread) .getContextClassLoader (.getResource path))))

;; read configuration values into map
(defn- read-configuration [path]
    (let [fpath (get-resource path)]
        (if (nil? fpath)
            default-config
            (read-string (slurp fpath)))))

;; properties (as map)
(def props (read-configuration "app.properties.clj"))

;; FUNCTION: get properties vaules, e.g. (conf :db :uri)
(defn- conf [& path]
  (get-in props (vec path)))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; CONFIGURATION VALUES
(def get-app-name (conf :app :name))
(def get-app-version (conf :app :version))
(def get-log-message (conf :app :log-message))
(def get-db-path (conf :app :db-path))
(def get-conf-file (conf :app :conf-file))

; KUBERNETES, OPENSHIFT, APIS
(def ^:private openshift-oauth-token (atom (or (System/getenv "OAUTH_TOKEN") (conf :apis :openshift-oauth-token))))
(def ^:private kubernetes-endpoint (atom (or (System/getenv "K8s_API_URL") (conf :apis :kubernetes :url))))
(def ^:private server-ip (atom (or (System/getenv "SERVER_IP") (conf :apis :server-ip))))
(def ^:private openshift-endpoint (atom (or (System/getenv "OPENSHIFT_API_URL") (conf :apis :openshift :url))))

(defn get-openshift-oauth-token [] (deref openshift-oauth-token))
(defn set-openshift-oauth-token [t] (reset! openshift-oauth-token t))

(defn get-kubernetes-endpoint [] (deref kubernetes-endpoint))
(defn set-kubernetes-endpoint [url] (reset! kubernetes-endpoint url))

(defn get-server-ip [] (deref server-ip))
(defn set-server-ip [ip] (reset! server-ip ip))

(defn get-openshift-endpoint [] (deref openshift-endpoint))
(defn set-openshift-endpoint [url] (reset! openshift-endpoint url))
