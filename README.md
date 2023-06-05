# Custom DOOM WADs (for the linux DOOMer)

Build status (main) - ![build badge](https://github.com/gowerm123/wadman/actions/workflows/go.yml/badge.svg)

WadMan is a WAD archive manager that automatically downloads from the [DoomWorld IdGames](https://www.doomworld.com/idgames/) archive. To install wadman, you will need `git`, and `go` v1.14+ installed. Then, you can clone this repository, and install.
```
    git clone https://github.com/Gowerm123/wadman.git
    cd wadman
    ./install.sh
```

## Commands
WadMan supports seven basic commands

**Install**: `wadman -i target / wadman --install target` \
**Uninstall**: `wadman -r target / wadman --remove target` \
**Query**: `wadman -q target / wadman --query target` \
**Run**: `wadman -p target / wadman --play target` \
**List**: `wadman -l / wadman --list` \
**Sync**: `wadman -sy / wadman --sync` 
- If you have an intact `$HOME/.wadman/wadmanifest.json` then the sync command can be used to re-install any archives that have been removed/lost unintentionally. \

**Set**: `wadman -s KEY VALUE / wadman --set KEY VALUE`
- You can use the Set command to set the LAUNCHER, LAUNCHARGS, IWAD, or MIRRORS configurations.
    - When setting IWad Aliases (IWAD) use the format ```wadman -s IWAD iwad=path/to/DOOM2.WAD``` where iwad is one of `plutonia`, `doom2`, `tnt`, `doom`