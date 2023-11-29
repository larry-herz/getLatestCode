# getLatestCode

## Description

This is a go program that will parse a directory tree, looking for git repos and perform a `git pull` on all repos found.

The routine requires that an environment variable titled `SRC` exists. This is the top level directory where teh search is to begin.

Any directory that contains a `.terraform` directory is excluded, as the `.terraform` directory is added to terraform modules on initialization.

Any repo that does not have a remote defined is skipped, and the information is displayed on the command line

Any repo that is unclean, that is has any outstanding commits, is also skipped to avoid any accidental overwrites of uncommitted code.
