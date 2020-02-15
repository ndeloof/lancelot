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
