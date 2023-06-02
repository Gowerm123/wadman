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
WadMan supports five basic commands

**Install**: `wadman -i target / wadman --install target` \
**Uninstall**: `wadman -u target / wadman --uninstall target` \
**Query**: `wadman -q target / wadman --query target` \
**Run**: `wadman -r target / wadman --run target` \
**List**: `wadman -l / wadman --list` \
**Set**: `wadman -s KEY VALUE / wadman --set KEY VALUE`
- You can use the Set command to set the LAUNCHER, LAUNCHARGS, IWAD, or MIRRORS configurations.
    - When setting IWad Aliases (IWAD) use the format \
    ```wadman -s IWAD doom2=path/to/DOOM2.WAD```