#!/bin/bash
#
# A simple bash script to sync the GitHub org.
# See: https://github.com/kubernetes/test-infra/tree/master/prow/cmd/peribolos
#
# sync_org.sh <kubernets-test-infra-dir> <path-to-github-token>
set -ex
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

TEST_INFRA_DIR=$1
TOKEN_FILE=$2

pushd .
cd ${TEST_INFRA_DIR}

bazel run //prow/cmd/peribolos -- --config-path ${DIR}/kubeflow/org.yaml \
	--github-token-path ${TOKEN_FILE} \
	--required-admins=jlewi,abhi-g,google-admin,googlebot,richardsliu,vicaire