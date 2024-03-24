package main

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"kefniark/billion/shared"
	v1 "kefniark/billion/v1"
	v2 "kefniark/billion/v2"
	v3 "kefniark/billion/v3"
)

const small = "data/measurements-20.txt"
const medium = "data/measurements-10000000.txt"
const large = "data/measurements-1000000000.txt"

func main() {
	// Start CPU profiling
	f, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Version to run
	// version1()
	// version2()
	version3(large)
}

// Version 3

func version3(src string) {
	fStat, _ := os.OpenFile("stats.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fStat.Close()
	fmt.Fprint(fStat, "-- Version 3:\n")

	startParsing := time.Now()
	v3.Parse(src)
	fmt.Println("Parsing: ", time.Since(startParsing))
	fmt.Fprintf(fStat, "Parsing: %v\n", time.Since(startParsing))

	startOutput := time.Now()
	f, _ := os.Create("output.out")
	defer f.Close()
	outputToFileV3(v3.CitiesStats, f)
	fmt.Println("Output: ", time.Since(startOutput))
	fmt.Fprintf(fStat, "Output: %v\n", time.Since(startOutput))
}

func outputToFileV3(data map[uint64]*shared.CityStatV3, out io.Writer) error {
	// sort cities
	cities := make([]string, 0, len(data))
	cityLookup := map[string]uint64{}
	for city := range data {
		cityLookup[data[city].Name] = city
		cities = append(cities, data[city].Name)
	}
	sort.Strings(cities)

	// output to file
	fmt.Fprint(out, "{")
	for i, city := range cities {
		if i > 0 {
			fmt.Fprint(out, ",\n")
		}
		hash := cityLookup[city]
		s := data[hash]

		mean := float64(s.Sum) / float64(s.Count)
		fmt.Fprintf(out, "%s=%.1f/%.1f/%.1f", city, float64(s.Min)/10.0, mean/10.0, float64(s.Max)/10.0)
	}
	fmt.Fprint(out, "}\n")
	return nil
}

// Version 2

func version2(src string) {
	fStat, _ := os.OpenFile("stats.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fStat.Close()
	fmt.Fprint(fStat, "-- Version 2:\n")

	startParsing := time.Now()
	v2.Parse(src)
	fmt.Println("Parsing: ", time.Since(startParsing))
	fmt.Fprintf(fStat, "Parsing: %v\n", time.Since(startParsing))

	startOutput := time.Now()
	f, _ := os.Create("output.out")
	defer f.Close()
	outputToFileV2(v2.CitiesStats, f)
	fmt.Println("Output: ", time.Since(startOutput))
	fmt.Fprintf(fStat, "Output: %v\n", time.Since(startOutput))
}

func outputToFileV2(data map[string]*shared.CityStatV2, out io.Writer) error {
	// sort cities
	cities := make([]string, 0, len(data))
	for city := range data {
		cities = append(cities, city)
	}
	sort.Strings(cities)

	// output to file
	fmt.Fprint(out, "{")
	for i, city := range cities {
		if i > 0 {
			fmt.Fprint(out, ",\n")
		}
		s := data[city]

		mean := float64(s.Sum) / float64(s.Count)
		fmt.Fprintf(out, "%s=%.1f/%.1f/%.1f", city, float64(s.Min)/10.0, mean/10.0, float64(s.Max)/10.0)
	}
	fmt.Fprint(out, "}\n")
	return nil
}

// Version 1

func version1(src string) {
	fStat, _ := os.OpenFile("stats.out", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fStat.Close()
	fmt.Fprint(fStat, "-- Version 1:\n")

	startParsing := time.Now()
	v1.Parse(src)
	fmt.Println("Parsing: ", time.Since(startParsing))
	fmt.Fprintf(fStat, "Parsing: %v\n", time.Since(startParsing))

	startOutput := time.Now()
	f, _ := os.Create("output.out")
	defer f.Close()
	outputToFileV1(v1.CitiesStats, f)
	fmt.Println("Output: ", time.Since(startOutput))
	fmt.Fprintf(fStat, "Output: %v\n", time.Since(startOutput))
}

func outputToFileV1(data map[string]*shared.CityStat, out io.Writer) error {
	// sort cities
	cities := make([]string, 0, len(data))
	for city := range data {
		cities = append(cities, city)
	}
	sort.Strings(cities)

	// output to file
	fmt.Fprint(out, "{")
	for i, city := range cities {
		if i > 0 {
			fmt.Fprint(out, ",\n")
		}
		s := data[city]
		mean := s.Sum / float64(s.Count)
		fmt.Fprintf(out, "%s=%.1f/%.1f/%.1f", city, s.Min, mean, s.Max)
	}
	fmt.Fprint(out, "}\n")
	return nil
}
