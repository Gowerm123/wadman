# Custom DOOM WADs (for the linux DOOMer)
WadMan is a wad archive manager that automatically downloads from the [DoomWorld IdGames](https://www.doomworld.com/idgames/) database. To install wadman, you will need `git` installed. Then, you can clone this repository, and install.
```
    git clone https://github.com/Gowerm123/WADMan.git
    cd wadman
    sudo ./install.sh BASE_INSTALL_DIR
```

## Commands
WadMan supports nine basic commands

 - `search QUERY <QUERYTYPE>` - Searches the IdGames archive for the specified QUERY,QUERYTYPE is optional, and defaults to filename. Possible options are filename, title, author, email, description, credits, editors, textfile.
 - `install QUERY <QUERYTYPE>` - First performs a `search QUERY <QUERYTYPE>` then installs the first found file. It is recommended that you search based on filename here to narrow down overlapping projects.
 - `list` - Lists all currently installed wad archives. Information printed is name of archive, installed directory, idGamesUrl, and Aliases.
 - `remove NAME` - Removes the archive with the given name. If two are found, the first will be deleted.
 - `run` - There are two ways to call `run`. You can either call `run ALIAS/NAME` or `run IWAD ALIAS/NAME`. Note that you must include the IWAD if you have not registered an IWAD to the given `ALIAS/NAME`.
 - `register NAME IWAD` - Assigns the IWAD to the archive entry in the `pkglist` associated with NAME. This is used for the `run` command so you do not have to specify IWADs everytime you load a PWAD.
 - `configure` - Runs you through a prompt to fill out the configuration file. The file is a simple JSON file found at `/usr/share/.wadmanConfig`
 - `help` - Prints this text
 - `alias TARGET ALIAS` - Assigns an alias to the given archive. This alias can be used when performing the `run` command.