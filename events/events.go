package events

import (
	"github.com/Techassi/growler/internal/workerpool"
)

func WorkerInit(pool *workerpool.WorkerPool) {
	fmt.Println(len(pool.Queue.Items))
}
