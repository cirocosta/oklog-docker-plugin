<h1 align="center">oklog-docker-plugin 📝 </h1>

<h5 align="center">Forwards your container logs to Docker</h5>

<br/>

`docker-oklog-driver` is a custom Docker logging plugin that implements the log driver interface and acts as an OKLog forwarder.

### Install

```
# Install the plugin `cirocosta/oklog-docker-plugin`
# and then make it aliased to `oklog` so that we can
# make use of it via the alias.
#
# INGESTER configures the host name to use to look 
# for ingester nodes.
docker plugin install \
        --alias oklog \
        cirocosta/oklog-docker-plugin \
        INGESTER=127.0.0.1:7651
```

