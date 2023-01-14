package main

import (
	"flag"
	"fmt"
    "sync"
    "time"

	pkg "github.com/DylanMeeus/GoAudio/wave"
)

var (
	input  = flag.String("i", "", "input file")
	output = flag.String("o", "", "output file")
	amp    = flag.Float64("a", 1.0, "amp mod factor")
)

const SAMPLES_SIZE int = 100000

func Map[T, U any](ts []T, f func(T) U) []U {
    us := make([]U, len(ts))
    for i := range ts {
        us[i] = f(ts[i])
    }
    return us
}

var res []pkg.Frame

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
	fmt.Println("Parsing wave file..")
	flag.Parse()
	infile := *input
	outfile := *output
	scale := *amp
	wave, err := pkg.ReadWaveFile(infile)
	if err != nil {
		panic("Could not parse wave file")
	}

	fmt.Printf("Read %v samples\n", len(wave.Frames))

    var wg sync.WaitGroup
    start_t := time.Now()
    res = make([]pkg.Frame, len(wave.Frames))
    var ires int = 0
    for i := 0; i < len(wave.Frames); i += (SAMPLES_SIZE - 1) {
        var to_compute []pkg.Frame = wave.Frames[i:min(i+(SAMPLES_SIZE), len(wave.Frames) - 1)]
        wg.Add(1)
        go changeAmplitude(i, to_compute, scale, &wg)
        ires += min(SAMPLES_SIZE - 1, len(wave.Frames)-i)
    }

    wg.Wait()
    end_t := time.Now()
    fmt.Println("Finished, computed", ires, "samples in", end_t.Sub(start_t))

	if err := pkg.WriteFrames(res, wave.WaveFmt, outfile); err != nil {
    	panic(err)
	}

	fmt.Println("done")
}

func changeAmplitude(startIndex int, samples []pkg.Frame, scalefactor float64, wg *sync.WaitGroup) {
    for i, s := range samples {
        res[startIndex + i] = pkg.Frame(float64(s) * scalefactor)
    }
    wg.Done()
}
