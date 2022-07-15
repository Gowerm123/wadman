# Custom DOOM WADs (for the linux DOOMer)

Build status (main) - ![build badge](https://github.com/gowerm123/wadman/actions/workflows/go.yml/badge.svg)

WadMan is a wad archive manager that automatically downloads from the [DoomWorld IdGames](https://www.doomworld.com/idgames/) database. To install wadman, you will need `git`, and `go` v1.17+ installed. Then, you can clone this repository, and install.
```
    git clone https://github.com/Gowerm123/wadman.git
    cd wadman
    sudo ./install.sh
```

Note, because wadman uses the `$SUDO_USER` environment variable to identify the user's home directory, `doas` is not supported.

## Commands
WadMan supports nine basic commands

 - `search QUERY` - Searches the IdGames archive for the specified QUERY.
 - `install QUERY` - Runs search, if one file is found, it is installed, if multiple files are found, they are displayed for you to choose.
 - `list` - Lists all currently installed wad archives. Information printed is name of archive, installed directory, idGamesUrl, and Aliases.
 - `remove NAME` - Removes the archive with the given name. If two are found, the first will be deleted.
 - `run` - There are two ways to call `run`. You can either call `run ALIAS/NAME` or `run IWAD ALIAS/NAME`. Note that you must include the IWAD if you have not registered an IWAD to the given `ALIAS/NAME`.
 - `register NAME IWAD` - Assigns the IWAD to the archive entry in the `pkglist` associated with NAME. This is used for the `run` command so you do not have to specify IWADs everytime you load a PWAD.
 - `configure` - Runs you through a prompt to fill out the configuration file. The file is a simple JSON file found at `$HOME/.config/wadman-config.json`
 - `help` - Prints this text
 - `alias TARGET ALIAS` - Assigns an alias to the given archive. This alias can be used when performing the `run` command.

 ## Aliases

 There are two types of aliases that will be refrenced. IWAD Aliases, and PWAD aliases.

 ### IWAD Aliases

You can configure IWAD aliases using the `configure` command. IWAD aliases allow you to reference IWAD files with a chosen term, instead of a full file path. These aliases work when using either form of the `run` command.

### PWAD Aliases

PWADs can be aliased as well. Please note that the run command only runs PWAD archives from the root directory of the archive. This can cause it to fail with some PWADs, for instance, Sunlust will need to be ran using `sunlust/sunlust` as the `ALIAS/NAME`