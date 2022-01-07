# ggman

![CI Status](https://github.com/tkw1536/ggman/workflows/CI/badge.svg)

A golang tool that can manage all your git repositories. 
Originally a rewrite of [GitManager](https://github.com/tkw1536/GitManager), but has diverged. 

## Why ggman?

When you only have a couple of git repositories that you work on it is perfectly feasible to manage them by using `git clone`, `git pull` and friends. 
However once the number of repositories grows beyond a small number this can become tedious:

- It is hard to find which folder a repository has been cloned to
- Getting an overview of what is cloned and what is not is hard
- It's not easily possible to perform actions on more than one repos at once, e.g. `git pull`

This is the problem `ggman` is designed to solve. 
It allows one to:

- Maintain and expand a local directory structure of multiple repositories
- Run actions (such as `git clone`, `git pull`) on groups of repositories

ggman is not limited to using specific providers such as GitHub or GitLab.
In particular all core functionality works either offline, or is only based only on "git" itself.

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
You need [`go`](https://golang.org) along with `make` installed on your machine. This project is tested on go `1.17`, older releasees may or may not work. 

After cloning this repository, you can then simply type `make` and executables will be generated inside the `dist/` directory. 

You can then either manually place them in your `$PATH`, or type `make install` to have it installed directly. 

#### Optional Dependencies

The compiled `ggman` binary does not depend on any external software. Nonetheless, having`git` in `$PATH` allows for more efficient operations in some cases. 

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
The root directory can be echoed using the command `ggman root`. 

For example, when `ggman` clones a repository `https://github.com/hello/world.git`, this would automatically end up in `$GGROOT/github.com/hello/world`. 
This works not only for `github.com` urls, but for any kind of url. 
To see where a repository would be cloned to (but not actually cloning it), use `ggman where <REPO>`. 

As of `ggman 1.12`, this translation of URLs into paths takes existing paths into account.
In particular, it re-uses existing sub-paths if they differ from the requested path only by casing.

For example, say the directory `$GGROOT/github.com/hello` exist and the user requests to clone `https://github.com/HELLO/world.git`.
Before 1.12, this clone would end up in `$GGROOT/github.com/HELLO/world`, resulting in two directories `$GGROOT/github.com/HELLO` and `$GGROOT/github.com/hello`. 
After 1.12, this clone will end up in `$GGROOT/github.com/hello/world`.
While this means placing of repositories needs to touch the disk (and check for existing directories), it results in less directory clutter.

By default, the first matching directory (in alphanumerical order) is used as opposed to creating a new one.
If a directory with the exact name exists, this is prefered over a case-insensitive match.

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

### the '--for' and '--here' arguments

When running multi-repository operations, it is possible to limit the operations to a specific subset of repositories. 
This is achieved by using the 'for' keyword along with a pattern. 
For example, `ggman --for 'github.com/*/example' ls` will list all repositories from `github.com` that are named `example`. 

Examples for supported patterns can be found in this table:

| Pattern            | Examples                                                            |
| ------------------ | ------------------------------------------------------------------- |
| `world`            | `git@github.com:hello/world.git`, `https://github.com/hello/world`  |*
| `hello/*`          | `git@github.com:hello/earth.git`, `git@github.com:hello/mars.git`   |
| `hello/m*`         | `git@github.com:hello/mars.git`, `git@github.com:hello/mercury.git` |
| `github.com/*/*`   | `git@github.com:hello/world.git`, `git@github.com:bye/world.git`    |
| `github.com/hello` | `git@github.com:hello/world.git`, `git@github.com:hello/mars.git`   |

Note that the `--for` argument also works for exact repository urls, e.g. `ggman --for 'https://github.com/tkw1536/GitManager' ls`. 
`--for` also works with absolute or relative filepaths to locally installed repositories. 

In addition to the `--for` argument, the `--here` argument also works.
It is a shortcut that is equivalent to calling `--for` with the remote of the repository located in the current working directory. 

In addition, the `--for` argument by default uses a fuzzy matching algorithm.
For example, the pattern `wrld` will also match a repository called `world`.
Fuzzy matching only works for patterns that do not contain a special glob characters (`*` and friends).
It is also possible to turn off fuzzy matching entirely by passing the `--no-fuzzy-filter` / `-n` argument.

### 'ggman comps' and 'ggman canon'

On `github.com` and multiple other providers, it is usually possible to clone repositories via multiple urls. 
For example, the repository `https://github.com/hello/world` can actually be cloned via:
- `https://github.com/hello/world.git` and
- `git@github.com:hello/world.git`

Usually the latter url is prefered over the former one in order to use SSH authentication instead of having to constantly having to type a password. 
For this purpose, ggman implements the concept of `canonical urls`, that is it treats the latter url as the main one and uses it to clone the repository. 
This behaviour can be customized by the user. 

A canonical url is generated from an original url using a so-called `CANSPEC` (canonical specification).
An example CANSPEC is `git@^:$.git`. 

CANSPECs generate canonical urls by first taking the original urls, and splitting them into path-like components. 
These components also perform some normalisation, such as removing common prefixes and suffixes. 
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

A couple of examples can be found below:

| CANSPEC           | Canonical URL of the components `server.com`, `user`, `repository` |
| ----------------- | ------------------------------------------------------------------ |
| `git@^:$.git`     | `git@server.com:user/repository.git`                               |
| `ssh://%@^/$.git` | `ssh://user@server.com/repository.git`                             |
| (empty)           | `server.com/user/repository`                                       |

To get the canonical url of a repository use `ggman canon <URL> <CANSPEC>`. 

To customize the behaviour globally, a so-called `CANFILE` can be used. 
This `CANFILE` should either be called `.ggman` in the users home directory, or be pointed to by the `GGMAN_CANFILE` environment variable. 

A `CANFILE` should consist of several lines.
Each line should contain either one or two space-seperated strings. 
The first one is a pattern (as used with the `for` keyword) and the second is a CANSPEC to apply for all repositories matching this pattern. 
Empty lines and those starting with '#', '\\' are treated as comments. 

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
This clonening not only works for the canonical repository url, but for any other url as well. 
For example:

```bash
ggman clone https://github.com/hello/world.git
```

will do the same as the above command. 

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
ggman web travisci
```

will open the current repo on travis-ci.com. 
The following base urls are supported:

- `circle` - CircleCI
- `travis` - TravisCI

### 'ggman find-branch'

git 2.28 introduced the `init.defaultBranch` option to set the name of the default branch of new repositories. 
However this does not affect existing repositories. 

To find repositories with an old branch, the `ggman find-branch` command can be used. 
It takes a single argument (a branch name), and finds all repositories that contain a branch with the given name. 

### 'ggman sweep'

After moving repositories around (for example using `ggman relocate`, or by manual operations) empty directories are often left behind. 
Because of the nature of ggman sometimes directories containing only empty directories may also be left behind. 
To easily find these directories, the `ggman sweep` command can be used. 

It takes no arguments, and lists all directories, which are not git repositories and are empty, or contain only empty directories.
These are listed in such an order that they can be deleted in order using `rmdir` and friends.

### 'ggman exec'

Sometimes it is useful to run an arbitray command over all the known git repositories.
This can be achieved using the `ggman exec` command.
It simply takes a command as an argument and runs it in each repository.

### Useful aliases

In addition to ggman certain aliases can be very useful. 

```bash
# ggcd allows 'cd'-ing into a directory given a repository name
# e.g ggcd github.com/hello/world will cd into the directory where the
# 'github.com/hello/world' repository is checked out. 
ggcd () {
	ggman -f $1 ls --exit-code --one && cd $(ggman -f $1 ls --exit-code 2>&1)
}
# ggcode is like ggcd, except it opens an editor (here vscode) instead of cding. 
ggcode () {
	ggman -f $1 ls --exit-code && code $(ggman -f $1 ls --exit-code 2>&1)
}
```

## Changelog

### 1.13.0 (Upcoming)

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
- Added help command and better subcommand help behaviour
- Add versioning information to help page

### 1.0.0 (Released [Feb 17 2019](https://github.com/tkw1536/ggman/releases/tag/v1.0.0))

- Initial release

## LICENSE

`ggman` is licensed under the terms of the MIT LICENSE, see [LICENSE](LICENSE). 
To view accompanying license information use `ggman license`. 
