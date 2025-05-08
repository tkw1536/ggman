package env

//spellchecker:words reflect github pkglib reflectx
import (
	"os"
	"reflect"

	"github.com/tkw1536/pkglib/reflectx"
)

//spellchecker:words ggman GGROOT CANFILE GGNORM

// Variables represents the values of specific environment variables.
// Unset variables are represented as the empty string.
//
// This object is used to prevent code in ggman to access the environment directly, which is difficult to test.
// Instead access goes through this layer of indirection which can be mocked during testing.
//
// The env struct-tag indicates which environment variable the value corresponds to.
type Variables struct {
	// HOME is the path to the users' home directory
	// This is typically stored in the 'HOME' variable on unix-like systems.
	HOME string

	// other environment variables
	PATH    string `env:"PATH"`
	GGROOT  string `env:"GGROOT"`
	CANFILE string `env:"GGMAN_CANFILE"`
	GGNORM  string `env:"GGNORM"`
}

// variableEnvNames holds a mapping from reflect-field-indexes in Variables to os.GetEnv() names.
var variablesEnvNames map[int]string

// initialize variablesEnvNames.
func init() {
	tVariables := reflect.TypeFor[Variables]()
	variablesEnvNames = make(map[int]string, tVariables.NumField())

	for field, index := range reflectx.IterFields(tVariables) {
		env, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}
		variablesEnvNames[index] = env
	}
}

// ReadVariables reads Variables from the operating system.
func ReadVariables() (v Variables) {
	// assign the os.Getenv() values
	rV := reflect.ValueOf(&v).Elem()
	for i, env := range variablesEnvNames {
		rV.Field(i).SetString(os.Getenv(env))
	}

	// set the HOME variable
	// errors result in an empty home
	v.HOME, _ = os.UserHomeDir()

	return
}
