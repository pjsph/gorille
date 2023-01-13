package main

import (
	"flag"
	"fmt"

	pkg "github.com/DylanMeeus/GoAudio/wave"
)

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	amp    = flag.Float64("a", 1.0, "amp mod factor")
)

func main() {
	fmt.Println("Parsing wave file..")
	flag.Parse()
	infile := *input
	//outfile := *output
	scale := *amp
	wave, err := pkg.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	fmt.Printf("Read %v samples\n", len(wave.Frames))

	// now try to write this file
    for i, s := range wave.Frames {
        ch := make(chan struct {int; pkg.Frame})
        go changeAmplitude(ch, i, s, scale)
        res := <- ch
        fmt.Println(res.int, res.Frame)
    }

	//scaledSamples := changeAmplitude(wave.Frames, scale)
	//if err := pkg.WriteFrames(scaledSamples, wave.WaveFmt, outfile); err != nil {
    //	panic(err)
	//}

	fmt.Println("done")
}

func changeAmplitude(ch chan struct {int; pkg.Frame}, i int, sample pkg.Frame, scalefactor float64) {
        ch <- struct {int; pkg.Frame}{i, pkg.Frame(float64(sample) * scalefactor)}
}
