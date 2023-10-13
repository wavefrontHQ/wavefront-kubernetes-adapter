#!/usr/bin/env bash
set -eou pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"

function create_pull_request() {
  local version="$1"
  local token="$2"
  local branch_name="$3"

  curl -fsSL -X 'POST' \
    -H 'Accept: application/vnd.github+json' \
    -H "Authorization: Bearer ${token}" \
    -H 'X-GitHub-Api-Version: 2022-11-28' \
    -d "{\"head\":\"${branch_name}\",\"base\":\"master\",\"title\":\"Bump version to ${version}\"}" \
    https://api.github.com/repos/wavefrontHQ/wavefront-kubernetes-adapter/pulls
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
  echo -e "\t-v bump version (required)"
  echo -e "\t-t github token (required)"
  exit 1
}

function main() {
  cd "${REPO_ROOT}"

  local VERSION=''
  local GITHUB_TOKEN=''

  while getopts ":v:t:" opt; do
    case $opt in
      v) VERSION="$OPTARG";;
      t) GITHUB_TOKEN="$OPTARG";;
      \?) print_usage_and_exit "Invalid option: -$OPTARG";;
    esac
  done

  check_required_argument "${VERSION}" "-v <VERSION> is required"
  check_required_argument "${GITHUB_TOKEN}" "-t <GITHUB_TOKEN> is required"

  local GIT_BUMP_BRANCH_NAME="bump-version-${VERSION}"
  git branch -D "$GIT_BUMP_BRANCH_NAME" &>/dev/null || true
  git checkout -b "$GIT_BUMP_BRANCH_NAME"

  make update-version NEW_VERSION="${VERSION}"

  create_pull_request "${VERSION}" "${GITHUB_TOKEN}" "${GIT_BUMP_BRANCH_NAME}"
}

main "$@"
