package main

import "flag"

func main() {
	source := flag.String("wav", "", "Path to WAV file to decode (or blank to read samples from stdin)")
	rate := flag.Int("rate", 44100, "Sample rate (when reading samples from stdin)")
	flag.Parse()

	_ = source
	_ = rate
}
