# internal-acls

Folder used to maintain google group ACLs used by the Kubeflow community.

The text files contain lists of folks that should be added
to various Google Groups in kubeflow.org that control
access to various shared resources.

The script `sync_groups.sh` can be used to sync groups
using the GAM CLI. Only administrators with appropriate
permissions will be able to sync groups.

To create a new group

```
GROUP_NAME=google-kubecon-europe
GROUP_EMAIL=${GROUP_NAME}@kubeflow.org 
DESCRIPTION="Some group description"

/home/jlewi/bin/gam/gam create group ${GROUP_EMAIL} who_can_join invited_can_join \
  name ${GROUP_NAME} description "${DESCRIPTION}" \
  allow_external_members true
```

## Setting Up GAM to run `sync_groups.sh`

1. You need a `xxxx@kubeflow.org` email address and you must be an admin in kubeflow.org gsuite. Ask an existing admin to invite you and give you permissions.
1. Follow instructions in https://github.com/jay0lee/GAM/wiki to set up the GAM CLI tool.
1. You can now run `sync_groups.sh`.
