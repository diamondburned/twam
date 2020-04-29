package main

import (
	"fmt"
	"math/rand"

	"github.com/diamondburned/twam/gifs"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	// Mole is actually a rat.
	MoleChar = "🐀"

	// Hole is a hole.
	HoleChar = "🕳️"

	// Skull is the troll attribute.
	SkullChar = "💀"
)

const (
	GridSize   = 4
	GridSum    = GridSize * GridSize
	DeadChance = 45 // 45%
)

type Hole struct {
	*gtk.Button
	Dead   bool // troll attribute
	Popped bool // whether or not the mole is up
}

func (h *Hole) SwitchPop() bool {
	if h.Dead {
		// don't pop if already dead
		return false
	}

	if h.Popped = !h.Popped; h.Popped {
		h.Button.SetLabel(MoleChar)
	} else {
		h.Button.SetLabel(HoleChar)
	}

	return true
}

// Wack a mole and return true, or false if there's no mole.
func (h *Hole) Wack() bool {
	if !h.Popped {
		// no mole
		return false
	}

	if rand.Float64()*100 < DeadChance {
		h.Dead = true
		h.Button.SetLabel(SkullChar)
	} else {
		h.Popped = false
		h.Button.SetLabel(HoleChar)
	}

	return true
}

type State struct {
	*gtk.Window
	Main *gtk.Stack

	GameOverBox   *gtk.Box
	GameOverImage *gtk.Image
	GameOverText  *gtk.Label
	GameOverReset *gtk.Button

	Frame *gtk.AspectFrame
	Grid  *gtk.Grid

	//    0  1  2
	// 0 [0][0][0]
	// 1 [0][0][0]
	// 2 [0][0][0]
	board [GridSize][GridSize]*Hole
	score int
}

func NewState() *State {
	s := &State{}
	s.Window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	s.Window.Show()

	s.Main, _ = gtk.StackNew()
	s.Main.Show()

	// prepare the game over screen
	s.GameOverImage, _ = gtk.ImageNew()
	s.GameOverImage.Show()
	s.GameOverImage.SetFromAnimation(gifs.Start())

	s.GameOverText, _ = gtk.LabelNew("Whack-a-mole!")
	s.GameOverText.Show()
	s.GameOverText.SetJustify(gtk.JUSTIFY_CENTER)

	s.GameOverReset, _ = gtk.ButtonNewWithLabel("Start!")
	s.GameOverReset.Show()
	s.GameOverReset.SetHAlign(gtk.ALIGN_CENTER)
	s.GameOverReset.Connect("clicked", s.Reset)

	s.GameOverBox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 25)
	s.GameOverBox.Show()
	s.GameOverBox.SetVAlign(gtk.ALIGN_CENTER)

	// prepare the grid
	s.Grid, _ = gtk.GridNew()
	s.Grid.Show()
	s.Grid.SetColumnHomogeneous(true)
	s.Grid.SetRowHomogeneous(true)
	s.Grid.SetHExpand(true)
	s.Grid.SetVExpand(true)

	s.Frame, _ = gtk.AspectFrameNew("", 0.5, 0.5, 1, false)
	s.Frame.SetSizeRequest(550, 550)
	s.Frame.Show()

	for x := range s.board {
		for y := range s.board[x] {
			button, _ := gtk.ButtonNew()
			button.Show()
			button.SetRelief(gtk.RELIEF_NONE)
			s.connect(button, x, y)

			s.board[x][y] = &Hole{Button: button}
			s.Grid.Attach(button, x, y, 1, 1)
		}
	}

	s.GameOverBox.Add(s.GameOverImage)
	s.GameOverBox.Add(s.GameOverText)
	s.GameOverBox.Add(s.GameOverReset)

	s.Frame.Add(s.Grid)
	s.Main.AddNamed(s.Frame, "game")
	s.Main.AddNamed(s.GameOverBox, "over")
	s.Window.Add(s.Main)

	// set game state
	s.Reset()
	s.Main.SetVisibleChildName("over")
	return s
}

func (s *State) Reset() {
	s.score = 0

	s.Main.SetVisibleChildName("game")
	s.DeltaScore(0)

	// reset the board
	for x := range s.board {
		for _, board := range s.board[x] {
			board.Dead = false
			board.Popped = false
			board.SetLabel(HoleChar)
		}
	}
}

func (s *State) connect(btn *gtk.Button, x, y int) {
	btn.Connect("button-press-event", func() {
		s.onPress(x, y)
	})
}

func (s *State) DeltaScore(score int) {
	s.score += score
	s.Window.SetTitle(fmt.Sprint("Score: ", s.score))
}

func (s *State) Update() {
	// If we're not in the game, then don't pop.
	if s.Main.GetVisibleChildName() != "game" {
		return
	}

	// try and pop or unpop
	if popped := s.PopRandom(); !popped {
		// This bug is interesting. We add 1 here otherwise we will always lose.
		// The reason is quite interesting.
		//
		// Imagine this scenario: all of the grid is dead except for one hole,
		// and Update() gets called. The UI thread is now running this function.
		//
		// This function then calls PopRandom(), which would then pop up a
		// mouse. All cool and good, the UI thread continues as usual.
		//
		// The user then sees the mouse as the UI updates, and then proceeds to
		// click it. BUT they're too slow! Update() is called again!
		//
		// PopRandom() then runs before the user can click, which sees that
		// there is no more empty hole. The user loses before they can even
		// react!
		//
		// To work around this issue, we add 1 into the dead count. This assumes
		// that the user will win when there is 1 hole left (which would
		// promptly be filled with a rat if we don't add 1). The game works as
		// intended with this.
		if s.GetDeadCount()+1 >= GridSum {
			// if we can't pop anymore and all squares are dead
			s.Win()
		} else {
			// if we can't pop anymore and all squares are moles
			s.Lose()
		}
	}
}

func (s *State) onPress(x, y int) {
	hole := s.board[x][y]

	// Lose if you hit the skull.
	if hole.Dead {
		s.Lose()
		return
	}

	// Deduct point if you hit a hole
	if !hole.Popped {
		s.DeltaScore(-1)

		// If negative score, then lose.
		if s.score < 0 {
			s.Lose()
		}

		return
	}

	hole.Wack()
	s.DeltaScore(1)
}

func (s *State) PopRandom() bool {
	unpopped := s.GetNotDead()
	if len(unpopped) == 0 {
		return false
	}

	pair := unpopped[rand.Intn(len(unpopped))]
	return s.board[pair[0]][pair[1]].SwitchPop()
}

func (s *State) Each(fn func(h *Hole, x, y int)) {
	for x := range s.board {
		for y := range s.board[x] {
			fn(s.board[x][y], x, y)
		}
	}
}

func (s *State) GetDeadCount() int {
	notDead := s.GetNotDead()
	return GridSum - len(notDead)
}

func (s *State) GetNotDead() [][2]int {
	var undead = make([][2]int, 0, GridSum)
	s.Each(func(h *Hole, x, y int) {
		if h.Dead {
			return
		}
		undead = append(undead, [2]int{x, y})
	})
	return undead
}

func (s *State) Win() {
	s.gameOver(gifs.Won(), "<big>You win!</big>")
}
func (s *State) Lose() {
	s.gameOver(gifs.Lost(), fmt.Sprintf(
		"<big>%s</big>\n<small>%s\n%s</small>",
		"You lost!",
		"Clicking on the skulls will instantly kill you.",
		"A negative score will instantly kill you.",
	))
}

func (s *State) gameOver(image *gdk.PixbufAnimation, markup string) {
	s.GameOverImage.SetFromAnimation(image)
	s.GameOverReset.SetLabel("Retry")
	s.GameOverText.SetMarkup(markup)
	s.Main.SetVisibleChildName("over")
	return
}
