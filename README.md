<a href="https://github.com/savedra1/clipse/actions"><img src="https://github.com/charmbracelet/bubbletea/workflows/build/badge.svg" alt="Build Status"></a> [![Last Commit](https://img.shields.io/github/last-commit/savedra1/clipse)](https://github.com/savedra1/clipse) [![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/savedra1/clipse.svg?color=green)](https://github.com/savedra1/clipse/issues) <br>

https://github.com/savedra1/clipse/assets/99875823/40af797c-2297-49b5-88ec-b8c04e8c829b

[![nix](https://img.shields.io/static/v1?label=Nix&message=1.0.9&color=blue)](https://search.nixos.org/packages?channel=unstable&show=clipse&from=0&size=50&sort=relevance&type=packages&query=clipse)
[![AUR](https://img.shields.io/aur/version/clipse.svg)](https://aur.archlinux.org/packages/clipse/)
<br>

### Table of contents

- [Features](#features)
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

This requires a system clipboard. I would recommend using `wl-clipboard` (Wayland) or `xclip` (X11/macOs) to get the best results. You can also use `xsel` and `termux-clipboard`, but these will not allow you to copy images. 


__[BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea)__

Does not require any additional dependency, but may require you to use a terminal environment that's compatible with [termenv](https://github.com/muesli/termenv).

# Features ✨

- Persistent history
- Supports text and image
- [Customizable UI theme](#Customization)
- [Customizable file paths](#Configuration)
- [Customizable maximum history limit](#Configuration)
- Filter items using a fuzzy find
- Image and text previews
- Mult-selection of items for copy and delete operations
- Bulk copy all active filter matches
- Pin items/pinned items view 
- Vim-like keybindings for navigation available
- [Run on any Unix machine](#Versatility) with single binary for the clipbboard monitor and TUI 
- Optional duplicates

### Customization 🧰 

A customizable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default but is editable from a `custom_theme.json` file that gets created when the program is run for the first time. See the [library](https://github.com/savedra1/clipse/blob/main/resources/library.md) for some example themes to give you inspiration.

An example `custom_theme.json` file: 

```json
{
	"UseCustom":          true,
	"TitleFore":          "#ffffff",
	"TitleBack":          "#6F4CBC",
	"TitleInfo":          "#3498db",
	"NormalTitle":        "#ffffff",
	"DimmedTitle":        "#808080",
	"SelectedTitle":      "#FF69B4",
	"NormalDesc":         "#808080",
	"DimmedDesc":         "#808080",
	"SelectedDesc":       "#FF69B4",
	"StatusMsg":          "#2ecc71",
	"PinIndicatorColor":  "#FFD700",
	"SelectedBorder":     "#3498db",
	"SelectedDescBorder": "#3498db",
	"FilteredMatch":      "#ffffff",
	"FilterPrompt":       "#2ecc71",
	"FilterInfo":         "#3498db",
	"FilterText":         "#ffffff",
	"FilterCursor":       "#FFD700",
	"HelpKey":            "#999999",
	"HelpDesc":           "#808080",
	"PageActiveDot":      "#3498db",
	"PageInactiveDot":    "#808080",
	"DividerDot":         "#3498db",
	"PreviewedText":      "#ffffff",
	"PreviewBorder":      "#3498db",
}
``` 

You can also easily specify source config like custom paths and max history limit in the apps `config.json` file. For more information see [Configuration](#configuration) section.  

### Versatility 🌐

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
  version = "1.0.9";

  src = fetchFromGitHub {
    owner = "savedra1";
    repo = "clipse";
    rev = "v${version}";
    hash = "sha256-Kpe/LiAreZXRqh6BHvUIn0GcHloKo3A0WOdlRF2ygdc=";
  };

  vendorHash = "sha256-Hdr9NRqHJxpfrV2G1KuHGg3T+cPLKhZXEW02f1ptgsw=";

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
Shout out to [@raininja](https://github.com/raininja) for creating and maintaining the AUR package!  

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
  wget -c https://github.com/savedra1/clipse/releases/download/v1.0.9/clipse_1.0.9_linux_arm64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Linux amd64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v1.0.9/clipse_1.0.9_linux_amd64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Linux 836</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v1.0.9/clipse_1.0.9_linux_836.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Darwin arm64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v1.0.9/clipse_1.0.9_darwin_arm64.tar.gz -O - | tar -xz
  ```
</details>

<details>
  <summary><b>Darwin amd64</b></summary>

  ```shell
  wget -c https://github.com/savedra1/clipse/releases/download/v1.0.9/clipse_1.0.9_darwin_amd64.tar.gz -O - | tar -xz
  ```
</details>


### Installing with Go

```shell

go install github.com/savedra1/clipse@v1.0.9

```

### Building from source 

```shell

git clone https://github.com/savedra1/clipse

cd clipse

go mod tidy

go build -o clipse

```

## Set up

As mentioned earlier, to get the most out of `clipse`, it's recommended to bind the two primary key commands to your system's config. The first key command is to open the clipboard history TUI:

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

exec-once = clipse -listen # run listener on startup

windowrulev2 = float,class:(clipse) # ensure you have a floating window class set if you want this behaviour
windowrulev2 = size 622 652,class:(clipse) # set the size of the window as necessary

bind = SUPER, V, exec,  <terminal name> --class clipse -e 'clipse' 

# Example: bind = SUPER, V, exec, alacritty --class clipse -e 'clipse'
```

[Hyprland reference](https://wiki.hyprland.org/Configuring/Window-Rules/)

### i3 

Add the following commands to your `~/.config/i3/config` file:

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
  - The debug log file
  - The clipboard UI theme file
- Setting a custom max history limit 
- Custom themes
- If duplicates are allowed

`clipse` looks for a base config file in `$CONFIG_DIR/clipse/config.json` _(`$CONFIG_DIR` being `$XDG_DATA_HOME` or `$HOME/.config`)_, and creates a default file if it does not find anything. The default config looks like this:
```json
{
    "historyFile": "clipboard_history.json",
    "maxHistory": 100,
    "allowDuplicates": false,
    "themeFile": "custom_theme.json",
    "tempDir": "tmp_files",
    "logFile": "clipse.log"
}
```

Note that all the paths provided (the theme, `historyFile`, and `tempDir`) are all relative paths. They are relative to the location of the config file that holds them. Thus, a file `config.json` at location `$HOME/.config/clipse/config.json` will have all relative paths defined in it relative to its directory of `$HOME/.config/clipse`.

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

clipse                # Open Clipboard TUI in persistent/debug mode

clipse -fc $PPID      # Open Clipboard TUI in 'force kill' mode 

clipse -listen        # Run a background listener process

clipse --listen-shell # Run a listener process in the current terminal (useful for debugging)

clipse -help          # Display menu option

clipse -v             # Get version

clipse -clear         # Wipe all clipboard history except for pinned items

clipse -clear-images  # Wipe all images from the history 

clipse -clear-text    # Wipe all text items from the clipboard history

clipse -clear-all     # Wipe entire clipboard history

clipse keep           # Keep the TUI open after selecting an item to copy (useful for debugging)

clipse -kill          # Kill any existing background processes

```

You can also view the full list of TUI key commands by hitting the `?` key when the `clipse` UI is open. 


## How it works 🤔

When the app is run for the first time it creates a `/home/$USER/.config/clipse` dir containing `config.json`, `clipboard_history.json`, `custom_theme.json` and a dir called `tmp_files` for storing image data. After the `clipse -listen` command is executed, a background process will be watching for clipboard activity and adding any changes to the `clipboard_history.json` file, unless a different path is specified in `config.json`.

The TUI that displays the clipboard history with the defined theme should then be called with the `clipse` command. Operations within the TUI are defined with the [BubbleTea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) framework, allowing for efficient concurrency and a smooth UX. `delete` operations will remove the selected item from the TUI view and the storage file, `select` operations will copy the item to the systems clipboard and exit the program.

The maximum item storage limit defaults at **100** but can be customized to anything you like in the `config.json` file.

## Contributing 🙏

I would love to receive contributions to this project and welcome PRs from everyone. The following is a list of example future enhancements I'd like to implement:

- [x] ~~Image previews in TUI view~~
- [x] ~~Pinned items~~
- [x] ~~Warn on deleting pinned items~~
- [x] ~~Color theme customizations for all UI elements~~
- [ ] Customisations for:
  - [x] ~~max history limit~~
  - [x] ~~config file paths~~
  - [x] ~~Duplicates allowed~~ 
  - [ ] key bindings
- [x] ~~Option to disable duplicate items~~
- [ ] Auto-forget based on where the text was copied
- [x] ~~Multi-select feature for copying multiple items at once~~
- [ ] Categorized pinned items with _potentially_ different tabs/views  
- [ ] System paste option _(building functionality to paste the chosen item directly into the next place of focus after the TUI closes)_
- [ ] Packages for:
  - [ ] apt 
  - [ ] dnf 
  - [ ] brew
  - [ ] other
- [ ] Theme/config adjustments made available via CLI
- [ ] Your custom theme for the [library](https://github.com/savedra1/clipse/blob/main/resources/library.md)
- [ ] debug mode _(eg `clipse --debug` / debug file / system alert on panic)_
- [ ] Cross compile binaries for `wl-clipboard`/`xclip` to remove dependency
- [x] TUI / theming enhancements:
  - [x] ~~Menu theme~~
  - [x] ~~Filter theme~~
  - [x] ~~Clear TUI view on select and close _(mirror close effect from `q` or `esc`)_~~
- [ ] Private mode _(eg `clipse --pause 1` )_

## FAQ

__My terminal window does not close on selection, even when using `clipse -fc $PPID`__ 

Some terminal environments reference system variables differently. For example, the fish terminal will need to use `$fish_pid` instead. To debug this error you can run `echo $PPID` to see what gets returned. To get the "close on selection" effect for macOs, you will need to install a different terminal environment like `Alacritty`._
<br>

__Is there risk of multiple parallel processes running?__
 
_No. The `clipse` command kills any existing TUI processes before opening up and the `clipse -listen`  command kills any existing background listeners before starting a new one._
<br>

__High CPU usage?__ 

When an image file has an irregular binary data pattern it can cause a lot of strain on the program when `clipse` reads its history file (even when the TUI is not open). If this happens, you will need to remove the image file from the TUI or by using `clipse -clear-images`. See issue #33 for an example.
<br>

__My copied entries are not recorded when starting the clipse listener on boot with a systemd service__ 
 
There may be a few ways around this but the workaround discovered in issue #41 was to use a `.desktop` file, stored in `~/.config/autostart/`. Eg:
  ```shell
  [Desktop Entry]
  Name=clipse
  Comment=Clipse event listener autostart.
  Exec=/home/usrname/Applications/bin/clipse/clipse_1.0.9_linux_amd64/clipse --listen %f
  Terminal=false
  Type=Application
  ```
<br>

__Copying images from a browser does not work correctly__  

Depending on the clipboard utility you are using (`wl-clipboard`/`xclip` etc) the data sent to the system clipboard is read differently when copying from browser locations. 
<br>
If using `wayland`, copying images from your browser should now work from most sites if using `v1.0.4` or later. This may copy the binary data as well as the metadata sting as a separate entry. Some sites/browsers may add the browser image data to the stdin in a way that `wl-clipboard` does not recognize.
<br>
If using `x11`, `MacOs` or other and copying browser images does not work, feel free to raise and issue (or a PR) detailing which sites/browser engines this does not work with for you.
  
<br>
