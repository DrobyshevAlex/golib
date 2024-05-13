package outbox

import (
	"context"
	"time"
)

type Sender interface {
	Send(ctx context.Context) (bool, error)
}

type Config struct {
	RetryInterval time.Duration
}

type Service struct {
	cnf    Config
	ch     chan struct{}
	sender Sender
}

func New(cnf Config, sender Sender) *Service {
	return &Service{
		cnf:    cnf,
		ch:     make(chan struct{}, 1),
		sender: sender,
	}
}

func (s Service) Send() {
	go func() {
		if len(s.ch) == 0 {
			s.ch <- struct{}{}
		}
	}()
}

func (s Service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-s.ch:
				if next, err := s.sender.Send(ctx); next {
					if s.cnf.RetryInterval > 0 && err != nil {
						time.Sleep(s.cnf.RetryInterval)
					}
					s.Send()
				}
			}
		}
	}
}
