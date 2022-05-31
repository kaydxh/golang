package gocv

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

var ZR Rect

func (r Rect) Scale(factor float32) Rect {
	ox := r.X + r.Width/2
	oy := r.Y + r.Height/2
	r.X = ox + int32(float32((r.X-ox))*factor)
	r.Y = oy + int32(float32((r.Y-oy))*factor)
	r.Width = int32(float32(r.Width) * factor)
	r.Height = int32(float32(r.Height) * factor)

	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rect) Empty() bool {
	return r.Width == 0 || r.Height == 0
}

func (r Rect) Intersect(s Rect) Rect {
	if r.X < s.X {
		r.X = s.X
	}
	if r.Y < s.Y {
		r.Y = s.Y
	}

	rMaxX := r.X + r.Width
	sMaxX := s.X + s.Width
	if rMaxX > sMaxX {
		rMaxX = sMaxX
	}
	r.Width = rMaxX - r.X

	rMaxY := r.Y + r.Height
	sMaxY := s.Y + s.Height
	if rMaxY > sMaxY {
		rMaxY = sMaxY
	}
	r.Height = rMaxY - r.Y

	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return ZR
	}
	return r
}

// Union returns the smallest rect that contains both r and s.
func (r Rect) Union(s Rect) Rect {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.X > s.X {
		r.X = s.X
	}
	if r.Y > s.Y {
		r.Y = s.Y
	}

	rMaxX := r.X + r.Width
	sMaxX := s.X + s.Width
	if rMaxX < sMaxX {
		rMaxX = sMaxX
	}
	r.Width = rMaxX - r.X

	rMaxY := r.Y + r.Height
	sMaxY := s.Y + s.Height
	if rMaxY < sMaxY {
		rMaxY = sMaxY
	}
	r.Height = rMaxY - r.Y
	return r
}

// In reports whether every point in r is in s.
// true means s is larger than r
func (r Rect) In(s Rect) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.X <= r.X && r.Y <= s.Y &&
		s.Width <= r.Width && r.Height <= s.Height
}
