# DevOpsSeptember2018
Stand-up Kubernetes Infra on Azure using Terraform, then code a self-replicating virus like app, and infect your cluster with it oO

## Steps to Deploy

* Ensure you're running Go version 1.10 or higher `go version`
* Deploy using

```bash
go get -u github.com/kim0/DevOpsSeptember2018 && cd ~/go/src/github.com/kim0/DevOpsSeptember2018
terraform init
terraform apply -auto-approve
mkdir -p ~/.kube && terraform output kube_config > ~/.kube/config
kubectl get po -o wide
go build -o virus
# Infect!
./virus
```

In a separate window, you might want to follow progress by

```bash
watch -n 0.2 'kubectl get po -o wide'
kubectl logs -f azure-loves-devops
```
