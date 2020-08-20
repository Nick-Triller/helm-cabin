# Helm Cabin

Helm Cabin is a web UI that visualizes Helm releases 
in a Kubernetes cluster. 

![](./screenshots/screenshot001.png)

Warning: Helm-Cabin shows all secrets that were created through Helm!

## Motivation
 
I wasn't able to find a simple web UI that visualizes the data managed by 
Tiller without hiding information behind additional abstractions. 
I decided to scratch my own itch and started Helm Cabin as part of 
[Hacktoberfest 2019](https://hacktoberfest.digitalocean.com/). 

## Features

- Supports Helm 2 & 3
- List all releases with any status (deleted, superseded, deployed, etc.)
- View revisions, rendered manifest, chart templates, chart values and chart files for any release

## Install

Install Helm Cabin with the provided chart. 

```bash
helm repo add helm-cabin https://nick-triller.github.io/helm-cabin/
helm repo update
helm upgrade --install --set-string helmVersion=3 helm-cabin helm-cabin/helm-cabin
```

Helm Cabin doesn't handle TLS itself. Please use a reverse proxy, 
e. g. [Traefik](https://traefik.io/), for TLS termination.

## Build from source

Helm Cabin uses [Magefile](https://github.com/magefile/mage) as task run. 
Check out the mage targets in `Magefile.go`, the scripts section in `web/package.json` 
and `build/Dockerfile` to see how the project is built. 

## Development

If you work in Helm 2 mode:

```bash
# Port-forward to tiller if you work in Helm 2 mode
kubectl port-forward svc/tiller-deploy -n kube-system 44134:44134
# Start the frontend through vue-cli-service. 
# The vue dev server proxies requests starting with /api to localhost:8080
# The frontend will be served on localhost:8081 by default.
npm --prefix web run serve
# Start the backend on localhost:8080 (assumes tiller is reachable at localhost:44134 in Helm 2 mode
# or a valid kube context is set in Helm 3 mode)
mage RunServer
```

## Project architecture

The project layout is based on 
[golang-standards/project-layout](https://github.com/golang-standards/project-layout).

### Backend

The backend periodically retrieves all releases including deleted releases and superseded revisions  
from Tiller (Helm 2) or Kubernetes secrets (Helm 3). 
The result is cached in memory. 

Helm Cabin tries to connect with Tiller via `tiller-deploy.kube-system.svc.cluster.local` by default. 
The tiller host can be overriden with the `tillerAddress` CLI flag.
Use port forwarding for local development.

If you work in Helm 3 mode, Helm Cabin uses the current context to 
connect to Kubernetes.

### Frontend

The frontend source is located in the web directory and uses VueJS. 
Take a look at the [frontend README](./web/README.md) for usage instructions. 
