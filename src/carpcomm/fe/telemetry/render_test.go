package telemetry

import "testing"

func TestRenderSecondsUsingClockTime(t *testing.T) {
	test := func(conf, s float64, expected string) {
		actual := renderSecondsUsingClockTime(conf, s)
		if actual != expected {
			t.Errorf("Actual: %s, expected: %s", actual, expected)
		}
	}

	const second = 1
	test(second, 10.0, "10 s")
	test(second, 300.0, "5 m 0 s")
	test(second, 317.0, "5 m 17 s")
	test(second, 9125.0, "2 h 32 m 5 s")
	test(second, 441125.0, "5 d 2 h 32 m 5 s")
	test(second, 31977125.0, "1 y 5 d 2 h 32 m 5 s")

	const minute = 60*second
	test(minute, 10.0, "0 m")
	test(minute, 300.0, "5 m")
	test(minute, 317.0, "5 m")
	test(minute, 9125.0, "2 h 32 m")
	test(minute, 441125.0, "5 d 2 h 32 m")
	test(minute, 31977125.0, "1 y 5 d 2 h 32 m")

	const hour = 60*minute
	test(hour, 10.0, "0 h")
	test(hour, 300.0, "0 h")
	test(hour, 317.0, "0 h")
	test(hour, 9125.0, "2 h")
	test(hour, 441125.0, "5 d 2 h")
	test(hour, 31977125.0, "1 y 5 d 2 h")

	const day = 24*hour
	test(day, 10.0, "0 d")
	test(day, 300.0, "0 d")
	test(day, 317.0, "0 d")
	test(day, 9125.0, "0 d")
	test(day, 441125.0, "5 d")
	test(day, 31977125.0, "1 y 5 d")

	const year = 365*day
	test(year, 10.0, "0 y")
	test(year, 300.0, "0 y")
	test(year, 317.0, "0 y")
	test(year, 9125.0, "0 y")
	test(year, 441125.0, "0 y")
	test(year, 31977125.0, "1 y")

	// Invalid cases.
	test(0.1, 0.001, "0.001000 s")
	test(0.1, -5.0, "-5.000000 s")
}