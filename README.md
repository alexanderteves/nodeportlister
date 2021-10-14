# nodeportlister

This small web application runs in your Kubernetes cluster and provides a quick and clean overview of configured nodeports:

```
NODEPORT(S)    SERVICE                      NAMESPACE
32202          jenkins                      development
32205          logging-backend              engineering
32204          logging-frontend             engineering
32203          php-backend                  web
30044          nginx                        web
```

## Usage

For deployment, install the Helm v3 chart inside `helm/`. Modify the ingress' hostname according to your setup in `helm/values.yaml`. The corresponding Docker image can be found [here](https://hub.docker.com/repository/docker/alexanderteves/nodeportlister).

## Compatibility

Binaries are compiled with `client-go` `v0.18.19` - if you run another Kubernetes version in your custer, your experience may vary.
