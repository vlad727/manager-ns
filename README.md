# Manager for kubernetes namespaces
Application creates limit ranges, resource quota and role binding for new namespaces kubernetes.
## How to build application
docker build -t webhook-app .
## How to deploy it to k8s
Templates for applicatoin should be contain deployment with image,validating webhook, cluster role binding, quota template and limit range template
