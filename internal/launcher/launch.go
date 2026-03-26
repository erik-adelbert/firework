package launcher

import (
	"iter"
	"slices"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/vec"
)

type Launcher struct {
	Looping bool

	State0  []*fw.Firework
	Actives []*fw.Firework
}

func New(fws []*fw.Firework, looping bool) *Launcher {
	show := new(Launcher)
	show.Init(fws, looping)

	return show
}

func (l *Launcher) Init(state0 []*fw.Firework, looping bool) {
	l.Looping = looping

	l.State0 = make([]*fw.Firework, len(state0))
	copy(l.State0, state0)

	l.Reset()
}

func (l *Launcher) Len() int {
	return len(l.State0)
}

func (l *Launcher) Add(fw *fw.Firework) {
	l.State0 = append(l.State0, fw)
}

func (l *Launcher) Activate(fw *fw.Firework) {
	l.Actives = append(l.Actives, fw)
}

func (l *Launcher) Reset() {
	for _, fw := range l.State0 {
		fw.Reset()
	}

	helper.Copy(&l.Actives, l.State0)
}

func (l *Launcher) Update(now time.Time, dt time.Duration) {
	actives := l.Actives[:0]

	for _, fwk := range l.Actives {
		fwk.Update(now, dt)

		if fwk.State != fw.Dead {
			actives = append(actives, fwk)
		}
	}

	l.Actives = actives

	if l.Looping && len(l.Actives) == 0 && len(l.State0) > 0 {
		l.Reset()
		l.Trigger(now)
	}
}

func (l *Launcher) Trigger(t0 time.Time) {
	for _, fw := range l.Actives {
		fw.Trigger(t0)
	}
}

type allParticles = iter.Seq2[int, vec.Vec]
type allTrails = iter.Seq2[int, allParticles]

func (l *Launcher) AllActiveTrails(fid int) allTrails {
	return func(yield func(int, allParticles) bool) {
		if fid >= len(l.Actives) {
			return
		}

		fw := l.Actives[fid]
		for i, p := range fw.AllParticles() {
			if !yield(i, p.AllTrailPoints()) {
				return
			}
		}
	}
}

func (l *Launcher) AllFireworks() iter.Seq[*fw.Firework] {
	return slices.Values(l.State0)
}
