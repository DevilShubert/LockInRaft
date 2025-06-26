package setting

// 标准数据库配置
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}
