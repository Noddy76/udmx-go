# udmx-go

A simple client library for the [uDMX](https://www.anyma.ch/research/udmx/) USB to DMX interface. Based on the C uDMX code on [GitHub](https://github.com/mirdej/udmx).

## Example Usage

```go
package main

import (
	"log"

	udmx "github.com/noddy76/udmx-go"
)

func main() {
	dmx, err := udmx.NewUdmx()
	if err != nil {
		log.Panicf("Error opening DMX : %v", err)
	}

	dmx.SetChannelRange(0, []byte{255, 255, 0, 0})

	err = dmx.Close()
	if err != nil {
		log.Panicf("Error closing DMX : %v", err)
	}
}
```
