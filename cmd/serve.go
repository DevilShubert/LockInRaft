package cmd

import (
	"fmt"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/http/rest"
	myraft "github.com/liuzheran/lockInRaft/pkg/raft"
	"github.com/liuzheran/lockInRaft/pkg/repository"
	"github.com/liuzheran/lockInRaft/pkg/service"
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
	dbConfig := setting.ProviderDBConfig()
	raftConfig := setting.ProviderRaftConfig()
	fmt.Println("从配置文件读取的结果：", dbConfig)
	fmt.Println("从配置文件读取的结果：", raftConfig)

	// 初始化配置
	repo := repository.NewLockRecordRepository()
	lockService := service.NewLockService(repo)

	// 初始化raft配置并启动raft节点
	raftConfig.RaftDir = raftConfig.RaftDir + "_" + raftConfig.RaftId
	raft, myFsm, err := myraft.NewRaft(raftConfig.RaftAddr, raftConfig.RaftId, raftConfig.RaftDir)
	if err != nil {
		fmt.Printf("Failed to initialize raft: %v", err)
		return
	}

	raftService := service.NewRaftService(raft, myFsm, raftConfig)
	if raftConfig.BootStrap {
		raftService.BootStrap()
	}

	// 初始化api配置
	lockApi := rest.NewLockApi(lockService, raftService)

	// 启动协程，监听来自Leader的变更，并使用原子修改isLeader的标志位
	go func() {
		for leader := range raft.LeaderCh() {
			if leader {
				atomic.StoreInt64(&(raftService.IsLeader), 1)
			} else {
				atomic.StoreInt64(&(raftService.IsLeader), 0)
			}
		}
	}()

	// TODO 启动 gin 服务
	ginEngine := gin.Default()

	// api/v1
	v1 := ginEngine.Group("/api/v1")
	v1.GET("/list", lockApi.List)
	v1.POST("/acquire", lockApi.Acquire)
	v1.POST("/release", lockApi.Release)
	v1.POST("/keepalive", lockApi.KeepAlive)

	// api/lock
	lockEndopint := ginEngine.Group("/api/lock")
	lockEndopint.POST("/addNode", lockApi.AddNode)
	lockEndopint.POST("/removeNode", lockApi.RemoveNode)
	lockEndopint.POST("/getClusterInfo", lockApi.GetClusterInfo)
	lockEndopint.POST("/getLeader", lockApi.GetLeader)

	ginEngine.Run()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
