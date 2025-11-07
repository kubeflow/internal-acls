---
name: Kubeflow GitHub Org Request
about: Request to join or modify Kubeflow GitHub membership

---

# Membership Request

## Instructions
Please read the [guidelines](https://www.kubeflow.org/docs/about/contributing/#joining-the-community) for joining the Kubeflow GitHub org before opening an issue.

### Provide links to your PRs or other contributions (2-3):

### List 2 existing members who are sponsoring your membership:

### Create a new membership pull request (PR):
- Fork the [kubeflow/internals-acls](https://github.com/kubeflow/internal-acls/) repo and clone it locally.
- Modify [github-orgs/kubeflow/org.yaml](github-orgs/kubeflow/org.yaml) to include your GitHub username in the `org.kubeflow.members` list.
- Test your code changes with:
    ```bash
    pytest github-orgs/test_org_yaml.py
    ```
    Confirm that the test run passed, and make sure to copy the test output so you can include it in your PR description (on GitHub).
- Push your changes to your `kubeflow/internals-acls` repository fork.
- Open a PR from the [kubeflow/internals-acls](https://github.com/kubeflow/internal-acls/) repo.
- Your PR will be reviewed by someone from the Kubeflow team. Work with your reviewers to address any outstanding issues.

## Additional Instructions
- After your PR is merged please wait at least 1 hour for changes to propagate.
- You will receive an email invite (to your GitHub associated email address) to join Kubeflow on GitHub. Follow the instructions on the email to accept your invitation.
- If after an hour you haven't received an invite to join the GitHub org (or your invite has expired) please open an issue with an [owner](https://github.com/kubeflow/internal-acls/blob/master/OWNERS) tagged to request follow-up.
