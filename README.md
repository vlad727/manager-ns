# Manager for kubernetes namespaces
Application creates limit range, resource quota and role binding for new namespaces also it set annotation for namespace "requester: system:serviceaccount:vlku4:vlku4"
## How to build application
docker build -t webhook-app .
## How to deploy it to k8s
Templates for applicatoin should be contain deployment with image,validating webhook, cluster role binding, quota template and limit range template
