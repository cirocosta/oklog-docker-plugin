<h1 align="center">oklog-docker-plugin 📝 </h1>

<h5 align="center">Forwards your container logs to Docker</h5>

<br/>

`docker-oklog-driver` is a custom Docker logging plugin that implements the log driver interface and acts as an OKLog forwarder.

### Install

```sh
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


# Run a container with `oklog` being used as the log
# plugin. 
docker run \
        --log-driver oklog \
        alpine echo lol

# Query oklog for the the last logs stored in the last
# 1m
oklog query -from 1m
cid=f034d50a48a964cc5ef7b4dd678cc155246de8209c32a25891cfe81c9ce296fc lol


# Run a container with `oklog` as the log plugin but
# capturing the labels specified so that we can query
# against them in the future.
docker run \
        --log-driver oklog \
        --log-opt labels=foo \
        --label foo=bar \
        alpine echo lol2

# Query all the logs from the last 3min
oklog query -from 3m
cid=f034d50a48a964cc5ef7b4dd678cc155246de8209c32a25891cfe81c9ce296fc lol
cid=fcf89d4cd2a6bcb33b50cb21fe29d4245efb473e94e8b427fb252ed298434dc1 foo=bar lol2

# Query all the logs from the last 3min but only
# those that container `foo=bar`
oklog query -from 3m -q foo=bar
cid=fcf89d4cd2a6bcb33b50cb21fe29d4245efb473e94e8b427fb252ed298434dc1 foo=bar lol2
```


