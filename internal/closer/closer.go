// Package closer provides functionality to manage closers that should be closed when the application stops.
package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

// Closer manages a set of cleanup functions that need to be executed.
// upon receiving a termination signal or when explicitly invoked.
type Closer struct {
	once    sync.Once
	done    chan struct{}
	mu      sync.Mutex
	closers []func() error
}

// NewCloser creates a new Closer instance. If signals are provided,
// it starts a goroutine that waits for any of those signals to be received.
func NewCloser(sig ...os.Signal) *Closer {
	clsr := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
		}()
	}

	return clsr
}

// Add provide functionality to add new closer to the storage.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.closers = append(c.closers, f...)
	c.mu.Unlock()
}

// Wait using to wait done signal.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll runs all closers that were added via the Add method.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.closers
		c.closers = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, fn := range funcs {
			go func(fn func() error) {
				errs <- fn()
			}(fn)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Printf("closing error: %v\n", err)
			}
		}
	})
}
