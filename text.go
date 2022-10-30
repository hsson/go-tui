package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/mitchellh/go-wordwrap"
)

type Paragraph struct {
	view
	text string

	AutoScroll          bool
	DisableLineWrapping bool
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		view: newView(),
	}
}

func (p *Paragraph) ID() string {
	return p.id
}

func (p *Paragraph) draw(_ *gocui.Gui, v *gocui.View) error {
	v.Autoscroll = p.AutoScroll
	text := p.text
	if !p.DisableLineWrapping {
		maxWidth, _ := p.size()
		text = wordwrap.WrapString(text, uint(maxWidth-1))
	}
	fmt.Fprint(v, text)
	return nil
}

func (p *Paragraph) SetText(text string) {
	defer p.redraw()
	p.text = text
}

func (p *Paragraph) Append(text string) {
	defer p.redraw()
	p.text += text
}

func (p *Paragraph) AppendLine(text string) {
	defer p.redraw()
	p.text += text + "\n"
}

func (p *Paragraph) NewLine() {
	defer p.redraw()
	p.text += "\n"
}

func (c *Paragraph) size() (int, int) {
	return c.width, c.height
}
