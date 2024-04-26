package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func DropFileData(path string) (string, int, int, int) {
	// path needs to include trailing slash!
	var dropAlias string
	var dropTimeLeft string
	var dropEmulation string
	var nodeNum string

	file, err := os.Open(strings.ToLower(path + "door32.sys"))
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	count := 0
	for _, line := range text {
		if count == 6 {
			dropAlias = line
		}
		if count == 8 {
			dropTimeLeft = line
		}
		if count == 9 {
			dropEmulation = line
		}
		if count == 10 {
			nodeNum = line
		}
		if count == 11 {
			break
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	timeInt, err := strconv.Atoi(dropTimeLeft) // return as int
	if err != nil {
		log.Fatal(err)
	}

	emuInt, err := strconv.Atoi(dropEmulation) // return as int
	if err != nil {
		log.Fatal(err)
	}
	nodeInt, err := strconv.Atoi(nodeNum) // return as int
	if err != nil {
		log.Fatal(err)
	}

	return dropAlias, timeInt, emuInt, nodeInt
}
