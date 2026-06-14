# Documentation: https://nix.dev/tutorials/packaging-existing-software
# https://nix.dev/tutorials/packaging-existing-software
# clone the nixpkgs repository so that you can have all the nix expressions at
# hand to learn from `git clone --depth 1 git@github.com:nixos/nixpkgs.git`
#

{ pkgs, ... }:

pkgs.buildGoModule {
  pname = "habits";
  version = "0.1.0";
  src = ./.;
}
