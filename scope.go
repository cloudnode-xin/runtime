package runtime

import (
	"context"
	"slices"
)

type Scope struct {
	parent   Finder
	children []Servicer
	name     string
}

type ScopeOption func(s *Scope)

func (s *Scope) Get(name string) (Servicer, error) {
	for _, c := range s.children {
		if c.Name() == name {
			return c, nil
		}
	}

	if s.parent != nil {
		return s.parent.Get(name)
	}

	return nil, ErrServiceNotExist
}

func (s *Scope) MustGet(name string) Servicer {
	c, err := s.Get(name)

	if err != nil {
		panic(err)
	}

	return c
}

func (s *Scope) Name() string {
	return s.name
}

func (s *Scope) IsHealthy() bool {
	return true
}

func (s *Scope) Use(svc ...Servicer) {
	for _, v := range svc {
		if scope, ok := v.(*Scope); ok {
			scope.parent = s
		}
	}

	s.children = append(s.children, svc...)
}

func (s *Scope) Load(f Finder) error {
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

	return nil
}

func (s *Scope) Start(f Finder, ctx context.Context) error {
	log := s.MustGet("#logger").(*Logger).New("runtime.start")

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
		log.Debugf("[%s].%s", getPath(s), getName(c))
	}

	return nil
}

func (s *Scope) Stop(f Finder) error {

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
		log.Debugf("[%s].%s", getPath(s), getName(c))
	}

	return nil
}

func NewScope(name string, opts ...ScopeOption) *Scope {
	s := &Scope{
		name: name,
	}

	for _, v := range opts {
		v(s)
	}

	return s
}
