package runtime

import "context"

type Servicer interface {
	HealthChecker
	// property
	Name() string

	// life cycle
	Load(f Finder) error
	Start(f Finder, ctx context.Context) error
	Stop(f Finder) error
}
