package radios

import (
  "fmt"
)

// Packet used by radios to store packets and associated information
type Packet struct {
  Data []byte
  Freq int
  FreqErr int
  Rssi float64
}

func (p Packet) String() string {
  return fmt.Sprtinf("Bad Packet Recevied - Freq %d, RSSI: %3.1f, FreqErr: %d, Data: [% x]", p.Freq, p.Rssi, p.FreqErr, p.Data)
}
