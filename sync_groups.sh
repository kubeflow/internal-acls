#!/bin/bash
#
# Sync groups using the gam CLI
set -ex

gam update group kf-demo-owners@kubeflow.org sync member file kf-demo-owners.members.txt
gam update group devrel-team@kubeflow.org sync member file devrel-team.members.txt