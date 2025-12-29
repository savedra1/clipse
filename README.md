[![Build](https://img.shields.io/github/actions/workflow/status/savedra1/clipse/go-test.yml)](https://github.com/savedra1/clipse/actions)
[![Last Commit](https://img.shields.io/github/last-commit/savedra1/clipse)](https://github.com/savedra1/clipse)
[![GitHub closed issues](https://img.shields.io/github/issues-closed-raw/savedra1/clipse.svg?color=green)](https://github.com/savedra1/clipse/issues)
<br>

<https://github.com/savedra1/clipse/assets/99875823/40af797c-2297-49b5-88ec-b8c04e8c829b>

[![nix](https://img.shields.io/static/v1?label=Nix&message=1.1.0&color=blue)](https://search.nixos.org/packages?channel=unstable&show=clipse&from=0&size=50&sort=relevance&type=packages&query=clipse)
[![AUR](https://img.shields.io/aur/version/clipse.svg)](https://aur.archlinux.org/packages/clipse/)
<br>
---

### Table of contents

- [Features](#features)
- [Installation & Set-up](#Installation-&-Set-up)
- [Configuration](#configuration)
- [All commands](#all-commands-)
- [Contributing](#contributing-)
- [FAQs](#faq)
---

### Release information

If moving to a new release of `clipse` please review the [changelog](https://github.com/savedra1/clipse/blob/main/CHANGELOG.md).

# About 📋

`clipse` is a configurable clipboard manager application written in Go with minimal dependency. The application is optimized for a Linux OS using a dedicated window manager, but `clipse` can also be used on any Unix-based system. Simply install the package and bind the open command to get your desired clipboard behavior. See the [Installation & Set-up](#Installation-&-Set-up) section for examples of this.

## Dependencies

**Wayland**

If using a Wayland-based display server, you must have [wl-clipboard](https://github.com/bugaevc/wl-clipboard) installed to handle both text and image data. To enable the auto-paste feature, you may also need to allow persistent permission to the `/dev/uinput` device. See [Configuration.Auto-Paste](#Auto-paste).

**X11**

For X11 builds, the [robotgo](https://github.com/go-vgo/robotgo) library is used to implement the auto-paste feature. Please check [this section](https://github.com/go-vgo/robotgo#for-everything-else) of their docs for the runtime dependencies you may need available.

**Darwin**

There are no known dependencies for Darwin.

# Features ✨

- Text and image support
- Highly customizable
- Image and text previews
- Item pinning
- Fuzzy filtering
- Multi-select
- Auto-paste
- Excluded apps/windows
- Bulk copy/output
- Portable (runs on any wayland/x11/darwin machine)
- CLI operations

# Installation & Set-up 🏗️

## Installation

<details>
  <summary><b>NixOS</b></summary>

  Due to how irregularly the stable branch of Nixpkgs is updated, you may find the unstable package is more up to date. The Nix package for `clipse` can be found [here](https://search.nixos.org/packages?channel=25.05&from=0&size=50&sort=relevance&type=packages&query=clipse)

  __Direct install__

  ```nix
  nix-env -iA nixpkgs.clipse # OS == NixOs
  nix-env -f channel:nixpkgs -iA clipse # OS != NixOs
  ```

  __Nix shell__

  ```nix
  nix shell -p clipse
  ```

  __System package__

  ```nix
  environment.systemPackages = [
      pkgs.clipse
    ];
```
</details>

<details>
  <summary><b>Arch</b></summary>

  Thank you [@raininja](https://github.com/raininja) for creating and maintaining the [AUR package](https://aur.archlinux.org/packages/clipse)!

  __Installing with yay__

  ```shell
  yay -S clipse
  ```

  __Installing from pkg source__

  ```shell
  git clone https://aur.archlinux.org/clipse.git && cd clipse && makepkg -si
  ```

</details>

<details>
  <summary><b>Fedora/Rhel</b></summary>

  Thank you [@RadioAndrea](https://github.com/RadioAndrea) for creating and maintaining the [COPR package](https://copr.fedorainfracloud.org/coprs/azandure/clipse/)!

  ```shell
  dnf copr enable azandure/clipse
  ```

</details>

<details>
  <summary><b>Installing with Go</b></summary>

  ```shell
  go install github.com/savedra1/clipse@v1.2.0
  ```

</details>

<details>
  <summary><b>Building from source</b></summary>

  Building `clipse` from source requires different build tags depending on your system's display server. To ensure the correct tags are uses, please use the Makefile to build the executable. E.g. `make wayland`

  ```shell
  git clone https://github.com/savedra1/clipse
  cd clipse
  go mod tidy
  make x11/darwin/wayland
  ```

  Once you have build the binary, you can install this to your executable path, E.g. `install -m 755 clipse /usr/bin || mv clipse /bin/`.

</details>


## Set up

To get the most out of `clipse`, it's recommended to bind the two primary key commands to your system's config. The first key command is to open the clipboard history TUI:

```shell
clipse
```

The second command doesn't need to be bound to a key combination, but rather to the system boot to run the background listener process on start-up:

```shell
clipse -listen
```

### Example window manager configurations

<details>
  <summary><b>Hyprland</b></summary>

  Add the following lines to your Hyprland config file:

  ```shell

  exec-once = clipse -listen # run listener on startup

  windowrulev2 = float,class:(clipse) # ensure you have a floating window class set if you want this behavior
  windowrulev2 = size 622 652,class:(clipse) # set the size of the window as necessary

  bind = SUPER, V, exec,  <terminal name> --class clipse -e 'clipse'

  # Example: bind = SUPER, V, exec, alacritty --class clipse -e 'clipse'
  ```

  [Hyprland reference](https://wiki.hypr.land/Configuring/Window-Rules/)

</details>

<details>
  <summary><b>i3</b></summary>

  Add the following commands to your `~/.config/i3/config` file:

  ```shell
  exec --no-startup-id clipse -listen                                                           # run listener on startup
  bindsym $mod+V exec --no-startup-id urxvt -e "$SHELL" -c "i3-msg 'floating enable' && clipse" # Bind floating shell with TUI selection to something nice
  ```

  [i3 reference](https://wiki.archlinux.org/title/i3)

</details>

<details>
  <summary><b>Sway</b></summary>

  Add the following config to your `~/.config/sway/config` file:

  ```shell
  exec clipse -listen                                                                        # run the background listener on startup

  for_window [app_id="clipse"] floating enable, move position center, resize set 80ppt 80ppt # style window to look nice

  bindsym $mod+V exec <terminal name> --class clipse -e clipse                               # bind floating shell with clipse TUI

  # Example: bindsym $mod+V exec alacritty --class clipse -e clipse
  ```

  [Sway reference](https://wiki.archlinux.org/title/Sway#Floating_windows)

</details>

<details>
  <summary><b>macOS</b></summary>

  The same concept applies to macOS, where the listen command is bound to the WM config. A separate application, like [skhd](https://github.com/asmvik/skhd) may be required to configure the custom keybind. See the below example using `yabai` WM and `skhd`.

  In the `yabairc` config:

  ```shell
  # yabairc

  clipse -listen
  ...
  yabai -m rule --add title="^Clipse$" manage=off

  ```

  In the `skhdrc` config:

  ```shell
  # skhdrc
  alt - v : alacritty -t "Clipse" --option window.dimensions.lines=40 --option window.dimensions.columns=70 -e clipse
  ```

  ---

  Without a WM, you can still achieve similar behavior by creating a listener service manually. E.g:

  Create the file `~/Library/LaunchAgents/clipse.plist` with the following content:
  ```
  <?xml version="1.0" encoding="UTF-8"?>
  <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"\>
  <plist version="1.0">
  <dict>
      <key>Label</key>
      <string>com.savedra1.clipse</string>
      <key>ProgramArguments</key>
      <array>
          <string>/path/to/clipse</string>
          <string>--listen-shell</string>
      </array>
      <key>RunAtLoad</key>
      <true/>
  </dict>
  </plist>
  ```

  Then in a terminal, activate the agent with: `launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/clipse.plist`. Right away and after your next login, you can check that `clipse` is running by executing `ps -e | grep '[c]lipse'`.

</details>

<br>

# Configuration

Your `configuration.json` file will initially be created with the following default values:

```json
{
    "allowDuplicates": false,
    "historyFile": "clipboard_history.json",
    "maxHistory": 100,
    "deleteAfter": 0,
    "logFile": "clipse.log",
    "pollInterval": 50,
    "maxEntryLength": 65,
    "themeFile": "custom_theme.json",
    "tempDir": "tmp_files",
    "keyBindings": {
        "choose": "enter",
        "clearSelected": "S",
        "down": "down",
        "end": "end",
        "filter": "/",
        "forceQuit": "Q",
        "home": "home",
        "more": "?",
        "nextPage": "right",
        "prevPage": "left",
        "preview": " ",
        "quit": "esc",
        "remove": "backspace",
        "selectDown": "shift+down",
        "selectSingle": "s",
        "selectUp": "shift+up",
        "togglePin": "p",
        "togglePinned": "tab",
        "up": "up",
        "yankFilter": "ctrl+s"
    },
    "imageDisplay": {
        "type": "basic",
        "scaleX": 9,
        "scaleY": 9,
        "heightCut": 2
    },
    "excludedApps": [
        "1Password",
        "Bitwarden",
        "KeePassXC",
        "LastPass",
        "Dashlane",
        "Password Safe",
        "Keychain Access"
    ],
    "autoPaste": {
        "enabled": false,
        "keybind": "ctrl+v",
        "buffer": 10
    },
    "enableMouse": true,
    "enableDescription": true
}
```

__If any values from this file are removed, they will not be readded when the program runs, but the default values will be used.__

## General Configuration

| Option              | Type   | Description                                                                                  |
| ------------------- | ------ | -------------------------------------------------------------------------------------------- |
| `allowDuplicates`   | bool   | Whether identical clipboard entries are allowed to appear more than once in history.         |
| `historyFile`       | string | File path used to persist clipboard history.                                                 |
| `maxHistory`        | int    | Maximum number of clipboard entries stored in history.                                       |
| `deleteAfter`       | int    | Time (in seconds) after which entries are automatically deleted; `0` disables auto-deletion. |
| `logFile`           | string | File path for application logs.                                                              |
| `pollInterval`      | int    | Clipboard polling interval in milliseconds (x11/Darwin).                                     |
| `maxEntryLength`    | int    | Maximum number of characters shown per clipboard entry.                                      |
| `themeFile`         | string | Path to the custom theme configuration file.                                                 |
| `tempDir`           | string | Directory used for image files.                                                              |
| `enableMouse`       | bool   | Enables mouse interaction in the UI.                                                         |
| `enableDescription` | bool   | Shows additional descriptive text for clipboard entries.                                     |
| `keyBindings`       | map    | Custom keybind definitions.                                                                  |
| `autoPaste`         | map    | Auto-paste options.                                                                          |
| `imageDisplay`      | map    | Image display options (basic/kitty/sixel).                                                   |
| `excludedApps`      | list   | List of App/Window names form which to exclude any copied data                               |

All the paths provided (the theme, `historyFile`, and `tempDir`) are all relative paths. They are relative to the location of the config file that holds them. E.g, a file `config.json` at location `$HOME/.config/clipse/config.json` will have all relative paths defined in it relative to its directory of `$HOME/.config/clipse`.

Absolute paths starting with `/`, paths relative to the user home dir using `~`, or any environment variables like `$HOME` and `$XDG_CONFIG_HOME` are also valid paths that can be used in this file instead.


## Key Bindings

| Option                      | Type   | Description                                 |
| --------------------------- | ------ | ------------------------------------------- |
| `keyBindings.choose`        | string | Confirms and selects the highlighted entry. |
| `keyBindings.clearSelected` | string | Clears all currently selected entries.      |
| `keyBindings.down`          | string | Moves selection down by one entry.          |
| `keyBindings.end`           | string | Jumps to the last entry in the list.        |
| `keyBindings.filter`        | string | Activates filtering/search mode.            |
| `keyBindings.forceQuit`     | string | Immediately exits the application.          |
| `keyBindings.home`          | string | Jumps to the first entry in the list.       |
| `keyBindings.more`          | string | Shows additional help or options.           |
| `keyBindings.nextPage`      | string | Moves to the next page of entries.          |
| `keyBindings.prevPage`      | string | Moves to the previous page of entries.      |
| `keyBindings.preview`       | string | Toggles preview of the selected entry.      |
| `keyBindings.quit`          | string | Exits the application gracefully.           |
| `keyBindings.remove`        | string | Deletes the selected entry.                 |
| `keyBindings.selectDown`    | string | Extends selection downward.                 |
| `keyBindings.selectSingle`  | string | Selects a single entry.                     |
| `keyBindings.selectUp`      | string | Extends selection upward.                   |
| `keyBindings.togglePin`     | string | Pins or unpins the selected entry.          |
| `keyBindings.togglePinned`  | string | Toggles display of pinned entries.          |
| `keyBindings.up`            | string | Moves selection up by one entry.            |
| `keyBindings.yankFilter`    | string | Copies the current filter text.             |



## Image Display

| Option                   | Type   | Description                                                                   |
| ------------------------ | ------ | ----------------------------------------------------------------------------- |
| `imageDisplay.type`      | string | Rendering mode used for displaying images. Allowed options: basic|kitty|sixel |
| `imageDisplay.scaleX`    | int    | Horizontal scaling factor for images.                                         |
| `imageDisplay.scaleY`    | int    | Vertical scaling factor for images.                                           |
| `imageDisplay.heightCut` | int    | Number of rows trimmed from image height.                                     |

Currently these are the supported options for `imageDisplay.type`:
 - `basic`
 - `kitty`
 - `sixel`

 The `scaleX` and `scaleY` options are the scaling factors for the images. Depending on the situation, you need to find suitable numbers to ensure the images are displayed correctly and completely. You can make adjustments based on [this implementation](https://github.com/savedra1/clipse/pull/138#issue-2530565414).


## Auto-paste

| Option              | Type   | Description                                                |
| ------------------- | ------ | ---------------------------------------------------------- |
| `autoPaste.enabled` | bool   | Enables automatic pasting after selecting an entry.        |
| `autoPaste.keybind` | string | Key combination used to trigger paste. (E.g. cmd+v)        |
| `autoPaste.buffer`  | int    | Delay buffer (in milliseconds) before auto-paste executes. |

When enabling the auto-paste feature, you may need to allow `clipse` the required permissions to the system keyboard, and ensure the relevant system APIs are available. The requirements are different based on your display server.

<details>
  <summary><b>Wayland</b></summary>

  On wayland, the only auto-paste requirement is for `clipse` to have access to the `/dev/uinput` device. There are a number of ways to do this. E.g. on NixOS you can add the following to your `configuration.nix`:

  ```nix
  # User Configurations
  users.users.${userConfig.username} = {
    isNormalUser = true;
    home = userConfig.homeDirectory;
    shell = pkgs.zsh; # Setting Zsh as the default shell
    extraGroups = [ "wheel" "networkmanager" "input" ];
  };

  # Create udev rule for uinput access
  services.udev.extraRules = ''
    KERNEL=="uinput", MODE="777", GROUP="input", OPTIONS+="static_node=uinput"
  '';
  ```

  To do this manually you can do something like:

  ```shell
  sudo groupadd input
  sudo usermod -aG input <username>
  sudo vi /etc/udev/rules.d/99-uinput.rules --> add 'KERNEL=="uinput", GROUP="input", MODE="0660"'
  sudo udevadm control --reload-rules
  sudo udevadm trigger
  ```

</details>

<details>
  <summary><b>X11</b></summary>

  As mentioned in [Dependencies](#dependencies), X11 builds utilize the [robotgo](https://github.com/go-vgo/robotgo) library to implement auto-paste. This shouldn't require any build dependencies, like GCC and Go, but certain `xlib` APIs may need to be installed in they are not already present Please see [this section](https://github.com/go-vgo/robotgo#for-everything-else) of the `robotgo` docs for more information.

</details>

<details>
  <summary><b>Darwin</b></summary>

  The only requirement for enabling auto-paste on Darwin could be to allow the `clipse` executable permissions to the system keyboard. This is usually not required, but you can do this via the System Settings GUI: `System Settings -> Accessibility -> '+' -> Select clipse binary file`.

</details>


## Theme

A customizable TUI allows you to easily match your system's theme. The app is based on your terminal's theme by default but is editable from the file specified under for `themeFile` (defaults to `custom_theme.json`). See the [library](https://github.com/savedra1/clipse/blob/main/.github/.resources/library.md) for some example themes to give you inspiration.

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

# All commands 💻

Additionally to the `clipse` TUI, there are a number of CLI commands for managing clipboard operations directly from the terminal.

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

clipse                   # Open Clipboard TUI in persistent/debug mode

clipse -listen           # Run a background listener process

clipse --listen-shell    # Run a listener process in the current terminal (useful for debugging)

clipse -help             # Display menu option

clipse -v                # Get version

clipse -clear            # Wipe all clipboard history except for pinned items

clipse -clear-images     # Wipe all images from the history

clipse -clear-text       # Wipe all text items from the clipboard history

clipse -clear-all        # Wipe entire clipboard history

clipse -clean            # sanitize existing text entries and remove orphaned image entries

clipse keep              # Keep the TUI open after selecting an item to copy (useful for debugging)

clipse -kill             # Kill any existing background processes

clipse -pause <arg>      # Pause clipboard monitoring for a specified duration. Example: `clipse -pause 5m` pauses for 5 minutes

clipse -output-all       # Print the entire clipboard history to stdout

clipse -enable-real-time # Enable real time updates to the TUI

```

You can also view the full list of TUI key commands by hitting the `?` key when the `clipse` UI is open.

---

## Contributing 🙏

Please see the following for what contribution suggestions. If you have an idea that's not listed, please create an issue/discussion around this first.

- [x] ~~Image previews in TUI view~~
- [x] ~~Pinned items~~
- [x] ~~Warn on deleting pinned items~~
- [x] ~~Color theme customizations for all UI elements~~
- Customizations for:
  - [x] ~~max history limit~~
  - [x] ~~config file paths~~
  - [x] ~~Duplicates allowed~~
  - [x] ~~key bindings~~
  - [x] ~~image preview display render type~~
- [x] ~~Option to disable duplicate items~~
- [x] ~~Multi-select feature for copying multiple items at once~~
- [ ] Categorized pinned items with _potentially_ different tabs/views
- [x] ~~System paste option _(building functionality to paste the chosen item directly into the next place of focus after the TUI closes)_~~
- Packages for:
  - [ ] apt
  - [x] ~~dnf~~
  - [ ] brew
  - [ ] other
- [ ] Your custom theme for the [library](https://github.com/savedra1/clipse/blob/main/resources/library.md)
- [ ] debug mode _(eg `clipse --debug` / debug file / system alert on panic)_
- [x] TUI / theming enhancements:
  - [x] ~~Menu theme~~
  - [x] ~~Filter theme~~
  - [x] ~~Clear TUI view on select and close _(mirror close effect from `q` or `esc`)_~~
- [x] ~~Private mode _(eg `clipse --pause 1` )_~~

---

## FAQ

<details> <summary><b>Clipse crashes when I enter certain characters into the search bar</b></summary>

See issue [#148](https://github.com/savedra1/clipse/issues/148). This is caused by the fuzzy find algo (within the BubbleTea TUI framework code) crashing when it encounters non-compatible characters in the history file, such as an irregular image binary pattern or a rare non-ascii text character. The fix is to remove the clipboard entry that contains the problematic character. I recommend pinning any items you do not want to lose and running clipse -clear.

</details>

<details>
  <summary><b>My copied entries are not recorded when starting the clipse listener on boot with a systemd service</b></summary>

  There may be a few ways around this, but the workaround discovered in issue [#41](https://github.com/savedra1/clipse/issues/41) was to use a `.desktop` file stored in `~/.config/autostart/`, for example:

  ```shell
  [Desktop Entry]
  Name=clipse
  Comment=Clipse event listener autostart.
  Exec=/home/usrname/Applications/bin/clipse/clipse_1.1.0_linux_amd64/clipse --listen %f
  Terminal=false
  Type=Application
  ```

</details>

<details>
  <summary><b>How it works (TLDR)</b></summary>

  When the app is run for the first time it creates a `$XDG_CONFIG_HOME/clipse` directory containing `config.json`, `clipboard_history.json`, custom_theme.json, and a tmp_files directory for storing image data.

  After `clipse -listen` is executed, a background process watches for clipboard activity and records changes in `clipboard_history.json` unless a different path is specified in `config.json`.

  The TUI is launched with the clipse command. It is built using the BubbleTea
  framework, enabling efficient concurrency and a smooth UX.

  Delete removes the selected item from both the UI and storage

  Select copies the item to the system clipboard and exits the program

</details>
