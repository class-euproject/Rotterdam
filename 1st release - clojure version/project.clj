;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; This code is being developed for the CLASS Project
;;
;; Copyright: Roi Sucasas Font, Atos Research and Innovation, 2018.
;;
;; This code is licensed under an Apache 2.0 license. Please, refer to the
;; LICENSE.TXT file for more information
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; KUBERNETES - REST API
;;     - v0.0.1 ....... 30.08.2018 CLASS project - kubernetes
;;     - v0.0.2 ....... 18.09.2018 basic operations added to api
;;     - v0.0.3 ....... 20.09.2018 basic operations working in K8s
;;     - v0.0.4 ....... 01.10.2018 gui updated; code reworked
;;     - v0.0.5 ....... 09.10.2018 responses updated
;;     - v0.0.6 ....... 19.11.2018 openshift implementation - ini
;;     - v0.0.8 ....... 27.11.2018 openshift implementation - errors fixed
;;     - v0.0.9 ....... 13.12.2018 integration with SLALite; basic rules engine
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(defproject class/kubernetes-api "0.0.9.4"
    :description "CLASS - Rotterdam-CaaS"
    :url "-"
    :min-lein-version "2.0.0"
    :dependencies [[org.clojure/clojure "1.9.0"]
                   [org.clojure/core.async "0.4.474"]         ; Eclipse Public License; https://github.com/clojure/core.async
                   [compojure "1.6.1"]                        ; https://github.com/weavejester/compojure
                   [ring/ring-json "0.4.0"]                   ; MIT License https://github.com/ring-clojure/ring-json
                   [ring/ring-defaults "0.3.2"]               ; MIT License https://github.com/ring-clojure/ring-defaults
                   [ring-cors/ring-cors "0.1.12"]             ; Eclipse Public License
                   [org.clojure/tools.logging "0.3.1"]        ; Eclipse Public License - Version 1.0 https://github.com/clojure/tools.logging
                   [log4j/log4j "1.2.17"                      ; Apache License, Version 2.0 http://logging.apache.org/log4j/1.2/
                    :exclusions [javax.mail/mail
                                 javax.jms/jms
                                 com.sun.jdmk/jmxtools
                                 com.sun.jmx/jmxri]]
                   [io.replikativ/konserve "0.5.0-beta4"]     ; Eclipse Public License - Version 1.0;  "0.5.0-beta4", https://github.com/replikativ/konserve
                   [clj-http "3.9.1"]                         ; MIT License         https://github.com/dakrone/clj-http/
                   [org.clojure/data.json "0.2.6"]]           ; EPL - Version 1.0   https://github.com/clojure/data.json
    :plugins [[lein-ring "0.12.0"]                            ; https://github.com/weavejester/lein-ring
              [lein-uberwar "0.2.0"]]
    :uberwar {:handler atos.class.restapi.handler/app}
    :ring {:handler atos.class.restapi.handler/app
           :port 18083
           :open-browser? false
           :resources-war-path "WEB-INF/classes/"}
    :test-paths   ["test"]
    :profiles {:dev {:dependencies [[javax.servlet/servlet-api "2.5"]
                                    [ring/ring-mock "0.3.0"]]}}
    ;; jvm configuration
    :jvm-opts ["-Xmx256M" "-Djava.net.preferIPv4Stack=true"])
