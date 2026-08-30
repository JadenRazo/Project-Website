#!/usr/bin/env bash
# Reports which Go packages do not build, and fails if that set grows.
#
# `go build ./...` does not pass in backend/ and has not for a long time: some
# packages under internal/messaging/ have duplicate declarations and references
# to a type that no longer exists. None of them is reachable from any of the
# eleven binaries in cmd/, so everything that actually ships builds and runs
# fine - which is exactly why the breakage survived unnoticed.
#
# Two bad options and one good one:
#
#   - Gate CI on `go build ./...`. Red from the first commit, permanently, so
#     the signal is worthless and gets ignored - which is how this repo got here.
#   - Do not check it at all. That is the status quo: the set can grow forever
#     and nobody finds out.
#   - Ratchet. Record today's broken set, fail if anything NEW breaks, and fail
#     if the file lists something that now builds, so the baseline can only
#     shrink. Debt is visible, bounded, and monotonically decreasing.
#
# This is a pure function of the commit: no network, no credentials, no clock.
#
# Usage:
#   scripts/ci/build-ratchet.sh            # check against the baseline
#   scripts/ci/build-ratchet.sh --write    # regenerate the baseline (do this
#                                          # only when deliberately recording a
#                                          # reduction, and say so in the PR)

set -uo pipefail

cd "$(dirname "$0")/../../backend" || exit 2

BASELINE="../scripts/ci/build-baseline.txt"
WRITE=0
[ "${1:-}" = "--write" ] && WRITE=1

# -e keeps listing after a package fails to load, instead of aborting.
#
# Packages holding only _test.go files are excluded: `go build` reports
# "no non-test Go files in ..." and exits non-zero for them, which is not a
# build failure. internal/performance is one, and the first run of this script
# flagged it as broken when it was fine. Their compilation is covered by the
# backend-test job, which builds test files.
mapfile -t packages < <(go list -e -f '{{if .GoFiles}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | sed '/^$/d' | sort)
mapfile -t testonly < <(go list -e -f '{{if and (not .GoFiles) (or .TestGoFiles .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | sed '/^$/d' | sort)

if [ "${#packages[@]}" -eq 0 ]; then
  echo "error: go list returned no buildable packages - is this the right directory?" >&2
  exit 2
fi

broken=()
for pkg in "${packages[@]}"; do
  if ! go build -o /dev/null "$pkg" >/dev/null 2>&1; then
    broken+=("$pkg")
  fi
done

printf '%s\n' "${broken[@]}" | sed '/^$/d' | sort > /tmp/broken.now

echo "packages: ${#packages[@]}   not building: $(wc -l < /tmp/broken.now)   test-only (skipped): ${#testonly[@]}"

if [ "$WRITE" = "1" ]; then
  cp /tmp/broken.now "$BASELINE"
  echo "wrote $(wc -l < "$BASELINE") package(s) to $BASELINE"
  exit 0
fi

if [ ! -f "$BASELINE" ]; then
  echo "error: $BASELINE is missing. Generate it with: scripts/ci/build-ratchet.sh --write" >&2
  exit 2
fi

sed -e 's/#.*//' -e '/^[[:space:]]*$/d' "$BASELINE" | sort > /tmp/broken.baseline

new=$(comm -23 /tmp/broken.now /tmp/broken.baseline)
fixed=$(comm -13 /tmp/broken.now /tmp/broken.baseline)

rc=0

if [ -n "$new" ]; then
  echo
  echo "::error::$(echo "$new" | wc -l) package(s) newly fail to build:"
  echo "$new" | sed 's/^/  /'
  echo
  echo "These are not pre-existing debt. Fix them, or if a package was"
  echo "knowingly left broken, add it to scripts/ci/build-baseline.txt and say why."
  rc=1
fi

if [ -n "$fixed" ]; then
  echo
  echo "::error::$(echo "$fixed" | wc -l) package(s) in the baseline now build:"
  echo "$fixed" | sed 's/^/  /'
  echo
  echo "Good news, but the baseline has to come down with it or it stops"
  echo "measuring anything. Run: scripts/ci/build-ratchet.sh --write"
  rc=1
fi

[ "$rc" = "0" ] && echo "ok: broken set unchanged ($(wc -l < /tmp/broken.baseline) package(s), none reachable from cmd/)"

if [ -n "${GITHUB_STEP_SUMMARY:-}" ]; then
  {
    echo "### Go package build health"
    echo
    echo "\`${#packages[@]}\` packages, \`$(wc -l < /tmp/broken.now)\` not building."
    echo
    if [ -s /tmp/broken.now ]; then
      echo "<details><summary>Packages that do not build</summary>"
      echo
      echo '```'
      cat /tmp/broken.now
      echo '```'
      echo
      echo "</details>"
    fi
  } >> "$GITHUB_STEP_SUMMARY"
fi

exit $rc
