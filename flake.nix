{
  description = "Golang development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-25.11";
  };

  outputs = { self, nixpkgs, ... }:
  let
    systems = ["x86_64-linux"];
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in
  {
      devShells = forAllSystems(system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        go = pkgs.go_1_25;
      in
      {
        default = pkgs.mkShell {
            packages = [go] ++ (with pkgs; [
                gopls # lsp
                golangci-lint # linter
                delve # debugger
            ]);
        };
      });
      formatter = forAllSystems(system: nixpkgs.legacyPackages.${system}.nixos-fmt);
  };
}
