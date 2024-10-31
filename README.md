# Backup

Tools for backup via rsync.

## Restore Backup

```bash
ssh-add ~/.ssh/id_ed25519_personal
ssh -A fire-k3s-dev.hm.benjamin-borbe.de
sudo -E -s
cd /var/lib/rancher/k3s/storage
DIR=pvc-0cb8996e-18b1-48ae-b959-c9c3cea1f8e8_dev_datadir-core-signal-finder-dwx-gold-0 \
rsync -av --progress --rsync-path="sudo rsync" "bborbe@hell.hm.benjamin-borbe.de:/backup/$(hostname -f)/current/var/lib/rancher/k3s/storage/${DIR}" .
```

```
ssh-add ~/.ssh/id_ed25519_personal
ssh -A fire-k3s-dev.hm.benjamin-borbe.de
sudo -E -s
cd /var/lib/rancher/k3s/storage
for DIR in $(ls -1d *-core-discord-*); do
rsync -av --progress --rsync-path="sudo rsync" "bborbe@hell.hm.benjamin-borbe.de:/backup/$(hostname -f)/current/var/lib/rancher/k3s/storage/${DIR}" .
done
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
