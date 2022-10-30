package tui

import (
	"github.com/google/uuid"
	"github.com/jroimartin/gocui"
)

type View interface {
	ID() string
	Parent() View
	Focus()
	DeFocus()

	onFocus()
	onFocusLost()
	draw(*gocui.Gui, *gocui.View) error
	size() (int, int)
	pos() (int, int)
	context(parent View, x, y, width, height int)
	setUI(ui *UI)
	getUI() *UI
	findUI() *UI
	redraw()
}

func newView() view {
	return view{
		id: uuid.NewString(),
	}
}

type view struct {
	id            string
	parent        View
	width, height int
	x, y          int
	ui            *UI
	isFocused     bool
}

func (v *view) ID() string {
	return v.id
}

func (v *view) Parent() View {
	return v.parent
}

func (v *view) draw(*gocui.Gui, *gocui.View) error {
	return nil
}

func (v *view) size() (int, int) {
	return v.width, v.height
}

func (v *view) pos() (int, int) {
	return v.x, v.y
}

func (v *view) setUI(ui *UI) {
	v.ui = ui
}

func (v *view) getUI() *UI {
	return v.ui
}

func (v *view) findUI() *UI {
	if ui := v.getUI(); ui != nil {
		return ui
	}
	p := v.Parent()
	if p == nil {
		return nil
	}
	return p.findUI()
}

func (v *view) Focus() {
	if ui := v.findUI(); ui != nil {
		ui.setFocused(v)
	}
}

func (v *view) DeFocus() {
	if ui := v.findUI(); ui != nil {
		ui.setFocused(nil)
	}
}

func (v *view) onFocus()     { v.isFocused = true }
func (v *view) onFocusLost() { v.isFocused = false }

func (v *view) redraw() {
	if ui := v.findUI(); ui != nil {
		ui.Redraw()
	}
}

func (v *view) context(parent View, x, y, width, height int) {
	v.parent = parent
	v.x, v.y = x, y
	v.width, v.height = width, height
}

func parentCount(v View) int {
	if v.Parent() == nil {
		return 0
	}
	return 1 + parentCount(v.Parent())
}
