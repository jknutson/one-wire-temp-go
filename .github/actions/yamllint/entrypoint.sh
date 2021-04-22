#!/bin/sh

set -e

# Copy the matcher to the host system; otherwise "add-matcher" can't find it.
cp /yamllint.json /github/workflow/yamllint.json
echo "::add-matcher::${RUNNER_TEMP}/_github_workflow/yamllint.json"

sh -c "yamllint $*"
