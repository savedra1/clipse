{ pkgs ? import <nixos-unstable> { config = { allowUnfree = true; }; } }:

  pkgs.mkShell {
    buildInputs = with pkgs; [ 
      unstable.clipse
  ];

}