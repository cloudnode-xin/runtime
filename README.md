# runtime

## Usage

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudnodexin/runtime"
	"github.com/cloudnodexin/runtime/logger"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigChan)

	s := runtime.New()
	s.Use(logger.Setup(logger.LevelString("info")))

	if err := s.Start(); err != nil {
		panic(err)
	}

	<-sigChan

	if err := s.Stop(); err != nil {
		panic(err)
	}
}

```