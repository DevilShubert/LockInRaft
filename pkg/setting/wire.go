package setting

import "github.com/spf13/viper"

func ProviderDBConfig() *DBConfig {
	// 使用viper来读取环境变量值
	dbconfig := &DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DbName:   viper.GetString("db.dbname"),
	}
	return dbconfig
}

func ProviderRaftConfig() *RaftConfig {
	raftConfig := &RaftConfig{
		BootStrap:   viper.GetBool("raft.bootStrap"),
		HttpAddr:    viper.GetString("raft.httpAddr"),
		RaftAddr:    viper.GetString("raft.raftAddr"),
		RaftId:      viper.GetString("raft.raftId"),
		RaftCluster: viper.GetString("raft.raftCluster"),
		RaftDir:     viper.GetString("raft.raftDir"),
	}
	return raftConfig
}

func ProviderHttpConfig() *HttpConfig {
	httpConfig := &HttpConfig{
		Port:           viper.GetInt("http.port"),
		ReadTimeout:    viper.GetInt("http.readTimeout"),
		WriteTimeout:   viper.GetInt("http.writeTimeout"),
		MaxHeaderBytes: viper.GetInt("http.maxHeaderBytes"),
	}
	return httpConfig
}
