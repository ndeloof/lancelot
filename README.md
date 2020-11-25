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

```

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

