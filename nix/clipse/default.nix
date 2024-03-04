
{ lib, buildGoModule, fetchFromGitHub }:

buildGoModule {
  pname = "clipse";
  version = "0.0.2";

  src = fetchFromGitHub {
    owner = "savedra1";
    repo = "clipse";
    rev = "v0.0.2";
    hash = "sha256-jZtaZszduD+xZu5hv7kzkAJKkD24l25YSQf6Z7l9Wiw=";
  };

  vendorHash = "sha256-GIUEx4h3xvLySjBAQKajby2cdH8ioHkv8aPskHN0V+w=";

  meta = with lib; {
    description = "A useful clipboard manager TUI for Unix.";
    homepage = "https://github.com/savedra1/clipse";
    license = licenses.gpl3Only;
    platforms = platforms.linux;
    maintainers = with maintainers; [ "savedra1" ];
    mainProgram = "clipse";
  };
}