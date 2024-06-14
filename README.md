<a href="https://github.com/savedra1/clipse/actions"><img src="https://github.com/charmbracelet/bubbletea/workflows/build/badge.svg" alt="Build Status"></a> [![Last Commit](https://img.shields.io/github/last-commit/savedra1/clipse)](https://github.com/savedra1/clipse) [![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/savedra1/clipse.svg?color=green)](https://github.com/savedra1/clipse/issues) <br>

https://github.com/savedra1/clipse/assets/99875823/c9c1998d-e96d-4c75-b7d9-060e39ac40ab

[![nix](https://img.shields.io/static/v1?label=Nix&message=0.0.71&color=blue)](https://search.nixos.org/packages?channel=unstable&show=clipse&from=0&size=50&sort=relevance&type=packages&query=clipse)
[![AUR](https://img.shields.io/aur/version/clipse.svg)](https://aur.archlinux.org/packages/clipse/)
<br>

### Table of contents

- [Why use clipse?](#why-use-clipse)
- [Installation](#installation)
- [Set up](#set-up)
- [Configuration](#configuration)
- [All commands](#all-commands-)
- [How it works](#how-it-works-)
- [Contributing](#contributing-)
- [FAQs](#faq)

### Release information

If moving to a new release of `clipse` please review the [changelog](https://github.com/savedra1/clipse/blob/main/CHANGELOG.md).


# About 📋
`clipse` is a configurable, TUI-based clipboard manager application written in Go with minimal dependency. Though the app is optimized for a Linux OS using a dedicated window manager, `clipse` can also be used on any Unix-based system. Simply install the package and bind the open command to get your desired clipboard behavior. Further instructions for setting this up can be found below.

### Dependency info and libraries used 
__[atotto/clipboard](https://github.com/atotto/clipboard)__

This requires a standard system clipboard, something like *one* of the following:
- wl-clipboard
- xclip
- xsel
- termux-clipboard

__[go-ps](github.com/mitchellh/go-ps)__

Does not require any additional dependency.

__[BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea)__

Does not require any additional dependency.

# Why use *clipse*? 

### 1. Configurability 🧰 

A customizable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default but is editable from a `.config/clipse/custom_theme.json` file that gets created when the program is run for the first time. Some example themes (based on my terminal)...

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

```json
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
    "StatusMsg": "#8FBCBB",
    "PinIndicatorColor": "#ff0000"
}

```

Simply leaving this file alone or setting the `useCustomTheme` value to `false` will give you a nice default theme... 

<p align="left">

  <img src="./resources/examples/default.png?raw=true" alt="Gruvbox" />

</p>

You can also easily specify source config like custom paths and max history limit in the apps `config.json` file. For more information see [Configuration](#configuration) section.  

### 2. Usability ✨

Easily recall, add, and delete clipboard history via a smooth TUI experience built with Go's excellent [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) library. A simple fuzzy finder, callable with the `/` key can easily match content from a theoretically unlimited amount of time in the past: 

<p align="left">

  <img src="./resources/examples/fuzzy.png?raw=true" alt="Gruvbox" />

</p>

Items can be pinned using the `p` key and toggled using `tab`, these pinned items will be not be removed from your history unless explicitly deleted with `backpace`/`x` in the TUI or `clipse -clear`. 

<p align="left">

  <img src="./resources/examples/pinned.png?raw=true" alt="Gruvbox" />

</p>

Content can also be added explicitly from the command-line using the `-a` or `-c` flags. This would allow you to easily pipe any CLI output directly into your history with commands like:

```shell

ls /home | clipse -a

``` 

```shell

clipse -a "a custom string value"

```

### 3. Efficiency 💥

Due to Go's inbuilt garbage collection system and the way the application is built, `clipse` is pretty selfless when it comes to CPU consumption and memory. The below image shows how little resources are required to run the background event listener used to continually update the history displayed in the TUI... 

<p align="left">

  <img src="./resources/examples/htop.png?raw=true" alt="Gruvbox" />

</p>

### 4. Versatility 🌐

The `clipse` binary, installable from the repository, can run on pretty much any Unix-based OS, though currently optimized for Linux. Being terminal-based also allows for easy integration with a window manager and configuration of how the TUI behaves. For example, binding a floating window to the `clipse` command as shown at the top of the page using [Hyprland window manager](https://hyprland.org/) on __NixOs__.

**Note that working with image files will require one of the following dependencies to be installed on your system**:

- Linux (X11) & macOS: [xclip](https://github.com/astrand/xclip)
- Linux (Wayland): [wl-clipboard](https://github.com/bugaevc/wl-clipboard)

# Setup & installation 🏗️

See below for instructions on getting clipse installed and configured effectively. 

## Installation 

### Installing on NixOs

Due to how irregularly the stable branch of Nixpkgs is updated, you may find the unstable package is more up to date. The Nix package for `clipse` can be found [here](https://search.nixos.org/packages?channel=24.05&from=0&size=50&sort=relevance&type=packages&query=clipse)

**Direct install**
```nix
nix-env -iA nixpkgs.clipse # OS == NixOs
nix-env -f channel:nixpkgs -iA clipse # OS != NixOs
```

**Nix shell**
```nix
nix shell -p clipse
```

**System package**
```nix
environment.systemPackages = [
    pkgs.clipse
  ];
```

If building `clipse` from the unstable branch as a system package, I would suggest referencing [this article](https://discourse.nixos.org/t/installing-only-a-single-package-from-unstable/5598) for best practice. The derivation can also be built from source using the following: 
```nix
{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule rec {
  pname = "clipse";
  version = "0.0.71";

  src = fetchFromGitHub {
    owner = "savedra1";
    repo = "clipse";
    rev = "v${version}";
    hash = "sha256-88GuYGJO5AgWae6LyMO/TpGqtk2yS7pDPS0MkgmJUQ4=";
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

### Installing on Arch
Shout out to [mcdenkijin](https://www.reddit.com/user/mcdenkijin/) for creating the AUR package!  

**Installing with yay**
```shell
yay -S clipse
```

**Installing from pkg source**
```shell
git clone https://aur.archlinux.org/clipse.git && cd clipse && makepkg -si
```

### Installing with wget

<details>
  <summary><b>Linux arm64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v0.0.6/clipse_0.0.71_linux_arm64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Linux amd64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v0.0.71/clipse_0.0.71_linux_amd64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Linux 836</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v0.0.71/clipse_0.0.71_linux_836.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Darwin arm64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v0.0.71/clipse_0.0.71_darwin_arm64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Darwin amd64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v0.0.71/clipse_0.0.71_darwin_amd64.tar.gz -O - | tar -xz
  ```
</details>


### Installing with Go

```shell

go install github.com/savedra1/clipse@v0.0.71

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

clipse

``` 

The second command doesn't need to be bound to a key combination, but rather to the system boot to run the background listener on start-up:

```shell

clipse -listen  

``` 

The above command creates a `nohup` process of `clipse --listen-shell`, which if called on its own will start a listener in your current terminal session instead. If `nohup` is not supported on your system, you can use your preferred method of running `clipse --listen-shell` in the background instead.

__Note: The following examples are based on bash/zsh shell environments. If you use something else like `foot` or `fish`, you may need to construct the command differently, referencing the relevant documentation.__

### Hyprland

Add the following lines to your Hyprland config file:

```shell

exec-once = clipse -listen                                                            # run listener on startup

windowrulev2 = float,class:(floating)                                                 # ensure you have defined a floating window class

bind = SUPER, V, exec,  <terminal name> --class floating -e <shell-env>  -c 'clipse'  # bind the open clipboard operation to a nice key. 

                                                                                      # Example: bind = SUPER, V, exec, alacritty --class floating -e zsh -c 'clipse'
```

[Hyprland reference](https://wiki.hyprland.org/Configuring/Window-Rules/)

### i3 

Add the following commands to your `.config/i3/config` file:

```shell

exec --no-startup-id clipse -listen                                                           # run listener on startup

bindsym $mod+V exec --no-startup-id urxvt -e "$SHELL" -c "i3-msg 'floating enable' && clipse" # Bind floating shell with TUI selection to something nice 

``` 

[i3 reference](https://wiki.archlinux.org/title/i3)

### Sway

Add the following config to your `~/.config/sway/config` file:

```shell

exec clipse -listen                                                                                                                     # run the background listener on startup

bindsym $mod+V exec <terminal name> -e sh -c "swaymsg floating enable, move position center; swaymsg resize set 80ppt 80ppt && clipse"  # Bind floating shell with TUI selection to something nice

```
[Sway reference](https://wiki.archlinux.org/title/sway#:~:text=To%20enable%20floating%20windows%20or,enable%20floating%20windows%2Fwindow%20assignments.)

### MacOs

The native terminal on MacOs will not close once the `clipse` program completes, even when using the `-fc` argument. You will therefore need to use a different terminal environment like [Alacritty](https://alacritty.org/) to achieve the "close on selection" effect. The bindings used to open the TUI will then need to be defined in your settings/window manager. 

### Other

Every system/window manager is different and hard to determine exactly how to achieve the more ‘GUI-like’ behavior. If using something not mentioned above, just refer to your systems documentation to find how to:

- Run the `clipse -listen` / `clipse --listen-shell` command on startup
- Bind the `clipse` command to a key that opens a terminal session (ideally in a window)

If you're not calling `clipse` with a command like `exec <terminal name> -e sh -c` and want to force the terminal window to close on selection of an item, use the `-fc` arg to pass in the `$PPID` variable so the program can force kill the shell session. EG `clipse -fc $PPID`. _Note that the $PPID variable is not available in every terminal environment, like fish terminal where you'd need to use $fish_pid instead._

## Configuration

The configuration capabilities of `clipse` will change as `clipse` evolves and grows. Currently, clipse supports the following configuration:
- Setting custom paths for:
  - The clipboard history file
  - The clipboard binaries directory (copied images and other binary data is stored in here)
  - The clipboard UI theme file
- Setting a custom max history limit 
- Custom themes

`clipse` looks for a base config file in `$HOME/.config/clipse/config.json`, and creates a default file if it does not find anything. The default config looks like this:
```json
{
    "historyFile": "clipboard_history.json",
    "maxHistory": 100,
    "themeFile": "custom_theme.json",
    "tempDir": "tmp_files"
}
```

Note that all the paths provided (the theme, `historyFile`, and `tempDir`) are all relative paths. They are relative to the location of the config file that holds them. Thus, a file `config.json` at location `$HOME/.config/clipse/config.json` will have all relative paths defined in it relative to its directory of `$HOME/.config/clipse/`.

Absolute paths starting with `/`, paths relative to the user home dir using `~`, or any environment variables like `$HOME` and `$XDG_CONFIG_HOME` are also valid paths that can be used in this file instead.

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

clipse -fc $PPID      # Open Clipboard TUI in 'force kill' mode 

clipse -listen        # Run a background listener process

clipse --listen-shell # Run a listener process in the current terminal (useful for debugging)

clipse -help          # Display menu option

clipse -v             # Get version

clipse -clear         # Wipe all clipboard history and current system clipboard value

clipse keep           # Keep the TUI open after selecting an item to copy (useful for debugging)

clipse -kill          # Kill any existing background processes

clipse                # Open Clipboard TUI in persistent/debug mode

```

You can view the full list of key bind commands when in the TUI by hitting the `?` key: 

<p align="left">

  <img src="./resources/examples/menu.png?raw=true" alt="Gruvbox" />

</p>

## How it works 🤔

When the app is run for the first time it creates a `/home/$USER/.config/clipse` dir containing `config.json`, `clipboard_history.json`, `custom_theme.json` and a dir called `tmp_files` for storing image data. After the `clipse -listen` command is executed, a background process will be watching for clipboard activity and adding any changes to the `clipboard_history.json` file, unless a different path is specified in `config.json`.

The TUI that displays the clipboard history with the defined theme should then be called with the `clipse` command. Operations within the TUI are defined with the [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) framework, allowing for efficient concurrency and a smooth UX. `delete` operations will remove the selected item from the TUI view and the storage file, `select` operations will copy the item to the systems clipboard and exit the program.

The maximum item storage limit defaults at **100** but can be customized to anything you like in the `config.json` file.

## Contributing 🙏

I would love to receive contributions to this project and welcome PRs from anyone and everyone. The following is a list of example future enhancements I'd like to implement:
- [ ] Image previews in TUI view
- [x] ~~Pinned items~~
- [ ] Customisations for:
  - [x] ~~max history limit~~
  - [x] ~~config file paths~~
  - [ ] key bindings
- [ ] Auto-forget feature based on where the text was copied
- [ ] Multi-select feature for copying multiple items at once
- [ ] Categorized pinned items with _potentially_ different tabs/views  
- [ ] System paste option _(building functionality to paste the chosen item directly into the next place of focus after the TUI closes)_
- [ ] Packages for apt, dnf, brew etc
- [ ] Theme/config adjustments made available via CLI
- [ ] Better debugging / debug mode _(eg `clipse --debug` / debug file / system alert on panic)_
- [ ] Use of a GUI library such as fyne/GIO (only with minimal CPU cost)
- [ ] Cross compile binaries for `wl-clipboard`/`xclip` to remove dependency
- [x] TUI / theming enhancements:
  - [x] ~~Menu theme~~
  - [x] ~~Filter theme~~
  - [x] ~~Clear TUI view on select and close _(mirror close effect from `q` or `esc`)_~~
- [ ] Private mode _(eg `clipse --pause 1` )_

## FAQ

- __My terminal window does not close on selection, even when using `clipse -fc $PPID`__ - _Some terminal environments reference system variables differently. For example, the fish terminal will need to use `$fish_pid` instead. To debug this error you can run `echo $PPID` to see what gets returned. The 'close on selection functionality is also not currently available for macOS as killing the terminals ppid does not close the window - it seems AppleScript is needed to achieve this._

- __Is there risk of multiple parallel processes running?__ - _No. The `clipse` command kills any existing TUI processes before opening up and the `clipse -listen`  command kills any existing background listeners before starting a new one._

- __High CPU usage?__ - When an image file has an irregular binary data pattern it can cause a lot of strain on the program when `clipse` reads its history file (even when the TUI is not open). If this happens, you will need to remove the image file from the TUI or buy using `clipse -clear` or `clipse -clear-images` to remove all files if you don't want to spend the time figuring out which one is causing the issue. See issue #33 for an example.

- __My copied entries are not recorded when starting the clipse listener on boot with a systemd service__ - There may be a few ways around this but the workaround discovered in issue #41 was to use a `.desktop` file, stored in `~/.config/autostart/`. Eg:
  ```shell
  [Desktop Entry]
  Name=clipse
  Comment=Clipse event listener autostart.
  Exec=/home/brayan/Applications/bin/clipse/clipse_0.0.71_linux_amd64/clipse --listen %f
  Terminal=false
  Type=Application
  ```

<br>
