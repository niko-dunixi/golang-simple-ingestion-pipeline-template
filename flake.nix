{
  description = "Simple intestion pipeline template"
  
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.05-small";
  inputs.nixpkgs-regression.url = "github:NixOS/nixpkgs/215d4d0fd80ca5163643b03a33fde804a29cc1e2";
  inputs.lowdown-src = { url = "github:kristapsdz/lowdown"; flake = false; };
  inputs.flake-compat = { url = "github:edolstra/flake-compat"; flake = false; };

  
}