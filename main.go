package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aronfan/xerrors"
)

var (
	e = flag.String("e", "", "pipe mode, such as \"shmkey=1&op=save\"")
	h = flag.Bool("h", false, "display help")
)

func main() {
	xerrors.SetSysInternalError(-1)

	flag.Parse()
	if *h {
		fmt.Println("share memory dumper")
		return
	}
	if *e != "" {
		code := int(0)
		if err := pipe(*e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			code = int(xerrors.Int(err))
		}
		os.Exit(code)
	} else {
		// display the interactive ui
	}
}
