# Rotterdam-CaaS

[![version](https://img.shields.io/badge/version-0.0.9--SNAPSHOT-blue.svg)]()

&copy; Atos Spain S.A. 2016

Rotterdam-CaaS (_version 0.0.9_) is a component of the European Project class (https://class-project.eu/).

-----------------------

[Description](#description)

[Installation Guide](#installation-guide)

[Usage Guide](#usage-guide)

[LICENSE](#license)

-----------------------

### Description


#### Requirements

#### 1. Installation in Kubernetes

1. Kubernetes

#### 2. Server installation (Linux / Windows)

1. Java version 8 or higher
2. leiningen
3. Connection to a Kubernetes REST API


-----------------------

### Installation Guide

- You can download the repository and launch the server application using **leiningen** (see "_Usage Guide - 2. Server installation (Linux / Windows)_")

- You can download the repository and create the docker image from the "clojure version" folder:

    ```bash
    sudo docker build -t rotterdam .
    sudo docker run -p 8082:8082 rotterdam
    ```

- You can also run **Rotterdam** in Docker by pulling the image from Docker Hub:

    ```bash
    docker pull atosclass/rotterdam-caas:0.0.9.4
    docker run [OPTIONS] atosclass/rotterdam-caas:0.0.9.4 [COMMAND] [ARG...]
    ```

- Finally, to run **Rotterdam** in Openshift, deploy the [image](https://hub.docker.com/r/atosclass/rotterdam-caas) from Docker Hub using the OKD UI. The following environment variables can be defined:

    - **KubernetesEndPoint** (e.g.) "http://192.168.7.28:8001"
    - **OpenshiftEndPoint** (e.g.) "https://192.168.7.28:8443"
    - **ServerIP** (e.g.) "192.168.7.28"
    - **OpenshiftOauthToken** (e.g.) "eyJhbGciOiJSUzI1 ... 3MiOiJrdWJlcm5ldGVzL3Nlc"

    1. In OKD Web UI, go to selected project / namespace, i.e. the _default_ namespace, and select `Add to project > Deploy Image`
	   ![rotterdam](https://github.com/class-rotterdam/Rotterdam/raw/master/images/deploy.png)
	
    2. Select `Image Name`
        - `atosclass/rotterdam`
          - Name: `rotterdam-caasXXX`
          - Environment Variables: `OpenshiftOauthToken`
    3. Deploy
    4. Go to new application / deplopyment and select `Create Route`
        - Rotterdam - Hostname: `rotterdam-cass.192.168.7.28.xip.io`

-----------------------

### Usage Guide

##### 1. Installation in Kubernetes

If the application was successfully deployed in K8s or Openshift, it can be accessed in port `8082`. This can be changed in the _Service_ yaml (parameters _ports-port_ and _externalIPs_).

##### 2. Server installation (Linux / Windows)

Launch application:

```bash
cd rotterdam_path
lein ring server
```

By default the API can be accessed in port `8082`. This can be modified in `project.clj` file, before launching the application.


---------------------------------

### LICENSE

`Rotterdam` is licensed under [Apache License, version 2](LICENSE.txt).

