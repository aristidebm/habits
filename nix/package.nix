# Documentation: https://nix.dev/tutorials/packaging-existing-software
# https://nix.dev/tutorials/packaging-existing-software
# clone the nixpkgs repository so that you can have all the nix expressions at
# hand to learn from `git clone --depth 1 git@github.com:nixos/nixpkgs.git`
#

{ pkgs, lib, ... }:

pkgs.buildGoModule {
  pname = "habits";
  version = "0.1.0";
  src = ./.;
  meta = {
    description = "A terminal-based habit tracking application with calendar views, multiple habit types, and note-taking capabilities.";
    mainPogram = "habits";
    homePage = "https://github.com/aristidebm/habits";
    licenses = [ lib.licenses.mit ];
    maintainers = [ ];
  };
}
