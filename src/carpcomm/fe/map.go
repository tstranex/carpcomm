// Author: Timothy Stranex <tstranex@carpcomm.com>
// Copyright 2013 Timothy Stranex

package main

import "carpcomm/scheduler"
import "image"
import "image/color"
import "image/draw"
import "math"

func minmax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func line(img draw.Image, a, b image.Point, c color.Color) {
	minx, maxx := minmax(a.X, b.X)
	miny, maxy := minmax(a.Y, b.Y)

	Δx := float64(b.X - a.X)
	Δy := float64(b.Y - a.Y)

	if maxx - minx > maxy - miny {
		d := 1
		if a.X > b.X {
			d = -1
		}
		for x := 0; x != b.X - a.X + d; x += d {
			y := int(float64(x) * Δy / Δx)
			img.Set(a.X + x, a.Y + y, c)
		}
	} else {
		d := 1
		if a.Y > b.Y {
			d = -1
		}
		for y := 0; y != b.Y - a.Y + d; y += d {
			x := int(float64(y) * Δx / Δy)
			img.Set(a.X + x, a.Y + y, c)
		}
	}
}

func orbitMap(points []scheduler.SatPoint) image.Image {
	width := 640
	height := 496
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	scale := 101.95

	var last image.Point
	for i, p := range points {
		// Mercator projection
		x := int(float64(width) * (p.LongitudeDegrees + 180.0) / 360.0)
		ϕ := math.Pi * p.LatitudeDegrees / 180.0
		y := height/2 - int(scale*math.Log(math.Tan(math.Pi/4 + ϕ/2)))
		if y < 0 {
			y = 0
		}
		if y >= height {
			y = height - 1
		}

		pt := image.Pt(x, y)

		if i == 0 {
			// Draw the cubesat.
			r := image.Rect(x-5, y-5, x+5, y+5)
			draw.Draw(img, r, image.White, image.Pt(0, 0), draw.Src)
		}

		// Avoid drawing a discontinous line.
		dx := int(math.Abs(float64(pt.X - last.X)))
		if i > 0 && dx < width - dx {
			line(img, last, pt, color.White)
		}

		last = pt
	}

	return img
}