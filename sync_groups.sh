#!/bin/bash
#
# Sync groups using the gam CLI
set -ex

KF_GROUPS=( "ci-team" "release-team" "google-kubeflow-admins" "google-team" "kf-demo-owners" "devrel-team" "modeldb-team" "code-search-team" "kubeflow-examples-gcr-writers" )
for g in "${KF_GROUPS[@]}"
do
	gam update group ${g}@kubeflow.org sync member file ${g}.members.txt
done
