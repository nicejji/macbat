package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Temperature     float32
	Cycles          int
	DesignCapacity  int
	MaxCapacity     int
	CurrentCapacity int
	IsCharging      bool
}

const CSI string = "\x1b["
const RESET string = CSI + "0m"
const BOLD string = CSI + "1m"
const NORMAL string = CSI + "22m"
const ITALIC string = CSI + "3m"
const BLUE string = CSI + "34m"
const RED string = CSI + "31m"
const YELLOW string = CSI + "33m"
const CYAN string = CSI + "36m"
const GRAY string = CSI + "90m"

func formatOpt(name string, value string, desc string, color string) string {
	return fmt.Sprintf(
		color+"%s\t%s  %s"+RESET,
		NORMAL+name,
		BOLD+value,
		GRAY+NORMAL+ITALIC+desc,
	)
}

func (info *Info) print() {
	health_percent :=
		float32(info.MaxCapacity) / float32(info.DesignCapacity) * 100
	charge_percent :=
		float32(info.CurrentCapacity) / float32(info.MaxCapacity) * 100
	charging_symbol := "⇣"
	if info.IsCharging {
		charging_symbol = "⇡"
	}
	fmt.Println(
		formatOpt(
			"Raw Health",
			fmt.Sprintf(
				"%.2f%% (%d/%d mAh)",
				health_percent,
				info.MaxCapacity,
				info.DesignCapacity,
			),
			"Raw battery health",
			BLUE,
		),
	)
	fmt.Println(
		formatOpt(
			"Temperature",
			fmt.Sprintf("%.2f °C", info.Temperature),
			"Battery temperature",
			RED,
		),
	)
	fmt.Println(
		formatOpt(
			"Cycles Count",
			fmt.Sprintf("%d", info.Cycles),
			"Cycles count",
			YELLOW,
		),
	)
	fmt.Println(
		formatOpt(
			"Charge Info",
			fmt.Sprintf(
				"%.2f%% (%d/%d mAh) %s",
				charge_percent,
				info.CurrentCapacity,
				info.MaxCapacity,
				charging_symbol,
			),
			"Charge percent",
			CYAN,
		),
	)
}

func main() {
	out, err := exec.Command(
		"ioreg",
		"-w",
		"0",
		"-r",
		"-c",
		"AppleSmartBattery",
	).Output()
	if err != nil {
		log.Fatal(err)
	}
	info := Info{}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines[2 : len(lines)-1] {
		parsed := strings.Split(strings.TrimSpace(line), "=")
		name := strings.Trim(strings.TrimSpace(parsed[0]), "\"")
		value := strings.TrimSpace(parsed[1])
		switch name {
		case "IsCharging":
			info.IsCharging = value == "Yes"
		case "AppleRawMaxCapacity":
			info.MaxCapacity, _ = strconv.Atoi(value)
		case "AppleRawCurrentCapacity":
			info.CurrentCapacity, _ = strconv.Atoi(value)
		case "DesignCapacity":
			info.DesignCapacity, _ = strconv.Atoi(value)
		case "CycleCount":
			info.Cycles, _ = strconv.Atoi(value)
		case "Temperature":
			temperature, _ := strconv.Atoi(value)
			info.Temperature = float32(temperature) / 100.0
		}
	}

	info.print()
}
