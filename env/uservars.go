package env

import (
	"strings"

	"github.com/tkw1536/goprogram/meta"
	"golang.org/x/exp/slices"
)

// UserVariable is a variable that is exposed to the user.
// See GetUserVariables() for a details.
type UserVariable struct {
	Key         string
	Description string
	Get         func(Env, meta.Info) string
}

var allVariables = []UserVariable{
	{
		Key:         "GGROOT",
		Description: "root folder all ggman repositories will be cloned to",
		Get:         func(env Env, info meta.Info) string { return env.Root },
	},
	{
		Key:         "PWD",
		Description: "current working directory",
		Get: func(env Env, info meta.Info) string {
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
		Get: func(e Env, i meta.Info) string {
			return e.Git.GitPath()
		},
	},

	{
		Key:         "GGMAN_VERSION",
		Description: "the version of ggman this version is",
		Get: func(e Env, i meta.Info) string {
			return i.BuildVersion
		},
	},
	{
		Key:         "GGMAN_TIME",
		Description: "the time this version of ggman was built",
		Get: func(e Env, i meta.Info) string {
			return i.BuildTime.String()
		},
	},
}

func init() {
	slices.SortFunc(allVariables, func(a, b UserVariable) bool {
		return a.Key < b.Key
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
