;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(ns atos.class.db.db
  (:require [atos.class.logs.logs :as LOG]
            [atos.class.config :as config]
            [konserve.filestore :refer [new-fs-store delete-store]]
            [konserve.memory :refer [new-mem-store]]
            [konserve.core :as k]
            [clojure.core.async :as async :refer [<!! <! go]]))

;; konserve:
;; FROM https://github.com/replikativ/konserve
;;
;; A simple document store protocol defined with core.async semantics to allow Clojuresque collection operations on associative
;; key-value stores, both from Clojure and ClojureScript for different backends. Data is generally serialized with edn semantics
;; or, if supported, as native binary blobs and can be accessed similar to clojure.core functions get-in,assoc-in and update-in.
;; update-in especially allows to run functions atomically and returns old and new value. Each operation is run atomically and
;; must be consistent (in fact ACID), but further consistency is not supported (Riak, CouchDB and many scalable solutions don't
;; have transactions over keys for that reason). This is meant to be a building block for more sophisticated storage solutions
;; (Datomic also builds on kv-stores). An append-log for fast writes is also implemented.
;;
;;
;; NOTE: We use the thread blocking operations <!! here only to synchronize with the REPL. <!! does not compose well with
;; async contexts, so prefer composing your application with go and <! instead.


;; FUNCTION: open-db
;; USAGE: (open-db "/tmp/store");
(defn open-db "Create db in file system"
  [f]
  (LOG/debug "Creating / opening file system database [" f "] ...")
  (<!! (new-fs-store f)))


;; FUNCTION: create-mem-db
;; USAGE: (create-mem-db);
(defn create-mem-db "Create in-memory database"
  []
  (LOG/debug "Creating in-memory database ...")
  (<!! (new-mem-store)))


;; FUNCTION: remove-db
;; USAGE: (remove-db "/tmp/store");
(defn remove-db "Delete db from file system"
  [f]
  (delete-store f))


;; FUNCTION: dbw
;; USAGE: (dbw my-db [:store-cars :car_373_ :contents] {:id 123 :scontent "0100"});
(defn dbw ""
  [db path content]
  (<!! (k/assoc-in db path content)))


;; FUNCTION: dbr
;; USAGE: (dbr my-db [:store-cars :car_373_ :contents]);
(defn dbr ""
  [db path]
  (<!! (k/get-in db path)))


;; FUNCTION: delete-table-db
;; USAGE: (delete-table-db my-db [:store-cars :car_373_ :contents]);
(defn delete-table-db "Deletes 'table' / keyword"
  [db k]
  (<!! (k/dissoc db k)))


;; FUNCTION: dbd
;; USAGE: (dbd my-db [:store-cars :car_373_ :contents]);
(defn dbd "Deletes path"
  [db path]
  (<!! (k/update-in db path (constantly {}))))


;; FUNCTION: dbu
;; USAGE: (dbu my-db [:store-cars :car_373_ :contents]);
(defn dbu "Updates path"
  [db path content]
  (<!! (k/update-in db path (constantly content))))

;; FUNCTION: dba
;; USAGE: (dbu my-db [:store-cars :car_373_ :contents]);
(defn dba "add to path"
  [db path content]
  (<!! (k/assoc-in db path content)))


;; FUNCTION: save-mem-table-to-disk
;; USAGE: (save-mem-table-to-disk db "/tmp/store" :table);
(defn save-mem-table-to-disk ""
  [db f k]
  (let [store   (open-db f)
        data    (dbr db (vector k))]
    (dbw store (vector k) data)))


;; FUNCTION: load-mem-table-from-disk
;; USAGE: (load-mem-table-from-disk db "/tmp/store" :table :table2);
(defn load-mem-table-from-disk ""
  [db f k k2]
  (let [store   (open-db f)
        data    (dbr db (vector k))]
    (dbw db (vector k2) data)))


;; FUNCTION: create-mem-from-files
;; USAGE: (create-mem-from-files "/tmp/store");
(defn create-mem-from-files "Create new in-memory database from files"
  [f]
  (LOG/warning "not implemented"))
