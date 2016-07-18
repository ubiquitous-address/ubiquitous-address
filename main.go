package main

import (
	"fmt"
	"log"
	"math"

	"github.com/golang/geo/s1"
)

// Earth's radius in kilometers
const (
	R       = 6371.0
	Radians = math.Pi / 180.0
)

//
type Degree struct {
	degrees, minutes, seconds int64
}

//
func (arc Degree) String() string {
	return fmt.Sprintf(`%03d°%02d'%02d"`, arc.degrees, arc.minutes, arc.seconds)
}

//
func (arc Degree) NS() string {
	d := 'N'
	if arc.degrees < 0 {
		d = 'S'
	}
	return fmt.Sprintf(`%s%c`, arc, d)
}

//
func (arc Degree) EW() string {
	d := 'E'
	if arc.degrees < 0 {
		d = 'W'
	}
	return fmt.Sprintf(`%s%c`, arc, d)
}

//
func Round(a float64) int64 {
	if a < 0 {
		return int64(math.Ceil(a - 0.5))
	}
	return int64(math.Floor(a + 0.5))
}

/*
   radius = (radius === undefined) ? 6371e3 : Number(radius);

   // φ2 = asin( sinφ1⋅cosδ + cosφ1⋅sinδ⋅cosθ )
   // λ2 = λ1 + atan2( sinθ⋅sinδ⋅cosφ1, cosδ − sinφ1⋅sinφ2 )
   // see http://williams.best.vwh.net/avform.htm#LL

   var δ = Number(distance) / radius; // angular distance in radians
   var θ = Number(bearing).toRadians();

   var φ1 = this.lat.toRadians();
   var λ1 = this.lon.toRadians();

   var φ2 = Math.asin(Math.sin(φ1)*Math.cos(δ) + Math.cos(φ1)*Math.sin(δ)*Math.cos(θ));
   var x = Math.cos(δ) - Math.sin(φ1) * Math.sin(φ2);
   var y = Math.sin(θ) * Math.sin(δ) * Math.cos(φ1);
   var λ2 = λ1 + Math.atan2(y, x);

   return new LatLon(φ2.toDegrees(), (λ2.toDegrees()+540)%360-180); // normalise to −180..+180°
*/

//
func move(d, brng, lng_, lat_ float64) (lng, lat float64) {
	// log.Printf("%8.2f:\t%.9f,%.9f", d, lng_, lat_)

	δ := d / R
	cosδ := math.Cos(δ)
	sinδ := math.Sin(δ)
	// log.Printf("δ \t%.9f\t%.9f\t%.9f", δ, cosδ, sinδ)

	θ := brng * Radians
	cosθ := math.Cos(θ)
	sinθ := math.Sin(θ)
	// log.Printf("θ \t%.9f\t%.9f\t%.9f", θ, cosθ, sinθ)

	φ1 := lat_ * Radians
	cosφ1 := math.Cos(φ1)
	sinφ1 := math.Sin(φ1)
	// log.Printf("φ1\t%.9f\t%.9f\t%.9f", φ1, cosφ1, sinφ1)

	λ1 := lng_ * Radians
	// log.Printf("λ1\t%.9f", λ1)

	φ2 := math.Asin(sinφ1*cosδ + cosφ1*sinδ*cosθ)
	// log.Printf("φ2\t%.9f\t%.9f", φ2, φ2-φ1)

	sinφ2 := math.Sin(φ2)

	x := cosδ - sinφ1*sinφ2
	// log.Printf("x \t%.9f\t%.9f", x, sinφ2)
	y := sinθ * sinδ * cosφ1
	// log.Printf("y \t%.9f", y)

	λ2 := λ1 + math.Atan2(y, x)
	// log.Printf("λ2\t%.9f", λ2)

	lat = s1.Angle(φ2).Degrees()
	lng = s1.Angle(λ2).Degrees()

	return
}

//
func arc(coord float64) (arc Degree) {
	arc.seconds = Round(coord * 3600.0)
	arc.degrees = arc.seconds / 3600
	arc.seconds = int64(math.Abs(float64(arc.seconds % 3600)))
	arc.minutes = arc.seconds / 60
	arc.seconds %= 60
	return
}

//
func main() {
	log.Printf("%+v", "/usr/share/dict/words")

	var lng, lat = 0.0, 0.0
	var increment = 0.25 // kilometers

	for n := 0; lat <= 89.99; n += 1 {
		lng = 0.0
		for w := 0; lng <= 89.99; w += 1 {
			lng, lat = move(increment, 90.0, lng, lat)
			log.Printf("==  %d,%d\t%13.9f,%13.9f\t%s %s", n, w, lng, lat, arc(lat).NS(), arc(lng).EW())
		}
		lng, lat = move(increment, 0.0, lng, lat)
	}
}
