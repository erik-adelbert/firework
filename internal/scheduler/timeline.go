package scheduler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"image/color"
	"io"
	"iter"
	"strconv"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/vec"
)

type Timeline []Event

type Event struct {
	NewFirework fw.NewFirework
	Color       color.RGBA
	Center      vec.Vec
	Time        time.Duration
}

var (
	ErrInvalidFile   = errors.New("invalid cvs file")
	ErrInvalidRecord = errors.New("invalid timeline")
)

type palette = map[string]color.RGBA
type fireworks = map[string]fw.NewFirework

func AllCSVTimeline(r io.Reader, lut palette, fires fireworks) (iter.Seq2[Event, error], error) {
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidFile, err)
	}

	var noEvent Event

	return func(yield func(Event, error) bool) {
		for i, rec := range records[1:] {
			err := fmt.Errorf(
				"%w: invalid record at line %d", ErrInvalidRecord, i+2,
			)

			if len(rec) < 3 {
				yield(noEvent, err)
				return
			}

			t, errT := time.ParseDuration(rec[0])
			x, errX := strconv.Atoi(rec[2])
			y, errY := strconv.Atoi(rec[3])
			f, okF := fires[rec[1]]

			switch {
			case errT != nil:
				err = fmt.Errorf("%w: %w", err, errT)
			case errX != nil:
				err = fmt.Errorf("%w: %w", err, errX)
			case errY != nil:
				err = fmt.Errorf("%w: %w", err, errY)
			case !okF:
				err = fmt.Errorf("%w: unknown firework %s", err, rec[1])
			default:
				err = nil
			}

			if err != nil {
				yield(noEvent, err)
				return
			}

			e := Event{
				Time: t,
				Center: vec.Vec{
					X: float64(x),
					Y: float64(y),
				},
				NewFirework: f,
			}

			if !yield(e, nil) {
				return
			}
		}
	}, nil
}
