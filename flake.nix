{
  description = "clipse dev shell";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }: let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
        go
        gnumake
      ];

      shellHook = ''
        echo "clipse dev shell ready"
        echo "  make wayland  → Wayland build"
        echo "  make x11      → X11 build"
      '';
    };
  };
}
