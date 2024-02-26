# About 
*clipse* is dependencyless, configurable TUI-based clipboard manager built in Go, and supports any Linux OS (especially Wayland). 

Why use *clipse*?
- Customizable TUI allows you to easily match your system's theme
- Fuzzy finder capabilities
- Easily recall, add and delete clipboard history with no explicit expiry
- Image file support
- Terminal-based means TUI can be easily called bound to a key, simulating a GUI without the memory overhead   

<div style="text-align: center;">
  <video width="900" height="600" controls allowfullscreen>
    <source src="resources/demo.mp4" type="video/mp4">
    Your browser does not support the video tag.
  </video>
</div>

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