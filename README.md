# internal-acls

Repository used to maintain group ACLs used by the Kubeflow community.

For Google Groups in kubeflow.org, refer to `google_groups` subfolder.

Run `pytest` after modifying `org.yaml`. Currently, this must be run
manually. This test run automatically in a future change.

## Joining the Kubeflow GitHub organization

To join the Kubeflow GitHub organization, complete the following steps:

* Read the [guidelines for joining the Kubeflow GitHub org](https://www.kubeflow.org/docs/about/contributing/#joining-the-community) before opening an issue.
* **Provide 2-3 links to PRs or other contributions.**
* **List 2 existing members who are sponsoring your membership.**
* **Test your PR by running the following:**

  ```bash
  cd github-orgs
  pytest test_org_yaml.py
  ```

  Include the output in the PR.

### Additional Instructions

After your PR is merged, wait at least 1 hour for changes to propagate.
If, after an hour, you haven't recieved an invite to join the Kubeflow
GitHub org, open an issue.
