# Backup

Tools for backup via rsync.

## List

http://backup.hell.hm.benjamin-borbe.de/setloglevel/4
http://backup.hell.hm.benjamin-borbe.de/list
http://backup.hell.hm.benjamin-borbe.de/trigger

## Build

```
VERSION=develop make build upload
```

## Keel

```
helm repo add keel https://charts.keel.sh
helm repo update
```

Install through Helm (with Helm provider enabled by default):

Keel must be installed into the same namespace as Tiller, typically kube-system

```
helm upgrade --install keel --namespace=kube-system keel/keel
```
