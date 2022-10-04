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
  in rec {
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

    packages.${system} = {
      pinit = pkgs.callPackage ./nix/pinit.nix {};
    };

    apps.${system}.default = {
      type = "app";
      program = "${packages.${system}.pinit}/bin/pinit";
    };

    overlays.default = final: prev: packages.${system};

    formatter.${system} = pkgs.treefmt;
  };
}
