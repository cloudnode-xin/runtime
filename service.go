package runtime

import (
	"context"
	"slices"
)

type Service struct {
	cancel   context.CancelFunc
	children []Servicer
}

func (s *Service) Get(name string) (Servicer, error) {
	for _, c := range s.children {
		if c.Name() == name {
			return c, nil
		}
	}

	return nil, ErrServiceNotExist
}

func (s *Service) MustGet(name string) Servicer {
	v, err := s.Get(name)

	if err != nil {
		panic(err)
	}

	return v
}

func (s *Service) Use(svc ...Servicer) {
	for _, v := range svc {
		if scope, ok := v.(*Scope); ok {
			scope.parent = s
		}
	}

	s.children = append(s.children, svc...)
}

func (s *Service) Start() error {
	if s.cancel != nil {
		return nil
	}

	log := s.MustGet("#logger").(*Logger).New("runtime.start")

	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())
	for _, c := range s.children {
		if scope, ok := c.(*Scope); ok {
			if err := scope.Load(scope); err != nil {
				return err
			}
			continue
		}

		if err := c.Load(s); err != nil {
			return err
		}
	}

	for _, c := range s.children {
		if scope, ok := c.(*Scope); ok {
			if err := scope.Start(scope, ctx); err != nil {
				return err
			}

			continue
		}

		if err := c.Start(s, ctx); err != nil {
			return err
		}
		log.Debugf("%s", getName(c))
	}

	return nil
}

func (s *Service) Stop() error {
	if s.cancel == nil {
		return nil
	}

	s.cancel()
	s.cancel = nil

	log := s.MustGet("#logger").(*Logger).New("runtime.stop")

	reverse := slices.Clone(s.children)
	slices.Reverse(reverse)

	for _, c := range reverse {

		if scope, ok := c.(*Scope); ok {
			if err := scope.Stop(scope); err != nil {
				return err
			}

			continue
		}

		if err := c.Stop(s); err != nil {
			return err
		}
		log.Debugf("%s", getName(c))
	}

	return nil
}

func New() *Service {
	s := &Service{}

	s.Use(logger())

	return s
}
