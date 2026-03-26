package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	base "github.com/erik-adelbert/firework/fireworks"
	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/launcher"
	sched "github.com/erik-adelbert/firework/internal/scheduler"
	"github.com/erik-adelbert/firework/internal/vec"
	"github.com/erik-adelbert/firework/pkg/epilepsy"
	"github.com/erik-adelbert/firework/tui"
)

func main() {
	noWarning := flag.Bool(
		"no-warning", false, "Skip the epilepsy warning screen",
	)

	demoFile := flag.String(
		"f", "", "Path to timeline file to load",
	)

	flag.Parse()

	if !*noWarning {
		if ok := epilepsy.Warn(); !ok {
			return
		}
	}

	h, w, err := terminalSize()
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		os.Exit(1)
	}

	var mod *tui.Model

	now := time.Now()

	switch {
	case *demoFile != "":
		show := loadShow(*demoFile, h, w, now)
		mod = tui.NewModel(show)
		mod.SetDemo()
	default:
		show := mkShow(h, w, now)
		mod = tui.NewModel(show)
	}

	prg := tea.NewProgram(mod)
	if _, err := prg.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func loadShow(path string, _, _ int, now time.Time) *sched.Scheduler {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	allTimeline, err := sched.AllCSVTimeline(f, fw.Palette(), base.Catalog())
	if err != nil {
		fmt.Printf("Error loading timeline: %v\n", err)
		os.Exit(1)
	}

	timeline := make([]sched.Event, 0, 100)
	for ev, err := range allTimeline {
		if err != nil {
			fmt.Printf("Error parsing timeline: %v\n", err)
			os.Exit(1)
		}

		timeline = append(timeline, ev)
	}
	allEvents := slices.Values(timeline)

	show := sched.New(
		launcher.New([]*fw.Firework{}, false),
		sched.Config{Looping: true},
	)
	show.SetTimeline(allEvents)
	show.Trigger(now)

	return show
}

func mkShow(h, w int, now time.Time) *launcher.Launcher {
	var fires []*fw.Firework

	for range 8 {
		x0 := rand.IntN(9*w/20-w/20) + w/20 // x0 in [w/10, 9w/10)
		y0 := rand.IntN(h/2-h/6) + h/6      // y0 in [h/6, h/2)

		o := vec.Vec{
			X: float64(x0),
			Y: float64(y0),
		}

		fires = append(
			fires,
			func() *fw.Firework {
				return []*fw.Firework{
					base.NewPeony(o),
					base.NewChrysanthemum(o),
					base.NewBrocade(o),
					base.NewPalm(o),
					base.NewWillow(o),
					base.NewFish(o),
					base.NewGlitter(o),
					base.NewSaturn(o),
					base.NewSphere(o),
					base.NewSun(o),
					base.NewKamuro(o),
					base.NewComet(o),
					base.NewLaser(o),
					base.NewDigit(o, rand.IntN(4)),
					base.NewFeather(o),
				}[rand.IntN(15)]
			}(),
		)
	}

	show := launcher.New(fires, true)
	show.Trigger(now)

	return show
}

func terminalSize() (h, w int, err error) {
	fd := os.Stdout.Fd()

	if !term.IsTerminal(fd) {
		return 0, 0, fmt.Errorf("stdout is not a terminal")
	}

	w, h, err = term.GetSize(fd)
	return
}
