package setting

import "github.com/spf13/viper"

func ProviderDBConfig() *DBConfig {
	// TODO 使用viper来读取环境变量值
	dbconfig := &DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DbName:   viper.GetString("db.dbname"),
	}
	return dbconfig
}
