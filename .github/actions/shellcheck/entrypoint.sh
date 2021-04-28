#!/bin/sh

set -e

# Copy the matcher to the host system; otherwise "add-matcher" can't find it.
cp /shellcheck.json /github/workflow/shellcheck.json
echo "::add-matcher::${RUNNER_TEMP}/_github_workflow/shellcheck.json"

echo $HOME

sh -c "shellcheck --format gcc $*"
