package events

import (
	"fmt"

	"github.com/Techassi/growler/internal/workerpool"
)

func WorkerProcess(pool *workerpool.WorkerPool) {
	fmt.Println(len(pool.Queue.Items))
}
