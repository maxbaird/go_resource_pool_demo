package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"go_resource_pool_demo/model"
	"go_resource_pool_demo/pool"
)

const maxAllocations = 10 ^ 6

var (
	r          = rand.New(rand.NewSource(time.Now().Unix()))
	p, _       = pool.New(1000)
	withPool   bool
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	wg         sync.WaitGroup
)

func main() {
	flag.BoolVar(&withPool, "wp", false, "use resource pool")
	flag.Parse()

	// ---- START CPU profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	// ---- END CPU profiing

	if withPool {
		testWithPool()
	} else {
		testWithoutPool()
	}
	wg.Wait()
	if withPool {
		fmt.Println(p)
	}

	// ---- START MEM profiing
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
	// ---- END MEM profiing
}

func testWithoutPool() {
	fmt.Println("Testing without resoruce pool")
	for i := 0; i < maxAllocations; i++ {
		gp := &model.GamePiece{}

		wg.Add(1)
		go func() {
			gp.Field3[0] = 0 // just do something
			// do something for a little bit to hold on to the object
			time.Sleep(time.Duration(r.Int()%100) * time.Millisecond)
			wg.Done()
		}()
	}
}

func testWithPool() {
	fmt.Println("Testing with resource pool")
	for i := 0; i < maxAllocations; i++ {
		gp := p.Alloc()

		wg.Add(1)
		go func() {
			gp.Field3[0] = 0 // just do something
			// do something for a little bit to hold on to the object
			time.Sleep(time.Duration(r.Int()%100) * time.Millisecond)
			// release object
			p.Release(gp)
			wg.Done()
		}()
	}
}
