# Exec a Pod via WebSocket
Learning how to exec a pod via WebSocket writing in Go.

## Building
First clone the repo, and execute `Makefile` using make tool:
```sh
$ git clone https://github.com/kairen/k8s-ws-exec.git $GOPATH/src/github.com/kairen/k8s-ws-exec
$ cd $GOPATH/src/github.com/kairen/k8s-ws-exec
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