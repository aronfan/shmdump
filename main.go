package main

import (
	"flag"
	"fmt"
)

var (
	e = flag.String("e", "", "pipe mode, such as \"shmkey=1&op=save\"")
	h = flag.Bool("h", false, "display help")
)

func main() {
	flag.Parse()
	if *h {
		fmt.Println("share memory dumper")
		return
	}
	if *e != "" {
		if err := pipe(*e); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	} else {
		// display the interactive ui
	}
}
