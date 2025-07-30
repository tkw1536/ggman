# shellcheck shell=sh
# spellchecker:words ggcd ggshow ggcode ggclone ggman ggdo ggcursor

# This file contains useful shell aliases by ggman. 
# The comment above says to check against sh, but this is only for IDE support. 
# In reality, it is checked against sh, bash, dash, ksh and busybox.

# gg is a function that takes a for-style ggman pattern and a command to execute.
# It first resolves the pattern, and then executes the command on the first matching repository. 
# The command is passed as the last argument, following any additional arguments.
#
# Usage: gg <pattern> <command> [args...]
ggdo () {
	if [ $# -lt 2 ]; then
		echo "Usage: gg <pattern> <command> [args...]" >&2
		return 1
	fi

	REPO="$(ggman --for "$1" ls --exit-code --one 2>&1)" || return $?
	shift
	
	echo "$REPO"
	"$@" "$REPO" || return $?
}

# To avoid conflicts with other aliases, don't override gg if it's already set. 
if ! command -v gg >/dev/null 2>&1; then
	alias gg=ggdo
fi

# ggcd allows 'cd'-ing into a directory given a repository name
# e.g ggcd github.com/hello/world will cd into the directory where the
# 'github.com/hello/world' repository is checked out. 
#
# This also works with short names, e.g. "ggcd world" will cd into the first
# repository matching "world".
ggcd () {
	ggdo "$1" cd
}

# ggcode is like 'gg $1 code'
ggcode () {
	PATTERN="$1" 
	shift
	ggdo "$PATTERN" code "$@" || return $?
}

# ggcursor is like 'gg $1 cursor'
ggcursor () {
	PATTERN="$1" 
	shift
	ggdo "$PATTERN" cursor "$@" || return $?
}

# ggshow is like ggcd, except that it runs ggman show on the output
ggshow () {
	REPO="$(ggman --for "$1" ls --exit-code --one 2>&1)" && ggman --for "$REPO" show --no-patch 2>&1 || return $?
}

# ggclone clones a repository if it does not yet exist, and then cds into the correct directory.
ggclone () {
	DEST="$(ggman --no-fuzzy-filter --for "$1" ls --one)"
	if [ "$DEST" = "" ]; then
		ggman clone "$@" || return $?
		DEST="$(ggman where "$1")"
	fi
	echo "$DEST"
	cd "$DEST" || return $?
}
