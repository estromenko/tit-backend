# tit-backend

## Dependencies

### Mandatory

- Golang version 18(Installation instruction [link](https://go.dev/doc/install))

## Getting started

### Development

- You need to install dev dependencies

```bash
$ make install-dev-requirements
```

- Run PostgreSQL database

```bash
$ docker run -d --rm --name tutorintech_postgres -v /srv/_tutorintech_postgres:/var/lib/postgresql/data -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:15-alpine
$ docker run --rm -it --link tutorintech_postgres:postgres -e PGPASSWORD=secret postgres:15-alpine createdb -h postgres -U postgres tutorintech
```

- Install tool for migrations and apply them:

```
$ make install-migrate
$ migrate -source file://./migrations -database postgres://postgres:secret@localhost:5432/tutorintech\?sslmode=disable up
```

- Setup local image registry:

```bash
$ docker run --rm --name registry --net host -d registry:2 
$ docker login localhost:5000
```

- Build required images and push them to the registry:

```bash
$ docker compose -f docker/dashboard/docker-compose.yml build
$ docker tag tit-dashboard:latest localhost:5000/tit-dashboard:latest
$ docker push localhost:5000/tit-dashboard:latest
```

- Install k3s

IMPORTANT NOTE: default k3s setup is provided with flannel CNI,
which does not support network policies required for disabling networking inside dashboards.
Thats why we use `calico` as alternative.

Install k3s excluding flannel CNI using command below:

```bash
$ curl -sfL https://get.k3s.io | K3S_KUBECONFIG_MODE="644" INSTALL_K3S_EXEC="--flannel-backend=none --cluster-cidr=192.168.0.0/16 --disable-network-policy" sh -
```

Then install calico:

```bash
$ k3s kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/tigera-operator.yaml
$ k3s kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/custom-resources.yaml
```

Add following content to `/etc/rancher/k3s/registries.yaml` file:
```yaml
mirrors:
  local:
    endpoint:
    - "http://localhost:5000"
```

- Add following line to `/etc/hosts` file:

```
127.0.0.1 dashboards.tutorin.tech
```

- Create `.env` file with following content:

```dotenv
DEBUG=true
SECRET_KEY=secret
DASHBOARD_IMAGE=localhost:5000/tit-dashboard:latest
DASHBOARD_INGRESS_DOMAIN=dashboards.tutorin.tech
```

- Run application. You can configure app by passing env variables directly or create .env 
file in project root.

```bash
$ make dev
```

### Create superuser

To create superuser locally use command below:

```bash
$ make createsuperuser
```

In production environment `createsuperuser` binary is available from anywhere

### Production

Before deployment make sure `cert-manager` installed in your cluster.
If it is not, use command below:
```bash
$ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml
```

To deploy the app in production environment you should use werf
(Installation instruction [link](https://werf.io/documentation/v1.2/#installing-werf)).

To build and deploy app use command below:
```bash
$ werf converge --repo registry.tutorin.tech/tit-backend --env prod
```

To override default configuration you can use `--set` flags or custom `values.yaml` file:
```bash
$ werf converge --repo registry.tutorin.tech/tit-backend --env prod \
    --set env.PG_HOST=localhost --set env.PG_NAME=postgres
$ # Alternatively
$ werf converge --repo registry.tutorin.tech/tit-backend --env prod \
    --values /path/to/custom/values.yaml
```
