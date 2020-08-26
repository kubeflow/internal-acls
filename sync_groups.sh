#!/bin/bash
#
# Sync groups using the gam CLI
set -ex

KF_GROUPS=(
  "calendar-admins"
  "ci-team"
  "ci-viewer"
  "code-search-team"
  "community-meeting-hosts"
  "devrel-team"
  "devstats"
  "drive-content-managers"
  "example-maintainers"
  "feast-team"
  "github-team"
  "google-kubeflow-admins"
  "google-team"
  "insight-fellowship"
  "kf-demo-owners"
  "kf-kcc-admins"
  "kubeflow-examples-gcr-writers"
  "modeldb-team"
  "release-team"
)
for g in "${KF_GROUPS[@]}"; do
  gam update group ${g}@kubeflow.org sync member file ${g}.members.txt
done
