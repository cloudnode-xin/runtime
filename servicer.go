package runtime

import "context"

type Servicer interface {
	// property
	Name() string
	IsHealthy() bool

	// life cycle
	Load(f Finder) error
	Start(f Finder, ctx context.Context) error
	Stop(f Finder) error
}
