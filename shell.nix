{ pkgs ? import <nixpkgs> {} }:

pkgs.stdenv.mkDerivation rec {
	name = "twam";
	version = "0.0.1";

	buildInputs = with pkgs; [
		gnome3.glib gnome3.gtk
	];

	nativeBuildInputs = with pkgs; [
		go
	];
}
