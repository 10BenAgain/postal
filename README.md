### Basic Info

Cross-platform TUI application for Pokemon Fire Red & Leaf Green mail glitch simulation and editing. Postal can read from `.pk3`, `.ek3`, and `.sav` files to load Pokemon data into the editor models. These files are generated from [PKHEX](https://projectpokemon.org/home/files/file/1-pkhex/) (except .sav) and are convenient ways to manage Pokemon data. 

To select a Pokemon from your save file and load it into the editor, simply tell postal where the `.sav` file is located
```
postal [path-to-save-file]
```
The same argument works for just an `.pk3` or .`ek3`
```
postal [path-to-pokemon-file]
```
If you don't the exact path, just run `postal` and file picker menu will appear where you can select the desired file. 
