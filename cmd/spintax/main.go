package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/nikandfor/spintax"
)

var (
	cnt  = flag.Bool("c", false, "print count of distinct possible texts instead of spinning")
	file = flag.Bool("f", false, "interpret argument as a file name with spintax not as spintax itself")
	iter = flag.Bool("a", false, "print all the possible texts instead of one random")
	help = flag.Bool("h", false, "print help and exit")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 || *help {
		fmt.Printf("usage: %v [-c] [-f] <spintax> // generates one random text out of given spintax template\n", os.Args[0])
		if *help {
			flag.PrintDefaults()
		}
		return
	}
	arg := flag.Arg(0)

	var templ string
	if *file {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			panic(err)
		}
		templ = string(data)
	} else {
		templ = arg
	}

	e := spintax.Parse(templ)

	if *cnt {
		fmt.Printf("%v\n", e.Count())
		return
	}
	if *iter {
		c := e.Iter()
		for l := range c {
			fmt.Printf("%v\n", l)
		}
		return
	}

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("%v\n", e.Spin())
}
