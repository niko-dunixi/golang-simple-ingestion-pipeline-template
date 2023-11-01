
# { pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-22.11") {} }:

# pkgs.mkShell {
#   packages = [
#     (pkgs.python3.withPackages (ps: [
#       ps.flask
#     ]))

#     pkgs.go_1_20
#     pkgs.curl
#     pkgs.jq
#   ];
# }

# templ = buildGoModule rec {
#   pname = "templ";
#   version = "0.3.4";

#   src = fetchFromGitHub {
#     owner = "knqyf263";
#     repo = "pet";
#     rev = "v${version}";
#     hash = "sha256-Gjw1dRrgM8D3G7v6WIM2+50r4HmTXvx0Xxme2fH9TlQ=";
#   };

#   vendorHash = "sha256-ciBIR+a1oaYH+H1PcC8cD8ncfJczk1IiJ8iYNM+R6aA=";

#   meta = with lib; {
#     description = "Simple command-line snippet manager, written in Go";
#     homepage = "https://github.com/knqyf263/pet";
#     license = licenses.mit;
#     maintainers = with maintainers; [ kalbasit ];
#   };
# }

# # {
# #   pkgs ? import (fetchTarball {
# #     url = "https://github.com/NixOS/nixpkgs/archive/4fe8d07066f6ea82cda2b0c9ae7aee59b2d241b3.tar.gz";
# #     sha256 = "sha256:06jzngg5jm1f81sc4xfskvvgjy5bblz51xpl788mnps1wrkykfhp";
# #   }) {}
# # }:
# # pkgs.mkShell rec {
# #    buildInputs = with pkgs; [
# #     go_1_20
# #     # cmake
# #     # boost
# #     # simgrid

# #     # # debugging tools
# #     # gdb
# #     # valgrind
# #    ];
# # }

# # # { stdenv, fetchurl }:


# # # stdenv.mkDerivation rec {
# # #   name = "golang-simple-ingestion-pipeline-template";
# # #   version = "1.0";

# # #   src = fetchurl {
# # #     url = "https://example.com/my-package-1.0.tar.gz"; # Replace with the actual URL of the tarball
# # #     sha256 = "<replace_with_sha256_checksum>"; # Replace with the actual SHA-256 checksum
# # #   };

# # #   meta = with stdenv.lib; {
# # #     description = "My package, a great software tool.";
# # #     license = licenses.mit;
# # #   };
# # # }

# # # # environment.systemPackages = [
# # # #   pkgs.go
# # # # ];
# # # # environment.systemPackages = [
# # # #   pkgs.podman
# # # # ];

# # # let
# # #   nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-22.11";
# # #   pkgs = import nixpkgs { config = {}; overlays = []; };
# # # in


# # # # pkgs.mkShell {
# # # #   buildInputs = [ pkgs.go_1_20 ];
# # # # }

# # # pkgs.mkShell {
# # #   packages = [
# # #     pkgs.go_1_20
# # #   ];
# # # }