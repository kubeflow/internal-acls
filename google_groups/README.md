# Google Groups Sync For kubeflow.org

This directory contains

* A Go Binary to auto sync Google group membership based on YAML configurations checked into GitHub
* The YAML configurations in groups/

## How this Works

* We use the [Directory API](https://developers.google.com/admin-sdk/directory/v1/quickstart/go) 
  to programmatically update Google groups based on YAML files

* The script runs using the gsuite account `autobot@kubeflow.org`

  * Password and recovery codes are stored in secret manager
    
    * **secret**: [projects/kf-infra-gitops/secrets/autobot-at-kubeflow-oauth-admin-api](https://console.cloud.google.com/security/secret-manager/secret/autobot-at-kubeflow-oauth-admin-api?project=kf-infra-gitops)

      ```
      gcloud --project=kf-infra-gitops secrets versions access latest --secret="autobot-at-kubeflow-oauth-admin-api"
      ```

* An OAuthClient ID is stored in gcs: `gs://kf-infra-gitops_secrets/autobot-at-kubeflow_client_secret.json`

* An OAuth2 refresh token is stored in secret manager to allow the script to run without human intervention

* When `groups` runs it uses a GSA to read the OAuth2 refresh token from secret manager and then uses it
  to authenticate as `autobot@kubeflow.org` to the calendar API

* A side car runs git-sync](https://github.com/kubernetes/git-sync) in a side car to synchronize the repo to a volume mount

* The `groups` program polls the location of the YAML file and when it detects a change (based on a hash of file contents) it runs a synchronization

  * The groups program will also periodically force a sync even if no changes are detected to deal with any
    drift for this to occur


* The account `autobot@kubeflow.org` is a groups admin for kubeflow.org

## To Manually Synchronize the Groups

In order to run the sync you need the following

  * You must be a Groups admin for kubeflow.org

    * You must use an @kubeflow.org email an @google.com account won't work.

  * You must have access to an OAuth2 ClientSecret

    * The Makefile command below reads the OAuth2 credentials from GCS
    * You must have access to that GCS file or else change it to use a different file

```
make sync
```

## Refreshing the OAuth2 Refresh Token For `autobot@kubeflow.org`

The OAuth2 refresh token is stored inside a secret in secret manager

  * **secret** [projects/kf-infra-gitops/secrets/autobot-at-kubeflow-oauth-admin-api](https://console.cloud.google.com/security/secret-manager/secret/autobot-at-kubeflow-oauth-admin-api?project=kf-infra-gitops)

To regenerate the refresh token

1. Destroy the latest version of the secret stored in secret manager

1. Run a sync manually

   ```
   run --input=./groups/*.yaml \
      --credentials-file=gs://kf-infra-gitops_secrets/autobot-at-kubeflow_client_secret.json \
      --secret=kf-infra-gitops/autobot-at-kubeflow-oauth-admin-api
   ```

   * You will be directed through the OAuth2 Web Flow
   * Be sure to login using the account **autobot@kubeflow.org**

     * The password and recovery codes for **autobot@kubeflow.org** are stored in secret manager

       * **secret** [projects/kf-infra-gitops/secrets/autobot-kubeflow-org-password](https://console.cloud.google.com/security/secret-manager/secret/autobot-kubeflow-org-password?project=kf-infra-gitops)


## Importing Settings

The groups binary has an `import` command which can be used to update the YAML files with the latest configuration
in Google Groups

* This is useful for intializing the configs for any Groups not currently controlled via GitOps

```
groups import \
  --domain=kubeflow.org \
  --output=./google_groups/groups \
  --credentials-file=gs://kf-infra-gitops_secrets/autobot-at-kubeflow_client_secret.json
```

## References

* https://developers.google.com/admin-sdk/directory/v1/quickstart/go