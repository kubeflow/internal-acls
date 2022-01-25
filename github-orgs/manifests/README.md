# Auto sync for GitHub org

* **project**: kubeflow-admin
* **cluster**: kf-admin-cluster
* **namespace**: github-admin

## Deploy new updates

Each change in this folder need to be deployed to take effect, only Kubeflow admins
have the permission to do so. Follow these steps to connect to the admin cluster:

```bash
# First time
gcloud container clusters get-credentials kf-admin-cluster --project kubeflow-admin --region us-central1-a
# Rename the context to make future usage easier
kubectl config rename-context gke_kubeflow-admin_us-central1-a_kf-admin-cluster kf-admin
# Next time, we can switch to this context via
kubectl config use-context kf-admin
```

To update the deployments:

```bash
git checkout master
git pull
make apply
```

## Trigger a github org sync manually

After connecting to `kf-admin` cluster like above:

```bash
make run-github-sync-once
```

## Github Token

We need a [GitHub token](https://github.com/settings/tokens) with admin:org privileges.

```bash
kubectl -n github-admin create secret generic github-org-admin-token-zijianjoy --from-literal=github_token=<TOKEN>
```

* We are currently using the token **Kubeflow github-admin token** owned by zijianjoy.

## Validate config map

We use a config map to provide the python code used to validate the config

```bash
make create-config-map
```
