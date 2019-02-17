package commands

// This file defines various resources used by all commands
// such as strings

// ErrorCodeCustom is a return code that can be used by custom commands
const ErrorCodeCustom int = 1
const stringRepoAlreadyExists string = "Unable to clone repository: Another git repository already exists in target location. "
const stringLinkDoesNotExist string = "Unable to link repository: Can not open source repository. "
const stringLinkAlreadyExists string = "Unable to link repository: Another directory already exists in target location. "
const stringLinkSamePath string = "Unable to link repository: Link source and target are identical. "

// ErrorUnknownCommand is the return code when an unknown command is called
const ErrorUnknownCommand int = 2
const stringUnknownCommand string = "Unknown command. Must be one of 'root', 'ls', 'lsr', 'where', 'canon', 'comps', 'fetch', 'pull', 'fix', 'clone', 'link', 'license'. "

// ErrorGeneralParsArgs is the return code when generic argument parsing fails
const ErrorGeneralParsArgs int = 3
const stringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. \nUse `ggman license` to view licensing information. "
const stringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "

// ErrorSpecificParseArgs is the return code when specific argument parsing fails
const ErrorSpecificParseArgs int = 4
const stringRootNoFor string = "Wrong number of arguments: 'root' takes no 'for' argument. "
const stringRootTakesNoArguments string = "Wrong number of arguments: 'root' takes no arguments. "
const stringLSArguments string = "Unknown argument: 'ls' must be called with either '--exit-code' or no arguments. "
const stringLSRArguments string = "Unknown argument: 'lsr' must be called with either '--canonical' or no arguments. "
const stringWhereNoFor string = "Wrong number of arguments: 'where' takes no 'for' argument. "
const stringWhereTakesOneArgument string = "Wrong number of arguments: 'where' takes exactly one arguments. "
const stringCanonNoFor string = "Wrong number of arguments: 'canon' takes no 'for' argument. "
const stringCanonTakesOneOrTwoArguments string = "Wrong number of arguments: 'canon' takes exactly one or exactly two arguments. "
const stringCompsNoFor string = "Wrong number of arguments: 'comps' takes no 'for' argument. "
const stringCompsTakesOneArgument string = "Wrong number of arguments: 'comps' takes exactly one argument. "
const stringFetchTakesNoArguments string = "Wrong number of arguments: 'fetch' takes no arguments. "
const stringPullTakesNoArguments string = "Wrong number of arguments: 'pull' takes no arguments. "
const stringFixArguments string = "Wrong number of arguments: Unknown argument: 'fix' must be called with either '--simulate' or no arguments."
const stringCloneNoFor string = "Wrong number of arguments: 'clone' takes no 'for' argument. "
const stringCloneTakesOneArgument string = "Wrong number of arguments: 'clone' takes exactly one argument. "
const stringLinkNoFor string = "Wrong number of arguments: 'link' takes no 'for' argument. "
const stringLinkTakesOneArgument string = "Wrong number of arguments: 'link' takes exactly one argument. "

// ErrorMissingConfig is the return code when the configuration is missing or invalid
const ErrorMissingConfig int = 5
const stringUnableParseRootDirectory string = "Unable to find GGROOT directory. "
const stringInvalidCanfile string = "Invalid CANFILE found. "

// ErrorInvalidRepo is the return code when an invalid repo name is passed
const ErrorInvalidRepo int = 6
const stringUnparsedRepoName string = "Unable to parse repository name. "
