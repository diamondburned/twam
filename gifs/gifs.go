package gifs

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
)

var (
	gifWon   *gdk.PixbufAnimation
	gifLost  *gdk.PixbufAnimation
	gifStart *gdk.PixbufAnimation
)

func Won() *gdk.PixbufAnimation {
	return gifLoad(&gifWon, "./gifs/won.gif")
}
func Lost() *gdk.PixbufAnimation {
	return gifLoad(&gifLost, "./gifs/lost.gif")
}
func Start() *gdk.PixbufAnimation {
	return gifLoad(&gifStart, "./gifs/start.gif")
}

func gifLoad(gif **gdk.PixbufAnimation, path string) *gdk.PixbufAnimation {
	if pb := *gif; pb != nil {
		return pb
	}

	p, err := gdk.PixbufAnimationNewFromFile(path)
	if err != nil {
		log.Fatalln("Failed to get start.gif:", err)
	}
	*gif = p
	return p
}
