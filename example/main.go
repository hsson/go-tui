package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hsson/go-tui"
)

const logfile = "log.txt"

func main() {
	setupLogging()
	ui, err := tui.New()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	c := tui.NewHorizontalContainer()
	p1 := tui.NewParagraph()
	p2 := tui.NewParagraph()
	p2.AutoScroll = true
	p3 := tui.NewParagraph()
	cl := tui.NewCheckList([]string{"Multiple Selection A", "Multiple Selection B"})

	cl.SetOnCheckChangedListener(func(cl *tui.CheckList, item string, index int, checked bool) {
		log.Printf("Checked %v %s (%d)", checked, item, index)
	})
	rl := tui.NewRadioList([]string{"Single Selection A", "Single Selection B"})
	rl.SetOnRadioChangedListener(func(cl *tui.RadioList, item string, index int, checked bool) {
		log.Printf("Checked %v %s (%d)", checked, item, index)
	})

	cv1 := tui.NewVerticalContainer()
	cv1.Border = true
	cv1.Add(p1, p2, p3)

	cv2 := tui.NewVerticalContainer()
	cv2.Border = true
	cv2.Add(cl, cl)

	c.Add(cv1, rl, cv2)
	c.Border = true

	ui.SetDefaultFocus(cl)

	p1.SetText("Hello...")
	p3.SetText("Toggle between focused inputs using Tab to go forward, or Alt-Tab to go backward.")

	ui.SetRoot(c)

	go func() {
		count := 0
		t := time.NewTicker(1 * time.Second)
		for range t.C {
			count++
			p2.AppendLine(fmt.Sprint(count))
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		p1.Append(" World!")

		m := tui.NewModal("Important Message", 50, rl)
		ui.ShowModal(m)

		time.Sleep(5 * time.Second)
		// ui.HideModal()
		m.Hide()

		rl.Focus()
	}()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func setupLogging() {
	o, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	tui.Logger = log.New(o, "", log.Ldate|log.Ltime|log.LUTC)
	log.SetOutput(o)
}
