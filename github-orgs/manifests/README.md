# Auto sync for GitHub org

* **project**: kubeflow-admin
* **cluster**: kf-admin-cluster 
* **namespace**: github-admin

## Github Token

We need a GitHub token with admin:org priveleges

```
kubectl -n github-admin create secret generic github-org-admin-token-bobgy --from-file=github_token=<PATH TO TOKEN>
```

* We are currently using the token **peribolos-kubeflow-org-admin** owned by Bobgy


## Validate config map

We use a config map to provide the python code used to validate the config

```
make create-config-map
```
