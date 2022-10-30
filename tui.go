package tui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var Logger *log.Logger

func init() {
	// By default, don't save logs
	Logger = log.New(io.Discard, "", 0)
}

type UI struct {
	g                  *gocui.Gui
	root               View
	started            bool
	modal              *Modal
	focused            View
	focusedBeforeModal View
	defaultFocus       View

	focusable       []View
	didAddFocusable map[string]struct{}
}

func New() (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, fmt.Errorf("init gocui: %w", err)
	}
	ui := &UI{g: g, didAddFocusable: make(map[string]struct{})}
	return ui, nil
}

func (ui *UI) Close() {
	ui.g.Close()
}

func (ui *UI) SetRoot(v View) {
	if ui.started {
		panic("can not set root after starting UI")
	}
	ui.root = v
}

func (ui *UI) SetDefaultFocus(v View) {
	ui.defaultFocus = v
}

func (ui *UI) ShowModal(m *Modal) {
	if ui.modal != nil {
		ui.HideModal()
	}
	defer ui.Redraw()
	ui.modal = m
	ui.focusedBeforeModal = ui.focused
	ui.setFocused(m.child)
}

func (ui *UI) HideModal() {
	defer ui.setFocused(ui.focusedBeforeModal)
	m := ui.modal
	if m == nil {
		return
	}
	c := m.child
	if c == nil {
		return
	}
	defer ui.Redraw()
	ui.modal = nil

	if err := ui.g.DeleteView(m.child.ID()); err != nil && err != gocui.ErrUnknownView {
		Logger.Fatal("hide modal: %w", err)
	}
}

func (ui *UI) Redraw() {
	if !ui.started || ui.root == nil {
		return
	}
	ui.g.Update(func(g *gocui.Gui) error {
		v, err := ui.g.View(ui.root.ID())
		if err != nil {
			return fmt.Errorf("get root view: %w", err)
		}
		if err := ui.root.draw(g, v); err != nil {
			return fmt.Errorf("draw root: %w", err)
		}
		if ui.modal != nil {
			if err := ui.modal.draw(ui); err != nil {
				return fmt.Errorf("draw modal: %w", err)
			}
		}
		return nil
	})
}

func (ui *UI) Run() error {
	if ui.root == nil {
		return fmt.Errorf("no root view specified")
	}
	ui.g.SetManager(gocui.ManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		Logger.Printf("layout (%d, %d)", maxX, maxY)
		ui.root.context(nil, 0, 0, maxX, maxY)
		ui.root.setUI(ui)
		v, err := g.SetView(ui.root.ID(), 0, 0, maxX, maxY)
		if err != nil && err != gocui.ErrUnknownView {
			return fmt.Errorf("set root view %w", err)
		}
		if err := ui.root.draw(g, v); err != nil {
			return fmt.Errorf("draw root: %w", err)
		}
		if ui.modal != nil {
			if err := ui.modal.draw(ui); err != nil {
				return fmt.Errorf("draw modal: %w", err)
			}
		}
		if ui.defaultFocus != nil && ui.focused == nil {
			Logger.Println("set default focus")
			ui.setFocused(ui.defaultFocus)
		}
		return nil
	}))
	if err := ui.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return fmt.Errorf("setup quit keybinding: %w", err)
	}
	if err := ui.g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		ui.toggleFocus(false)
		return nil
	}); err != nil {
		return fmt.Errorf("setup toggle focus forward keybinding: %w", err)
	}
	if err := ui.g.SetKeybinding("", gocui.KeyTab, gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		ui.toggleFocus(true)
		return nil
	}); err != nil {
		return fmt.Errorf("setup toggle focus backward keybinding: %w", err)
	}
	ui.started = true
	if err := ui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (ui *UI) Quit() {
	ui.g.Update(func(g *gocui.Gui) error {
		ui.started = false
		return gocui.ErrQuit
	})
}

func (ui *UI) ReleaseFocus() {
	ui.setFocused(nil)
}

func (ui *UI) setFocused(v View) {
	var didRedraw bool
	if ui.focused != nil {
		defer ui.Redraw()
		ui.g.SetCurrentView("")
		didRedraw = true
		ui.focused.onFocusLost()
	}
	if v == nil {
		return
	}
	if !didRedraw {
		defer ui.Redraw()
	}

	ui.focused = v
	v.onFocus()
	for _, view := range ui.g.Views() {
		if strings.HasPrefix(view.Name(), v.ID()) {
			ui.g.SetCurrentView(view.Name())
			return
		}
	}
}

func (ui *UI) registerFocusable(v View) {
	if _, ok := ui.didAddFocusable[v.ID()]; ok {
		return
	}
	ui.didAddFocusable[v.ID()] = struct{}{}

	ui.focusable = append(ui.focusable, v)
}

func (ui *UI) toggleFocus(backwards bool) {
	if ui.modal != nil {
		return
	}
	if len(ui.focusable) == 0 {
		return
	}
	if ui.focused == nil {
		ui.setFocused(ui.focusable[0])
		return
	}
	nextIdx := 0
	for i, f := range ui.focusable {
		if f.ID() == ui.focused.ID() {
			nextIdx = i
			if backwards {
				nextIdx -= 1
			} else {
				nextIdx += 1
			}
		}
	}

	if nextIdx < 0 {
		nextIdx = len(ui.focusable) - 1
	} else if nextIdx >= len(ui.focusable) {
		nextIdx = 0
	}

	ui.setFocused(ui.focusable[nextIdx])
}
