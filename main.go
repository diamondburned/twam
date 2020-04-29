package main

import (
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const CSS = `
	headerbar {
		min-height: 45px;
	}

	window {
		font-size: 30px;
		color: black;
		background-color: white;
	}

	grid * {
		font-size: 1.75em;
	}
`

func main() {
	gtk.Init(&os.Args)

	game := NewState()
	game.Connect("destroy", gtk.MainQuit)

	css, _ := gtk.CssProviderNew()
	css.LoadFromData(CSS)

	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, css, uint(gtk.STYLE_PROVIDER_PRIORITY_USER))

	// Bind the main game loop at 3Hz.
	glib.TimeoutAdd(1000/3, func() bool {
		game.Update()
		return true
	})

	gtk.Main()
}
