#!/bin/sh

set -e

# Copy the matcher to the host system; otherwise "add-matcher" can't find it.
cp /goss.json /github/workflow/goss.json
echo "::add-matcher::${RUNNER_TEMP}/_github_workflow/goss.json"

sh -c "goss $*"
