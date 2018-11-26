#!/bin/bash
#
# Script to update the IAM policy.
#
# The script checks that the etag matches the actual etag.
# If it doesn't we know the iam policy is out of sync
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


usage() {
  echo "Usage: update_iam_policy --project=PROJECT"
}

getIamPolicy() {
  local POLICYFILE=$1
  gcloud projects --format=yaml get-iam-policy ${PROJECT} > ${POLICYFILE}
}

updateIamPolicy() {
  local PROJECT=$1
  local POLICY_FILE=$2

  cd ${DIR}
  pushd .

  if [ -z ${PROJECT} ]; then
    echo PROJECT is empty
    exit 1
  fi

  if [ -z ${POLICY_FILE} ]; then
    echo POLICY_FILE is empty
    exit 1
  fi

  EXPECTED_ETAG=$(yq -r ".etag" ${POLICY_FILE})

  # File to store policy fetched via gcloud
  NAME=$(basename ${POLICY_FILE})
  LIVE_POLICY=$(tempfile -p ${NAME}.live)

  getIamPolicy ${LIVE_POLICY}

  ACTUAL_ETAG=$(yq -r ".etag" ${LIVE_POLICY})

  echo "Current etag in ${POLICY_FILE} is ${EXPECTED_ETAG}"
  echo "Current IAM policy has etag ${ACTUAL_ETAG}"

  if [ ${ACTUAL_ETAG} != ${EXPECTED_ETAG} ]; then
    echo "Expected etag doesn't match actual etag."
    echo "Ensure ${POLICY_FILE} is in sync and then"
    echo "update ${POLICY_FILE} with current etag ${ACTUAL_ETAG}"
    exit 1
  fi

  # Update the policy
  gcloud projects set-iam-policy ${PROJECT} ${POLICY_FILE}

  getIamPolicy ${LIVE_POLICY}
  NEW_ETAG=$(yq -r ".etag" ${LIVE_POLICY})
  yq -y -r ".etag=\"${NEW_ETAG}\"" ${POLICY_FILE} > ${POLICY_FILE}.new
  mv ${POLICY_FILE}.new ${POLICY_FILE}

  popd 
}

main() {

  # List of required parameters
  names=(project)

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

  POLICY_FILE=${project}.iam.policy.yaml
  if [ ! -f ${POLICY_FILE} ]; then
    echo "Policy file ${POLICY_FILE} doesn't exist"
    exit 1
  fi

  updateIamPolicy ${project} ${POLICY_FILE}
}

parseArgs "$*"
main
