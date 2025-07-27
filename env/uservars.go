package env

//spellchecker:words slices strings
import (
	"slices"
	"strings"
)

//spellchecker:words GGROOT ggman workdir

// UserVariable is a variable that is exposed to the user.
// See GetUserVariables() for a details.
type UserVariable struct {
	Key         string
	Description string
	Get         func(Env) string
}

var allVariables = []UserVariable{
	{
		Key:         "GGROOT",
		Description: "root folder all ggman repositories will be cloned to",
		Get:         func(env Env) string { return env.Root },
	},
	{
		Key:         "PWD",
		Description: "current working directory",
		Get: func(env Env) string {
			workdir, err := env.Abs(".")
			if err != nil {
				return env.Workdir
			}
			return workdir
		},
	},

	{
		Key:         "GIT",
		Description: "path to the native git",
		Get: func(e Env) string {
			return e.Git.GitPath()
		},
	},
}

func init() {
	slices.SortFunc(allVariables, func(a, b UserVariable) int {
		return strings.Compare(a.Key, b.Key)
	})
}

// GetUserVariables returns all user variables in consistent order.
//
// This function is untested because the 'ggman env' command is tested.
func GetUserVariables() []UserVariable {
	return slices.Clone(allVariables)
}

// GetUserVariable gets a single user variable of the given name.
// The name is checked case-insensitive.
// When the variable does not exist, returns ok = False.
//
// This function is untested because the 'ggman env' command is tested.
func GetUserVariable(name string) (variable UserVariable, ok bool) {
	for _, v := range allVariables {
		if strings.EqualFold(v.Key, name) {
			return v, true
		}
	}
	return
}
