// usage:
// uvarint.exe >> pprof.log
// go tool pprof gb.prof
//png
//quit

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	. "github.com/vipally/binary/internal/bench"
)

var cases = BenchCases()

func main() {
	f, err := os.Create("gb.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	n := flag.Int("n", -1, fmt.Sprintf("sel number, max %d", len(cases)-1))
	flag.Parse()

	doCase(*n, true, time.Now())
}

func doCase(n int, head bool, start time.Time) {
	doCnt := 20000000
	if head {
		fmt.Printf("\n===============\n")
		fmt.Printf("time = %s\n", start.Format("2006-01-02 15:04:05"))
		fmt.Printf("buildtime = %s\n", time.BuildTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("doCnt = %d\n", doCnt)
		fmt.Printf("%-10s", "StdEnCode")
		fmt.Printf("%-10s", "StdDecode")
		fmt.Printf("%-10s", "Encode")
		fmt.Printf("%-10s", "Decode")
		fmt.Printf("\n")
	}

	_doCnt := doCnt
	dur, speed, _ := DoBenchUvarint(BenchStdWrite, UvarintCases, _doCnt)
	fmt.Printf("%-10s", dur.String())
	dur, speed, _ = DoBenchUvarint(BenchStdRead, UvarintStdBytes, _doCnt)
	fmt.Printf("%-10s", dur.String())
	dur, speed, _ = DoBenchUvarint(BenchEncode, UvarintCases, _doCnt)
	fmt.Printf("%-10s", dur.String())
	dur, speed, _ = DoBenchUvarint(BenchDecode, UvarintBytes, _doCnt)
	fmt.Printf("%-10s", dur.String())

	dur, speed = dur, speed
}
