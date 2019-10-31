# growler

Another crawler for the web. But this one is different! It is written in Go(lang) and by me...

## Contents

- [Overview](#overview)
- [Usage](#usage)
	- [Basic](#basic)
	- [Advanced](#advanced)
- [API](#api)
- [Current Development](#current-development)

## Overview

More coming soon.

## Usage

### Basic

Windows:
```shell
.\growler --url https://example.com
```

Linux:
```shell
./growler --url https://example.com
```

### Advanced

Possible arguments

Argument | Explanation
--- | ---
`--url`| *Required.* The url used as the entry point (seed) to start crawling

More coming soon.

## API

### Queue

```golang
q, err := queue.NewQueue(1000)
if err != nil {
	panic(err)
}
```

To initialize the queue call `NewQueue(max int)` with `max` being the max count of items which can be in queue at one time. `NewQueue` will return an error if `max` is smaller than 0.

```shell
panic: -1 is smaller than 0. Provide a value greater than 0
```

### WorkerPool
#### Initialize
```golang
p := workerpool.NewWorkerPool(10, q, action)
```

To initialize the worker pool call `NewWorkerPool(workers int, queue queue.Queue, func(interface{}) interface{})`. The argument `workers` sets the count of total workers being used to crawl the urls. The argument `queue` sets the queue which should be used to poll jobs from. The last arguments provides a function which will be executed by the workers, so this is your actual "work".

`NewWorkerPool()` returnes the created pool.

For an example work function check out [this](https://github.com/Techassi/growler/blob/master/internal/crawl/crawl.go) file.

#### Start

```golang
p.Start()
```

This starts the worker pool and the workers get created and start polling jobs from the queue.

### Lifecycle Events
```golang
err := p.On("init", workerInit)
if err != nil {
	panic(err)
}
```

You can attach event function (handlers) on lifecycle events emitted by workers through the worker pool. Call `On(event string, action func(pool *workerpool.WorkerPool))` to subscribe to an event called `event`. Every time this event occurres `action` will get called.

Possible events
- init: Triggred when the worker is being created / initialized.

## Current Development

- Lifecycle Events
