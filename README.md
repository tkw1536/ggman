# ggman

[![Build Status](https://travis-ci.org/tkw1536/ggman.svg?branch=master)](https://travis-ci.org/tkw1536/ggman)

A golang script that can manage multiple git repositories locally. 
Golang rewrite of [GitManager](https://github.com/tkw1536/GitManager). 

## the `ggman` command

The `ggman` command is implemented in golang and can be compiled using the standard golang tools. 
In addition, a `Makefile` is provided. 

The command is split into several sub-commands, which are described below. 
`ggman` has the following general exit behaviour:

| Exit Code          | Description                                                         |
| ------------------ | ------------------------------------------------------------------- |
| 0                  | Everything went ok                                                  |
| 1                  | Command Parsing went ok, but a subcommand-specific error occured    |
| 2                  | The user asked for an unknown subcommand                            |
| 3                  | Command-independent argument parsing failed, e.g. an invalid 'for'  |
| 4                  | Command-dependent argument parsing failed                           |
| 5                  | The `GGROOT` directory does not exist                               |
| 6                  | Unable to parse a repository name                                   |


### 'ggman root' and 'ggman where'

It manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
The root folder defaults to `~/Projects` but can be customized using the `$GGROOT` environment variable. 
The root directory can be echoed using the command `ggman root`. 

For example, when `ggman` clones a repository `https://github.com/hello/world.git`, this would automatically end up in `$GGROOT/github.com/hello/world`. 
This works not only for `github.com` urls, but for any kind of url. 
To see where a repository would be cloned to (but not actually cloning it), use `ggman where <REPO>`. 

### 'ggman ls'

While creating this folder structure when cloning new repositories, `ggman` can run operation on any other folder structure contained within the `GGROOT` directory. 
For this purpose the `ggman ls` command lists all repositories that have been found in this structure. 

For easier integration into scripts, `ggman ls` supports an `--exit-code` argument. 
If this is given, the command will return exit code 0 iff at least one repository is found, and exit code 1 otherwise. 

### the 'for' keyword
When running multi-repository operations, it is possible to limit the operations to a specific subset of repositories. 
This is achieved by using the 'for' keyword along with a pattern. 
For example, `ggman for 'github.com/*/frontend' ls` will list all repositories from `github.com` that are named `frontend`. 

Examples for supported patterns can be found in this table:

| Pattern            | Examples                                                            |
| ------------------ | ------------------------------------------------------------------- |
| `world`            | `git@github.com:hello/world.git`, `https://github.com/hello/world`  |
| `hello/*`          | `git@github.com:hello/earth.git`, `git@github.com:hello/mars.git`   |
| `hello/m*`         | `git@github.com:hello/mars.git`, `git@github.com:hello/mercury.git` |
| `github.com/*/*`   | `git@github.com:hello/world.git`, `git@github.com:bye/world.git`    |
| `github.com/hello` | `git@github.com:hello/world.git`, `git@github.com:hello/mars.git`   |

Note that the `for` keyword also works for exact repository urls, e.g. `ggman for 'https://github.com/tkw1536/GitManager' ls`. 
