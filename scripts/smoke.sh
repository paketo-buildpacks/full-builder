#!/usr/bin/env bash

set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACKDIR="$(cd "${PROGDIR}/.." && pwd)"

# shellcheck source=SCRIPTDIR/.util/tools.sh
source "${PROGDIR}/.util/tools.sh"

# shellcheck source=SCRIPTDIR/.util/print.sh
source "${PROGDIR}/.util/print.sh"

function main() {
  local name

  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --help|-h)
        shift 1
        usage
        exit 0
        ;;

      --name|-n)
        name="${2}"
        shift 2
        ;;

      "")
        # skip if the argument is empty
        shift 1
        ;;

      *)
        util::print::error "unknown argument \"${1}\""
    esac
  done

  if [[ ! -d "${BUILDPACKDIR}/smoke" ]]; then
      util::print::warn "** WARNING  No Smoke tests **"
  fi

  if [[ -z "${name:-}" ]]; then
    name="testbuilder"
  fi

  tools::install
  builder::create "${name}"
  image::pull::lifecycle "${name}"
  tests::run "${name}"
}

function usage() {
  cat <<-USAGE
smoke.sh [OPTIONS]

Runs the smoke test suite.

OPTIONS
  --help        -h         prints the command usage
  --name <name> -n <name>  sets the name of the builder that is built for testing
USAGE
}

function tools::install() {
  util::tools::pack::install \
    --directory "${BUILDPACKDIR}/.bin"
}

function image::pull::lifecycle() {
  local name lifecycle_image
  name="${1}"

  lifecycle_image="index.docker.io/buildpacksio/lifecycle:$(
    docker inspect "${name}" \
      | jq -r '.[0].Config.Labels."io.buildpacks.builder.metadata"' \
      | jq -r '.lifecycle.version'
  )"

  util::print::title "Pulling lifecycle image..."
  docker pull "${lifecycle_image}"
}

function builder::create() {
  local name
  name="${1}"

  pack create-builder "${name}" --config "${BUILDPACKDIR}/builder.toml"
}

function tests::run() {
  local name
  name="${1}"

  util::print::title "Run Buildpack Runtime Integration Tests"

  testout=$(mktemp)
  pushd "${BUILDPACKDIR}" > /dev/null
    if GOMAXPROCS="${GOMAXPROCS:-4}" go test -count=1 -timeout 0 ./smoke/... -v -run Smoke --name "${name}" | tee "${testout}"; then
      util::tools::tests::checkfocus "${testout}"
      util::print::success "** GO Test Succeeded **"
    else
      util::print::error "** GO Test Failed **"
    fi
  popd > /dev/null
}

main "${@:-}"
