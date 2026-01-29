package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Mapping โครงสร้างให้ตรงกับ yaml
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // เพิ่ม
	TimeZone string // เพิ่ม
}

func LoadConfig() (*Config, error) {
	// บอก Viper ว่าไฟล์ชื่ออะไร อยู่ที่ไหน
	viper.SetConfigName("config") // ชื่อไฟล์ config (ไม่ต้องใส่นามสกุล .yaml)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // ให้หาที่ root folder

	// เผื่ออยาก override ด้วย Environment Variable (เช่น SERVER_PORT=8080)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// อ่านไฟล์
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// แปลงข้อมูลใส่ Struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
