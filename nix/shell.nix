{ go, pkgs, ... }:

pkgs.mkShell {
  packages = [
    go
  ]
  ++ (with pkgs; [
    gopls # lsp
    golangci-lint # linter
    delve # debugger
  ]);
}
