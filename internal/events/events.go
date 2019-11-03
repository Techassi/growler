package events

import (
	"os"
	"fmt"

	"github.com/Techassi/growler/internal/workerpool"
)

var f *File
var err error

func init() {
	f, err = os.Create("links.txt")
	if err != nil {
		panic(err)
	}
}

func WorkerFinish(pool workerpool.Event) {
	_, e := f.WriteString(fmt.Sprintf("Hello\n"))
	f.Sync()
}
