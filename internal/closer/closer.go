package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

type Closer struct {
	once    sync.Once
	done    chan struct{}
	mu      sync.Mutex
	closers []func() error
}

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

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.closers = append(c.closers, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

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

		for i := 0; i < cap(errs); i += 1 {
			if err := <-errs; err != nil {
				log.Printf("closing error: %v\n", err)
			}
		}
	})
}
