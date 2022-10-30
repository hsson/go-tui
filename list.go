package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type OnCheckChangedListener func(cl *CheckList, item string, index int, checked bool)

type CheckList struct {
	view
	items            []string
	current          int
	checked          map[int]bool
	onCheckedChanged OnCheckChangedListener
}

func NewCheckList(items []string) *CheckList {
	return &CheckList{
		view:             newView(),
		items:            items,
		checked:          make(map[int]bool, len(items)),
		onCheckedChanged: func(cl *CheckList, item string, index int, checked bool) {},
	}
}

func (cl *CheckList) SetOnCheckChangedListener(listener OnCheckChangedListener) {
	cl.onCheckedChanged = listener
}

func (cl *CheckList) draw(g *gocui.Gui, v *gocui.View) error {
	if ui := cl.findUI(); ui != nil {
		ui.registerFocusable(cl)
	}

	x, y := v.Cursor()
	Logger.Println("cl cursor", x, y)
	g.DeleteKeybindings(v.Name())
	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if cl.current < len(cl.items)-1 {
			cl.current++
		}
		return nil
	}); err != nil {
		return fmt.Errorf("set down keybinding: %w", err)
	}

	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if cl.current > 0 {
			cl.current--
		}
		return nil
	}); err != nil {
		return fmt.Errorf("set up keybinding: %w", err)
	}

	g.SetKeybinding(v.Name(), gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, idx := v.Cursor()
		if idx >= 0 || idx < len(cl.items) {
			item := cl.items[idx]
			cl.checked[idx] = !cl.checked[idx]
			checked := cl.checked[idx]
			cl.onCheckedChanged(cl, item, idx, checked)
		}

		return nil
	})

	if err := v.SetCursor(0, cl.current); err != nil {
		return fmt.Errorf("set cursor: %w", err)
	}
	v.SelFgColor = gocui.ColorMagenta
	if cl.isFocused {
		v.Highlight = true
	} else {
		v.Highlight = false
	}

	for i, li := range cl.items {
		if cl.checked[i] {
			fmt.Fprintf(v, "[x] %s\n", li)
		} else {
			fmt.Fprintf(v, "[ ] %s\n", li)
		}
	}
	return nil
}

type OnRadioChangedListener func(cl *RadioList, item string, index int, checked bool)

type RadioList struct {
	view
	items            []string
	current          int
	checked          int
	onCheckedChanged OnRadioChangedListener
}

func NewRadioList(items []string) *RadioList {
	return &RadioList{
		view:    newView(),
		items:   items,
		checked: -1,
	}
}

func (rl *RadioList) SetOnRadioChangedListener(listener OnRadioChangedListener) {
	rl.onCheckedChanged = listener
}

func (rl *RadioList) draw(g *gocui.Gui, v *gocui.View) error {
	if ui := rl.findUI(); ui != nil {
		ui.registerFocusable(rl)
	}
	g.DeleteKeybindings(v.Name())
	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if rl.current < len(rl.items)-1 {
			rl.current++
		}
		return nil
	}); err != nil {
		return fmt.Errorf("set down keybinding: %w", err)
	}

	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if rl.current > 0 {
			rl.current--
		}
		return nil
	}); err != nil {
		return fmt.Errorf("set up keybinding: %w", err)
	}

	g.SetKeybinding(v.Name(), gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, idx := v.Cursor()
		if idx >= 0 || idx < len(rl.items) {
			current := rl.items[idx]
			prevChecked := rl.checked
			if prevChecked == idx {
				// Only uncheck Current
				rl.checked = -1
				rl.onCheckedChanged(rl, current, idx, false)
			} else if prevChecked == -1 {
				// Only check Current
				rl.checked = idx
				rl.onCheckedChanged(rl, current, idx, true)
			} else {
				// Switch selection
				prev := rl.items[prevChecked]
				rl.checked = idx
				rl.onCheckedChanged(rl, prev, prevChecked, false)
				rl.onCheckedChanged(rl, current, idx, true)
			}
		}
		return nil
	})

	if err := v.SetCursor(0, rl.current); err != nil {
		return fmt.Errorf("set cursor: %w", err)
	}
	v.SelFgColor = gocui.ColorMagenta
	if rl.isFocused {
		v.Highlight = true
	} else {
		v.Highlight = false
	}

	for i, li := range rl.items {
		if rl.checked == i {
			fmt.Fprintf(v, "(x) %s\n", li)
		} else {
			fmt.Fprintf(v, "( ) %s\n", li)
		}
	}
	return nil
}
