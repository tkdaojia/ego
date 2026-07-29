package config

type Grpc struct {
	Host  string `mapstructure:"host" json:"host" yaml:"host"`
	Addr  int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	Token string `mapstructure:"token" json:"token" yaml:"token"`
	Open  bool   `mapstructure:"open" json:"open" yaml:"open"`
}
