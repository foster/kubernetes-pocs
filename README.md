# Building

1) Build app/

        cd app/
        docker build . -t app

2) Build sidecar/

        cd sidecar/
        docker build . -t sidecar

# Running

1) Deploy app

        kubectl create -f app-deployment.yml

2) Debug

        kubectl get deployments
        kubectl describe deployments
        kubectl get pods -l app=app
        kubectl logs app-deployment-664b98d96b-fw5pn app
        kubectl logs -l app=app -c app
        kubectl exec app-deployment-5788b7f744-d5wbd -c sidecar -it sh

3) Scale

        kubectl scale deployment.v1.apps/app-deployment --replicas=1

4) Clean up

        kubectl delete -f app-deployment.yml 
