# About 
`clipse` is a dependencyless, configurable TUI-based clipboard manager built with Go. Though the app is optimised for a linux OS with a dedicated window manager, `clipse` can also be used on Macos and Windows. Simply install the package and bind the open command to get [your desired clipboard behaviour](https://www.youtube.com/watch?v=ZE2F8Mj0_I0). Instructions for setting this up below.

[Click here to see a video demo for clipse](https://www.youtube.com/watch?v=ZE2F8Mj0_I0)

## Why use *clipse*?
1.  Customizable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default, but is them customisable from a `.config/clipse/custom_theme.json` file that gets created when the program is ran for the first time. Some example themes (based on my terminal)...

**Nord**
<p align="left">
  <img src="./resources/examples/nord.png?raw=true" alt="Nord" />
</p>

**Dracula**
<p align="left">
  <img src="./resources/examples/dracula.png?raw=true" alt="Dracula" />
</p>

**Gruvbox**
<p align="left">
  <img src="./resources/examples/gruvbox.png?raw=true" alt="Gruvbox" />
</p>

an example `.config/clipse/custom_theme.json`: 
```
{
    "useCustomTheme": true,
    "DimmedDesc": "#4C566A",
    "DimmedTitle": "#4C566A",
    "FilteredMatch": "#A3BE8C",
    "NormalDesc": "#81A1C1",
    "NormalTitle": "#B48EAD",
    "SelectedDesc": "#A3BE8C",
    "SelectedTitle": "#A3BE8C",
    "SelectedBorder": "#88C0D0",
    "SelectedDescBorder": "#88C0D0",
    "TitleFore": "#D8DEE9",
    "Titleback": "#3B4252",
    "StatusMsg": "#8FBCBB"
}

```
Simply leaving this file alone or setting the will give you a default theme 

- Fuzzy finder capabilities
- Easily recall, add and delete clipboard history with no explicit expiry
- Image file support
- Terminal-based means TUI can be easily called bound to a key, simulating a GUI without the memory overhead   




# Setup & installation

# Themes

## Considerations

## TODO
- README (with videos)
- Nix package
- Makefile installation
- blog posts

## Future Considerations
- Use with rofi front end? (Can add a section for this is readme but not essential) - Needs an arg to just list the history essentially 
- Sytem Paste option
- Extra config added to json file, adjustable with CLI args
- Better debugging? (currently can use > to a custom file)
- Use of a GUI library such as fyne/GIO
- Increased customisation (keybinds/themes)
