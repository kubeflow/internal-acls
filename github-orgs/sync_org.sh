#!/bin/bash
#
# A simple bash script to sync the GitHub org.
# See: https://github.com/kubernetes/test-infra/tree/master/prow/cmd/peribolos
#
# sync_org.sh <kubernets-test-infra-dir> <path-to-github-token> <confirm>
set -ex
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

parseArgs() {
  # Parse all command line options
  while [[ $# -gt 0 ]]; do
    # Parameters should be of the form
    # --{name}=${value}
    echo parsing "$1"
    if [[ $1 =~ ^--(.*)=(.*)$ ]]; then
    	name=${BASH_REMATCH[1]}
    	value=${BASH_REMATCH[2]}
     	eval ${name}="${value}"
    elif [[ $1 =~ ^--(.*)$ ]]; then
		name=${BASH_REMATCH[1]}
		value=true
		eval ${name}="${value}"
    else
    	echo "Argument $1 did not match the pattern --{name}={value} or --{name}"
    fi
    shift
  done
}


action=$1
shift
parseArgs $*

usage() {
	echo "sync_org.sh <command> --test_infra_dir=<kubernetes-test-infra-dir> --token_file=<path-to-github-token> [--confirm] [--admins]"
}

if [ -z ${test_infra_dir} ]; then
	usage
	exit 1
fi

if [ -z ${token_file} ]; then
	usage
	exit 1
fi

if [ -z ${confirm} ]; then
	echo confirm not set defaulting to dryrun mode
    confirm=false
fi
pushd .
cd ${test_infra_dir}

if ${admins}; then
	FIX_ADMINS=--fix-org-members
fi

if [ "${action}" == "dump" ]; then
	#bazel run //prow/cmd/peribolos -- --help
	bazel run //prow/cmd/peribolos -- --dump=kubeflow \
		--github-token-path ${token_file} > --config-path ${DIR}/kubeflow/org.yaml
elif [ "${action}" == "sync" ]; then
	# Fix teams to create delete any teams
	bazel run //prow/cmd/peribolos -- \
		--fix-teams \
		--fix-team-members \
		--fix-org-members ${FIX_ADMINS} --config-path ${DIR}/kubeflow/org.yaml \
		--github-token-path ${token_file} \
		--required-admins=james-jwu \
		--required-admins=google-admin \
		--required-admins=googlebot \
		--confirm=${confirm}
    echo "Note: if dryrun=true you might get errors updating groups if the group doesn't exist"
else
  echo "command=${action} is not a valid command; valid commands are dump and sync"
  exit 1
fi
