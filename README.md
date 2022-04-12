# Build & Deploy
* Build the application docker image,
  ```
  make build
  ```  
* Create a K8S cluster to host the application. For example, to build a cluster using [kind](https://kind.sigs.k8s.io/),
  ```
  kind create cluster --name wolverine
  kind get kubeconfig > ~/.kube/wolverine-kubeconfig
  export KUBECONFIG=~/.kube/wolverine-kubeconfig
  chmod 600 ~/.kube/wolverine-kubeconfig
  kubectl cluster-info
  ```
* If using kind load the docker image from your host into the cluster. This makes available for kubernetes to pull. 
  ```
  make kind-load-image
  ```  
* Deploy the application into Kubernetes,
  ```
  make deploy
  ```
* Check the application logs,
  ```
  kubectl logs -f -l app=wolverine
  ```

# Deploy Prometheus
Install Prometheus using helm (requires helm 3).
```
$ helm version
version.BuildInfo{Version:"v3.8.1", GitCommit:"5cb9af4b1b271d11d7a97a71df3ac337dd94ad37", GitTreeState:"clean", GoVersion:"go1.17.5"}

$ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
$ helm repo update
$ helm install prometheus prometheus-community/prometheus
```

To view the Prometheus UI,
```
$ export POD_NAME=$(kubectl get pods --namespace default -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 9090
```
Then visit http://localhost:9090/

Prometheus is configured to scrape metrics from pods with the annotation `prometheus.io/scrape: "true"`. Two further annotations tell Prometheus what port and path to use. For example, the Wolverine pod spec includes,
```
prometheus.io/scrape: "true"
prometheus.io/path: /metrics
prometheus.io/port: "2112"
```

# Deploy Grafana
Install Grafana using helm (requires helm 3).
```
$ helm repo add grafana https://grafana.github.io/helm-charts
$ helm repo update
$ helm install -f k8s/grafana_values.yaml grafana grafana/grafana
```
The file `k8s/grafana_values.yaml` includes details of the Prometheus Datasource (from which Grafana can pull metrics) and enables the loading of dashboards from ConfigMaps.

Create the Wolverine dashboard using,
```
$ kubectl apply -f k8s/wolverine-grafana-dashboard.yaml
```

To view the Grafana UI:

Get the 'admin' user password using,
```
$ kubectl get secret --namespace default grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

Port-forward to the UI port in Grafana, 
```
$ export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=grafana" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 3000
```

To view the Wolverine dashboard,
* Open http://localhost:3000. Login with the usename 'admin' and the password retrieved above.
* Go to Dasboards -> Browse, and click on "Wolverine".
* Select the monitored URL from the "url" dropdown.

The dashboard has two charts, "Response Time" and "URL Status".

The "Response Time" chart reports the 95th percentile response time for successful (response code 200) HEAD requests.

The "URL Status" chart reports wheather the URL is up (response code 200) or down (not a 200 resonse code). This is simple "Status History" chart so it shows a green bar when the URL is up and a red bar when its down.

See the [screenshots](screenshots/) folder for a sample dashboard.
