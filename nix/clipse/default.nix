{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule rec {
  pname = "clipse";
  version = "1.0.6";

  src = fetchFromGitHub {
    owner = "savedra1";
    repo = "clipse";
    rev = "v${version}";
    hash = "sha256-nkB7HUleEDSZTbV6u+my6DM76OlCA5r2qmXxmPPGwfI=";
  };

  vendorHash = "sha256-QEBRlwNS8K44chB3fMOJZxYnIaWMnuDySIhKfF7XtxM=";

  meta = {
    description = "Useful clipboard manager TUI for Unix";
    homepage = "https://github.com/savedra1/clipse";
    license = lib.licenses.mit;
    mainProgram = "clipse";
    maintainers = [ lib.maintainers.savedra1 ];
  };
}
