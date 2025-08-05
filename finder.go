package runtime

type Finder interface {
	Get(name string) (Servicer, error)
	MustGet(name string) Servicer
}
