package protojson

func WithMashalMultiline(multiline bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.Multiline = multiline
	})
}

func WithMashalUseProtoNames(useProtoNames bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.UseProtoNames = useProtoNames
	})
}

func WithMashalUseEnumNumbers(useEnumNumbers bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.UseEnumNumbers = useEnumNumbers
	})
}

func WithMashalEmitUnpopulated(emitUnpopulated bool) MarshalerOption {
	return MarshalerOptionFunc(func(m *Marshaler) {
		m.EmitUnpopulated = emitUnpopulated
	})
}
