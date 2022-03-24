To run (with default kind context as target):

```
make docker-build
make deploy
kubectl apply -f config/samples/tokenmap.yaml
```

You should see this output in namespace `test-sdk-system`

```console
➜  test-sdk kubectl apply -f config/samples/tokenmap.yaml 
tokenmap.k8ssandra.io/tokenmap-sample created
➜  test-sdk kubectl get all
NAME                                               READY   STATUS              RESTARTS   AGE
pod/cluster0-node-0                                0/1     Completed           0          1s
pod/cluster0-node-1                                0/1     ContainerCreating   0          1s
pod/cluster0-node-2                                0/1     ContainerCreating   0          1s
pod/test-sdk-controller-manager-7f8bb84999-6f292   1/1     Running             0          12s

NAME                                                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/test-sdk-controller-manager-metrics-service   ClusterIP   10.96.217.226   <none>        8443/TCP   12s

NAME                                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/test-sdk-controller-manager   1/1     1            1           12s

NAME                                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/test-sdk-controller-manager-7f8bb84999   1         1         1       12s
➜  test-sdk kubectl logs pod/cluster0-node-0
Replacing 10.64.1.69 with cluster0-node-0 10.244.2.3%                                                                                                                                                                           
➜  test-sdk kubectl logs pod/cluster0-node-1
Replacing 10.64.6.59 with cluster0-node-1 10.244.1.5%                                                                                                                                                                           
➜  test-sdk kubectl logs pod/cluster0-node-2
Replacing 10.64.7.4 with cluster0-node-2 10.244.6.4%                                                                                                                                                                            
➜  test-sdk
```