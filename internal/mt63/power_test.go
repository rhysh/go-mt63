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
