package radios

// Packet used by radios to store packets and associated information
type Packet struct {
  Data []byte
  Freq int
  FreqErr int
  Rssi float64
}
