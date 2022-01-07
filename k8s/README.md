# Cludo Server (cludod) Kubernetes Manifests

It is easy to install the cludo server (cludod) using either [kustomize](https://kustomize.io/) or [helm](https://helm.sh/).

## Kustomize

* **Note**: You must provide a `cludod.yaml` that contains your real secrets. At the moment `kustomize` looks for a `.gitignored` file here: `k8s/kustomize/base/files-secrets/secret-cludod.yaml`. You can find an example of what this file should look like here: `k8s/kustomize/base/files-secrets/cludod-EXAMPLE.yaml`

* With recent`kubectl` releases

```sh
kubectl kustomize k8s/kustomize/base
kubectl kustomize k8s/kustomize/overlays/development
kubectl kustomize k8s/kustomize/overlays/staging
kubectl kustomize k8s/kustomize/overlays/production
```

* With `kustomize` 4+

```sh
kustomize build k8s/kustomize/base
kustomize build k8s/kustomize/overlays/development
kustomize build k8s/kustomize/overlays/staging
kustomize build k8s/kustomize/overlays/production
```

## Helm

* **Note**: You must provide a `cludod.yaml` that contains your real secrets. At the moment `helm` looks for a `.gitignored` file here: `k8s/helm/cludod/secret-cludod.yaml`. You can find an example of what this file should look like, in the kustomize folder, here: `k8s/kustomize/base/files-secrets/cludod-EXAMPLE.yaml`

* With `helm` 3+

```
helm install cludod ./k8s/helm/cludod/
helm uninstall cludod
```
