package mt63

import (
	"fmt"
	"math"
)

type Bandwidth int

const (
	Bw500 = Bandwidth(500)
	Bw1k  = Bandwidth(1000)
	Bw2k  = Bandwidth(2000)
)

type Interleaving int

const (
	Short = Interleaving(32)
	Long  = Interleaving(64)
)

type Carrier int

const (
	baseFreq = 500.0
)

type MT63 struct {
	Bandwidth    Bandwidth
	Interleaving Interleaving
}

func (m *MT63) String() string {
	bw := map[Bandwidth]string{
		Bw500: "500",
		Bw1k:  "1k",
		Bw2k:  "2k",
	}[m.Bandwidth]
	inter := map[Interleaving]string{
		Short: "S",
		Long:  "L",
	}[m.Interleaving]
	if bw == "" || inter == "" {
		return "[invalid]"
	}
	return "mt63-" + bw + inter
}

func (m *MT63) Frequency(c Carrier) float64 {
	if 0 <= c && c <= 63 {
		return float64(m.Bandwidth)/float64(64)*float64(c) + baseFreq
	}
	panic(fmt.Sprintf("invalid carrier %d", c))
}

type Decoder struct {
	Mode       *MT63
	SampleRate float64
}

func (d *Decoder) PowerAt(c Carrier, vals []float64) float64 {
	radPerSecond := d.Mode.Frequency(c) * (2 * math.Pi)
	var sinSum, cosSum float64

	carrierCycle := 64 / float64(d.Mode.Bandwidth)

	for i, val := range vals {
		t := float64(i) / d.SampleRate
		if t >= carrierCycle { // integer cycle count for all carriers
			break
		}
		sin, cos := math.Sincos(t * radPerSecond)
		sinSum += sin * val
		cosSum += cos * val
	}
	return sinSum*sinSum + cosSum*cosSum
}
