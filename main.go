package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"io/ioutil"

	"strings"

	"flag"

	"github.com/tatsushid/go-fastping"
)

var (
	result PairList

	filename = flag.String("h", "host.txt", "determine the location of host-list file")
	times    = flag.Int("t", 5, "determine the times of ping")
	reserve  = flag.Bool("r", false, "reserve result order")
)

func main() {
	flag.Parse()

	hosts, err := readHostsFromFile(*filename)
	if err != nil {
		fmt.Printf("Read host file error:%+v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, host := range hosts {
		if len(strings.TrimSpace(host)) == 0 {
			continue
		}

		fmt.Printf("ping %s ...\n", host)
		wg.Add(1)
		go doPing(host, *times, &wg)
	}

	wg.Wait()
	fmt.Println("\n------------ RESULT ------------")
	if *reserve {
		sort.Sort(sort.Reverse(result))
	} else {
		sort.Sort(result)
	}

	for _, p := range result {
		fmt.Printf("%-20s %.2fms\n", p.Key, p.Value)
	}
}

func doPing(host string, times int, wg *sync.WaitGroup) {
	defer wg.Done()

	if avg, err := pingIP(host, times); err != nil {
		fmt.Printf("%-20s -> ERR:%+v\n", host, err)
	} else {
		result = append(result, Pair{host, avg})
	}
}

func readHostsFromFile(p string) (hosts []string, err error) {
	b, e := ioutil.ReadFile(p)
	if e != nil {
		err = e
		return
	}

	hosts = strings.Split(string(b), "\n")
	return
}

func pingIP(host string, times int) (avg float64, err error) {
	ra, e := net.ResolveIPAddr("ip4:icmp", host)
	if e != nil {
		err = e
		return
	}

	total := float64(0)
	p := fastping.NewPinger()
	p.AddIPAddr(ra)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		//fmt.Printf("%s %s %s\n", host, addr.String(), rtt)
		total = total + float64(rtt.Nanoseconds())
	}

	for i := 0; i < times; i++ {
		if err = p.Run(); err != nil {
			return
		}
	}

	avg = float64(total) / float64(times) / 1000000
	return
}

type Pair struct {
	Key   string
	Value float64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
