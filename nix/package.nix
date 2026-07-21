# Documentation: https://nix.dev/tutorials/packaging-existing-software
# https://nix.dev/tutorials/packaging-existing-software
# clone the nixpkgs repository so that you can have all the nix expressions at
# hand to learn from `git clone --depth 1 git@github.com:nixos/nixpkgs.git`

# Check here https://github.com/NixOS/nixpkgs/blob/4975a89f37fd5e13ddd0c6539886a052ce2efc52/pkgs/build-support/go/module.nix
# to see how buildGoModule is implemented
{
    lib,
    buildGoModule,
}:

let
    version = "0.1.0";
in
buildGoModule {
  # use pname instead of name to follow this RFC https://github.com/NixOS/rfcs/pull/35
  pname = "habits";
  inherit version;
  src = ./..;
  meta = {
    description = "A terminal-based habit tracking application with calendar views, multiple habit types, and note-taking capabilities.";
    mainProgram = "habits";
    homePage = "https://github.com/aristidebm/habits";
    licenses = [ lib.licenses.mit ];
    maintainers = [ "Aristide Bamazi" ];
  };
  postInstall = ''
    # We need to rename the module after postInstall because
    # buildGoModule name the program after the directory that contains the main.go file
    mv $out/bin/cmd $out/bin/habits
  '';
  # This helps getting the correct vendorHash after build failure with this
  # vendorHash = lib.fakeHash;
  vendorHash = "sha256-0FPXA2/E92L/0P4c4VZ5c2aAOZwL03nKgPZHDceJMas=";
}
