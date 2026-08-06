package config

type System struct {
	Addr       string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Path       string `mapstructure:"path" json:"path" yaml:"path"`
	DbType     string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	Tabpre     string `mapstructure:"tabpre" json:"tabpre" yaml:"tabpre"`
	Host       string `mapstructure:"host" json:"host" yaml:"host"`
	Webroute   string `mapstructure:"webroute" json:"webroute" yaml:"webroute"`
	Webname    string `mapstructure:"webname" json:"webname" yaml:"webname"`
	Webversion string `mapstructure:"webversion" json:"webversion" yaml:"webversion"`
	Cookiename string `mapstructure:"cookiename" json:"cookiename" yaml:"cookiename"`
	Appid      string `mapstructure:"appid" json:"appid" yaml:"appid"`
	Zrlogin    bool   `mapstructure:"zrlogin" json:"zrlogin" yaml:"zrlogin"`
}
