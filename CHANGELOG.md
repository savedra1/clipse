## v0.0.6 -> v0.0.71

- feat: removed need for `$PPID` arg. made optional with `-fc` flag
- feat: pinned items view
- feat: custom paths for source files and temp dir
- bug fix: high CPU usage when image files copies
- bug fix: `clipse -clear` not deleting all temp files
- feat: replace `<BINARY FILE>` image indicator with 📷

## v0.0.71 -> v0.1.0

- feat: additional `clear` commands to save pinned items/images
- feat: multi-select for copy and delete
- feat: multi-select copy from active filter
- feat: warn on deleting a pinned item in the ui
- feat: custom theme support for all ui components
- feat: previews for text and images
- feat: added debug logging
- ci: `golangci-lint` added to build workflow (thank you @ccoVeille)
- bug fix: updated description to show local time rather than UTC
- bug fix: removed duplicate `No items.` status message when clipboard empty
- optimization: improved the listener's go routine pattern to save CPU usage
- optimization: refactored the core codebase to make fewer calls to external files

## v1.0.0 -> v1.0.3

- bug fix: toggle pin status message showing opposite event
- bug fix: duplicated images sharing the same reference file
- feat: optional duplicates

## v1.0.3 -> v1.0.7

- feat: added a separate Wayland listener client to access data directly from the stdin using `wl-clipboard --watch`.
- feat: significantly improved CPU usage if using Wayland
- fix: not able to copy images from a browser if using Wayland
- bug fix: images copied form stdin and from their temp file no longer share the same byte length for wayland. This lead to a bug where the initial image would not be 'de-duplicated' and would sometimes cause rendering issues. Implemented a fix where all no images can now be duplicated, even if `duplicatesAllowed` is set to `true`.
- bug fix: images not keeping pinned status after being chosen on Wayland


## v1.0.7 -> v1.0.8

- bug fix: image binary data sometimes parsing as a string on Wayland
- bug fix: inconsistent viewport start position
- bug fix: inconsistent confirmation list start position

## v1.0.9 -> v1.1.0

- bug fix: EDT timezone bug when saving images
- feat: custom keybinding
- feat: custom image preview rendering with `kitty` shell

## v1.1.0 -> v1.2.0
- bug fix: custom keybinds not registering
- feat: pause listener
- feat: automatic deletion after a specified amount of time
- feat: realtime UI updates
- feat: Exclude certain apps
- feat: wayland exclude sensitive content
- feat: auto-paste
- feat: exit code for custom apis
- bug fix: images not copying on MacOs
- bug fix: high cpu usage on x11 and MacOs
- feat: allow exit preview with quit key and allow copy from preview
- feat: force quit keybind
- bug fix: ensure existing processes are killed correctly
- feat: output-all
- feat: configurable max item length
- feat: feat: optional description timestamp
- feat: optional mouse actions
- chore: remove `attoto clipboard` dependency
- bc: remove `-fc` option

## v1.2.0 -> v1.2.1
- bug fix: Image preview persists when viewport closed with quit key
- bug fix: Correctly handle paste operations for entries with prefix '-' on Wayland
- bug fix: Shell processes being killed when exe is symlinked to a path containing substring 'sh'
- feat: Allow keybinds to take multiple values
