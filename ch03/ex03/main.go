// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

var plots = map[string]func(float64, float64) float64{
	"sin(r)/r": func(x, y float64) float64 {
		r := math.Hypot(x, y)
		return math.Sin(r) / r
	},
	"eggbox": func(x, y float64) float64 {
		return math.Sin(x) * math.Sin(y) / 4
	},
	"moguls": func(x, y float64) float64 {
		return (math.Sin(x) + math.Sin(y)) / 25
	},
	"saddle": func(x, y float64) float64 {
		return (math.Pow(x, 3.0) - 3*x*math.Pow(y, 2.0)) / 5000
	},
}
var plot = flag.String("p", "sin(r)/r", "plot, one of sin(r)/r, eggbox, moguls, saddle")

func main() {
	flag.Parse()
	f, ok := plots[*plot]
	if !ok {
		fmt.Fprintf(os.Stderr, "unsupported plot %q\n", *plot)
		os.Exit(1)
	}
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j, f)
			bx, by, bz := corner(i, j, f)
			cx, cy, cz := corner(i, j+1, f)
			dx, dy, dz := corner(i+1, j+1, f)
			if math.IsNaN(ax) || math.IsNaN(ay) ||
				math.IsNaN(bx) || math.IsNaN(by) ||
				math.IsNaN(cx) || math.IsNaN(cy) ||
				math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			color := "#000000"
			if az > 0 && bz > 0 && cz > 0 && dz > 0 {
				color = "#ff0000"
			} else if az < 0 && bz < 0 && cz < 0 && dz < 0 {
				color = "#0000ff"
			}
			fmt.Printf("<polygon style='stroke: %s;' points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				color, ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int, f func(float64, float64) float64) (float64, float64, float64) {
	// Find point (x,y) at corner of cell(i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas(sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}
