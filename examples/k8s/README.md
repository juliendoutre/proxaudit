# Kubernetes sidecar HTTP(s) proxy

This guide walks you through using proxaudit as a Kubernetes sidecar container.

This tutorial makes use of the [Kind project](https://kind.sigs.k8s.io/) to run a local k8s cluster. You can install Kind on MacOS with:
```shell
brew install kind
kind --version
```

For other platforms, please check their [quick start page](https://kind.sigs.k8s.io/docs/user/quick-start/).

Then you can get started by creating a cluster and switching context with:
```shell
kind create cluster
kubectl cluster-info --context kind-kind
```

> Feel free to use your own Kubernetes cluster instead of Kind if you want, this example should work the same.

In order to deploy the example, run:
```shell
kubectl apply -f ./pod.yaml
```

This creates a pod with 3 containers:
- an init container running `mkcert` to generate a CA key and certificate. It writes them to a volume that will be shared with other containers.
- a sidecar container running proxaudit, reading the CA key and certificate from the volume previously populated by the `mkcert` container.
- an ubuntu `box`. In order for HTTPS interception to work smoothly, the container must have access to the CA certificate through the volume previously populated by the `mkcert` container **and** install it in its trustore. Hence the `apt update && apt upgrade --yes && apt install --yes curl mkcert && mkcert -install` command. `&& sleep infinity` is appended to make sure the container does not stop immediatly and lets users a chance to `kube exec` in. `HTTP_PROXY` environment variables are configured to point to `localhost` since in Kubernetes, pods share the same network namespace.

You can check the pod started successfully with:
```shell
kubectl get pods
```

Then try to perform some HTTP requests from the `box` container:
```shell
kubectl exec box -c box -it  -- curl https://google.com
```

Streaming the `proxaudit` sidecar container logs should display the intercepted requests:
```shell
kubectl logs box -c proxaudit
```

When you're done, you can delete the test resource with:
```shell
kubectl delete pod box
```

For Kind users, don't forget to delete the local cluster with:
```shell
kind delete cluster
```
