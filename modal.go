package tui

import (
	"fmt"
	"math"

	"github.com/jroimartin/gocui"
)

type Modal struct {
	child                View
	title                string
	modalWidthPercentage int
}

func NewModal(title string, w int, child View) *Modal {
	if w < 0 || w > 100 {
		panic("invalid percentage")
	}
	if child == nil {
		panic("child must not be nil")
	}
	return &Modal{
		title:                title,
		modalWidthPercentage: w,
		child:                child,
	}
}

func (m *Modal) draw(ui *UI) error {
	g := ui.g
	maxX, maxY := g.Size()
	xP := float64(m.modalWidthPercentage) / 100.0
	yP := 0.10
	width := int(math.Round(float64(maxX) * xP))
	height := int(math.Round(float64(maxY) * yP))
	if height < 2 {
		height = 2
	}
	x0 := maxX/2 - width/2
	y0 := maxY/2 - height/2
	x1 := maxX/2 + width/2
	y1 := maxY/2 + height/2

	v, err := g.SetView(m.child.ID(), x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return fmt.Errorf("set modal child view: %w", err)
	}
	v.Clear()
	v.Frame = true
	v.Title = m.title
	m.child.context(nil, x0, y0, width, height)
	m.child.setUI(ui)
	if err := m.child.draw(g, v); err != nil {
		return fmt.Errorf("draw child view: %w", err)
	}
	return nil
}

func (m *Modal) Hide() {
	if m.child == nil {
		return
	}
	if ui := m.child.findUI(); ui != nil {
		ui.HideModal()
	}
}
