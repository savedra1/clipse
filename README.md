<a href="https://github.com/savedra1/clipse/actions"><img src="https://github.com/charmbracelet/bubbletea/workflows/build/badge.svg" alt="Build Status"></a>

<p align="centre">
  <img src="./resources/examples/demo.gif?raw=true" width=50% alt="gif" />
</p>

# About 

`clipse` is a dependency-less, configurable TUI-based clipboard manager built with Go. Though the app is optimised for a linux OS with a dedicated window manager, `clipse` can also be used on any Unix-based system. Simply install the package and bind the open command to get [your desired clipboard behavior](https://www.youtube.com/watch?v=ZE2F8Mj0_I0). Further instructions for setting this up can be found below.

[Click here to see a video demo for clipse](https://www.youtube.com/watch?v=ZE2F8Mj0_I0)

## Why use *clipse*?

### 1. Configurability 

Customizable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default but is customisable from a `.config/clipse/custom_theme.json` file that gets created when the program is ran for the first time. Some example themes (based on my terminal)...

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

An example `.config/clipse/custom_theme.json`: 

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

Simply leaving this file alone or setting the `useCustomTheme` value to `false` will give you a nice pink and where default theme... 

<p align="left">

  <img src="./resources/examples/default.png?raw=true" alt="Gruvbox" />

</p>

### 2. Usability 

Easily recall, add, and delete clipboard history via a dreamy TUI built with Go's excellent [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) library. A simple fuzzy finder, callable with the `/` key can easily match content from a theoretically unlimited amount of time in the past: 

<p align="left">

  <img src="./resources/examples/fuzzy.png?raw=true" alt="Gruvbox" />

</p>

Items can be permanently deleted from the list by simply hitting `backspace` when the item is selected, as seen in the [demo video](https://youtu.be/ZE2F8Mj0_I0), and can be added explicitly from the command line using the `-a` flag. This would allow you to easily pipe any CLI output directly into your history with commands like:

```shell

ls /home | clipse -a

``` 

```shell

clipse -a "a custom string value"

```

### 3. Efficiency 

Due to Go's inbuilt garbage collection system and the way the application is built, `clipse` is pretty selfless when it comes to CPU consumption and memory. The below image shows how little resources are required to poser the background event listener used to continually update the history displayed in the TUI... 

<p align="left">

  <img src="./resources/examples/htop.png?raw=true" alt="Gruvbox" />

</p>

### 4. Versatility

The `clipse` binary, installable from the repo, can run on pretty much any Unix-based OS and will require zero external dependencies. Being terminal-based also allows for easy integration with a window manager and configuration of how the TUI behaves. For example, binding a floating window to the `clipse` command as shown in [my example](https://youtu.be/ZE2F8Mj0_I0) using [Hyprland window manager](https://hyprland.org/) on __NixOs__.

**Note that working with image files will require one of the following dependencies**:

- Linux (X11) & MacOs: [xclip](https://github.com/astrand/xclip)
- Linux (Wayland): [wl-clipboard](https://github.com/bugaevc/wl-clipboard)

# Setup & installation

## Installation

### Installing on NixOs

TBC

### Installing with Go

```shell

go install https://github.com/savedra1/clipse@latest

```

### Building from source

```shell

git clone https://github.com/savedra1/clipse

cd clipse

go build -o clipse

```

## Set up

As mentioned earlier, to get the most out of `clipse` you'll want to bind the two key commands to your systems config. The first key command is to open the clipboard history TUI:

```shell

clipse $PPID

```

Passing in the `$PPID` variable as an arg to the main command allows the TUI to close the terminal session in which it's hosted on the `choose` event, simulating a full GUI experience. Without passing in `$PPID`, your TUI selection will still be copied to your system's clipboard, however, the terminal session will not close automatically.    

The second command doesn't need to be bound to a key combination, but rather to the system boot to run the background listener on start-up:

```shell

clipse -listen  

``` 

The above command creates a `nohup` process of `clipse --listen-shell`, which if called on its own will start a listener in your current terminal session instead.

### Hyprland
Add the following lines to your Hyprland config to acheive the optimal TUI behaviour:
```shell
exec-once = clipse -listen # run listener on startup
windowrulev2 = float,class:(floating) # ensure you have defined a floating window class
bind = SUPER, V, exec,  <terminal name> --class floating -clipse $PPID # bind the open clipboard operation to a nice key
```

### i3 
TBC

## All commands

```shell

clipse $PPID # Open Clipboard TUI

clipse -listen # Run a background listener process

clipse --listen-shell # Run a listener process in the current terminal

clipse -help # Display menu option

clipse -v # Get version

clipse -clear # Wipe all clipboard history

clipse -kill # Kill any existing background processes

clipse # Open Clipboard TUI in persistant/debug mode

```

## How it works

When the app is run for the first time it creates a `/home/$USER/.config/clipse` dir with a `clipboard_hostory.json` file, a `custom_theme.json` file, and `tmp_files` folder for storing image data. After the `clipse -listen` command is executed, a background process will be watching for clipboard activity and adding any changes to the `clipboard_hstory.json` file. 

The TUI that displays the clipboard history should then be called with the `clipse $PPID` command. Passing in the terminal's PPID is irregular, but allows the terminal-based app to close itself from within the program itself, simulating the behavior of a full GUI without the memory overhead. A worthy trade-off in my opinion. 

Operations within the TUI are defined with the [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) framework, allowing for efficient concurrency and a smooth UX. `Delete` operations will remove the selected item form the TUI view and the storage file, `select` operastions will copy the item to the systems clipboard and close the terminal window in which the session is currently hosted.  

The maximum item storage limit is currently hardcoded at **100**. However there are plans to make this configurable in the future.

## Contributing

I would love to receive contributions to this project and welcome PRs from anyone and everyone. The following is a list of example future enhancements I'd like to implement:
- System Paste option (building functionality to paste the chosen item directly into th next place of focus after the TUI closes)
- Theme adjustments made available via CLI 
- Better debugging
- Use of a GUI library such as fyne/GIO (only with minimal CPU cost)
- Custom key binds added to config file 

## TODO
- Nix package
- Apt package
- Makefile installation
