package config

import (
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func NewMySQLConfig(v *viper.Viper) mysql.Config {
	v.SetDefault("MYSQL_USER", "root")
	v.SetDefault("MYSQL_PASSWORD", "")
	v.SetDefault("MYSQL_NET", "tcp")
	v.SetDefault("MYSQL_ADDR", "127.0.0.1:3306")
	v.SetDefault("MYSQL_DATABASE", "nice_dcv_db")

	return mysql.Config{
		User:                 v.GetString("MYSQL_USER"),
		Passwd:               v.GetString("MYSQL_PASSWORD"),
		Net:                  v.GetString("MYSQL_NET"),
		Addr:                 v.GetString("MYSQL_ADDR"),
		DBName:               v.GetString("MYSQL_DATABASE"),
		AllowNativePasswords: true,
	}
}
