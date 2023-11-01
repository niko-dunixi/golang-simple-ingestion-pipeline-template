# environment.systemPackages = [
#   pkgs.go
# ];
# environment.systemPackages = [
#   pkgs.podman
# ];

let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-22.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShell {
  packages = with pkgs; [
    go
    podman
  ];
}