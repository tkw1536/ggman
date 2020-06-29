package constants

// StringUsage is a usage string
const StringUsage string = `ggman version %s
(built %s)

Usage:
    ggman [help|--help|-h] [version|--version|-v] [for|--for|-f FILTER] COMMAND [ARGS...]

    help, --help, -h
        Print this usage dialog and exit
    
    version|--version|-v
		Print version message and exit. 
	
    for FILTER, --for FILTER, -f FILTER
        Filter the list of repositories to apply command to by FILTER. 
	
    COMMAND [ARGS...]
	    Command to call. One of %s. See individual commands for more help. 

ggman is licensed under the terms of the MIT License. Use 'ggman license'
to view licensing information. 
`

// StringVersion is a version string
const StringVersion string = `ggman version %s, built %s`

// StringNeedOneArgument is an error message when unable to parse arguments
const StringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information. "

// StringNeedTwoAfterFor is an error message when unable to parse arguments
const StringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "
