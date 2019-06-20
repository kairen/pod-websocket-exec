# Exec a Pod via WebSocket
Learning how to exec a pod via WebSocket written in Go.

## Building
First clone the repo, and execute `Makefile` using make tool:
```sh
$ git clone https://github.com/kubedev/websocket-exec.git $GOPATH/src/github.com/kubedev/websocket-exec
$ cd $GOPATH/src/github.com/kubedev/websocket-exec
$ make
```

## Running
To exec the command on pod as below:
```sh
$ kubectl run nginx --image nginx --restart=Never
$ ./k8s-ws-exec -p nginx -c nginx -t -i --command sh
# ls
bin   dev  home  lib64	mnt  proc  run	 srv  tmp  var
boot  etc  lib	 media	opt  root  sbin  sys  usr
```
