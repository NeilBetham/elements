package main

import (
  "os"
  "time"
  "log"
  "flag"
  "periph.io/x/periph/host"

  "github.com/NeilBetham/elements/radios"
  "github.com/NeilBetham/elements/protocol"
  "github.com/NeilBetham/elements/reporting"
  "github.com/NeilBetham/elements/config"
)

func main() {
  if _, err := host.Init(); err != nil {
    os.Exit(1)
  }

  configPathPtr := flag.String("config", "elements_config.yml", "The config yaml to use")
  flag.Parse()

  if *configPathPtr == "bad_yaml" {
    log.Fatal("Please specify a config file")
  }
  config, err := config.ReadConfig(*configPathPtr)
  if err != nil{
    log.Fatal("Error reading config: %s", err)
  }

  reporter := reporting.NewReporter(config)

  r, err := radios.NewRFM69("/dev/spidev0.0", "GPIO4", "GPIO5")
  if err != nil{
    log.Fatal("Failed to open spi port ro radio: %s", err)
  }

  timeout := (2562500 + 200000) * time.Microsecond

  log.Printf("Waiting for packets...")

  ph := protocol.NewProtocolHandler(0)
  ph.NextHop()
  nextHop := ph.NextHop()
  log.Printf("Hopping to %v", nextHop)
  r.SetFreq(uint32(nextHop.Freq))

  for {
    packet, timeout, _ := r.ReceiveData(timeout)
    shouldHop, reading := ph.HandlePacket(packet, timeout)

    if reading.Valid {
      log.Printf("Reading: %s", reading)
      err := reporter.ReportReading(reading)
      if err != nil {
        log.Printf("Error reporting reading: %s", err)
      }
    }

    if shouldHop {
      nextHop := ph.NextHop()
      log.Printf("Hopping to %v", nextHop)
      r.SetFreq(uint32(nextHop.Freq))
    }
  }

  os.Exit(0)
}
