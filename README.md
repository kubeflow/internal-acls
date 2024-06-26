# internal-acls

Repository used to maintain group ACLs used by the Kubeflow community.

For Google Groups in kubeflow.org, refer to `google_groups` subfolder.

To check if everything is fine after modifying `org.yaml`, you can run
`pytest` locally from `github-orgs` directory.

## Joining Kubeflow GitHub organization

To join the Kubeflow GitHub organization, complete the following steps:

* Read the [guidelines for joining the Kubeflow GitHub org](https://www.kubeflow.org/docs/about/contributing/#joining-the-community)
  before opening an issue.
* **Provide 2-3 links to PRs or other contributions.**
* **List 2 existing members who are sponsoring your membership.**
* **Test your PR by running the following:**

  ```
  cd github-orgs
  pytest test_org_yaml.py
  ```

**Additional Instructions**

After your PR is merged, wait at least 1 hour for changes to propagate.
If, after an hour, you haven't recieved an invite to join the Kubeflow
GitHub org, open an issue.

You can contact build cop in the #buildcop channel of [kubeflow.slack.com](https://kubeflow.slack.com).
