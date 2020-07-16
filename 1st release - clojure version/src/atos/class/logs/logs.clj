;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.logs.logs
  (:require [clojure.tools.logging :as logging]
            [atos.class.config :as config]))

;; DEFINITION: use to check not nil
(def ^:private not-nil? (complement nil?))

;; FUNCTION: print-log
(defn- print-log "prints log content on screen console"
  [debug-level txt]
  {:pre [(not-nil? debug-level) (not-nil? txt)]}
  (cond
      (= debug-level "DEBUG")   (logging/debug txt)
      (= debug-level "INFO")    (logging/info txt)
      (= debug-level "ERROR")   (logging/error txt)
      (= debug-level "WARNING") (logging/warn txt)
      :else                     (logging/trace txt)))

;; FUNCTION: pr-log
(defn- pr-log
  [l-type & txt]
  (print-log l-type (apply str config/get-log-message txt)))

;; FUNCTION: debug
(defn debug [& txt] {:pre [(not-nil? txt)]}
  (apply pr-log "DEBUG" txt))

;; FUNCTION: info
(defn info [& txt] {:pre [(not-nil? txt)]}
  (apply pr-log "INFO" txt))

;; FUNCTION: error
(defn error [& txt] {:pre [(not-nil? txt)]}
  (apply pr-log "ERROR" txt))

;; FUNCTION: warning
(defn warning [& txt] {:pre [(not-nil? txt)]}
  (apply pr-log "WARNING" txt))

;; FUNCTION: trace
(defn trace [& txt] {:pre [(not-nil? txt)]}
  (apply pr-log "TRACE" txt))

;; FUNCTION: create error map
(defn get-error-stacktrace [e]
  (error "> Caught exception: [" (.getMessage e) "], stackTrace: \n    " (clojure.string/join "\n    " (map str (.getStackTrace e)))))

;; FUNCTION: create-map-error
(defn create-map-error "creates a map with Exception info"
  [e]
  {:code        "ERROR"
   :message     (str "caught exception: " (.getMessage e))
   :stacktrace  (str "StackTrace: " (clojure.string/join "\n " (map str (.getStackTrace e))))
   :response    false})

;; FUNCTION:
(defn log-exception "creates a map with Exception info"
  [e]
  (error "ERROR: caught exception: " (.getMessage e)
              "\n       stackTrace: " (clojure.string/join "\n " (map str (.getStackTrace e)))))
