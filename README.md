# getignore

Description
=============

getignore is a command line tool to get .gitignore file from [github/gitignore](https://github.com/github/gitignore).

INSTALLATION
=============

```
$ go get github.com/sudix/getignore
```

COMMANDS
=============

### update

Clone `github/gitignore` repository into a work directory($HOME/.getignore/gitignore_files).

```
$ getignore update
```

### list

List up available .gitignore files.

```
$ getignore list
Actionscript.gitignore
Ada.gitignore
Agda.gitignore
Android.gitignore
AppEngine.gitignore
AppceleratorTitanium.gitignore
ArchLinuxPackages.gitignore
Autotools.gitignore
        .
        .
        .
```

If a query argument is given, only files whose name contain the query are listed.

```
$ getignore list PHP
CakePHP.gitignore
FuelPHP.gitignore
```

### get

Copy specified .gitignore file to current directory.
A file name must be given.(case sensitive)

```
$ getignore get Go
# This copies Go.gitignore to .gitignore.
```
