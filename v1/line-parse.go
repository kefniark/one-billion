package v1

import (
	"bufio"
	"os"
	"strconv"

	"kefniark/billion/shared"
)

var CitiesStats = map[string]*shared.CityStat{}

const symbolSeparator = byte(';')
const symbolDot = byte('.')
const symbolLineBreak = byte('\n')

func Parse(name string) error {
	readFile, err := os.Open(name)
	if err != nil {
		return err
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	i := 0
	for fileScanner.Scan() {
		i++
		bytes := fileScanner.Bytes()
		parseLine(bytes, len(bytes))
		// if i >= 5 {
		// 	break
		// }
	}

	return nil
}

func parseLine(buffer []byte, size int) {
	start, separator, _ := 0, 0, 0
	// fmt.Println("Size", buffer, size)

	for i := range size {
		if buffer[i] == symbolSeparator {
			separator = i
		} else if buffer[i] == symbolDot {
			// dot = i
		} else if buffer[i] == symbolLineBreak || i == size-1 {
			end := i - 1
			if i == size-1 {
				end = i + 1
			}
			// fmt.Println(
			// 	start, separator, dot, i,
			// 	string(buffer[start:separator]), string(buffer[separator+1:dot]), string(buffer[dot+1:end]),
			// )
			city := string(buffer[start:separator])
			temp, _ := strconv.ParseFloat(string(buffer[separator+1:end]), 64)

			stat, ok := CitiesStats[city]
			if ok {
				stat.Min = min(stat.Min, temp)
				stat.Max = max(stat.Max, temp)
				stat.Sum += temp
				stat.Count++
			} else {
				stat = &shared.CityStat{Min: temp, Max: temp, Sum: temp, Count: 1}
			}

			CitiesStats[city] = stat
			start = i + 1
		}
	}
}
