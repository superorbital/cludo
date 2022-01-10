# Cludo Server (cludod) Kubernetes Manifests

It is easy to install the cludo server (cludod) using either [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/), [kustomize](https://kustomize.io/) or [helm](https://helm.sh/).

## kustomize (and kubectl)

* **Note**: You must provide a `cludod.yaml` that contains your real secrets. At the moment `kustomize` looks for a `.gitignored` file here: `k8s/kustomize/base/files-secrets/secret-cludod.yaml`. You can find an example of what this file should look like here: `k8s/kustomize/base/files-secrets/cludod-EXAMPLE.yaml`

* With recent `kubectl` releases

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

## helm

* With `helm` 3+

```
helm repo add superorbital https://helm.superorbital.io/
helm install cludod superorbital/cludod
```
