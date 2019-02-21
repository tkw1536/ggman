package constants

// StringUsage is a usage string
const StringUsage string = `Usage: ggman [help|--help|-h] [for|--for|-f FILTER] COMMAND [ARGS...]

	help, --help, -h
		Print this usage dialog and exit
	
	for FILTER, --for FILTER, -f FILTER
		Filter the list of repositories to apply command to by FILTER. 
	
	COMMAND [ARGS...]
		Command to call. One of 'root', 'ls', 'lsr', 'where', 'canon',
		'comps', 'fetch', 'pull', 'fix', 'clone', 'link', 'license'.
		See individual commands for more help. 

ggman is licensed under the terms of the MIT License. Use 'ggman license'
to view licensing information. 
`

// StringNeedOneArgument is an error message when unable to parse arguments
const StringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information. "

// StringNeedTwoAfterFor is an error message when unable to parse arguments
const StringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "
