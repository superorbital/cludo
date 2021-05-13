# Cludo Server Kubernetes Manifests

* In addition to the other manifests in this directory you will need to create a that Kubernetes can use to connect to the private Docker Registry.

  * `kubectl create secret docker-registry docker-registry --docker-server=docker.io --docker-username=MYUSER --docker-password=MYPW --docker-email=MYEMAIL
