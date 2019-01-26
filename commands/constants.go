package commands

// This file defines various resources used by all commands
// such as strings

// ErrorCodeCustom is a return code that can be used by custom commands
const ErrorCodeCustom int = 1

// ErrorUnknownCommand is the return code when an unknown command is called
const ErrorUnknownCommand int = 2
const stringUnknownCommand string = "Unknown command. Must be one of 'root', 'ls', 'where'. "

// ErrorGeneralParsArgs is the return code when generic argument parsing fails
const ErrorGeneralParsArgs int = 3
const stringNeedOneArgument string = "Unable to parse arguments: Need at least one argument. "
const stringNeedTwoAfterFor string = "Unable to parse arguments: At least two arguments needed after 'for' keyword. "

// ErrorSpecificParseArgs is the return code when specific argument parsing fails
const ErrorSpecificParseArgs int = 4
const stringRootNoFor string = "Wrong number of arguments: 'root' takes no 'for' argument. "
const stringRootTakesNoArguments string = "Wrong number of arguments: 'root' takes no arguments. "
const stringLSArguments string = "Unknown argument: 'ls' must be called with either '--exit-code' or no arguments. "
const stringWhereNoFor string = "Wrong number of arguments: 'where' takes no 'for' argument. "
const stringWhereTakesOneArgument string = "Wrong number of arguments: 'where' takes exactly one arguments. "

// ErrorNoRoot is the return code when the root directory can not be found
const ErrorNoRoot int = 5
const stringUnableParseRootDirectory string = "Unable to find GGROOT directory. "

// ErrorInvalidRepo is the return code when an invalid repo name is passed
const ErrorInvalidRepo int = 6
const stringUnparsedRepoName string = "Unable to parse repository name. "
