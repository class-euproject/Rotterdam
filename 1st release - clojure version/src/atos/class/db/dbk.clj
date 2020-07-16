;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.db.dbk
  (:require [atos.class.logs.logs :as LOG]
            [atos.class.config :as config]
            [atos.class.db.db :as db]))


(def TASK-DB-SCHEMA {
  ; string
  :name         "dname"
  ; string
  :deployment   "dname"
  ; map
  :depl-content {}  ; k8s deployment content
  ; map
  :serv-content {}  ; k8s service content
  ; vector
  :service      ["sname1", "sname2"]
  ; vector
  :volumenClaim ["vname1", "vname2"]
  ; vector
  :ports        [1, 2, 3]
  ; string
  :status       "deleted"
  ; string
  :url          ""
})


;; FUNCTION: create-urls
(defn- create-urls ""
  [url ports tname]
  (str tname "." (config/get-server-ip) ".xip.io"))
  ;(for [x ports]
  ;  (str url ":" (x :port))))


;; FUNCTION: save-task
(defn save-task ""
  [tname snames vnames ports dcontent scontent url]
  (LOG/debug "DB>> save-task>> params>> " tname ", " snames ", " vnames ", " ports ", " dcontent ", " scontent)
  (let [dbstore (db/open-db config/get-db-path)
        res (db/dbw dbstore [:tasks (keyword tname)]
              { :name         tname
                :deployment   tname
                :depl-content dcontent  ; map
                :serv-content scontent  ; map
                :service      snames    ; vector
                :volumenClaim vnames    ; vector
                :ports        ports     ; vector
                :status       "created"
                :urls         (str "http://" tname "." (config/get-server-ip) ".xip.io")})] ; (create-urls url ports tname)})]
    (LOG/info "DB>> save-task>> result>> " res)
    res))


;; FUNCTION: delete-task
(defn delete-task ""
  [tname]
  (LOG/debug "DB>> delete-task>> params>> " tname)
  (let [dbstore (db/open-db config/get-db-path)
        res (db/dbd dbstore [:tasks (keyword tname)])]
    (LOG/debug "DB>> delete-task>> result>> " res)
    res))


;; FUNCTION: read-task
(defn read-task ""
  [tname]
  (LOG/debug "DB>> read-task>> params>> " tname)
  (let [dbstore (db/open-db config/get-db-path)
        res (db/dbr dbstore [:tasks (keyword tname)])]
    (LOG/info "DB>> read-task>> result>> " res)
    res))


;; FUNCTION: update-task
(defn update-task ""
  [tname field value]
  (LOG/debug "DB>> update-task>> params>> " tname ", " field ", " value)
  (let [dbstore (db/open-db config/get-db-path)]
    (db/dbu dbstore [:tasks (keyword tname) field] value)
    (let [res (db/dbr dbstore [:tasks (keyword tname)])]
      (LOG/debug "DB>> update-task>> result>> " res)
      res)))


;; FUNCTION: get-db-tasks
(defn get-db-tasks ""
  []
  (let [dbstore (db/open-db config/get-db-path)]
    (db/dbr dbstore [:tasks])))


;; FUNCTION: get-available-port
(defn get-available-port ""
  []
  {})
