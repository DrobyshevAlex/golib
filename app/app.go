package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Func func(ctx context.Context)

type closerQueue struct {
	queue []Func
	m     sync.Mutex
	wg    sync.WaitGroup
}

func (c *closerQueue) add(f Func) {
	c.m.Lock()
	defer c.m.Unlock()
	c.queue = append(c.queue, f)
}

func (c *closerQueue) close(ctx context.Context) {
	c.m.Lock()
	defer c.m.Unlock()
	c.wg.Add(len(c.queue))
	for _, f := range c.queue {
		go func(f Func) {
			defer c.wg.Done()
			f(ctx)
		}(f)
	}
	c.wg.Wait()
}

type App struct {
	wg     sync.WaitGroup
	closed closerQueue
	gf     closerQueue
}

func (a *App) Run(ctx context.Context, f func(ctx context.Context)) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		f(ctx)
	}()
}

func (a *App) OnClose(f func(ctx context.Context)) {
	a.closed.add(f)
}

func (a *App) OnGracefulShutdown(f func(ctx context.Context)) {
	a.gf.add(f)
}

func (a *App) close(ctx context.Context) {
	a.closed.close(ctx)
}

func (a *App) gracefulShutdown(ctx context.Context) {
	a.gf.close(ctx)
}

func Run(ctx context.Context, run func(ctx context.Context, app *App) error) error {
	appCtx, cancel := context.WithCancel(ctx)
	app := App{}

	ch := make(chan error)

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		err := run(appCtx, &app)
		if err != nil {
			ch <- err
		}
	}()

	grace := make(chan os.Signal)
	signal.Notify(grace, syscall.SIGINT, syscall.SIGTERM)

	var err error
	app.wg.Add(1)
	go func() {
		select {
		case err = <-ch:
			panic(err)
		case <-grace:
		}
		app.wg.Done()

		app.gracefulShutdown(ctx)

		cancel()
	}()

	app.wg.Wait()

	app.close(ctx)

	return err
}
