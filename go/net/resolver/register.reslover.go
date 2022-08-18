package resolver

import "sync"

var (
	resolversMu sync.RWMutex
	// m is a map from scheme to resolver builder.
	m = make(map[string]Builder)
	// defaultScheme is the default scheme to use.
	defaultScheme = "passthrough"
)

// TODO(bar) install dns resolver in init(){}.

// Register registers the resolver builder to the resolver map. b.Scheme will be
// used as the scheme registered with this builder.
//
func Register(b Builder) {
	resolversMu.Lock()
	defer resolversMu.Unlock()

	if b == nil {
		panic("register builder is nil")
	}
	if _, ok := m[b.Scheme()]; ok {
		panic("double register scheme " + b.Scheme())
	}
	m[b.Scheme()] = b
}

// Get returns the resolver builder registered with the given scheme.
//
// If no builder is register with the scheme, nil will be returned.
func Get(scheme string) Builder {
	resolversMu.Lock()
	defer resolversMu.Unlock()

	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}

// SetDefaultScheme sets the default scheme that will be used. The default
// default scheme is "passthrough".
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe. The scheme set last overrides
// previously set values.
func SetDefaultScheme(scheme string) {
	defaultScheme = scheme
}

// GetDefaultScheme gets the default scheme that will be used.
func GetDefaultScheme() string {
	return defaultScheme
}
