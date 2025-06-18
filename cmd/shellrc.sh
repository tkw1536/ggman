# ggcd allows 'cd'-ing into a directory given a repository name
# e.g ggcd github.com/hello/world will cd into the directory where the
# 'github.com/hello/world' repository is checked out. 
#
# This also works with short names, e.g. "ggcd world" will cd into the first
# repository matching "world".
ggcd () {
	REPO="$(ggman -f "$1" ls --exit-code --one 2>&1)" && echo "$REPO" && cd "$REPO"
}

# ggshow is like ggcd, except that is displays the 
ggshow () {
	REPO="$(ggman -f "$1" ls --exit-code --one 2>&1)" && ggman -f "$REPO" show --no-patch 2>&1
}

# ggcode is like ggcd, except it opens an editor (here vscode) instead of cding. 
ggcode () {
	REPO="$(ggman -f "$1" ls --exit-code --one 2>&1)" && echo "$REPO" && code "$REPO"
}

# ggclone clones a repository if it does not yet exist, and then ccds into the correct directory.
ggclone () {
	DEST="$(ggman --no-fuzzy-filter -f "$1" ls --one)"
	if [ "$DEST" = "" ]; then
		ggman clone "$@" || return $?
		DEST="$(ggman where "$1")"
	fi
	echo "$DEST"
	cd "$DEST" 
}
