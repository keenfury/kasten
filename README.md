## Kasten Assesement

This code base will expose Prometheus metrics for:

- /powerball (type Couter)
- pod count of Kubernetes cluster provided by Kasten (type Gauge)

### Setup

- git clone this repo
- for security sake, k8s config is not provided, you will need to add a kube_config file with the correct config information to the cluster, place in the root directory of the project
- go build (version: +1.16.0)
- ./kasten (start server)

### Usage

The server will start on port 8080

Available endpoints:

- /powerball: will give you your random powerball numbers
- /metrics: prometheus metrics (under: "a_kasten_assessment" subsystem)