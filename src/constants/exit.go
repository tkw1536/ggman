package constants

// ErrorCodeCustom is a subcommand-specific error
const ErrorCodeCustom int = 1

// ErrorUnknownCommand indicates an unknown command
const ErrorUnknownCommand int = 2

// ErrorGeneralParsArgs when generic argument parsing fails
const ErrorGeneralParsArgs int = 3

// ErrorSpecificParseArgs is the return code when specific argument parsing fails
const ErrorSpecificParseArgs int = 4

// ErrorMissingConfig is the return code when the configuration is missing or invalid
const ErrorMissingConfig int = 5

// ErrorInvalidRepo is the return code when an invalid repo name is passed
const ErrorInvalidRepo int = 6
