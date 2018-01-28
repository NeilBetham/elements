package main

import (
  "os"
  "time"
  "log"
  "periph.io/x/periph/host"

  "github.com/NeilBetham/elements/radios"
  "github.com/NeilBetham/elements/protocol"
)

func main() {
  if _, err := host.Init(); err != nil {
		os.Exit(1)
	}

  r, _ := radios.NewRFM69("/dev/spidev0.0", "GPIO4", "GPIO5")

  timeout := 2562500 * time.Microsecond

  log.Printf("Waiting for packets...")

  ph := protocol.NewProtocolHandler(0)
  r.SetFreq(uint32(ph.NextHop()))

  for {
    data, rssi, _ := r.ReceiveData(timeout)
    log.Printf("Recevied packet - RSSI: %3.1f,  Data:[% x]", rssi, data)
    shouldHop := ph.HandlePacket(data)
    if shouldHop {
      r.SetFreq(uint32(ph.NextHop()))
    }
  }

  os.Exit(0)
}
