package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/http/api"
	"github.com/liuzheran/lockInRaft/pkg/setting"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	// 使用viper 加载配置
	db := setting.ProviderDBConfig()
	fmt.Println("从配置文件读取的结果：", db)

	// 一系列初始化配置（且顺序是从底到高？）
	// repo := repository.NewLockRecordRepository(db)
	// service := service.NewLockService(repo)

	// TODO 启动 gin 服务
	ginEngine := gin.Default()
	ginEngine.GET("/api/v1/list", func(c *gin.Context) {
		api.List(c)
	})
	// ginEngine.POST('/api/v1/lock')
	// ginEngine.POST('/api/v1/release')
	// ginEngine.GET('/api/v1/renew')
	ginEngine.Run()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
