// usage:
// pprof.exe >> pprof.log
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
		fmt.Printf("%-30s", "Case")
		fmt.Printf("%-10s", "GobEncode")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "GobDecode")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "GobSize")
		fmt.Printf("%-10s", "StdWrite")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "StdRead")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "StdSize")

		fmt.Printf("%-10s", "EncodeY")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "DecodeY")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "BinYSize")
		fmt.Printf("%-10s", "EncodeN")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "DecodeN")
		fmt.Printf("%-12s", "Speed")
		fmt.Printf("%-10s", "BinNSize")
		fmt.Printf("%-10s", "Cost")
		fmt.Printf("%-10s", "TotalTime")
		fmt.Println("")
		//fmt.Printf("%-30s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s%-10s\n", "Case", "GobEncode", "GobDecode", "StdWrite", "StdRead", "EncodeY", "DecodeY", "EncodeN", "DecodeN", "Cost", "TotalTime")
	}

	if n < 0 {
		for i := 0; i < len(cases); i++ {
			doCase(i, false, start)
		}
		fmt.Printf("finish time = %s\nCost = %s\n", time.Now().Format("2006-01-02 15:04:05"), Duration(time.Now().Sub(start)).String())
	} else {
		if n >= len(cases) {
			println("max case", len(cases)-1)
			return
		}
		v := cases[n]
		st := time.Now()
		fmt.Printf("%-30s", v.Name)
		_doCnt := doCnt / v.Length
		dur, speed, size := DoBench(BenchGobEncode, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		dur, speed, size = DoBench(BenchGobDecode, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		fmt.Printf("%-10s", size.String())

		dur, speed, size = DoBench(BenchStdWrite, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		dur, speed, size = DoBench(BenchStdRead, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		fmt.Printf("%-10s", size.String())
		dur, speed, size = DoBench(BenchEncode, v.Data, _doCnt, true)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		dur, speed, size = DoBench(BenchDecode, v.Data, _doCnt, true)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		fmt.Printf("%-10s", size.String())
		dur, speed, size = DoBench(BenchEncode, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		dur, speed, size = DoBench(BenchDecode, v.Data, _doCnt, false)
		fmt.Printf("%-10s", dur.String())
		fmt.Printf("%-12s", speed.String())
		fmt.Printf("%-10s", size.String())
		durCost := Duration(time.Now().Sub(st))
		fmt.Printf("%-10s", durCost.String())
		durAll := Duration(time.Now().Sub(start))
		dur, speed, size = dur, speed, size
		fmt.Printf("%-10s", durAll.String())
		fmt.Println("")
	}
}
