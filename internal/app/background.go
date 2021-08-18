package app

import "fmt"

func (s *server) background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				s.logger.WithError(fmt.Errorf("%s", err)).Error("background email error")
			}
		}()

		fn()
	}()
}