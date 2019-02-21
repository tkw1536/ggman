package commands

// This file defines various resources used by all commands
// such as strings
const stringUsage string = `Usage: ggman [help|--help|-h] [for|--for|-f FILTER] COMMAND [ARGS...]

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

const stringRepoAlreadyExists string = "Unable to clone repository: Another git repository already exists in target location. "
const stringLinkDoesNotExist string = "Unable to link repository: Can not open source repository. "
const stringLinkAlreadyExists string = "Unable to link repository: Another directory already exists in target location. "
const stringLinkSamePath string = "Unable to link repository: Link source and target are identical. "

const stringUnknownCommand string = "Unknown command. Must be one of 'root', 'ls', 'lsr', 'where', 'canon', 'comps', 'fetch', 'pull', 'fix', 'clone', 'link', 'license'. "

const stringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. \nUse `ggman license` to view licensing information. "
const stringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "

const stringCmdNoFor string = "Wrong number of arguments: '%s' takes no 'for' argument. "
const stringRootTakesNoArguments string = "Wrong number of arguments: 'root' takes no arguments. "
const stringLSArguments string = "Unknown argument: 'ls' must be called with either '--exit-code' or no arguments. "
const stringLSRArguments string = "Unknown argument: 'lsr' must be called with either '--canonical' or no arguments. "
const stringWhereTakesOneArgument string = "Wrong number of arguments: 'where' takes exactly one arguments. "
const stringCanonTakesOneOrTwoArguments string = "Wrong number of arguments: 'canon' takes exactly one or exactly two arguments. "
const stringCompsTakesOneArgument string = "Wrong number of arguments: 'comps' takes exactly one argument. "
const stringFetchTakesNoArguments string = "Wrong number of arguments: 'fetch' takes no arguments. "
const stringPullTakesNoArguments string = "Wrong number of arguments: 'pull' takes no arguments. "
const stringFixArguments string = "Wrong number of arguments: Unknown argument: 'fix' must be called with either '--simulate' or no arguments."
const stringCloneTakesOneArgument string = "Wrong number of arguments: 'clone' takes exactly one argument. "
const stringLinkTakesOneArgument string = "Wrong number of arguments: 'link' takes exactly one argument. "

const stringUnableParseRootDirectory string = "Unable to find GGROOT directory. "
const stringInvalidCanfile string = "Invalid CANFILE found. "

const stringUnparsedRepoName string = "Unable to parse repository name. "
