package mt63

import (
	"fmt"
	"testing"
)

func TestCarrierFrequency(t *testing.T) {
	testcase := func(v *MT63, c Carrier, want float64) func(t *testing.T) {
		return func(t *testing.T) {
			if have := v.Frequency(c); have != want {
				t.Errorf("Carrier %d of %s; %f != %f", c, v, have, want)
			}
		}
	}

	mt63_2kL := &MT63{Bandwidth: Bw2k, Interleaving: Long}

	t.Run("", testcase(mt63_2kL, 0, 500))
	t.Run("", testcase(mt63_2kL, 1, 531.25))
	t.Run("", testcase(mt63_2kL, 4, 625))
}

func TestModeName(t *testing.T) {
	testcase := func(v *MT63, want string) func(t *testing.T) {
		return func(t *testing.T) {
			if have := fmt.Sprint(v); have != want {
				t.Errorf("mode %#v; %s != %s", v, have, want)
			}
		}
	}

	t.Run("", testcase(&MT63{Bandwidth: Bw500, Interleaving: Short}, "mt63-500S"))
	t.Run("", testcase(&MT63{Bandwidth: Bw2k, Interleaving: Long}, "mt63-2kL"))
}

func TestPowerAt(t *testing.T) {
	// Testing 500 Hz square wave "carrier" at 2000 Hz sample rate
	mode := &MT63{Bandwidth: Bw2k, Interleaving: Long}
	dec := &Decoder{Mode: mode, SampleRate: 2000}

	square1 := []float64{ // 64 ms
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,

		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
	}
	square2 := []float64{ // 64 ms with phase reversal
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,
		-1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1,

		1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1,
		1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1,
		1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1,
		1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1,
	}
	square3 := []float64{
		-1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1,
		-1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1,
		-1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1,
		-1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1,
	}
	square4 := []float64{
		1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1,
		1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1,
		1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1,
		1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1, 1, -1, -1, 1,
	}

	testLarge := func(wave []float64) func(t *testing.T) {
		return func(t *testing.T) {
			have := dec.PowerAt(0, wave)
			want := float64(64 * 64 / 2)
			if have < 0.9*want || have > 1.1*want {
				t.Errorf("want large magnitude; %f != %f", have, want)
			}
		}
	}

	testSmall := func(wave []float64) func(t *testing.T) {
		return func(t *testing.T) {
			have := dec.PowerAt(0, wave)
			high := float64(64 * 64 / 2)
			if have > 0.01*high {
				t.Errorf("want small magnitude; %f != %f", have, 0.0)
			}
		}
	}

	// We expect PowerAt to look at only a 32 ms window of data (for the 2 kHz
	// variant). We expect it to be indifferent to the phase of the signal.
	//
	// But, we expect it to detect very low power when it's centered on a phase
	// reversal.

	t.Run("double-length square wave", testLarge(square1))
	t.Run("square head", testLarge(square1[:64]))
	t.Run("square tail", testLarge(square1[64:]))
	t.Run("square middle", testLarge(square1[32:]))

	t.Run("double-length phase reversal", testLarge(square2))
	t.Run("phase reversal head", testLarge(square2[:64]))
	t.Run("phase reversal tail", testLarge(square2[64:]))
	t.Run("phase reversal", testSmall(square2[32:]))

	t.Run("90° shift", testLarge(square3))
	t.Run("270° shift", testLarge(square4))
}
