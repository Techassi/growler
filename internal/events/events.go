package events

import (
	"os"
	"fmt"

	m "github.com/Techassi/growler/internal/models"
)

var f *os.File
var err error

func init() {
	f, err = os.Create("events.txt")
	if err != nil {
		panic(err)
	}
}

func WorkerFinish(event m.Event) {
	_, e := f.WriteString(fmt.Sprintf("Worker %s finished\n", event.Worker.ID.String()))
	if e == nil {
		f.Sync()
	}
}
