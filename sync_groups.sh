#!/bin/bash
#
# Sync groups using the gam CLI
set -ex

gam update group devrel-team@kubeflow.org sync member file devrel-team.members.txt