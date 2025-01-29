package geo

import "math"

const (
	Meters RadiusType = 6378137
	Km     RadiusType = 6378
)

type RadiusType int

func (r RadiusType) float() float64 {
	return float64(r)
}

type Cord struct {
	Lat float64
	Lon float64
}

func HaversineDistance(p, q Cord, r RadiusType) float64 {
	toRadians := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	lat1Rad, lon1Rad := toRadians(p.Lat), toRadians(p.Lon)
	lat2Rad, lon2Rad := toRadians(q.Lat), toRadians(q.Lon)

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return r.float() * c
}
