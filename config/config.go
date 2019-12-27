package config
import(
  "os"
  "gopkg.in/yaml.v3"
)

type Config struct {
  Server struct {
    Host string `yaml:"host"`
    Port string `yaml:"port"`
    Ssl bool `yaml:"ssl"`
    SslVerify bool `yaml:"ssl_verify"`
    StationId string `yaml:"station_id"`
  } `yaml:"server"`
  Credentials struct {
    ApiKey string `yaml:"api_key"`
  } `yaml:"credentials"`
}


func ReadConfig(path string) (Config, error) {
  var cfg Config

  f, err := os.Open(path)
  if err != nil {
    return cfg, err
  }
  defer f.Close()
  decoder := yaml.NewDecoder(f)
  err = decoder.Decode(&cfg)
  return cfg, err
}
