{
  description = "default";

  inputs = {
    nixpkgs-stable.url = "nixpkgs/nixos-22.05";
    nixpkgs-unstable.url = "nixpkgs/nixos-unstable";
  };

  outputs = {...} @ inputs: let
    system = "x86_64-linux";
    lib = inputs.nixpkgs-stable.lib;

    unstable-overlay = _: _: {unstable = import inputs.nixpkgs-unstable {inherit system;};};

    pkgs = import inputs.nixpkgs-stable {
      inherit system;
      overlays = [unstable-overlay];
    };
  in {
    devShell.${system} = pkgs.mkShell {
      packages = with pkgs; [
        go_1_18
        gopls # language server
        delve # debugger
        go-tools # linter

        # Formatting
        treefmt
        alejandra
        gofumpt
      ];
    };
  };
}
