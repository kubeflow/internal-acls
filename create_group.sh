#!/bin/bash
#
# Helper script to create a group
set -xe

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

usage() {
  echo "Usage: checkout_repos --group=<name for group> --description=<description>"
}

main() {

  cd "${DIR}"

  # List of required parameters
  names=(group description)

  missingParam=false
  for i in ${names[@]}; do
    if [ -z ${!i} ]; then
      echo "--${i} not set"
      missingParam=true
    fi
  done

  if ${missingParam}; then
    usage
    exit 1
  fi

  GROUP_EMAIL=${group}@kubeflow.org  

  if [ ! -f ${group}.members.txt ]; then
  	touch ${group}.members.txt
  fi

  gam create group ${GROUP_EMAIL} who_can_join invited_can_join \
  	name ${group} description "${description}" \
  	allow_external_members true  

}

parseArgs $*
main
