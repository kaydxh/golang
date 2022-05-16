package gocv

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

var ZR Rect

func (r Rect) Scale(factor float32) Rect {
	ox := r.X + r.Width/2
	oy := r.Y + r.Height/2
	r.X = ox + int(float32((r.X-ox))*factor)
	r.Y = ox + int(float32((r.Y-oy))*factor)
	r.Width = int(float32(r.Width) * factor)
	r.Height = int(float32(r.Height) * factor)

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
	if r.Height > s.Height {
		r.Height = s.Height
	}
	if r.Width > s.Width {
		r.Width = s.Width
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return ZR
	}
	return r
}
