package protocol

import (
  "log"
  "time"
  "github.com/NeilBetham/elements/crc"
  "github.com/NeilBetham/elements/radios"
)

// ProtocolHandler handles receipt of packets and channel hopping
type ProtocolHandler struct {
  crc.CRC

  hopTime time.Duration
  hopIndex int
  hopPattern []int
  channels []int

  goodPkts int
  badPkts int
  resync bool
}

// NewProtocolHandler sets up a new protocol handler
func NewProtocolHandler(stationNumber int) (ph ProtocolHandler){
  ph.CRC = crc.NewCRC("CCITT-16", 0, 0x1021, 0)

  ph.hopTime = time.Duration(2562500 + (stationNumber * 62500)) * time.Microsecond

  ph.hopIndex = 0

  ph.channels = []int{
    902355835, 902857585, 903359336, 903861086, 904362837, 904864587,
    905366338, 905868088, 906369839, 906871589, 907373340, 907875090,
    908376841, 908878591, 909380342, 909882092, 910383843, 910885593,
    911387344, 911889094, 912390845, 912892595, 913394346, 913896096,
    914397847, 914899597, 915401347, 915903098, 916404848, 916906599,
    917408349, 917910100, 918411850, 918913601, 919415351, 919917102,
    920418852, 920920603, 921422353, 921924104, 922425854, 922927605,
    923429355, 923931106, 924432856, 924934607, 925436357, 925938108,
    926439858, 926941609, 927443359,
  }

  ph.hopPattern = []int{
    50, 18, 40, 24, 7, 46, 31, 12, 35, 21, 2, 28, 43, 15, 4, 26, 37, 9,
    48, 20, 1, 29, 41, 13, 47, 6, 23, 33, 44, 0, 16, 38, 25, 8, 30, 49,
    36, 11, 19, 32, 3, 42, 27, 14, 34, 5, 39, 10, 22, 45, 17,
  }

  ph.goodPkts = 0
  ph.badPkts = 0
  ph.resync = true
  return
}

// HandlePacket handles incomming packets and decides if a hop should happen
func (ph *ProtocolHandler) HandlePacket(pkt radios.Packet, timeout bool) (hop bool){
  // Davis ISS transmits LSB first
  for index, data := range pkt.Data {
    pkt.Data[index] = swapBitOrder(data)
  }

  // If the checksum is valid then hop
  if ph.Checksum(pkt.Data) != 0 {
    log.Printf("Bad Packet Recevied - Freq %d, RSSI: %3.1f, FreqErr: %d, Data: [% x]", pkt.Freq, pkt.Rssi, pkt.FreqErr, pkt.Data)
    if !ph.resync {
      ph.badPkts++
    }
  } else if timeout {
    if !ph.resync {
      ph.badPkts++
    }
  } else {
    if ph.resync {
      ph.resync = false
    }

    log.Printf("Packet Recevied - Freq %d, RSSI: %3.1f, FreqErr: %d, Data: [% x]", pkt.Freq, pkt.Rssi, pkt.FreqErr, pkt.Data)

    ph.badPkts = 0
    ph.goodPkts++
  }

  // If we start to accumulate bad packets it's time for a resync
  if ph.badPkts >= 5 {
    log.Printf("Resync needed")
    ph.resync = true
    ph.badPkts = 0
  }

  if ph.resync {
    hop = false
  } else {
    hop = true
  }
  return
}

// NextHop gets the next hopping frequency and increments the hop index
func (ph *ProtocolHandler) NextHop() (freq int){
  freq = ph.channels[ph.hopPattern[ph.hopIndex]]
  ph.hopIndex++
  if ph.hopIndex > 50 {
    ph.hopIndex = 0
  }
  return
}

// CurrentChannel returns the current channel
func (ph *ProtocolHandler) CurrentChannel() (freq int){
  hopIndex := ph.hopIndex
  if hopIndex == 0 {
    hopIndex = 50
  } else {
    hopIndex--
  }

  return ph.channels[ph.hopPattern[hopIndex]]
}

func swapBitOrder(b byte) byte {
  b = ((b & 0xF0) >> 4) | ((b & 0x0F) << 4)
  b = ((b & 0xCC) >> 2) | ((b & 0x33) << 2)
  b = ((b & 0xAA) >> 1) | ((b & 0x55) << 1)
  return b
}
