# default.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation {
  name = "build-essential";
  buildInputs = with pkgs; [
    nixclip
  ];
}
