# default.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation {
  name = "build-essential";
  buildInputs = with pkgs; [
    gcc
    glibc
    binutils
    coreutils
    make
    bash
    patch
    findutils
    grep
    sed
    tar
    gzip
    bzip2
    xz
    unzip
    curl
    wget
    git
    perl
  ];
}
