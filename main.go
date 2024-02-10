package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	temperature      float32
	cycles           int
	design_capacity  int
	max_capacity     int
	current_capacity int
	is_charging      bool
}

func (info *Info) print() {
	health_percent := float32(info.max_capacity) / float32(info.design_capacity)
	charging_symbol := "⇣"
	if info.is_charging {
		charging_symbol = "⇡"
	}
	fmt.Printf("\x1b[34mhealth\t\t%.2f%s (%d/%d mAh)\x1b[0m\n", health_percent*100, "%", info.max_capacity, info.design_capacity)
	fmt.Printf("\x1b[31mtemp\t\t%.2f °C\x1b[0m\n", info.temperature)
	fmt.Printf("\x1b[33mcycles\t\t%d\x1b[0m\n", info.cycles)
	fmt.Printf("\x1b[36mcharge\t\t%.2f%s %s\x1b[0m\n", float32(info.current_capacity)/float32(info.max_capacity)*100, "%", charging_symbol)
}

func main() {
	out, err := exec.Command("ioreg", "-w", "0", "-r", "-c", "AppleSmartBattery").Output()
	if err != nil {
		log.Fatal(err)
	}
	info := Info{}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines[2 : len(lines)-1] {
		parsed := strings.Split(strings.TrimSpace(line), "=")
		name := strings.Trim(strings.TrimSpace(parsed[0]), "\"")
		value := strings.TrimSpace(parsed[1])
		if name == "IsCharging" {
			info.is_charging = value == "Yes"
		}
		if name == "AppleRawMaxCapacity" {
			max_capacity, _ := strconv.Atoi(value)
			info.max_capacity = max_capacity
		}
		if name == "AppleRawCurrentCapacity" {
			current_capacity, _ := strconv.Atoi(value)
			info.current_capacity = current_capacity
		}
		if name == "DesignCapacity" {
			design_capacity, _ := strconv.Atoi(value)
			info.design_capacity = design_capacity
		}
		if name == "CycleCount" {
			cycles, _ := strconv.Atoi(value)
			info.cycles = cycles
		}
		if name == "Temperature" {
			temperature, _ := strconv.Atoi(value)
			info.temperature = float32(temperature) / 100.0
		}
	}

	info.print()
}
