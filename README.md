# k8s-teamhack demo

## build and run locally using docker

1. Run `complexity-service`
```
docker build -t ${DOCKER_USERNAME}/complexity:0.1 -f ./Dockerfile.complexity . 
docker run -it -p 3000:3000 --name complexity-service ${DOCKER_USERNAME}/complexity:0.1
```
2. Get `complexity-service` "public" IP Address:
```
docker inspect complexity-service | grep IPAddress
```

> **NOTE**: Make sure that there is `bridge` connection set up in your local Docker environment:
> `docker network ls`:
> ```
> NETWORK ID     NAME           DRIVER    SCOPE
> 5236b3958156   bridge         bridge    local
> ```

3. Run `tasks-service` and pass IP Address of `complexity-service` as environmental variable

```
docker build -t ${DOCKER_USERNAME}/tasks:0.1 -f Dockerfile.tasks .
docker run -it -p 8080:8080 --env APP_COMPLEXITY_SERVICE_URL=http://172.17.0.2:3000/complexity --name tasks-service ${DOCKER_USERNAME}/tasks:0.1
```

4. Run example calls from `hack/example.http`

## deploy on k8s cluster using kind

### Prerequisites
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/)

```
kubectl apply -f ${YAML_MANIFEST_PATH}.yaml
```


