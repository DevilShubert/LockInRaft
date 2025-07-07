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
		BootStrap:   viper.GetBool("raft.bootstrap"),
		HttpAddr:    viper.GetString("raft.http_addr"),
		RaftAddr:    viper.GetString("raft.raft_addr"),
		RaftId:      viper.GetString("raft.raft_id"),
		RaftCluster: viper.GetString("raft.raft_cluster"),
		RaftDir:     viper.GetString("raft.raft_dir"),
	}
	return raftConfig
}
