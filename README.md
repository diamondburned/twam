# twam

A troll wack-a-mole game made for AP CSP in Golang.

<img src="https://imgur.com/a/uOVnBqK" alt="video" />

## Disclaimer

**DO NOT USE THIS CODE AS REFERENCE. I REPEAT: DO NOT.**

## What's "troll" about this?

It works slightly different than regular whack-a-mole:

1. There can be multiple moles.
2. Holes can be destroyed, which would be marked with a skull.
3. Destroyed holes won't be counted in the randomization process.
4. These multiple moles can go down at random.
5. There are points. You lose if it goes below zero (0).


## There's nothing "troll" about that!

It's less "troll" and more annoying, although I think it's pretty fun as well.

The fact that dead holes don't get counted towards the randomization process
means that the game gets progressively harder as you play it.

When there's only one single hole left, a mole will randomly spawn and
disappear. This is very infuriating and "troll," as the player would try and
click it, but the mole would've already gone down, thus deducting their points.

## What is this crappy game for? Where's my gtkcord?!

This is just written for the PT of my AP CSP class of 2020.

## Can I use this code to learn or refer to it?

No. It's bad, and my license says no, so you can't use it. Period.

## Playing

### Dependencies

`go` version 1.14 or newer and `gtk` version 3. Nix users can do `nix-shell`.

### Running

```sh
git clone https://github.com/diamondburned/twam && cd twam
go run .
```

## Credits

This project uses 3 GIFs borrowed from other people, which are stored in `gifs/`.

- **won.gif**: https://en.wikipedia.org/wiki/File:Default_Dancing_Stick_Figure.gif
- **lost.gif**: https://tenor.com/view/stick-figure-blue-mad-gif-10086091
- **start.gif**: https://tenor.com/view/whack-amole-gif-8064490
