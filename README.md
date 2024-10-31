# Backup

Tools for backup via rsync.

## Restore Backup

```bash
ssh-add ~/.ssh/id_ed25519_personal
ssh -A fire-k3s-dev.hm.benjamin-borbe.de
sudo -E -s
cd /var/lib/rancher/k3s/storage/pvc-37f07c6e-8358-4202-abb8-334b8d21e86b_dev_datadir-frontend-candle-0
rsync --progress --rsync-path="sudo rsync" bborbe@hell.hm.benjamin-borbe.de:/backup/fire-k3s-dev.hm.benjamin-borbe.de/current/var/lib/rancher/k3s/storage/pvc-37f07c6e-8358-4202-abb8-334b8d21e86b_dev_datadir-frontend-candle-0/bolt.db .
```

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
