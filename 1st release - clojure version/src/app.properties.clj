{
  :app {
    ;; REST-API properties
    :name "class-k8s-rest-api"
    :version "0.0.9.4"
    :conf-file "app.properties.clj"
    ;; USED in logs
    :log-message "Rotterdam-CaaS REST API [0.0.9.4] "
    ;; database path
    :db-path "/tmp/store"
  }
  :apis {
    ;; server
    :server-ip "192.168.7.28"
    ;; oauth-token
    :openshift-oauth-token "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJjbGFzcyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJjbGFzcy1hcHAtdG9rZW4tajZrdnMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiY2xhc3MtYXBwIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiMWMzMDg5MTYtZjE1Yi0xMWU4LWI2NDItMDA1MDU2OTg2MDU5Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmNsYXNzOmNsYXNzLWFwcCJ9.s9ZvaDq7TKf4m3X36swQUD8g09HIRlTWptoiZ0jlFkx-1JZKuJu22LzXgP6J4PeCh3OJg16CuYyxAgRrogYERKJVf5FD1WBb0DQH5qVmtlAC5LtV823G5admXmZ8AMGbzn5TNELTECw6IOfE5PU3baRdRMvXVyli96HZpnUpVwqV8t5vh0HcoySn4vM7L6dT8ug9cJfqtfS5v_2WTfiLxr_cKhlyXOV_uRBWVcBBLL6M2UHZ9Ex00AYoA6xwhLvhgJ3yi8h3CflrueIKMlhsBI5aYrollszuKNqxJu-ctQCaRuslgkg6Ew5_SYQbivQShjpYi1nNgnZFr4LqrGVjEQ"
    ;; KUBERNETES connection data
    :kubernetes {
      :url  "http://192.168.7.28:8001"
    }
    :openshift {
      :url  "https://192.168.7.28:8443"
    }
  }
}
