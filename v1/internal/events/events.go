package events

import (
	"fmt"

	m "github.com/Techassi/growler/internal/models"
)

func WorkerFinish(event m.Event) {
	fmt.Printf("Worker %s finished\n", event.Worker.ID.String())
}
