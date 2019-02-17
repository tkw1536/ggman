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
| 5                  | Invalid configuration                                               |
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
For example, `ggman for 'github.com/*/example' ls` will list all repositories from `github.com` that are named `example`. 

Examples for supported patterns can be found in this table:

| Pattern            | Examples                                                            |
| ------------------ | ------------------------------------------------------------------- |
| `world`            | `git@github.com:hello/world.git`, `https://github.com/hello/world`  |
| `hello/*`          | `git@github.com:hello/earth.git`, `git@github.com:hello/mars.git`   |
| `hello/m*`         | `git@github.com:hello/mars.git`, `git@github.com:hello/mercury.git` |
| `github.com/*/*`   | `git@github.com:hello/world.git`, `git@github.com:bye/world.git`    |
| `github.com/hello` | `git@github.com:hello/world.git`, `git@github.com:hello/mars.git`   |

Note that the `for` keyword also works for exact repository urls, e.g. `ggman for 'https://github.com/tkw1536/GitManager' ls`. 

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

### 'ggman clone' and 'ggman link'

To clone a new repoistory into the respective location, use `ggman clone` with the name of the repository as the argument, for example:

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




## LICENSE

`ggman` is licensed under the terms of the MIT LICENSE, see [LICENSE](LICENSE). 