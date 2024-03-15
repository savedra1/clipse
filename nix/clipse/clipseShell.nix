{ pkgs ? import <nixpkgs> {} }:

  pkgs.mkShell {
    buildInputs = with pkgs; [ 
      go
  ];
  shellHook = '' 
    go install https://github.com/savedra1/clipse@latest
  '';
}