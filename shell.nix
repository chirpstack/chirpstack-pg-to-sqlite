{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-24.05.tar.gz") {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.opkg-utils
    pkgs.sqlite
    pkgs.docker
  ];
}
