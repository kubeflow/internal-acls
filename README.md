# internal-acls

Repository used to maintain group ACLs used by the Kubeflow community.

For Google Groups in kubeflow.org, refer to `google_groups` subfolder.

Currently after modifying org.yaml, pytest should be manually run in
github-orgs directory to verify the change. This test will be run
automatically in a future change.

## Joining Kubeflow GitHub organization

**Please**
* read the [guidelines](https://www.kubeflow.org/docs/about/contributing/#joining-the-community) for joining the Kubeflow GitHub org before opening an issue
* **provide links to PRs or other contributions (2-3)**
* **list 2 existing members who are sponsoring your membership**
* **test your PR**
  Run

  ```
  cd github_orgs
  pytest test_org_yaml.py
  ```
  Include the output in the PR

**Additional Instructions**

After your PR is merged please wait at least 1 hour for changes to propogate. 

If after an hour you haven't recieved an invite to join the GitHub org please open an issue.

You can contact build cop in #buildcop in kubeflow.slack.com
