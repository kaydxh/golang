package taskqueue

var (
	m TaskerMap
)

// Register registers the tasker to the tasker map. b.Name will be
// used as the name registered with this builder.
//
func Register(b Tasker) {
	if b == nil {
		panic("register tasker is nil")
	}

	tasker := Get(b.Scheme())
	if tasker != nil {
		panic("double register tasker " + b.Scheme())
	}
	m.Store(b.Scheme(), b)
}

// Get returns the tasker registered with the given scheme.
//
// If no tasker is register with the scheme, nil will be returned.
func Get(scheme string) Tasker {

	b, has := m.Load(scheme)
	if !has {
		return nil
	}

	return b
}
