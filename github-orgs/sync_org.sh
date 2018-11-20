#!/bin/bash
#
# A simple bash script to sync the GitHub org.
# See: https://github.com/kubernetes/test-infra/tree/master/prow/cmd/peribolos
#
# sync_org.sh <kubernets-test-infra-dir> <path-to-github-token> <confirm>
set -ex
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

TEST_INFRA_DIR=$1
TOKEN_FILE=$2
CONFIRM=$3
usage() {
	echo "sync_org.sh <kubernets-test-infra-dir> <path-to-github-token>"
}

if [ -z ${TEST_INFRA_DIR} ]; then
	usage
	exit 1	
fi

if [ -z ${TOKEN_FILE} ]; then
	usage
	exit 1	
fi

if [ -z ${CONFIRM} ]; then
	echo CONFIRM not set defaulting to dryrun mode
    CONFIRM=false
fi	
pushd .
cd ${TEST_INFRA_DIR}

bazel run //prow/cmd/peribolos -- --fix-org-members --config-path ${DIR}/kubeflow/org.yaml \
	--github-token-path ${TOKEN_FILE} \
	--required-admins=jlewi \
	--required-admins=abhi-g \
	--required-admins=google-admin \
	--required-admins=googlebot \
	--required-admins=richardsliu \
	--required-admins=vicaire \
	--confirm=${CONFIRM}