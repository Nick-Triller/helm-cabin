# Helm Cabin

Helm Cabin is a web UI that visualizes Helm releases 
in a Kubernetes cluster. 

## Motivation
 
I wasn't able to find a simple web UI that visualizes the data managed by 
Tiller without hiding information behind additional abstractions. 
I decided to scratch my own itch and started Helm Cabin as part of 
[Hacktoberfest 2019](https://hacktoberfest.digitalocean.com/). 

## Install

Install Helm Cabin with the provided chart. 

```bash
helm upgrade --install TODO
```

Helm Cabin doesn't handle TLS itself. Please use a reverse proxy, 
e. g. [Traefik](https://traefik.io/), for TLS termination.

## Build

TODO

## Project architecture

The project layout is based on 
[golang-standards/project-layout](https://github.com/golang-standards/project-layout).

### Backend

The backend periodically retrieves all releases including deleted releases and superseded revisions  
from Tiller. 
The result is cached in memory. 
Obviously, polling is a suboptimal solution if the number of 
releases (or revisions) is large because a lot of data is transferred on each poll.
Code of the Helm client is used to communicate with Tiller.

Helm Cabin tries to connect with Tiller via `tiller-deploy.svc.kube-system.cluster.local` first. 
If no connection could be established, it is assumed Helm Cabin runs outside a Kubernetes cluster. 
Helm Cabin tries to port-forward using the current Kubernetes context in this case. 
Port-forwarding mode is intended for development only.

### Frontend

The frontend source is located in the web directory and uses VueJS. 
Take a look at the [frontend README](./web/README.md) for usage instructions. 
