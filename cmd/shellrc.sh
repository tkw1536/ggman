# ggcd allows 'cd'-ing into a directory given a repository name
# e.g ggcd github.com/hello/world will cd into the directory where the
# 'github.com/hello/world' repository is checked out. 
#
# This also works with short names, e.g. "ggcd world" will cd into the first
# repository matching "world".
ggcd () {
	ggman -f "$1" ls --exit-code --one && cd "$(ggman -f "$1" ls --exit-code --one 2>&1)"
}

# ggcode is like ggcd, except it opens an editor (here vscode) instead of cding. 
ggcode () {
	ggman -f "$1" ls --exit-code --one && code "$(ggman -f "$1" ls --exit-code --one 2>&1)"
}

# ggclone clones a repository if it does not yet exist, and then ccds into the correct directory.
ggclone () {
	DEST=$(ggman --no-fuzzy-filter -f "$1" ls --one)
	if [ -n "$DEST" ] && ! [ -d "$DEST" ]; then
		ggman clone "$@" || return 1
	fi
	echo "$DEST"
	cd "$DEST" 
}