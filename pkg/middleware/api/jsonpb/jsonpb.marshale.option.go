package apijsonpb

func WithUseProtoNames(useProtoNames bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.useProtoNames = useProtoNames
	})
}

func WithUseEnumNumbers(useEnumNumbers bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.useEnumNumbers = useEnumNumbers
	})
}

func WithEmitUnpopulated(emitUnpopulated bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.emitUnpopulated = emitUnpopulated
	})
}

func WithDiscardUnknown(discardUnknown bool) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.discardUnknown = discardUnknown
	})
}

func WithIndent(indent string) JSONPbOption {
	return JSONPbOptionFunc(func(c *JSONPb) {
		c.opts.indent = indent
	})
}
