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
```sh
Install: wadman -i <target> / wadman --install <target> 
Uninstall: wadman -u <target> / wadman --uninstall <target> 
Search: wadman -s <target> / wadman --search <target>
Run: wadman -r <targe> / wadman --run <target>
List: wadman -l / wadman --list

```