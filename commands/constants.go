package commands

// This file defines various resources used by all commands
// such as strings

// ErrorUnknownCommand is the return code when an unknown command is called
const ErrorUnknownCommand int = 1
const stringUnknownCommand string = "Unknown command. Must be one of 'ls', 'where'. "

// ErrorGeneralParsArgs is the return code when generic argument parsing fails
const ErrorGeneralParsArgs int = 2
const stringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. "
const stringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "

// ErrorSpecificParseArgs is the return code when specific argument parsing fails
const ErrorSpecificParseArgs int = 3
const stringLSTakesNoArguments string = "Too many arguments: 'ls' takes no arguments. "
const stringWhereNoFor string = "Wrong number of arguments: 'where' takes no 'for' argument. "
const stringWhereTakesOneArgument string = "Wrong number of arguments: 'where' takes exactly one arguments. "

// ErrorNoRoot is the return code when the root directory can not be found
const ErrorNoRoot int = 4
const stringUnableParseRootDirectory string = "Unable to find GGROOT directory. "

// ErrorInvalidRepo is the return code when an invalid repo name is passed
const ErrorInvalidRepo int = 5
const stringUnparsedRepoName string = "Unable to parse repository name. "
