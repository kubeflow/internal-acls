#!/bin/bash
#
# Sync groups using the gam CLI
set -ex

gam update group ci-team@kubeflow.org sync member file ci-team.members.txt
gam update group release-team@kubeflow.org sync member file release-team.members.txt
gam update group google-kubeflow-admins@kubeflow.org sync member file google-kubeflow-admins.members.txt
gam update group google-team@kubeflow.org sync member file google-team.members.txt
gam update group kf-demo-owners@kubeflow.org sync member file kf-demo-owners.members.txt
gam update group devrel-team@kubeflow.org sync member file devrel-team.members.txt
gam update group modeldb-team@kubeflow.org sync member file modeldb-team.members.txt
gam update group code-search-team@kubeflow.org sync member file code-search-team.members.txt