# Building

1) Build app/

    cd app/
    docker build . -t app

2) Build sidecar/

    cd sidecar/
    docker build . -t sidecar-ctl

2) Deploy app

    cd ../
    kubectl create -f app-deployment.yml

3) Debug

    kubectl get deployments
    kubectl describe deployments
    kubectl get pods -l app=app
    kubectl logs app-deployment-664b98d96b-fw5pn app
    kubectl exec -it debugger bash

4) Clean up

    kubectl delete -f app-deployment.yml 
