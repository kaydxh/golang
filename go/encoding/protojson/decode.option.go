package protojson

// If AllowPartial is set, input for messages that will result in missing
func WithUnMashalAllowPartial(allowPartial bool) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.AllowPartial = allowPartial
	})
}

// If DiscardUnknown is set, unknown fields are ignored.
func WithUnmashalDiscardUnknown(discardUnknown bool) UnmarshalerOption {
	return UnmarshalerOptionFunc(func(m *Unmarshaler) {
		m.DiscardUnknown = discardUnknown
	})
}
