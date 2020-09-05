#!/bin/bash

# This script is used to generate license_notices.go.
# Call it using 'go generate'. 
set -e

INPACKAGE="$1"
OUTPACKAGE="$2"
OUTFILE="$3"

LICENSE_CSV=""
LEGAL_TEXT=""

function get_license_text() {
    TEMPDIR=$(mktemp -d)
    GOBIN="$(go env GOPATH)/bin"

    pushd "$TEMPDIR" > /dev/null 2>&1
    go mod init temp > /dev/null 2>&1
    go get github.com/google/go-licenses
    go install github.com/google/go-licenses
    popd > /dev/null 2>&1

    LICENSE_CSV=$("$GOBIN/go-licenses" csv "$INPACKAGE")

    "$GOBIN/go-licenses" save "$INPACKAGE" --save_path="$TEMPDIR/legal"

    pushd "$TEMPDIR/legal"  > /dev/null 2>&1

    N=$'\n'
    for modfile in $(find ./ -name 'LICENSE' | sort); do
        module=${modfile#"./"}
        module=${module%"/LICENSE"}
        LEGAL_TEXT="$LEGAL_TEXT$N================================================================================$N"
        LEGAL_TEXT="$LEGAL_TEXT${N}Go Module $module$N$N"
        LEGAL_TEXT="$LEGAL_TEXT$N$(cat "$modfile")"
        LEGAL_TEXT="$LEGAL_TEXT$N================================================================================$N"
    done

    popd  > /dev/null 2>&1
    rm -rf "$TEMPDIR"
}

function prepend() {
    while read line; do echo "${1}${line}"; done;
}

get_license_text
LEGAL_TEXT="${LEGAL_TEXT/\`/\`+\"\`\"+\`}"
LICENSE_CSV=$(echo "$LICENSE_CSV" | prepend '// ' | sort)


rm -f "$OUTFILE"
cat << EOF > "$OUTFILE"
package $OUTPACKAGE

// This file was generated automatically at $(date -u +"%Y-%m-%dT%H:%M:%SZ") using 'make_license_notices.sh'.
// Do not edit manually, as changes may be overwritten.

// StringLicenseNotices are legal notices required by the license.
//
$LICENSE_CSV
const StringLicenseNotices = \`${LEGAL_TEXT}\`


EOF

gofmt -w "$OUTFILE"