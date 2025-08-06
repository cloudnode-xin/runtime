package runtime

import (
	"context"
	"slices"
)

type Scope struct {
	parent   Finder
	name     string
	children []Servicer
	options  []ScopeOption
}

type ScopeOption func(s *Scope)

func (s *Scope) root() *Service {
	parent := s.parent
	for {
		rt, ok := parent.(*Service)
		if ok {
			return rt
		}

		sc, ok := parent.(*Scope)
		if ok {
			parent = sc.parent
		} else {
			return nil
		}
	}
}

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
	for _, c := range s.children {
		if !c.IsHealthy() {
			return false
		}
	}

	return true
}

func (s *Scope) Use(v Servicer) error {
	for _, c := range s.children {
		if c.Name() == v.Name() {
			return ErrServiceMultiple
		}
	}

	if health, ok := v.(*healthCheckService); ok {
		root := s.root()
		root.Use(health)
		health.Load(root)

		return nil
	}

	if scope, ok := v.(*Scope); ok {
		scope.parent = s
	}

	s.children = append(s.children, v)
	return nil
}

func (s *Scope) Load(f Finder) error {
	for _, opt := range s.options {
		opt(s)
	}

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

		log.Debugf("[%s].%s", getPath(s), getName(c))
		if err := c.Start(s, ctx); err != nil {
			return err
		}
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

		log.Debugf("[%s].%s", getPath(s), getName(c))
		if err := c.Stop(s); err != nil {
			return err
		}
	}

	return nil
}

func NewScope(name string, opts ...ScopeOption) *Scope {
	s := &Scope{
		name:    name,
		options: opts,
	}

	return s
}
