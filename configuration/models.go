package configuration

import "time"

type Server struct {
	Host    string `yaml:"host"`
	Port    int16 `yaml:"port"`
	Gateway string `yaml:"gateway"`
	ServerListFilePath	string	`yaml:"server_list_file_path"`
	ServerListFilePathForStaging	string	`yaml:"server_list_file_path_for_staging"`
}

type PartnerApi struct{
	Uri	string	`yaml:"uri"`
	SecretKey	string `yaml:"secret_key"`
	Timeout	time.Duration	`yaml:"timeout"`
	Concurrent	int	`yaml:"concurrent"`
}

type Conf struct {
	ConfigureFile string
	Server	Server	`yaml:"server"`
	PartnerApi	PartnerApi	`yaml:"partner_api"`
}
