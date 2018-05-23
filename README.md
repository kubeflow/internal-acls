# internal-acls

Repository used to maintain group ACLs used by the Kubeflow community.

The text files contain lists of folks that should be added
to various Google Groups in kubeflow.org that control
access to various shared resources.

The script `sync_groups.sh` can be used to sync groups
using the GAM CLI. Only administrators with appropriate
permissions will be able to sync groups.
