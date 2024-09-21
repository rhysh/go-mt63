package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/rhysh/go-mt63/internal/mt63"
)

type modeFlag mt63.MT63

var _ flag.Value

func (f *modeFlag) String() string {
	return (*mt63.MT63)(f).String()
}

func (f *modeFlag) Set(v string) error {
	switch v {
	default:
		return fmt.Errorf("invalid mode: %q", v)
	case "2kL":
		f.Bandwidth = mt63.Bw2k
		f.Interleaving = mt63.Long
	case "2kS":
		f.Bandwidth = mt63.Bw2k
		f.Interleaving = mt63.Short
	case "1kL":
		f.Bandwidth = mt63.Bw1k
		f.Interleaving = mt63.Long
	case "1kS":
		f.Bandwidth = mt63.Bw1k
		f.Interleaving = mt63.Short
	case "500L":
		f.Bandwidth = mt63.Bw500
		f.Interleaving = mt63.Long
	case "500S":
		f.Bandwidth = mt63.Bw500
		f.Interleaving = mt63.Short
	}
	return nil
}

func main() {
	source := flag.String("wav", "", "Path to WAV file to decode (or blank to read samples from stdin)")
	rate := flag.Int("rate", 44100, "Sample rate (when reading samples from stdin)")
	mode := &mt63.MT63{Bandwidth: mt63.Bw2k, Interleaving: mt63.Long}
	flag.Var((*modeFlag)(mode), "mode", "MT63 variant, such as '2kL' or '500S'")
	flag.Parse()

	_ = source

	dec := &mt63.Decoder{Mode: mode, SampleRate: float64(*rate)}

	var samples []float64
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		v, err := strconv.ParseFloat(sc.Text(), 64)
		if err != nil {
			fmt.Printf("Could not parse floating point input %q: %v\n", sc.Text(), err)
			os.Exit(1)
		}
		samples = append(samples, v)
	}

	symbolHz := int(dec.Mode.Bandwidth) / 100
	const looksPerSymbol = 4 // Inspect power levels four times per symbol, to decode clock

	carrierCycle := 64 / float64(dec.Mode.Bandwidth)
	windowLen := int(carrierCycle * dec.SampleRate)

	deltaWindow := make([]float64, windowLen)
	for i := 0; ; i++ {
		symbol := float64(i) / looksPerSymbol
		startTime := symbol / float64(symbolHz)
		startSample := int(startTime * dec.SampleRate)
		if startSample >= len(samples) {
			break
		}
		if startSample+windowLen*2 > len(samples) {
			break
		}
		doubleWindow := samples[startSample : startSample+windowLen*2]
		currWindow := doubleWindow[:windowLen]
		nextWindow := doubleWindow[windowLen:]
		for i := range deltaWindow {
			deltaWindow[i] = currWindow[i] - nextWindow[i]
		}
		for c := mt63.Carrier(0); c < 64; c++ {
			currAmpl := math.Sqrt(dec.PowerAt(c, currWindow))
			nextAmpl := math.Sqrt(dec.PowerAt(c, nextWindow))
			deltaAmpl := math.Sqrt(dec.PowerAt(c, deltaWindow))
			alignment := 0.0
			if sumAmpl := currAmpl + nextAmpl; sumAmpl > 0 {
				alignment = deltaAmpl / sumAmpl

			}
			fmt.Printf("%f %d %f %f %f %f\n", symbol, c, currAmpl, nextAmpl, deltaAmpl, alignment)
		}
	}
}
