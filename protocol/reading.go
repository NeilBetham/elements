package protocol

import (
  "log"
  "time"
  "github.com/NeilBetham/elements/radios"
)

type Reading struct {
  StationID int
  Sensor Sensor
  SensorName string
  Value float64

  WindSpeed float64
  WindDir float64
}

func (r Reading) String() string {
  return fmt.Sprintf(
    "Reading for %s, station: %d, value %f",
    r.SensorName,
    r.StationID,
    r.Value
  )
}

func ParsePacket(pkt radios.Packet) (rd Reading){
  rd.StationID = (pkt.Data[0] & 0x0f) + 1
  rd.Sensor = (pkt.Data[0] & 0xf0) >> 4
  rd.SensorName = fmt.Printf("%s", rd.Sensor)

  switch rd.Sensor {
  case SuperCapVoltage:
    rd.Value = convertSuperCapVoltage(pkt.Data[3:6])
  case UVIndex:
    rd.Value = 0
  case RainRate:
    rd.Value = convertRainRate(pkt.Data[3:6])
  case SolarRadiation:
    rd.Value = 0
  case Light:
    rd.Value = convertLight(pkt.Data[3:6])
  case Temperature:
    rd.Value = convertTemperature(pkt.Data[3:6])
  case WindGustSpeed:
    rd.Value = convertGustSpeed(pkt.Data[3:6])
  case Humidity:
    rd.Value = convertHumidity(pkt.Data[3:6])
  case RainClicks:
    rd.Value = convertRainClicks(pkt.Data[3:6])
  default:
    return fmt.Sprintf("Uknown Reading Type: %0x", rType)
  }

  rd.WindSpeed = float64(pkt.Data[1])
  rd.WindDir = float64(9 + ((float64(pkt.Data[2]) * 342) / 255))
}

func convertSuperCapVoltage(data []byte) float64 {
  return float64((uint(data[0]) * 4 + ((uint(data[1]) & 0xC0) / 64)) / 100)
}

func convertRainRate(data []byte) float64 {
  if data[0] == 0xff{
    return 0
  } else if data[1] & 0x40 == 0 {
    return float64(((uint(data[1]) & 0x30) / 16 * 250) + data[0])
  } else if data[1] & 0x40 == 0x40 {
    return float64(((uint(data[1]) & 0x30) / 16 * 250) + data[0]) / 16
  }
  log.Printf("Unkown Rain Rate State")
  return 0
}

func convertLight(data []byte) float64 {
  return float64(uint(data[0]) * 4) + ((uint(byte[1]) & 0xc0) / 64)
}

func convertTemperature(data []byte) float64 {
  return float64(uint(data[0]) * 256 + uint(data[1])) / 160
}

func convertGustSpeed(data []byte) float64 {
  return float64(data[0])
}

func convertHumidity(data []byte) float64 {
  return float64(((uint(data[1]) & 0xf0) * 16) + uint(data[0])) / 10
}

func convertRainClicks(data []byte) float64 {
  return float64(uint(data[0]) & 0x7f)
}

type Sensor byte

const (
  SuperCapVoltage Reading = 2
  UVIndex         Reading = 4
  RainRate        Reading = 5
  SolarRadiation  Reading = 6
  Light           Reading = 7
  Temperature     Reading = 8
  WindGustSpeed   Reading = 9
  Humidity        Reading = 0xA
  RainClicks      Reading = 0xE
)

func (r Sensor) String(rType byte) string {
  switch rType {
  case SuperCapVoltage:
    return "SuperCapVoltage"
  case UVIndex:
    return "UVIndex"
  case RainRate:
    return "RainRate"
  case SolarRadiation:
    return "SolarRadiation"
  case Light:
    return "Light"
  case Temperature:
    return "Temperature"
  case WindGustSpeed:
    return "WindGustSpeed"
  case Humidity:
    return "Humidity"
  case RainClicks:
    return "RainClicks"
  default:
    return fmt.Sprintf("Uknown Reading Type: %0x", rType)
  }
}
