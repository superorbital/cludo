# Cludo Server Kubernetes Manifests

* **Note**: You will need edit `cludod-secret.yaml` and add in your `cludod` configuration YAML before applying it.
  * You will also need to create an additional secret that Kubernetes can use to connect to the private Docker Registry.

  * `kubectl create secret docker-registry docker-registry --docker-server=docker.io --docker-username=MYUSER --docker-password=MYPW --docker-email=MYEMAIL
