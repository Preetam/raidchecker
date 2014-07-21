package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/PreetamJinka/raidcheck"
)

func main() {
	hostsStr := flag.String("hosts", "", "comma separated list of hosts")
	flag.Parse()

	var hosts []string
	if len(*hostsStr) > 0 {
		hosts = strings.Split(*hostsStr, ",")
	}

	wg := sync.WaitGroup{}
	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			checkHost(host)
			wg.Done()
		}(host)
	}

	wg.Wait()
}

func checkHost(host string) {
	printMessage(host, "connecting")

	ls := exec.Command("ssh", "root@"+host, "cat /proc/mdstat")
	out, err := ls.Output()
	if err != nil {
		printMessage(host, "error: "+err.Error())
		return
	}

	degraded := raidcheck.CheckDegraded(string(out))
	printMessage(host, fmt.Sprintf("degraded RAID: %v", degraded))
}

func printMessage(host, msg string) {
	fmt.Printf("[%s] %s\n", host, msg)
}
