package anticheat

import (
	"math"
	"time"

	"github.com/territory-run/api/internal/models"
)

type Result struct {
	TrustScore int
	Status     models.ActivityStatus
	Reason     string
}

// Analyze performs basic Phase 1 anti-cheat checks on GPS points.
func Analyze(points []models.GPSPoint, maxSpeed, maxAccuracy float64, minDurationSec int) Result {
	if len(points) < 2 {
		return Result{TrustScore: 0, Status: models.ActivityRejected, Reason: "too few points"}
	}

	duration := points[len(points)-1].Timestamp - points[0].Timestamp
	if duration < int64(minDurationSec) {
		return Result{TrustScore: 30, Status: models.ActivityQuarantined, Reason: "duration too short"}
	}

	score := 100
	for i := 1; i < len(points); i++ {
		prev, cur := points[i-1], points[i]
		if cur.Accuracy > maxAccuracy && prev.Accuracy > maxAccuracy {
			score -= 2
			continue
		}

		dt := float64(cur.Timestamp - prev.Timestamp)
		if dt <= 0 {
			continue
		}
		dist := haversineM(prev.Lat, prev.Lon, cur.Lat, cur.Lon)
		speed := dist / dt
		if speed > maxSpeed {
			return Result{TrustScore: 0, Status: models.ActivityRejected, Reason: "impossible speed detected"}
		}
		if speed > maxSpeed*0.85 {
			score -= 5
		}
	}

	if score < 50 {
		return Result{TrustScore: score, Status: models.ActivityQuarantined, Reason: "low gps quality"}
	}
	return Result{TrustScore: score, Status: models.ActivityVerified, Reason: ""}
}

func haversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * R * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// FilterPoints removes outliers before territory processing.
func FilterPoints(points []models.GPSPoint, maxSpeed, maxAccuracy float64) []models.GPSPoint {
	if len(points) == 0 {
		return points
	}
	out := []models.GPSPoint{points[0]}
	for i := 1; i < len(points); i++ {
		cur := points[i]
		if cur.Accuracy > maxAccuracy {
			continue
		}
		prev := out[len(out)-1]
		dt := float64(cur.Timestamp - prev.Timestamp)
		if dt <= 0 {
			continue
		}
		dist := haversineM(prev.Lat, prev.Lon, cur.Lat, cur.Lon)
		if dist/dt > maxSpeed {
			continue
		}
		out = append(out, cur)
	}
	return out
}

func DurationSec(points []models.GPSPoint) int {
	if len(points) < 2 {
		return 0
	}
	return int(points[len(points)-1].Timestamp - points[0].Timestamp)
}

func DistanceM(points []models.GPSPoint) float64 {
	var total float64
	for i := 1; i < len(points); i++ {
		total += haversineM(points[i-1].Lat, points[i-1].Lon, points[i].Lat, points[i].Lon)
	}
	return total
}

func NowUnix() int64 {
	return time.Now().Unix()
}
