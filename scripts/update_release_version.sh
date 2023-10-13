#!/usr/bin/env bash
set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "${REPO_ROOT}"

function get_bumped_version() {
  local version="$1"
  local bump_component="$2"

  semver-cli inc "${bump_component}" "${version}"
}

function check_required_argument() {
  local required_arg="$1"
  local failure_msg="$2"

  if [[ -z "${required_arg}" ]]; then
    print_usage_and_exit "${failure_msg}"
  fi
}

function print_usage_and_exit() {
  echo "Failure: $1"
  echo "Usage: $0 [flags] [options]"
  echo -e "\t-v version to bump (required)"
  echo -e "\t-s semver component to bump (required, ex: major, minor, patch)"
  exit 1
}

function main() {
  local VERSION=''
  local BUMP_COMPONENT=''
  local NEXT_RELEASE_VERSION=''
  local FUTURE_RELEASE_VERSION=''

  while getopts ":v:s:" opt; do
    case $opt in
      v) VERSION="$OPTARG";;
      s) BUMP_COMPONENT="$OPTARG";;
      \?) print_usage_and_exit "Invalid option: -$OPTARG";;
    esac
  done

  check_required_argument "${VERSION}" "-v <VERSION> is required"
  check_required_argument "${BUMP_COMPONENT}" "-s <BUMP_COMPONENT> is required"

  echo "Current release version: ${VERSION}"

  NEXT_RELEASE_VERSION="$(get_bumped_version "${VERSION}" "${BUMP_COMPONENT}")"
  echo "Next release version: ${NEXT_RELEASE_VERSION}"

  # update the version in the image tag
  sed -i.bak "s/wavefront-hpa-adapter:${VERSION}/wavefront-hpa-adapter:${NEXT_RELEASE_VERSION}/g" "${REPO_ROOT}"/deploy/manifests/05-custom-metrics-apiserver-deployment.yaml
  rm -f "${REPO_ROOT}"/deploy/manifests/05-custom-metrics-apiserver-deployment.yaml.bak

  echo "${NEXT_RELEASE_VERSION}" >"${REPO_ROOT}"/release/VERSION

  FUTURE_RELEASE_VERSION="$(get_bumped_version "${NEXT_RELEASE_VERSION}" "patch")"
  echo "Future release version: ${FUTURE_RELEASE_VERSION}"

  # update the version in the Makefile
  sed -i.bak "s/VERSION?=[0-9.]*/VERSION?=${FUTURE_RELEASE_VERSION}/g" "${REPO_ROOT}"/Makefile
  rm -f "${REPO_ROOT}"/Makefile.bak
}

main "$@"
