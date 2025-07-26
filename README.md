# ggman

<!-- spellchecker:words ggman ggroot shellrc ggcd ggclone ggshow ggcode ggnorm wrld cspec gopath godoc goprogram unsynced jessevdk struct POSIX pflag localgodoc CANSPEC CANFILE worktree reclone testutil subpackage -->

![CI Status](https://github.com/tkw1536/ggman/workflows/CI/badge.svg)

A golang tool that can manage all your git repositories. 

## What Is ggman?

When you only have a couple of git repositories that you work on it is perfectly feasible to manage them by using `git clone`, `git pull` and friends. 
However once the number of repositories grows beyond a small number this can become tedious:

- It is hard to find which folder a repository has been cloned to
- Getting an overview of what is cloned and what is not is hard
- It's not easily possible to perform actions on more than one repo at once, e.g. `git pull`

This is the problem `ggman` is designed to solve. 
It allows one to:

- Maintain and expand a local directory structure of multiple repositories
- Run actions (such as `git clone`, `git pull`) on groups of repositories

## Why ggman?

While similar tools exist these commonly have a lot of downsides:

- they enforce a flat directory structure;
- they are limited to one repository provider (such as GitHub or GitLab); or
- they are only available from within an IDE or GUI.

ggman considers these as major downsides. 
The goals and principles of ggman are:

- to be command-line first;
- to be simple to install, configure and use;
- to encourage an obvious hierarchical directory structure, but remain fully functional with any directory structure;
- to remain free of provider-specific code; and
- to not store any repository-specific data outside of the repositories themselves (enabling the user to switch back to only git at any point).

## Setting up and using 'ggman'

Setting up ggman consists of two steps:

1. Getting the 'ggman' binary
2. Configuring the installation

### getting the `ggman` binary

There are two primary options for getting and installing ggman. 


#### Download a pre-built binary. 
You can download a pre-built binary from the [releases page](https://github.com/tkw1536/ggman/releases/latest) on GitHub. 
This page includes releases for Linux, Mac OS X and Windows. 
Afterwards simply place it somewhere in your `$PATH`.

#### Install from Source

If you do not trust pre-built binaries you can build `ggman` from source. 
You need [`go`](https://golang.org) (version `1.21` or newer) along with `make` installed on your machine. 
After cloning this repository, you can simply type `make install` to have ggman installed automatically. 

#### Optional Dependencies

The compiled `ggman` binary does not depend on any external software. 
However having `git` in `$PATH` allows for more efficient operations in some cases. 

Furthermore, if you have a custom `.ssh/config` on your system, the `git`-less setup may not be fully supported.
You should install a native `git` executable on your system.

To check if a running ggman installation has found a `git` executable, run `ggman env --raw git`.
If one is found, it will print the path to it.
### configuring `ggman`

ggman is easy to configure. 

It manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
This root folder defaults to `~/Projects` but can be customized using the `$GGROOT` environment variable. 

Once the GGROOT environment variable is set, ggman is ready to be used. 

### first steps

To validate that the configuration is correct, you can list all repositories that ggman detects.
To do this type:

```
ggman ls
```

If you want to find a specific installed repository you can provide a `--for` argument.
For example:

```
ggman --for https://github.com/tkw1536/ggman ls
```

will print the location of a particular repository.

To clone a new repository into the root directory, you can use the `ggman clone` command.
For example:

```
ggman clone https://github.com/tkw1536/ggman
```

This will clone the repository https://github.com/tkw1536/ggman to the folder `$GGROOT/github.com/tkw1536/ggman` using ssh keys. 

If you want to move all existing repositories into the standardized structure you can use:

```
ggman relocate
```

If you want to only see what would be moved where you can instead use

```
ggman relocate --simulate
```

A more thorough documentation on the commands above and how the URL to path mapping works can be found in the thorough documentation below. 

### Optional: Additional shell aliases

In addition to ggman (described in detail below) certain aliases can also be very useful. 
They can be installed into your `.zshrc` or `.bashrc` by adding the following line:

```bash
eval "$(ggman shellrc)"
```

#### ggcd

`ggcd` allows 'cd'-ing into a directory given a repository name.
For example, `ggcd github.com/hello/world` will cd into the directory where the `github.com/hello/world` repository is checked out. 
This also works with any pattern matching a repository, e.g. `ggcd world` will cd into the first repository matching `world`.

#### ggclone

`ggclone` behaves similar to `ggman clone` and `ggcd`.
It takes the exact same arguments as `ggman clone`, but when the repository already exists in the default location does not clone it again.
Furthermore, after finishing the clone, automatically `cd`s into the cloned repository.

#### ggshow

ggshow is like ggcd, except that it prints the target directory and also shows the most recent `HEAD` commit.
This requires a locally installed git.

#### ggcode

ggcode is like ggcd, except it opens an editor (here vscode) instead of cd-ing.

## the `ggman` command

The `ggman` command is implemented in golang and can be compiled using standard golang tools. 
In addition, a `Makefile` is provided. 

The command is split into several sub-commands, which are described below. 
`ggman` has the following general exit behavior:

| Exit Code          | Description                                                           |
| ------------------ | --------------------------------------------------------------------- |
| 0                  | Everything went ok                                                    |
| 1                  | Command Parsing went ok, but a subcommand-specific error occurred     |
| 2                  | The user asked for an unknown subcommand                              |
| 3                  | Command-independent argument parsing failed, e.g. an invalid '--for'  |
| 4                  | Command-dependent argument parsing failed                             |
| 5                  | Invalid configuration                                                 |
| 6                  | Unable to parse a repository name                                     |


### 'ggman root' and 'ggman where'

It manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
The root folder defaults to `~/Projects` but can be customized using the `$GGROOT` environment variable. 
The root directory can be echoed using the command alias `ggman root`. 

For example, when `ggman` clones a repository `https://github.com/hello/world.git`, this would automatically end up in `$GGROOT/github.com/hello/world`. 
This works not only for `github.com` urls, but for any kind of url. 
To see where a repository would be cloned to (but not actually cloning it), use `ggman where <REPO>`. 

As of `ggman 1.12`, this translation of URLs into paths takes existing paths into account.
In particular, it re-uses existing sub-paths if they differ from the requested path only by casing.

For example, say the directory `$GGROOT/github.com/hello` exists and the user requests to clone `https://github.com/HELLO/world.git`.
Before 1.12, this clone would end up in `$GGROOT/github.com/HELLO/world`, resulting in two directories `$GGROOT/github.com/HELLO` and `$GGROOT/github.com/hello`. 
After 1.12, this clone will end up in `$GGROOT/github.com/hello/world`.
While this means placing of repositories needs to touch the disk (and check for existing directories), it results in less directory clutter.

By default, the first matching directory (in alphanumerical order) is used as opposed to creating a new one.
If a directory with the exact name exists, this is preferred over a case-insensitive match.

This normalization behavior can be controlled using the `GGNORM` environment variable.
It has three values:
- `smart` (use first matching path, prefer exact matches, default behavior);
- `fold` (fold paths, but do not prefer exact matches); and
- `none` (always use exact paths, legacy behavior)

### 'ggman ls'

While creating this folder structure when cloning new repositories, `ggman` can run operation on any other folder structure contained within the `GGROOT` directory. 
For this purpose the `ggman ls` command lists all repositories that have been found in this structure. 

For easier integration into scripts, `ggman ls` supports an `--exit-code` argument. 
If this is given, the command will return exit code 0 iff at least one repository is found, and exit code 1 otherwise.

Furthermore, the flag `--one` or the flags `--count` / `-n` can be given to limit the number of results.
This is useful in specific scripting circumstances.

### the '--for', '--here' and '--path' arguments

When running multi-repository operations, it is possible to limit the operations to a specific subset of repositories. 
This is achieved by using the 'for' keyword along with a pattern. 
For example, `ggman --for 'github.com/*/example' ls` will list all repositories from `github.com` that are named `example`. 

Examples for simple supported patterns can be found in this table:

| Pattern            | Examples                                                            |
| ------------------ | ------------------------------------------------------------------- |
| `world`            | `git@github.com:hello/world.git`, `https://github.com/hello/world`  |
| `hello/*`          | `git@github.com:hello/earth.git`, `git@github.com:hello/mars.git`   |
| `hello/m*`         | `git@github.com:hello/mars.git`, `git@github.com:hello/mercury.git` |
| `github.com/*/*`   | `git@github.com:hello/world.git`, `git@github.com:bye/world.git`    |
| `github.com/hello` | `git@github.com:hello/world.git`, `git@github.com:hello/mars.git`   |

Patterns are generally applied against URL components (see below for details on how the splitting works).
For example, to match the pattern `hello/*`, it is first split into the patterns `hello` and `*`.
These are then matched individually against the components of the URL.
The matched components have to be sequential, but don't have to be at either end of the URL.

Each component pattern can be one of the following:
- A case-insensitive `fnmatch.3` pattern with `*`, `?` and `[]` holding their usual meanings;
- A fuzzy string match, meaning the characters in the pattern have to occur in order of the characters in the string.
When no special fnmatch characters are found, the implementation assumes a fuzzy match. 
Fuzzy matching can also be explicitly disabled by passing the global `--no-fuzzy-filter` argument.

A special case is when a pattern begins with `^` or ends with `$` (or both).
Then any fuzzy matching is disabled, and any matches must start at the beginning  (in the case `^`) or end at the end  (in the case `$`) of the URL (or both).
For example `hello/world` matches both `git@github.com:hello/world.git` and `hello.com/world/example.git`, but `hello/world$` only matches the former.

Note that the `--for` argument also works for exact repository urls, e.g. `ggman --for 'https://github.com/tkw1536/ggman' ls`. 
`--for` also works with absolute or relative filepaths to locally installed repositories. 

In addition, the `--for` argument by default uses a fuzzy matching algorithm.
For example, the pattern `wrld` will also match a repository called `world`.
Fuzzy matching only works for patterns that do not contain a special glob characters (`*` and friends).
It is also possible to turn off fuzzy matching entirely by passing the `--no-fuzzy-filter` / `-n` argument.

In addition to the `--for` argument, you can also use the `--path` argument.
Instead of taking a pattern, it takes a (relative or absolute) filesystem path and matches all repositories under it.
This also works when not inside `GGROOT`.
The `--path` argument can be provided multiple times.

The `--here` argument is an alias for `--path .`, meaning it matches only the repository located in the current working directory, or repositories under it. 

### 'ggman comps' and 'ggman canon'

On `github.com` and multiple other providers, it is usually possible to clone repositories via multiple urls. 
For example, the repository `https://github.com/hello/world` can actually be cloned via:
- `https://github.com/hello/world.git` and
- `git@github.com:hello/world.git`

Usually the latter url is preferred over the former one in order to use SSH authentication instead of having to constantly having to type a password. 
For this purpose, ggman implements the concept of `canonical urls`, that is it treats the latter url as the main one and uses it to clone the repository. 
This behavior can be customized by the user. 

A canonical url is generated from an original url using a so-called `CANSPEC` (canonical specification).
An example CANSPEC is `git@^:$.git`. 

CANSPECs generate canonical urls by first taking the original urls, and splitting them into path-like components. 
These components also perform some normalization, such as removing common prefixes and suffixes. 
A few examples examples can be found in this table. 

| URL                            | Components                       |
| ------------------------------ | -------------------------------- |
| `git@github.com/user/repo`     | `github.com`, `user`, `repo`     |
| `github.com/hello/world.git`   | `github.com`, `hello`, `world`   |
| `user@server.com:repo.git`     | `server.com`, `user`, `repo`     |

To see exactly which components a URL has, use `ggman comps <URL>`.

After this, the canonical url is generated by parsing each character of the CANSPEC.
By default, a character of the CANSPEC simply ends up in the canonical url. 
However, two characters are treated differently:
- `%` is replaced by the second unused component of the URI (commonly a username)
- `$` is replaced by all remaining components in the URI joined with a '/'. Also stops all special character processing afterwards. 
If `$` does not exist in the cspec, it is assumed to be at the end of the CANSPEC.
Furthermore, the special CANSPEC `$$` always the original url unchanged.

A couple of examples can be found below:

| CANSPEC           | Canonical URL of the components `server.com`, `user`, `repository` |
| ----------------- | ------------------------------------------------------------------ |
| `git@^:$.git`     | `git@server.com:user/repository.git`                               |
| `ssh://%@^/$.git` | `ssh://user@server.com/repository.git`                             |
| (empty)           | `server.com/user/repository`                                       |

To get the canonical url of a repository use `ggman canon <URL> <CANSPEC>`. 

To customize the behavior globally, a so-called `CANFILE` can be used. 
This `CANFILE` should either be called `.ggman` in the users home directory, or be pointed to by the `GGMAN_CANFILE` environment variable. 

A `CANFILE` should consist of several lines.
Each line should contain either one or two space-separated strings. 
The first one is a pattern (as used with the `for` keyword) and the second is a CANSPEC to apply for all repositories matching this pattern. 
Empty lines and those starting with '#', '\\' are treated as comments.

An example CANFILE might be:

```
# for anything on git.example.com, clone with https
^git.example.com https://$.git

# for anything on git2.example.com leave the urls unchanged
^git2.example.com $$

# by default, clone via ssh
git@^:$.git
```


To resolve a canonical url with a CANFILE, simply omit the `CANSPEC` attribute of `ggman canon`. 

### 'ggman fix'

To fix an existing remote of a repository use `ggman fix`. 
This updates remotes of all matching repositories to their canonical form using the `CANFILE`. 
Optionally, you can pass a `--simulate` argument to `ggman fix`. 
Instead of storing any urls, it will only print what is being done to STDOUT. 

### 'ggman lsr'

To list the remotes of all installed repositories, use `ggman lsr`. 
It takes an optional argument `--canonical` which, if provided, cause ggman to print canonical urls instead of the provided ones. 

### 'ggman fetch' and 'ggman pull'

To fetch data for all repositories, or to run git pull, use `ggman fetch` and `ggman pull` respectively. 

### 'ggman clone', 'ggman link' and `ggman relocate`

To clone a new repository into the respective location, use `ggman clone` with the name of the repository as the argument, for example:

```bash
ggman clone git@github.com:hello/world.git
```

which will clone the the hello world repository into  `$GGROOT/github.com/hello/world`. 
This cloning not only works for the canonical repository url, but for any other url as well. 
For example:

```bash
ggman clone https://github.com/hello/world.git
```

will do the same as the above command.


When it is not desired that the canonical URL should be used, pass the `--exact-url` flag:

```bash
ggman clone --exact-url https://github.com/hello/world.git
```

This will clone using the exact url into the same folder as above. 

If ggman has access to a real `git` executable, it is also possible to pass additional arguments to it. 
For example:

```bash
ggman clone --exact-url https://github.com/hello/world.git -- --branch dev --depth 2
```

will execute the command ```git clone git@github.com:hello/world.git --branch dev --depth 2``` under the hood. 
The extra "--" is needed to allow ggman to separate the internal flags from the external flags. 

However sometimes for various reasons a repository needs to live in a non-standard location outside of `GGROOT`. 
For example, in the case of `go` packages these need to live within `$GOPATH`. 
In this case, it is sometimes useful to symlink these repositories into the existing directory structure. 
For this purpose, the `ggman link` command exists. 
This takes the path to the local clone of an existing repository, which will then be linked into the existing structure. 
For example

```bash
ggman link $HOME/go/src/github.com/hello/world
```

would link the repository in `$HOME/go/src/github.com/hello/world` into the right location. 
Here, this corresponds to `$GGROOT/github.com/hello/world`. 

Furthermore, sometimes a repository changes it's remote url and should be moved to the correct location. 
For this purpose the `ggman relocate` command can be used. 
It is called without arguments. 

### `ggman here`, `ggman web` and `ggman url`

```bash
ggman here
```
prints the current ggman-controlled repository. 
In addition, the command takes an optional `--tree` argument. 
When provided, also prints the location relative to the current git worktree root. 

Similarly, 

```bash
ggman web
```
attempts to open the url of the current repository in a web-browser. 
For this purpose it uses the CANSPEC `https://^/$`, which may not work with all git hosts. 
It also takes an optional `--tree`, which behaves similar and above and optionally opens a url pointing to the current folder. 

```bash
ggman url
```
is the same as `ggman web`, except that it only prints the URL to stdout. 

`ggman web` and `ggman url` also take an optional base url. 
If it is provided, the first component of the url is replace with the given base. 
ggman also supports a number of "default" base urls. 

For example:

```
ggman web godoc
```

will open the current repo on [pkg.go.dev](https://pkg.go.dev). 
To see a list of supported default urls, use `ggman web --list-bases`.

### 'ggman find-branch'

git 2.28 introduced the `init.defaultBranch` option to set the name of the default branch of new repositories. 
However this does not affect existing repositories. 

To find repositories with an old branch, the `ggman find-branch` command can be used. 
It takes a single argument (a branch name), and finds all repositories that contain a branch with the given name.

### 'ggman find-file'

Sometimes it is useful to find specific files inside repository directories.
This can be used to e.g. detect repositories of a specific language.

For this purpose the `ggman find-file` command can be used. 
It takes a single argument (a file name), and finds all repository directories that contain a file with the given path. 
For example, use `ggman find-file package.json` to find all repositories with a `package.json`.

### 'ggman sweep'

After moving repositories around (for example using `ggman relocate`, or by manual operations) empty directories are often left behind. 
Because of the nature of ggman sometimes directories containing only empty directories may also be left behind. 
To easily find these directories, the `ggman sweep` command can be used. 

It takes no arguments, and lists all directories, which are not git repositories and are empty, or contain only empty directories.
These are listed in such an order that they can be deleted in order using `rmdir` and friends.

### 'ggman exec'

Sometimes it is useful to run an arbitrary command over all the known git repositories.
This can be achieved using the `ggman exec` command.
It simply takes a command as an argument and runs it in each repository.

### 'ggman env'

To debug and inspect the current environment of the ggman command the `ggman env` command can be used.
The environment exposed consists of a set of variables representing the state.
Use `ggman env --list` to see a list of variables.
Use `ggman env --describe` to see their corresponding human-readable descriptions.
Use `ggman env --raw` to print the raw values of these variables in the same order.
Use `ggman env` without any arguments to print escaped (variable, value) pairs.

You can also provide a list of variables as arguments to list only those variables.
For instance: `ggman env --raw GGROOT` will print the unescaped, raw value of the GGROOT environment variable.
Variables are matched case-insensitive.

### Command Aliases

ggman comes with the following builtin aliases:

- `ggman git` behaves exactly like `ggman exec -- git`
- `ggman show` behaves exactly like `ggman exec -- git -c core.pager= show HEAD` to show the most recent HEAD commit of a repository 
- `ggman require` behaves exactly like `ggman clone --force`
- `ggman root` behaves exactly like `ggman env --raw GGROOT` (for backwards compatibility)

## Changelog

### 1.25.0 (Upcoming)

- change various short form options for consistency (global flags are upper case, local flags are lower case)
- rename `--here` flag of `ggman clone` to `--plain` (to avoid conflicts with the global `--here` flag)
- `ggman clone` and `ggman exec`: require `--` to separate flags to external commands
- tests: check overlap between command and global flags
- move go import paths to custom domain
- fix bug in shellrc (see #13, thanks @janezicmatej)
- update dependencies
- drop building a universal mac binary

### 1.24.0 (Released [May 12 2025](https://github.com/tkw1536/ggman/releases/tag/v1.24.0))

- add optional `--remote` url to `ggman url` and `ggman web`
- add special CANSPEC `$$` to return URL unchanged
- add anchors to patterns to match exactly at the start or end of a repository url
- update tool dependencies to latest
- internal improvements
- add a bunch of linting
- update some error messages in line with new `goprogram`
- update `goprogram` and `pkglib` dependencies  

### 1.23.1 (Released [Mar 30 2025](https://github.com/tkw1536/ggman/releases/tag/v1.23.1))

- update `ggman clone` help page wording 

### 1.23.0 (Released [Mar 30 2025](https://github.com/tkw1536/ggman/releases/tag/v1.23.0))

- update to `go 1.24`
- add `--exact-url` flag to `ggman clone`
- update dependencies to latest
- replace unneeded dependencies by standard library
- update to `goprogram` 0.7.0
- use new `go tool` during development
- migrate lint configuration to migrate
- minor performance improvements
- explicitly handle output errors


### 1.22.0 (Released [Sep 22 2024](https://github.com/tkw1536/ggman/releases/tag/v1.22.0))

- update to `go 1.23`
- add `--count` flag to `ggman ls`
- fix typo in Makefile

### 1.21.0 (Released [May 30 2024](https://github.com/tkw1536/ggman/releases/tag/v1.21.0))

- add `ggman find-file` command
- Various internal performance tweaks
- make spellchecker happy
- Update to `go 1.22`
- Update to goprogram `0.5.0`
- Update dependencies

### 1.20.1 (Released [Jun 1 2023](https://github.com/tkw1536/ggman/releases/tag/v1.20.1))

- fix `ggman relocate` behavior on Windows

### 1.20.0 (Released [Jun 1 2023](https://github.com/tkw1536/ggman/releases/tag/v1.20.0))

- add `ggman show` alias to show the head commit of a repository
- add `ggshow` utility
- rework some of the shell aliases
- `Makefile`: Make sure to always build with cgo disabled
- rework filter scoring to take position of match into account
- improve `ggman relocate` error handling in some edge cases 
- fix a lot of typos
- update to `goprogram` 0.4.0
- update various other dependencies

### 1.19.0 (Released [Apr 4 2023](https://github.com/tkw1536/ggman/releases/tag/v1.19.0))

- update to go 1.20
- `ggman exec` display output in parallel when running in parallel
- rename `--local` flag to `--here` in `ggman clone`
- improve `ggman relocate` behavior with symlinks
- update to `goprogram` 0.3.5
- minor bugfixes and CI updates
- update copyright year

### 1.18.2 (Released [Jul 14 2022](https://github.com/tkw1536/ggman/releases/tag/v1.18.2))

- fix another `ggclone` alias issue

### 1.18.1 (Released [Jul 10 2022](https://github.com/tkw1536/ggman/releases/tag/v1.18.1))

- fix `ggclone` alias

### 1.18.0 (Released [Jul 10 2022](https://github.com/tkw1536/ggman/releases/tag/v1.18.0))

- move aliases to new `ggman shellrc` command
  - add new `ggclone` alias
- build universal mac executables

### 1.17.0 (Released [May 30 2022](https://github.com/tkw1536/ggman/releases/tag/v1.17.0))

- update to new `goprogram`
	- format messages accordingly
	- remove `BeforeRegister` method from commands
	- remove unneeded pointer receivers
- add new `--scores` flag to `ggman ls`
- dependency updates

### 1.16.0 (Released [Apr 18 2022](https://github.com/tkw1536/ggman/releases/tag/v1.16.0))

- update to new `goprogram` version (refactors positional argument parsing)
- add `--list-bases` flag to `ggman web` and `ggman url`
- don't try to open invalid URLs in `ggman web` and `ggman url`
- add `ggman require` alias
- make documentation strings more consistent

### 1.15.0 (Released [Apr 8 2022](https://github.com/tkw1536/ggman/releases/tag/v1.15.0))

- add `--from-file` argument that reads `--for` arguments from a file
- add `ggman env` command to print information about ggman
- use native git when available in `ggman fetch`, `ggman pull`
- make built-in git `fetch`, `pull` and `clone` progress on standard error
- automatically upload releases from GitHub actions
- minor fixes

### 1.14.0 (Released [Mar 27 2022](https://github.com/tkw1536/ggman/releases/tag/v1.14.0))

- move to go `1.18`
- refactor `program` package into external `goprogram` package, using type-parameters and not depend on ggman
- add `--to` and `--local` flags to `ggman clone`
- fix `ggman pull` not respecting input / output streams
- README and documentation rework
- minor internal improvements

### 1.13.2 (Released [Jan 18 2022](https://github.com/tkw1536/ggman/releases/tag/v1.13.2))

- add new `--path` global argument to match repos under a specific path

### 1.13.1 (Released [Jan 16 2022](https://github.com/tkw1536/ggman/releases/tag/v1.13.1))

- add support for command aliases and add various aliases
- fix typos in README

### 1.13.0 (Released [Jan 14 2022](https://github.com/tkw1536/ggman/releases/tag/v1.13.0))

- add `ggman sweep` command
- add `ggman exec` command
- sort matches against fuzzy filters by score
- ensure shell escaping when generating scripts using `--simulate`
- prepare URLs to accept custom aliases
- internal testing improvements
- refactor main program initialization

### 1.12.0 (Released [Dec 23 2021](https://github.com/tkw1536/ggman/releases/tag/v1.12.0))

- add `GGNORM` variable: when placing repositories locally, take casing of existing paths into account
- add `--dirty` and `--clean` filter arguments
- add `--synced` and `--unsynced` filter arguments
- add `--tarnished` and `--pristine` filter arguments
- internal testing improvements

### 1.11.1 (Released [Sep 20 2021](https://github.com/tkw1536/ggman/releases/tag/v1.11.1))

- use `go1.17` for building and tests
- improved checking for local urls when running `ggman clone`
- minor internal improvements

### 1.11.0 (Released [May 8 2021](https://github.com/tkw1536/ggman/releases/tag/v1.11.0))

- add `--clone` and `--reclone` flag to `ggman url`
- add fuzzy matching support for repository patterns (can be disabled using `--no-fuzzy-filter`) 
- `ggman ls`: add `--one` argument to list at most one repository
- `ggman clone`: complain when trying to clone a local path
- internal code improvements and bugfixes

### 1.10.0 (Released [Mar 20 2021](https://github.com/tkw1536/ggman/releases/tag/v1.10.0))

- add `--force` flag to `ggman clone` to ignore errors when a cloned repository already exists.
- use `github.com/jessevdk/go-flags` to allow unknown options in argument parsing
- rewrite and extend help page generator
- embed license info using `go:embed`
- internal code improvements to `program` struct and text wrapping
- remove TODOs that are no longer required

### 1.9.0 (Released [Feb 21 2021](https://github.com/tkw1536/ggman/releases/tag/v1.9.0))

- move to go `1.16`
- `--for` now also matches filepaths
- add a new utility method to cleanup repeated code
- move `util` and `testutil` packages into new `internal` subpackage

### 1.8.0 (Released [Jan 31 2021](https://github.com/tkw1536/ggman/releases/tag/v1.8.0))

- add `force-repo-here` flag to `ggman web` and `ggman url` to force a repository even when there is none
- include go version when calling `ggman version`
- use `go 1.15` in tests
- improve `util/record.go` implementation
- update copyright year

### 1.7.2 (Released [Nov 29 2020](https://github.com/tkw1536/ggman/releases/tag/v1.7.2))

- `--for` now matches remote URL instead of clone path
- bug fixes

### 1.7.1 (Released [Nov 27 2020](https://github.com/tkw1536/ggman/releases/tag/v1.7.1))

- rewrite `--here` and `--for` filter flags
- minor bug fixes

### 1.7.0 (Released [Nov 22 2020](https://github.com/tkw1536/ggman/releases/tag/v1.7.0))

- add `--here` flag as a convenience to filter the current repository
- add `ggman relocate` command to move repositories to where they should be
- `ggman clone`: Only create parent folder to clone repository
- improved windows support

### 1.6.0 (Released [Oct 4 2020](https://github.com/tkw1536/ggman/releases/tag/v1.6.0))

- rework, document and test new package structure
- rework git repository scanning & filtering to take place in parallel
- add internal utility to automatically re-generate license notices
- add `branch` flag to `ggman web` and `ggman url`
- bugfix: Add missing `ggman find-branch` documentation
- replace `flag` package by POSIX-compatible `pflag` package
- cleanup runtime version handling

### 1.5.0 (Released [Aug 29 2020](https://github.com/tkw1536/ggman/releases/tag/v1.5.0))

- added `ggman find-branch` command
- upgrade to new version of go-git
- internal optimization and documentation
- rewrite and optimize internal url handling
- rewrite internal handling of git commands
- moved from Travis to GitHub Actions

### 1.4.1 (Released [Jul 23 2020](https://github.com/tkw1536/ggman/releases/tag/v1.4.1))

- stop compressing binary releases with upx
- added 'godoc' and 'localgodoc `BASE` urls to `ggman web` and `ggman url`
- rewrite handling of CanFile
- added more tests

### 1.4.0 (Released [Jul 1 2020](https://github.com/tkw1536/ggman/releases/tag/v1.4.0))

- added a `BASE` url to `ggman web` and `ggman url`
- refactored internal flag handling
- refactored `Makefile`
- moved to `travis-ci.com`
- added a CHANGELOG to the README

### 1.3.0 (Released [Aug 29 2019](https://github.com/tkw1536/ggman/releases/tag/v1.3.0))

- added `ggman web` which opens the current repository in a web browser
- added `ggman url` which prints the (web) url of the current repository
- added `ggman here` command which prints the current repository

### 1.2.0 (Released [Jul 10 2019](https://github.com/tkw1536/ggman/releases/tag/v1.2.0))

- `ggman link` now creates absolute symlinks instead of exactly echoing the path the user entered
- Added a bash alias to the README

### 1.1.0 (Released [Apr 11 2019](https://github.com/tkw1536/ggman/releases/tag/v1.1.0))

- Use external `git clone` command when available and allow passing options to it
- Added help command and better subcommand help behavior
- Add versioning information to help page

### 1.0.0 (Released [Feb 17 2019](https://github.com/tkw1536/ggman/releases/tag/v1.0.0))

- Initial release

## LICENSE

`ggman` is licensed under the terms of the MIT LICENSE, see [LICENSE](LICENSE). 
To view accompanying license information use `ggman license`. 
