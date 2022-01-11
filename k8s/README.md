# Cludo Server (cludod) Kubernetes Manifests

It is easy to install the cludo server (cludod) using either [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/), [kustomize](https://kustomize.io/) or [helm](https://helm.sh/).

## kustomize (and kubectl)

* **Note**: You must provide a `cludod.yaml` that contains your real secrets. At the moment `kustomize` looks for a `.gitignored` file here: `k8s/kustomize/base/files-secrets/secret-cludod.yaml`. You can find an example of what this file should look like here: `k8s/kustomize/base/files-secrets/cludod-EXAMPLE.yaml`

* View manifests with recent `kubectl` releases
  * Note that there are multiple overlays. Apply the one appropriate to your use case.

```sh
kubectl kustomize k8s/kustomize/overlays/development
```

* View manifests with `kustomize` 4+

```sh
kustomize build k8s/kustomize/overlays/development
```

* Apply and delete manifests with recent `kubectl` releases

```sh
kubectl apply -k k8s/kustomize/overlays/development/
kubectl delete -k k8s/kustomize/overlays/development/
```

## helm

* Apply and delete manifest with `helm` 3+

```
helm repo add superorbital https://helm.superorbital.io/
helm install cludod superorbital/cludod
helm uninstall cludod
```
