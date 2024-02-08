{ pkgs ? import <nixpkgs> {} }:

  pkgs.mkShell {
    buildInputs = with pkgs; [ 
      go
      gnumake
      xorg.libX11.dev
      xorg.libXft
      xorg.libXinerama
      xorg.libXcursor
      gnugo
      bison
      flex
      fontforge
      makeWrapper
      pkg-config
      gnumake
      gcc
      libiconv
      autoconf
      automake
      libtool # freetype calls glibtoolize
      xorg.xrandr
  ];

    shellHook = ''
        go get fyne.io/fyne/v2@latest
        go get fyne.io/fyne/v2/storage/repository@v2.4.3                                                
        go get fyne.io/fyne/v2/internal/painter@v2.4.3
        go get fyne.io/fyne/v2/widget@v2.4.3
        go get fyne.io/fyne/v2/internal/painter@v2.4.3
        go get fyne.io/fyne/v2/internal/driver/glfw@v2.4.3
        go get fyne.io/fyne/v2/app@v2.4.3
    '';
}