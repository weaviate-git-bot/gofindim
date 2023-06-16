{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-22.11";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  outputs = { self, nixpkgs, devenv, systems, nixpkgs-unstable, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
            unstablePkgs = nixpkgs-unstable.legacyPackages.${system};

          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {
                  # https://devenv.sh/reference/options/
                  packages = with unstablePkgs;[
                    go
                    gnumake
                    stdenv.cc
                    pkg-config
                  ];
                  env = with unstablePkgs;{
                    CGO_CFLAGS_ALLOW = "'-Xpreprocessor|-Xcompiler|-D__CORRECT_ISO_CPP_STRING_H_PROTO|-D_MT|-D_DLL'|gcc";
                    CGO_CPPFLAGS = "-I${glibc.dev}/include";
                    CGO_LDFLAGS = "${ lib.concatMapStringsSep " " (p: "-L${p}/lib") [glib]}";
                    CGO_ENABLED = "1";
                    PKG_CONFIG_PATH = "${opencv}/lib/pkgconfig:$PKG_CONFIG_PATH";
                  };

                  enterShell = ''
                  '';
                }
              ];
            };
            manual = with unstablePkgs;
              let
                packages = [
                ];
                lib-path = pkgs.lib.makeLibraryPath
                  (packages ++ [
                    glibc
                  ]);
              in

              pkgs.mkShell {
                packages = packages;
                nativeBuildInputs = [
                  glibc
                  stdenv.cc.cc
                  sqlite
                  opencv
                  pkg-config
                  go
                  faiss
                ];
                buildInputs = [
                  glibc
                ];
                # export CGO_CFLAGS_ALLOW="'-Xpreprocessor|-Xcompiler|-D__CORRECT_ISO_CPP_STRING_H_PROTO|-D_MT|-D_DLL'|gcc"
                shellHook = ''
                  export CGO_CPPFLAGS="-I${glibc.dev}/include"
                  export LD_LIBRARY_PATH="${stdenv.cc.cc}/lib"
                  export CGO_LDFLAGS="-L${lib-path} -L${glibc}/lib -L${opencv}/lib -L${faiss}/lib"
                '';
              };
          });
    };
}
