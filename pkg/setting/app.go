package setting

// 标准数据库配置
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type RaftConfig struct {
	BootStrap   bool
	HttpAddr    string
	RaftAddr    string
	RaftId      string
	RaftCluster string
	RaftDir     string
}
