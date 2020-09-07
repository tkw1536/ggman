#!/bin/bash

# make_license_notices.sh
#
# Generate some go code for license notices. 
#
# (c) Tom Wiesing 2020 available under the terms of the MIT License. 
#
# It should be used by means of a go-generate command, such as:
# //go:generate bash make_license_notices.sh INPACKAGE OUTPACKAGE OUTFILE DECLARATION

set -e

###########################
# Read parameters         #
###########################

INPACKAGE=$1
OUTPACKAGE="$2"
OUTFILE="$3"
DECLARATION="${4:-LicenseNotices}"

###########################
# Environment and globals #
###########################

TEMPDIR=$(mktemp -d)
GOBIN="$(go env GOPATH)/bin"

N=$'\n'

###########################
# Install Dependencies    #
###########################

pushd "$TEMPDIR" > /dev/null 2>&1
go mod init temp > /dev/null 2>&1
go get github.com/google/go-licenses
go install github.com/google/go-licenses
popd > /dev/null 2>&1

###########################
# Bootstrap comment       #
###########################

LEGAL_COMMENT=$(cat <<EOF
// ${DECLARATION} contains legal and license information of external software included in this program. 
// These notices consist of a list of dependencies along with their license information. 
// This string is intended to be displayed to the enduser on demand. 
//
// Even though the value of this variable is fixed at compile time it is omitted from this documentation. 
// Instead the list of go modules, along with their licenses, is listed below. 
//
EOF
)
LEGAL_COMMENT="$LEGAL_COMMENT${N}"

############################
# Function to record info  #
############################

declare -A pkg_license
declare -A pkg_url

function record_license_info() {
    PKG=$1
    LICENSE=$2
    URL=$3

    pkg_license[$PKG]="$LICENSE"
    pkg_url[$PKG]="$URL"
}


###########################
# Function to write text  #
###########################

function write_overview_text() {
    PKG=$1
    LICENSE=$2
    URL=$3

    LEGAL_TEXT="${LEGAL_TEXT}${N} - $PKG ($LICENSE"
    if [ "${URL}" != "Unknown" ]; then
        LEGAL_TEXT="${LEGAL_TEXT}; see ${URL}"
    fi
    LEGAL_TEXT="${LEGAL_TEXT})"
}

function write_const_text() {
    PKG=$1
    LICENSE=$2
    URL=$3
    LICENSETEXT=$4

    LEGAL_TEXT="${LEGAL_TEXT}${N}================================================================================${N}"
    LEGAL_TEXT="${LEGAL_TEXT}${N}Module $PKG${N}"
    LEGAL_TEXT="${LEGAL_TEXT}${N}Licensed under the Terms of the $LICENSE License. ${N}"
    if [ "${URL}" != "Unknown" ]; then
        LEGAL_TEXT="${LEGAL_TEXT}See also $URL. ${N}"
    fi
    LEGAL_TEXT="${LEGAL_TEXT}${N}${N}${LICENSETEXT}${N}"
    LEGAL_TEXT="${LEGAL_TEXT}${N}================================================================================${N}"
}

function write_comment_text() {
    PKG=$1
    LICENSE=$2
    URL=$3
    LICENSETEXT=$4

    LEGAL_COMMENT="${LEGAL_COMMENT}// ${N}"
    LEGAL_COMMENT="${LEGAL_COMMENT}// Module ${PKG//\// }${N}"
    LEGAL_COMMENT="${LEGAL_COMMENT}// ${N}"
    LEGAL_COMMENT="${LEGAL_COMMENT}// Module $PKG is licensed under the Terms of the $LICENSE License. ${N}"
    if [ "${URL}" != "Unknown" ]; then
        LEGAL_COMMENT="${LEGAL_COMMENT}// See also $URL. ${N}"
    fi

    LEGAL_COMMENT="${LEGAL_COMMENT}// ${N}"
    while IFS= read -r licenseline; do
        LEGAL_COMMENT="${LEGAL_COMMENT}//  ${licenseline}${N}"
    done <<< "$LICENSETEXT"

}
###########################
# Read License CSV        #
###########################

LEGAL_TEXT="The following go packages are imported: ${N}"
for line in $("$GOBIN/go-licenses" csv "$INPACKAGE/..." | sort); do
    # parse the csv, skip $INPACKAGE
    IFS=',' read -ra info <<< "$line"
    if [ "${info[0]}" == "$INPACKAGE" ]; then
        continue
    fi

    write_overview_text "${info[0]}" "${info[2]}" "${info[1]}"
    record_license_info "${info[0]}" "${info[2]}" "${info[1]}"
done
LEGAL_TEXT="${LEGAL_TEXT}${N}"

###########################
# Generate License Text   #
###########################

"$GOBIN/go-licenses" save "$INPACKAGE/..." --save_path="$TEMPDIR/legal"

pushd "$TEMPDIR/legal"  > /dev/null 2>&1
for licensefile in $(find ./ -name 'LICENSE' | sort); do
    pkg=${licensefile#"./"}
    pkg=$(dirname "$pkg")
    if [ "$pkg" == "$INPACKAGE" ]; then
        continue
    fi

    license=${pkg_license[$pkg]}
    licensetext=$(cat $licensefile)
    url=${pkg_url[$pkg]}

    write_const_text "$pkg" "$license" "$url" "$licensetext"
    write_comment_text "$pkg" "$license" "$url" "$licensetext"
done
popd  > /dev/null 2>&1

LEGAL_COMMENT=$(cat <<EOF
${LEGAL_COMMENT}// 
// Generation
//
// The variable (and this documentation) have been automatically generated using
//
//  make_license_notices.sh $*
//
// This variable was last updated at $(date -u +"%Y-%m-%dT%H:%M:%SZ"). 
EOF
)

###########################
# Cleanup                 #
###########################

rm -rf "$TEMPDIR"


#########################
# Write outfile         #
#########################

# escape `s
LEGAL_TEXT="${LEGAL_TEXT//\`/\`+\"\`\"+\`}"

rm -f "$OUTFILE"
cat << EOF > "$OUTFILE"
package $OUTPACKAGE

// ===========================================================================================================
// This file was generated automatically at $(date -u +"%Y-%m-%dT%H:%M:%SZ") using make_license_notices.sh $*.
// Do not edit manually, as changes may be overwritten.
// ===========================================================================================================


${LEGAL_COMMENT}
var ${DECLARATION} string
func init() {
    ${DECLARATION} = \`${LEGAL_TEXT}\`
}
EOF

gofmt -w "$OUTFILE"