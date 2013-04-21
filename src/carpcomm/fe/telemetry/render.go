package telemetry

import "carpcomm/pb"
import "log"
import "fmt"
import "math"
import "time"


var unitSymbol map[pb.TelemetryDatumSchema_Unit]string =
	map[pb.TelemetryDatumSchema_Unit]string{
	pb.TelemetryDatumSchema_KELVIN: "K",
	pb.TelemetryDatumSchema_VOLT: "V",
	pb.TelemetryDatumSchema_AMPERE: "A",
	pb.TelemetryDatumSchema_WATT: "W",
	pb.TelemetryDatumSchema_HERTZ: "Hz",
	pb.TelemetryDatumSchema_RADIAN_PER_SECOND: "rad/s",
	pb.TelemetryDatumSchema_SECOND: "s",
}

func bestScaleAndPrefix(unit *pb.TelemetryDatumSchema_Unit,
	values ...float64) (scale float64, prefix string) {
	// Heuristic (can be improved):
	// Use the smallest value.

	m := math.Abs(values[0])
	for _, v := range values {
		v = math.Abs(v)
		if v > m {
			m = v
		}
	}

	prefixes := []struct {
		string
		float64
	}{
		{"T", 1e12},
		{"G", 1e9},
		{"M", 1e6},
		{"k", 1e3},
		{"", 1e0},
		{"m", 1e-3},
		{"μ", 1e-6},
		{"n", 1e-9},
		{"p", 1e-12}}

	for _, s := range prefixes {
		if m >= s.float64 {
			return s.float64, s.string
		}
	}
	return 1.0, ""
}


func kelvinToDegrees(k float64) float64 {
	return k - 273.15
}

func shouldUseDegrees(k float64) bool {
	k = kelvinToDegrees(k)
	return k > -200.0 && k < 200.0
}

func shouldUseClockTime(s float64) bool {
	return s > 60.0
}

func renderSecondsUsingClockTime(confidence, s float64) string {
	if s < 1.0 {
		log.Printf("Error: invalid input to " +
			"renderSecondsUsingClockTime: %f", s)
		return fmt.Sprintf("%f s", s)
	}

	sizes := [5]int{1, 60, 60*60, 24*60*60, 365*24*60*60}
	units := [5]string{"s", "m", "h", "d", "y"}

	var values [5]int

	n := int64(s)
	for i := 4; i >= 0; i-- {
		values[i] = int(n / int64(sizes[i]))
		n = n % int64(sizes[i])
	}

	min_i := 0
	for i := 0; i < 5; i++ {
		if confidence >= float64(sizes[i]) {
			min_i = i
		}
	}

	max_i := 5
	for i := 4; i > min_i; i-- {
		if values[i] == 0 {
			max_i = i
		} else {
			break
		}
	}

	r := ""
	for i := max_i-1; i >= min_i; i-- {
		r += fmt.Sprintf("%d %s ", values[i], units[i])
	}
	return r[:len(r)-1]
}

func numDecimalDigits(confidence *float64, scale float64) int {
	if confidence == nil {
		return 1
	}
	for i := 0; i < 6; i++ {
		if scale <= *confidence {
			return i
		}
		scale = scale / 10.0
	}
	return 6
}

func renderDouble(unit *pb.TelemetryDatumSchema_Unit,
	confidence *float64, v float64) string {
	scale, prefix := bestScaleAndPrefix(unit, v)

	u := ""
	if unit != nil {
		u = unitSymbol[*unit]
		if *unit == pb.TelemetryDatumSchema_KELVIN &&
			shouldUseDegrees(v) {
			v = kelvinToDegrees(v)
			scale = 1.0
			u = "°C"
		}
		if u == "" {
			log.Printf("Error: Missing unit string for %v.", *unit)
		}
	}

	if unit != nil && *unit == pb.TelemetryDatumSchema_SECOND {
		if shouldUseClockTime(v) {
			c := 1.0
			if confidence != nil {
				c = *confidence
			}
			return renderSecondsUsingClockTime(c, v)
		}
	}

	n := numDecimalDigits(confidence, scale)
	format := fmt.Sprintf("%%.%df %%s%%s", n)
	return fmt.Sprintf(format, v/scale, prefix, u)
}

func renderInterval(unit *pb.TelemetryDatumSchema_Unit,
	confidence *float64, min, max float64) string {
	scale, prefix := bestScaleAndPrefix(unit, min, max)

	u := ""
	if unit != nil {
		u = unitSymbol[*unit]
		if *unit == pb.TelemetryDatumSchema_KELVIN &&
			shouldUseDegrees(min) || shouldUseDegrees(max) {
			min = kelvinToDegrees(min)
			max = kelvinToDegrees(max)
			scale = 1.0
			u = "°C"
		}
	}

	n := numDecimalDigits(confidence, scale)
	format := fmt.Sprintf("[%%.%df, %%.%df) %%s%%s", n, n)
	return fmt.Sprintf(format, min/scale, max/scale, prefix, u)
}

func renderBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func renderInt64(i int64) string {
	return fmt.Sprintf("%d", i)
}

const TimeFormat = "2006-01-02 15:04 MST"
func renderTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(TimeFormat)
}


func EngNotationHz(confidence, v float64) string {
	u := pb.TelemetryDatumSchema_HERTZ
	return renderDouble(&u, &confidence, v)
}

func EngNotationWatt(confidence, v float64) string {
	u := pb.TelemetryDatumSchema_WATT
	return renderDouble(&u, &confidence, v)
}


type LabelValue struct {
	Label string
	Value string
	Timestamp string
}

func getName(s *pb.TelemetryDatumSchema, lang string) string {
	if len(s.Name) == 0 {
		log.Printf("No names for datum %s.", *s.Key)
		return "Unknown"
	}
	for _, n := range s.Name {
		if n.Lang != nil && *n.Lang == lang {
			return *n.Text
		}
	}
	log.Printf("Missing name in %s for datum %s.", lang, *s.Key)
	return *s.Name[0].Text
}

func renderValue(s *pb.TelemetryDatumSchema,
	d pb.TelemetryDatum, lang string) string {
	switch *s.Type {
	case pb.TelemetryDatumSchema_BOOL:
		return renderBool(*d.Boolean)
	case pb.TelemetryDatumSchema_INT64:
		return renderInt64(*d.Int64)
	case pb.TelemetryDatumSchema_DOUBLE:
		return renderDouble(s.Unit, s.Confidence, *d.Double)
	case pb.TelemetryDatumSchema_INTERVAL:
		return renderInterval(s.Unit, s.Confidence,
			*d.IntervalMin, *d.IntervalMax)
	case pb.TelemetryDatumSchema_TIMESTAMP:
                return renderTimestamp(*d.UnixTimestamp)
	}
	log.Printf("Rendered unrecognized datum")
	return ""
}

func RenderTelemetry(schema pb.TelemetrySchema,	data []pb.TelemetryDatum,
	lang string) (r [][]LabelValue) {
	keyToDatum := make(map[string]pb.TelemetryDatum)
	for _, d := range data {
		existing, ok := keyToDatum[*d.Key]
		// If there are multiple datums for the same key, use the most
		// recent.
		if !ok || *existing.Timestamp < *d.Timestamp {
			keyToDatum[*d.Key] = d
		}
	}
	var lastGroup int32 = 0
	g := make([]LabelValue, 0)
	for _, s := range schema.Datum {
		d, ok := keyToDatum[*s.Key]
		if !ok {
			continue
		}

		// Insert a separator between groups.
		var group int32
		if s.DisplayGroup != nil {
			group = *s.DisplayGroup
		}
		if len(g) > 0 && lastGroup != group {
			r = append(r, g)
			g = make([]LabelValue, 0)
		}
		lastGroup = group

		g = append(g, LabelValue{
			getName(s, lang),
			renderValue(s, d, lang),
			renderTimestamp(*d.Timestamp)})
	}
	if len(g) > 0 {
		r = append(r, g)
	}
	return r
}