package scheduler

import (
	"container/heap"
	"iter"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/launcher"
)

type Scheduler struct {
	*launcher.Launcher

	startTime time.Time
	duration  time.Duration

	timeline iter.Seq[Event]

	cues queue

	running bool
}

type Config struct {
	Duration  time.Duration
	Looping   bool
	TimeScale float64
}

func New(launcher *launcher.Launcher, config Config) *Scheduler {
	s := &Scheduler{
		Launcher: launcher,
		duration: config.Duration,
	}

	heap.Init(&s.cues)

	s.Looping = config.Looping

	return s
}

func (s *Scheduler) AddCue(cue *Cue) {
	heap.Push(&s.cues, cue)
}

func (s *Scheduler) SetTimeline(all iter.Seq[Event]) {
	s.timeline = all

	var total time.Duration
	for t := range all {
		s.AddCue(&Cue{
			Time:        t.Time,
			NewFirework: t.NewFirework,
			Center:      t.Center,
		})

		total = max(total, t.Time)
	}

	if s.duration == 0 {
		s.duration = total + 10*time.Second
	}
}

func (s *Scheduler) Init(_ []*fw.Firework, looping bool) {
	s.Looping = looping
}

func (s *Scheduler) AllFireworks() iter.Seq[*fw.Firework] {
	return helper.EmptySeq[*fw.Firework]()
}

func (s *Scheduler) Reset() {
	s.startTime = time.Time{}
	s.running = false

	s.Launcher.Reset()
}

func (s *Scheduler) Trigger(now time.Time) {
	s.startTime = now
	s.running = true
	s.Launcher.Trigger(now)
}

func (s *Scheduler) Update(now time.Time, dt time.Duration) {
	if !s.running {
		return
	}

	elapsed := now.Sub(s.startTime)

	for s.cues.Len() > 0 && s.cues[0].Time <= elapsed {
		cue := heap.Pop(&s.cues).(*Cue)

		fw := cue.NewFirework(cue.Center)
		fw.Reset()
		fw.Trigger(now)

		s.Activate(fw)
	}

	s.Launcher.Update(now, dt)

	if elapsed >= s.duration && len(s.Actives) == 0 {
		s.Reset()

		if s.Looping {
			s.SetTimeline(s.timeline)
			s.Trigger(now)
		}
	}
}

type queue []*Cue

func (pq queue) Len() int {
	return len(pq)
}

func (pq queue) Less(i, j int) bool {
	if pq[i].Time == pq[j].Time {
		return pq[i].Priority < pq[j].Priority
	}

	return pq[i].Time < pq[j].Time
}

func (pq queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *queue) Push(x any) {
	n := len(*pq)
	item := x.(*Cue)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *queue) Pop() any {
	old := *pq
	n := len(old)

	item := old[n-1]
	item.index = -1

	*pq = old[0 : n-1]
	return item
}
