package program

// SubCommand represents a command that can be run with the program
type SubCommand func(args *GGArgs) (retval int, err string)
