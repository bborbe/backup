# Backup

Tools for backup via rsync.

## Build

```
BRANCH=develop make build upload
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
