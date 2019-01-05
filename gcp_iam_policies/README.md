# IAM Policies

Store IAM policies for various GCP projects.

## To start a new policy file

```
gcloud projects --format=yaml get-iam-policy ${PROJECT} > ${PROJECT}.iam.policy.yaml
```
## To update the iam policy for a project

1. Update the YAML file for the project
1. Run the update command

   ```
    ./update_iam_policy.sh --project=${PROJECT}
   ```

