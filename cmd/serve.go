package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuzheran/lockInRaft/pkg/http/rest"
	"github.com/liuzheran/lockInRaft/pkg/infra/db"
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
	ctx := context.Background()
	// 使用viper 加载配置
	dbConfig := setting.ProviderDBConfig()
	raftConfig := setting.ProviderRaftConfig()
	httpConfig := setting.ProviderHttpConfig()
	fmt.Println("dbConfig 从配置文件读取的结果：", dbConfig)
	fmt.Println("raftConfig 从配置文件读取的结果：", raftConfig)
	fmt.Println("httpConfig 从配置文件读取的结果：", httpConfig)

	db, err := db.ProvideLockDB(dbConfig)
	if err != nil {
		panic("初始化数据库时有问题")
	}

	// 初始化Repo
	lockRecordRepo := repository.NewLockRecordRepository()

	// 初始化raft配置并启动raft节点
	raftConfig.RaftDir = raftConfig.RaftDir + "_" + raftConfig.RaftId
	fmt.Println("raftConfig.RaftDir: ", raftConfig.RaftDir)
	os.MkdirAll(raftConfig.RaftDir, 0700)

	raft, myFsm, err := myraft.NewRaft(raftConfig.RaftAddr, raftConfig.RaftId, raftConfig.RaftDir)
	if err != nil {
		fmt.Printf("Failed to initialize raft: %v", err)
		return
	}

	raftManager := service.NewRaftManager(raft, myFsm, raftConfig)
	// 如果设置为Leader且是项目第一次启动才需要BootStrap
	if raftConfig.BootStrap {
		raftManager.BootStrap()
	}
	// 初始化CacheManager
	cacheManager := service.NewCacheManager(db, lockRecordRepo)
	// 初始化众多service
	lockService := service.NewLockService(cacheManager, raftManager)
	// 初始化api配置
	lockApi := rest.NewLockApi(lockService)

	// 启动协程，监听来自Leader的变更，并使用原子修改isLeader的标志位
	go func() {
		for leader := range raft.LeaderCh() {
			if leader {
				atomic.StoreInt64(lockService.RaftManager.IsLeaderPtr(), 1)
			} else {
				atomic.StoreInt64(lockService.RaftManager.IsLeaderPtr(), 0)
			}
		}
	}()

	// TODO 启动 gin 服务(改写使用HTTP.SERVER启动方式)
	ginEngine := gin.Default()

	// api/v1
	v1 := ginEngine.Group("/api/v1")
	v1.GET("/list", lockApi.List)
	v1.POST("/acquire", lockApi.Acquire)
	v1.POST("/release", lockApi.Release)
	v1.POST("/keepalive", lockApi.KeepAlive)
	v1.POST("/cache", lockApi.RebuildCache)

	// api/lock
	lockEndopint := ginEngine.Group("/api/lock")
	lockEndopint.POST("/addNode", lockApi.AddNode)
	lockEndopint.POST("/removeNode", lockApi.RemoveNode)
	lockEndopint.GET("/getClusterInfo", lockApi.GetClusterInfo)
	lockEndopint.GET("/getLeader", lockApi.GetLeader)

	s := &http.Server{
		// Addr需要写成 :port 格式
		Addr:           fmt.Sprintf(":%d", httpConfig.Port),
		Handler:        ginEngine,
		ReadTimeout:    time.Duration(httpConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(httpConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: httpConfig.MaxHeaderBytes,
	}
	s.ListenAndServe()

	go lockService.Elect(ctx)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
