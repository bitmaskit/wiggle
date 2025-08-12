package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

var IdleTimeout = 60 * time.Second // Default idle timeout if not set by env or args

const (
	WiggleDuration = 1 * time.Second
	PollInterval   = 250 * time.Millisecond // min 25ms
	Amplitude      = 6                      // Size of Wiggle: Min 1
	WiggleStep     = 10 * time.Millisecond  // Speed of wiggle movement
)

func init() {
	// Try environment variable first
	if idleTimeout, ok := os.LookupEnv("IDLE_TIMEOUT"); ok {
		if idle, err := strconv.Atoi(idleTimeout); err == nil && idle > 0 {
			IdleTimeout = time.Duration(idle) * time.Second
			return
		}
	}

	// Try command-line argument
	if len(os.Args) > 1 {
		if idle, err := strconv.Atoi(os.Args[1]); err == nil && idle > 0 {
			IdleTimeout = time.Duration(idle) * time.Second
			return
		}
	}

	// Fallback to default
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("Mouse Wiggler started. Idle timeout: %v, Wiggle: %v, Poll: %v, Amplitude: %dpx\n",
		IdleTimeout, WiggleDuration, PollInterval, Amplitude)
	fmt.Println("Press Ctrl+C to exit.")

	monitor(ctx)
}

func monitor(ctx context.Context) {
	// Establish initial position
	xPrev, yPrev := robotgo.Location()
	lastMove := time.Now()

	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			x, y := robotgo.Location()
			if x != xPrev || y != yPrev {
				lastMove = time.Now()
				xPrev, yPrev = x, y
				continue
			}

			if time.Since(lastMove) >= IdleTimeout {
				// Inactivity detected; perform wiggle then reset timer
				wiggleOnce(ctx)
				// After wiggle, reset lastMove to now and update previous position
				xPrev, yPrev = robotgo.Location()
				lastMove = time.Now()
			}
		}
	}
}
func wiggleOnce(ctx context.Context) {
	ox, oy := robotgo.Location()
	deadline := time.Now().Add(WiggleDuration)

	const numberOfPointsOnCircle = 16
	positions := make([]struct{ dx, dy int }, numberOfPointsOnCircle)
	for i := 0; i < numberOfPointsOnCircle; i++ {
		angle := 2 * 3.14159265 * float64(i) / float64(numberOfPointsOnCircle)
		positions[i].dx = int(float64(Amplitude) * math.Cos(angle))
		positions[i].dy = int(float64(Amplitude) * math.Sin(angle))
	}

	for idx := 0; time.Now().Before(deadline); {
		select {
		case <-ctx.Done():
			robotgo.Move(ox, oy)
			return
		case <-time.After(WiggleStep):
			p := positions[idx]
			robotgo.Move(ox+p.dx, oy+p.dy)
			idx = (idx + 1) % len(positions)
		}
	}
	robotgo.Move(ox, oy)
}
