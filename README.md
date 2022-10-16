# logrus2logtail
Logtail hooks for Logrus 

## Installation
```
go get github.com/icoom-lab/logrus2logtail
```

## Usage 

```go
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/icoom-lab/logrus2logtail"
)

func main() {	

	hook, err := logrus2logtail.NewHook("enter_token_here")
	if err != nil {
		log.Error("Error:", err)
	}
	log.AddHook(hook)

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.Debug("This is a debug message")
}
```

## License 

**logrus2logtail** is released under [the MIT license](LICENSE).
