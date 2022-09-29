# APKCTL: CLI for API Platform for K8s

Command Line Tool for APK, API Management Platform for Kubernetes to address the evergrowing industrial need to have a native experience when it comes to API Management in Kubernetes.

## Getting Started 
<hr>

To get started with APKCTL, you will need to have the prerequisites listed below. Once, that's sorted out you are good to use APKCTL.

### Prerequisites
<hr>

1. Make sure you have [Rancher Desktop](https://rancherdesktop.io/), [Minikube](https://minikube.sigs.k8s.io/docs/start/) or any other container orchestration platform installed in your machine.

2. Install the following (note that the tested versions are mentioned)

- **Go**: tested using version 1.19.1
- **Kubernetes**: tested using v1.24.3 (envoy gateway tested version - 1.24.1)
- **Docker**: tested using Docker version 20.10.17-rd
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- **Dockerd** is used as the container runtime (If your using Rancher,  you can verify the container runtime by navigating to Preferences -> Container Runtime 

### Clone and Setup APKCTL
<hr>

- Fork and clone/directly clone [this repo](https://github.com/BLasan/APKCTL-Demo).

- Change directory to root directory of cloned repo.

- Run the following command to build the ctl.

```go build apkctl.go```

- In order to try out the basic scenarios within APKCTL, you’ll need to deploy a sample service. You can simply deploy a sample service by executing the following command.

```kubectl apply -f sample-resources/BackendService.yaml```

- Command to install the platform.

```./apkctl install platform```

### API Create and Deploy using API Definition
<hr>

Use the following command to create an API using the sample [petstore swagger](https://github.com/BLasan/APKCTL-Demo/blob/main/sample-resources/swagger.yaml) and deploy to a particular namespace (deployed to default namespace if not provided).

```./apkctl create api swaggertest -f sample-resources/swagger.yaml --version 1.0.0 -n default```

### API Create and Deploy using Service URL
<hr>

Use the following command to create an API using the service URL (note that an API definition need not be provided for this scenario).

```./apkctl create api serviceurltest --service-url http://httpbin.default.svc.cluster.local:80```

### Invoking the Deployed API
<hr>

Once the deployment is succesful, you can try out the deployed API. First, look for the external IP of the gateway service (execute `kubectl get all --all-namespaces` command and look for the external IP of the service with name envoy-eg). Then invoke the API using the external IP as shown below.

```curl --verbose --header "Host: www.example.com" http://<EXTERNAL-IP>:8080/api/v3/pet/3```

### API Create Command with Dry Run Option
<hr>

You can use the API create command with  `--dry-run` option to generate the API project and store it in your file system. Note that the API does not get deployed with this command.

```./apkctl create api swaggertest -f sample-resources/swagger.yaml --version 1.0.0 -n default --dry-run```

### API Delete
<hr>

Use the following command to delete an API.

```./apkctl delete api swaggertest```

### Get APIs
<hr>

You can list down the APIs that are deployed using the below command.

```./apkctl get apis -n default```

### Clean up
<hr>

Clean up by executing the following commands once you are done.

- Command to delete the sample service.

```kubectl delete -f sample-resources/BackendService.yaml```

- Command to uninstall the platform.

```./apkctl uninstall platform ```

## List of Commands
<hr>

- apkctl install platform
- apkctl create api \<API Name\> --version \<API Version\> -n \<Namespace\> -f \<API Definition File\> --service-url \<Service URL\> --dry-run
- apkctl delete api \<API Name\>
- apkctl get apis -n \<Namespace> -o \<Outptut Format\>
- apkctl uninstall platform
