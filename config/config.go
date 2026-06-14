package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port      string `yaml:"port"`
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"server"`
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	RabbitMQ struct {
		URL string `yaml:"url"`
	} `yaml:"rabbitmq"`
	PayPal struct {
		ClientID  string `yaml:"client_id"`
		Secret    string `yaml:"secret"`
		ReturnURL string `yaml:"return_url"`
		CancelURL string `yaml:"cancel_url"`
	} `yaml:"paypal"`
	Alipay struct {
		AppID        string `yaml:"app_id"`
		PrivateKey   string `yaml:"private_key"`
		AliPublicKey string `yaml:"alipay_public_key"`
		NotifyURL    string `yaml:"notify_url"`
		ReturnURL    string `yaml:"return_url"`
	} `yaml:"alipay"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
