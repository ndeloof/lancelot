# Lancelot docker API proxy

![logo](logo.png)

## About

There's a few case (most notably: Continuous Integration) where a containerized process need 
to access docker infrastructure. A common practice is to bind mount host's `docker.sock`, but 
this comes with a high risk as, doing so, containerized process get granted access to the host 
with unrestricted privileges.

Lancelot is a docker API proxy which let you expose a Docker-like API endpoint, but is actually 
highly constrained to only let client use a subset of the Docker API, blocks some arguments 
(like `privileged`) and will also force a few options to restrict hosts resources usage. 

It also track resources created by client so it can easily cleanup after usage, removing the 
classic CI issue of volumes|images|containers accumulating on host. When closed, Lancelot 
removes all resources.

## Design

Lancelot is a HTTP server to replicate the Moby endpoints. It fully parse the API payload to
copy supported fields _only_ into Moby API Client structs, and add it's own ones.

```console

 +---------------+                   +----------------+                      +----------------+
 | docker client |--------> tcp:2375 | lancelot proxy | -------> docker.sock | docker engine  |
 +---------------+                   +----------------+                      +----------------+

                                        //  only allow API/parameters subset
                                        //  tag resources as com.docker.lancelot

```

Lancelot is a full Docker engine API implementation, it fully parse the payload objects and rebuild
them to access the actual engine using docker sdk. Doing so we enfore Lancelot does only expose a
curated set of attributes, compared to other proxying solution which just filter the REST paths.

Also, Lancelot injects a few parameters and label to track resource usage and enforce limits. 
Typically, Lancelot is ran with a control group parameter (defaults to self) that it will apply
to all container as `parent_cgroup`. Doing so, we enfore that code accessing Docker engine
through a Lancelot proxy won't be able to globaly consume more than configured memory/cpu.

### Moby API subset

Lancelot implements a subset of the Moby API, only methods that are required for the usages considered
safe without risk to impact other users or alter the docker host.

### Resource tracking

Lancelot track resources you created by it's API. Doing so it can allow use of `--volumes-from` as long
the target is a container you created, and will block otherwise.

## Usage

Lancelot is a work in progress. This section is set for illustration of our end goal
1. Setup a lancelot proxy as a "safe(r) docker session"
`docker run -v /var/run/docker.sock:/var/run/docker.sock -P 2375:2375 ndeloof/lancelot`
2. Configure your docker client to use lancelot instance as docker host
`export DOCKER_HOST=tcp://localhost:2375
3. Enjoy docker in a safe(r) context
`docker run ...`
`docker run --privileged` => Error
4. Stop lancelot proxy and get all resources released

## Why the name?

Arthurian legend's _Lancelot du Lac_](https://en.wikipedia.org/wiki/Lancelot) combines the idea for
being strict and inflexible regarding his role, and has a name with aquatic-life reference as required
for any Docker-related project.

_Lancelot_ is also the name for [famous beer drafted in britanny](http://brasserie-lancelot.bzh/) I enjoy.

## About Security

Lancelot does **not** make your Docker engine multi-tenant. Don't use it to
pretend isolate arbitrary third-party code into sandboxes. Anyway, there's
many scenarios where a legitimate container has to run with access to the
underlying Docker engine, and you don't want to grant it full access.
Authz plugins are a partial solution (see [Authobot](https://github.com/ndeloof/authobot)),
and Lancelot can also be used in combination with other solutions as an additional 
layer to implement [defense in depth](https://en.wikipedia.org/wiki/Defense_in_depth_%28computing%29).