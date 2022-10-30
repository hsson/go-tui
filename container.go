package tui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type Container struct {
	view
	horizontal bool
	Views      []View
	Border     bool
}

var _ View = (*Container)(nil)

func NewHorizontalContainer() *Container {
	return &Container{
		view:       newView(),
		horizontal: true,
	}
}

func NewVerticalContainer() *Container {
	return &Container{
		view:       newView(),
		horizontal: false,
	}
}

func (c *Container) ID() string {
	return c.id
}

func (c *Container) Parent() View {
	return c.parent
}

func (c *Container) size() (int, int) {
	return c.width, c.height
}

func (c *Container) Add(views ...View) {
	c.Views = append(c.Views, views...)
}

func (c *Container) draw(g *gocui.Gui, _ *gocui.View) error {
	if len(c.Views) == 0 {
		return nil
	}
	parents := parentCount(c)
	indent := "  " + strings.Repeat("  ", parents)

	maxX, maxY := c.size()
	posX, posY := c.pos()

	var isInOtherContainer bool
	if c.Parent() != nil {
		_, ok := c.Parent().(*Container)
		isInOtherContainer = ok
	}

	childCount := len(c.Views)
	xSize := maxX
	ySize := maxY

	if c.horizontal {
		xSize = maxX / childCount
	} else {
		ySize = maxY / childCount
	}

	Logger.Printf("%ssize (%d, %d)", indent, xSize, ySize)
	for i, view := range c.Views {
		xOffset, yOffset := posX, posY
		if c.horizontal {
			xOffset += i * (xSize + 1)
		} else {
			yOffset += i * (ySize + 1)
		}

		x0, y0 := xOffset, yOffset
		x1 := xOffset + xSize
		y1 := yOffset + ySize

		if c.horizontal {
			y1 = posY + maxY
		}
		if !isInOtherContainer {
			y1 -= 1
		}
		if i == len(c.Views)-1 {
			x1 = posX + maxX
			if !c.horizontal {
				y1 = posY + maxY
			}
			if !isInOtherContainer {
				x1 -= 1
			}
		}

		Logger.Printf("%sdraw %d -> x0: %d y0: %d x1: %d, y1: %d", indent, i, x0, y0, x1, y1)
		id := fmt.Sprintf("%s-%d", view.ID(), i)
		v, err := g.SetView(id, x0, y0, x1, y1)
		if err != nil && err != gocui.ErrUnknownView {
			return fmt.Errorf("set view %d %q: %w", i, id, err)
		}
		v.Clear()
		v.Frame = c.Border
		width := x1 - x0
		height := y1 - y0
		view.context(c, x0, y0, width, height)
		if err := view.draw(g, v); err != nil {
			return fmt.Errorf("draw view %d %q: %w", i, id, err)
		}
	}
	return nil
}
