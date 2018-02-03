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

  timeout := (2562500 + 200000) * time.Microsecond

  log.Printf("Waiting for packets...")

  ph := protocol.NewProtocolHandler(0)
  ph.NextHop()
  nextHop := ph.NextHop()
  log.Printf("Hopping to %v", nextHop)
  r.SetFreq(uint32(nextHop))

  for {
    packet, timeout, _ := r.ReceiveData(timeout)
    shouldHop := ph.HandlePacket(packet, timeout)

    if shouldHop {
      nextHop := ph.NextHop()
      log.Printf("Hopping to %v", nextHop)
      r.SetFreq(uint32(nextHop))
    }
  }

  os.Exit(0)
}
