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

## v1.0.3 -> v1.0.5

- feat: added a separate Wayland listener client to access data directly from the stdin using `wl-clipboard --watch`. 
- feat: significantly improved CPU usage if using Wayland 
- fix: not able to copy images from a browser if using Wayland
- bug fix: images copied form stdin and from their temp file no longer share the same byte length for wayland. This lead to a bug where the initial image would not be 'de-duplicated' and would sometimes cause rendering issues. Implemented a fix where all no images can now be duplicated, even if `duplicatesAllowed` is set to `true`. 


