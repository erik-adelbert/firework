package scheduler

import (
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/vec"
)

type Cue struct {
	NewFirework fw.NewFirework
	Center      vec.Vec
	Time        time.Duration
	Priority    int
	index       int
}
