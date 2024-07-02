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