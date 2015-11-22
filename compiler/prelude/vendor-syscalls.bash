#!/bin/bash

# FIXME
JS="$(cat $HOME/plasma/browsix/lib-dist/lib/syscall-api/syscall-api.js)"

cat >syscalls.go <<EOF
// AUTO GENERATED - DO NOT EDIT

package prelude

const syscalls = \`
$JS
\`
EOF
