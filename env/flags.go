package env

// Flags represent flags for the main ggman program
type Flags struct {
	Filters       []string `short:"f" long:"for" value-name:"filter" description:"Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url"`
	NoFuzzyFilter bool     `short:"n" long:"no-fuzzy-filter" description:"Disable fuzzy matching for filters"`

	Here bool     `short:"H" long:"here" description:"Filter the list of repositories to apply COMMAND to only contain repository in the current directory or subtree. Alias for '-p .'"`
	Path []string `short:"P" long:"path" description:"Filter the list of repositories to apply COMMAND to only contain repositories in or under the specified path. May be used multiple times"`

	Dirty bool `short:"d" long:"dirty" description:"List only repositories with uncommited changes"`
	Clean bool `short:"c" long:"clean" description:"List only repositories without uncommited changes"`

	Synced   bool `short:"s" long:"synced" description:"List only repositories which are up-to-date with remote"`
	UnSynced bool `short:"u" long:"unsynced" description:"List only repositories not up-to-date with remote"`

	Tarnished bool `short:"t" long:"tarnished" description:"List only repositories which are dirty or unsynced"`
	Pristine  bool `short:"p" long:"pristine" description:"List only repositories which are clean and synced"`
}
