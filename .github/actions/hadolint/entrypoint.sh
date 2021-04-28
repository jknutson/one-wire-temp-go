#!/bin/sh

set -e

# Copy the matcher to the host system; otherwise "add-matcher" can't find it.
cp /hadolint.json /github/workflow/hadolint.json
echo "::add-matcher::${RUNNER_TEMP}/_github_workflow/hadolint.json"

sh -c "/bin/hadolint $*"
