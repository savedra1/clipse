<a href="https://github.com/savedra1/clipse/actions"><img src="https://github.com/charmbracelet/bubbletea/workflows/build/badge.svg" alt="Build Status"></a>   ![GitHub last commit](https://img.shields.io/github/last-commit/savedra1/clipse)

<p align="centre">
  <img src="./resources/examples/demo.gif?raw=true" width=65% alt="gif" />
</p>

<br>

<summary>Table of contents</summary>

- [Why use clipse?](#why-use-clipse)
- [Installation](#installation)
- [Set up](#set-up)
- [All commands](#all-commands-💻)
- [How it works](#how-it-works-🤔)
- [Contributing](#contributing-🙏)

<br>

# About 📋
`clipse` is a dependencyless, configurable TUI-based clipboard manager application. Though the app is optimised for a linux OS with a dedicated window manager, `clipse` can also be used on any Unix-based system. Simply install the package and bind the open command to get [your desired clipboard behavior](https://www.youtube.com/watch?v=ZE2F8Mj0_I0). Further instructions for setting this up can be found below.

[Click here to see a video demo for clipse](https://www.youtube.com/watch?v=ZE2F8Mj0_I0)

# Why use *clipse*? 

### 1. Configurability 🧰 

A customisable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default but is editable from a `.config/clipse/custom_theme.json` file that gets created when the program is run for the first time. Some example themes (based on my terminal)...

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

Simply leaving this file alone or setting the `useCustomTheme` value to `false` will give you a nice default theme... 

<p align="left">

  <img src="./resources/examples/default.png?raw=true" alt="Gruvbox" />

</p>

### 2. Usability ✨

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

### 3. Efficiency 💥

Due to Go's inbuilt garbage collection system and the way the application is built, `clipse` is pretty selfless when it comes to CPU consumption and memory. The below image shows how little resources are required to poser the background event listener used to continually update the history displayed in the TUI... 

<p align="left">

  <img src="./resources/examples/htop.png?raw=true" alt="Gruvbox" />

</p>

### 4. Versatility 🌐

The `clipse` binary, installable from the repo, can run on pretty much any Unix-based OS and will require zero external dependencies. Being terminal-based also allows for easy integration with a window manager and configuration of how the TUI behaves. For example, binding a floating window to the `clipse` command as shown in [my example](https://youtu.be/ZE2F8Mj0_I0) using [Hyprland window manager](https://hyprland.org/) on __NixOs__.

**Note that working with image files will require one of the following dependencies to be installd on your system**:

- Linux (X11) & MacOs: [xclip](https://github.com/astrand/xclip)
- Linux (Wayland): [wl-clipboard](https://github.com/bugaevc/wl-clipboard)

# Setup & installation 🏗️

See below for instructions on getting clipse installed at configured effectively. 

## Installation 

### Installing on NixOs

As a new package, `clipse` is still currently on the `Unstable` branch of `nixpkgs`. You can use the following methods to install...

**Direct install** __(recommended)__ 
```shell
nix-env -f channel:nixpkgs-unstable -iA clipse
```

**Nix shell**
```shell
nix shell github:NixOS/nixpkgs#clipse
```

**System package**
Building `clipse` as a system package may depend on your nix environemnt. I would suggest referencing this article for best practice. The derivation can also be built fomr source using the following: 
```c
{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule rec {
  pname = "clipse";
  version = "0.0.6";

  src = fetchFromGitHub {
    owner = "savedra1";
    repo = "clipse";
    rev = "v${version}";
    hash = "sha256-DLvYTPlLkp98zCzmbeL68B7mHl7RY3ee9rL30vYm5Ow=";
  };

  vendorHash = "sha256-GIUEx4h3xvLySjBAQKajby2cdH8ioHkv8aPskHN0V+w=";

  meta = {
    description = "Useful clipboard manager TUI for Unix";
    homepage = "https://github.com/savedra1/clipse";
    license = lib.licenses.mit;
    mainProgram = "clipse";
    maintainers = [ lib.maintainers.savedra1 ];
  };
}
```

### Installing with wget
 
**Linux arm64**:
```shell
wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.6_linux_arm64.tar.gz -O - | tar -xz 
```

**Linux amd64**:
```shell
wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.6_linux_amd64.tar.gz -O - | tar -xz 
```

**Linux 836**:
```shell
wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.6_linux_836.tar.gz -O - | tar -xz 
```

**Darwin arm64**:
```shell
wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.6_darwin_arm64.tar.gz -O - | tar -xz 
```

**Darwin amd64**:
```shell
wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.6_darwin_amd64.tar.gz -O - | tar -xz 
```


### Installing with Go

```shell

go install github.com/savedra1/clipse@v0.0.6

```

### Building from source 

```shell

git clone https://github.com/savedra1/clipse

cd clipse

go mod tidy

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

The above command creates a `nohup` process of `clipse --listen-shell`, which if called on its own will start a listener in your current terminal session instead. If `nohup` is not supported on your system, you can use your preferred method of running `clipse --listen-shell` in the background instead.

### Hyprland

Add the following lines to your Hyprland config file:

```shell

exec-once = clipse -listen                                                                 # run listener on startup

windowrulev2 = float,class:(floating)                                                      # ensure you have defined a floating window class

bind = SUPER, V, exec,  <terminal name> --class floating -e <shell-env>  -c 'clipse $PPID' # bind the open clipboard operation to a nice key. 

                                                                                           # Example: bind = SUPER, V, exec, alacritty --class floating -e zsh -c 'clipse $PPID'

```

[Hyprland reference](https://wiki.hyprland.org/Configuring/Window-Rules/)

### i3 

Add the following commands to your `.config/i3/config` file:

```shell

exec --no-startup-id clipse -listen                                                                 # run listener on startup

bindsym $mod+V exec --no-startup-id urxvt -e "$SHELL" -c "i3-msg 'floating enable' && clipse $PPID" # Bind floating shell with TUI selection to something nice 

``` 

[i3 reference](https://wiki.archlinux.org/title/i3)

### Sway

Add the following config to your `~/.config/sway/config` file:

```shell

exec clipse -listen                                                                                                                           # run the background listener on startup

bindsym $mod+V exec <terminal name> -e sh -c "swaymsg floating enable, move position center; swaymsg resize set 80ppt 80ppt && clipse $PPID"  # Bind floating shell with TUI selection to something nice

```
[Sway reference](https://wiki.archlinux.org/title/sway#:~:text=To%20enable%20floating%20windows%20or,enable%20floating%20windows%2Fwindow%20assignments.)

### Other

Every system/window manager is different and hard to determine exactly how to achieve the more ‘GUI-like’ behaviour. If using something not mentioned above, just refer to your systems documentation to find how to:

- Run the `clipse -listen` / `clipse --listen-shell` command on startup
- Bind the `clipse $PPID` command to a key that opens a terminal session (ideally in a window)

## All commands 💻

`clipse` is more than just a TUI. It also offers a number of CLI commands for managing clipboard content directly from the terminal. 

```shell

# Operational commands 

clipse -a <arg>       # Adds <arg> directly to the clipboard history without copying to system clipboard (string

clipse -a             # Adds any standard input directly to the clipboard history without copying to the system clipboard.

                      # For example: echo "some data" | clipse -a

clipse -c <arg>       # Copy the <arg> to the system clipboard (string). This also adds to clipboard history if currently listening. 

clipse -c             # Copies any standard input directly to the system clipboard.

                      # For example: echo "some data" | clipse -c

clipse -p             # Prints the current clipboard content to the console. 

                      # Example: clipse -p > file.txt

# TUI management commands

clipse $PPID          # Open Clipboard TUI

clipse -listen        # Run a background listener process

clipse --listen-shell # Run a listener process in the current terminal

clipse -help          # Display menu option

clipse -v             # Get version

clipse -clear         # Wipe all clipboard history and current system clipboard value

clipse -kill          # Kill any existing background processes

clipse                # Open Clipboard TUI in persistent/debug mode

```

You can view the full list of key bind commands when in the TUI by hitting the `?` key: 

<p align="left">

  <img src="./resources/examples/menu.png?raw=true" alt="Gruvbox" />

</p>

## How it works 🤔

When the app is run for the first time it creates a `/home/$USER/.config/clipse` dir with a `clipboard_hostory.json` file, a `custom_theme.json` file, and `tmp_files` folder for storing image data. After the `clipse -listen` command is executed, a background process will be watching for clipboard activity and adding any changes to the `clipboard_hstory.json` file. 

The TUI that displays the clipboard history should then be called with the `clipse $PPID` command. Passing in the terminal's PPID is irregular, but allows the terminal-based app to close itself from within the program itself, simulating the behavior of a full GUI without the memory overhead. A worthy trade-off in my opinion. 

Operations within the TUI are defined with the [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) framework, allowing for efficient concurrency and a smooth UX. `Delete` operations will remove the selected item from the TUI view and the storage file, `select` operations will copy the item to the systems clipboard and close the terminal window in which the session is currently hosted.  

The maximum item storage limit is currently hardcoded at **100**. However, there are plans to make this configurable in the future.

## Contributing 🙏

I would love to receive contributions to this project and welcome PRs from anyone and everyone. The following is a list of example future enhancements I'd like to implement:
- Pinned items view
- Image previews in TUI view rather than `<BINARY FILE>`
- Customisable max history limit
- System paste option (building functionality to paste the chosen item directly into the next place of focus after the TUI closes)
- Packages for apt, dnf, brew etc  
- Theme adjustments made available via CLI 
- Better debugging
- Use of a GUI library such as fyne/GIO (only with minimal CPU cost)
- Custom key binds added to config file 
- Customisable config file storage paths 

<br><br>
## TODO
Pinned/favourite items
Solution for auto-closing terminal window on MacOs
