# Digital Ocean K8s challenge

This is a project for [Kubernetes Challenge](https://www.digitalocean.com/community/pages/kubernetes-challenge).

## Requirements:
* Account in Digital Ocean and valid ReadWrite API token
* [Pulumi](https://www.pulumi.com/) binary (already logged)
* Golang > 1.14
* make
Optional:
* doctl

## Quick start
```
DIGITALOCEAN_TOKEN=<TOKEN> make
```
The above command will:

* Create managed k8s in Digital Ocean
* Retreive kubeconfig and export in to pulumi state with key `kubeconfig`
* Deploy a kubegres controller based on official manifests with several changes for HA
* Deploy a kubegres managed postgresql cluster with reasonable defaults (multiple replicas, backup schedule, run as postgres user)

> You can use pulumi commands directly without `make`. Please see the Makefile.

To delete the k8s cluster
```
DIGITALOCEAN_TOKEN=<TOKEN> make clean
```
> Do not forget remove volumes from digital Ocean if you do not need it

## Demo
There is a demo video for the project on [Youtube](https://www.youtube.com/watch?v=9GhI34LqLhs) 

*Please turn on subtitles!*

1. Observe a k8s cluster managed by Digital Ocean and a Kubegres cluster
2. Connect to Posgresql cluster and generate some data
3. Increase replicas count and watch events
4. Provoke auto failover (via node draining)
5. trigger a backup

## Notes
The kubegres CRD for Pulumi can be regenerate via `make crd` command. crd2pulumi binary required.
