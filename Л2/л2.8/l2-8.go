package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

const ntpServer = "pool.ntp.org"

func main() {
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к NTP серверу: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Время согласно NTP серверу: %s\n", ntpTime)
}
