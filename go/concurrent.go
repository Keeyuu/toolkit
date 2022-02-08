package tool

type ConcurrentFunc func()

type Concurrent struct {
	taskChannel  chan ConcurrentFunc
	downChannel  chan struct{}
	maxGoroutine int
}

func NewConcurrent(maxGoroutine int) *Concurrent {
	return (&Concurrent{
		taskChannel:  make(chan ConcurrentFunc, maxGoroutine),
		downChannel:  make(chan struct{}, maxGoroutine),
		maxGoroutine: maxGoroutine,
	}).run()
}

func (c *Concurrent) run() *Concurrent {
	for i := 0; i < c.maxGoroutine; i++ {
		go func() {
			for {
				select {
				case <-c.downChannel:
					return
				case f := <-c.taskChannel:
					f()
				}
			}
		}()
	}
	return c
}

func (c *Concurrent) Add(f ConcurrentFunc) {
	c.taskChannel <- f
}

func (c *Concurrent) Down() {
	for i := 0; i < c.maxGoroutine; i++ {
		c.downChannel <- struct{}{}
	}
}
