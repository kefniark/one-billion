package v2

import (
	"math"
	"os"
	"sync"

	"kefniark/billion/shared"
)

var CitiesStats = map[string]*shared.CityStatV2{}
var mutex = &sync.Mutex{}

const symbolSeparator = byte(';')
const symbolDot = byte('.')
const symbolNegative = byte('-')
const symbolLineBreak = byte('\n')

func Parse(name string) error {
	var wg sync.WaitGroup

	for i := range 8 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			entries, _ := parseByChunk(name, 4*1024*1024, i, 8)
			mutex.Lock()
			for city, stat := range entries {
				data, ok := CitiesStats[city]
				if ok {
					data.Min = min(data.Min, stat.Min)
					data.Max = max(data.Max, stat.Max)
					data.Sum += stat.Sum
					data.Count += stat.Count
				} else {
					data = stat
				}
				CitiesStats[city] = stat
			}
			mutex.Unlock()
		}(i)
	}
	wg.Wait()

	return nil
}

func parseByChunk(name string, bufferSize int64, workerId int, workerTotal int64) (map[string]*shared.CityStatV2, error) {
	var workerStats = map[string]*shared.CityStatV2{}

	readFile, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer readFile.Close()

	buffer := make([]byte, bufferSize+40)

	off := int64(0)
	chunkId := int64(workerId)
	for {
		if chunkId > 0 {
			off = bufferSize*chunkId - 40
		}
		size, err := readFile.ReadAt(buffer, off)
		if chunkId == 0 {
			size -= 40
		}
		if err != nil {
			return workerStats, err
		}

		start := 0
		if chunkId > 0 {
			start = 40
			for start > 0 {
				if buffer[start] == symbolLineBreak {
					break
				}
				start--
			}
			// fmt.Println("Chunk", chunkId, "Start", start, string(buffer[0:start]), "|", string(buffer[start:25]))
		}

		// fmt.Println("Chunk", chunkId, string(buffer))
		for i := start; i < size; i++ {
			if buffer[i] == symbolLineBreak {
				parseLine(workerStats, buffer, start, i)
				start = i + 1
			}
		}
		chunkId += workerTotal
	}
	// return nil
}

func parseLine(stats map[string]*shared.CityStatV2, buffer []byte, start, end int) {
	separator, dot := 0, 0
	// fmt.Println("ParseLine", string(buffer[start:end]))

	for i := start; i < end; i++ {
		if buffer[i] == symbolSeparator {
			separator = i
		} else if buffer[i] == symbolDot {
			dot = i
		} else if i == end-1 {
			city := string(buffer[start:separator])
			temp := parseTemperature(buffer, separator, dot, end)
			// fmt.Println("City", city, "Temp", temp)
			stat, ok := stats[city]
			if ok {
				stat.Min = min(stat.Min, temp)
				stat.Max = max(stat.Max, temp)
				stat.Sum += temp
				stat.Count++
			} else {
				stat = &shared.CityStatV2{Min: temp, Max: temp, Sum: temp, Count: 1}
			}
			stats[city] = stat
		}
	}
}

func parseTemperature(buffer []byte, separator, dot, end int) int {
	negative := false
	temp := 0
	for index := separator + 1; index < end; index++ {
		if buffer[index] == symbolNegative {
			negative = true
			continue
		}
		if index == dot {
			continue
		}

		pos := end - index - 1
		if index < dot {
			pos -= 1
		}

		// use ascii code offest to get the int, and power of 10 to get the correct position
		temp += (int(buffer[index]) - 48) * int(math.Pow10(pos))
	}
	if negative {
		temp *= -1
	}
	return temp
}
