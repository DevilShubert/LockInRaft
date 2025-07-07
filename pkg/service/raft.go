package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	myraft "github.com/liuzheran/lockInRaft/pkg/raft"
	"github.com/liuzheran/lockInRaft/pkg/setting"
)

type RaftService interface {
	BootStrap()
	AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) error
}

type raftService struct {
	IsLeader   int64
	Raft       *raft.Raft
	MyFsm      *myraft.MyEmptyFsm
	RaftConfig *setting.RaftConfig
}

func NewRaftService(raft *raft.Raft, myFsm *myraft.MyEmptyFsm, raftConfig *setting.RaftConfig) *raftService {
	return &raftService{
		IsLeader:   0,
		Raft:       raft,
		MyFsm:      myFsm,
		RaftConfig: raftConfig,
	}
}

func (r *raftService) BootStrap() {
	// 获取集群配置（从stable中查找）
	servers := r.Raft.GetConfiguration().Configuration().Servers

	// 如果配置存在则跳过（因为正常情况下会自动恢复）
	if len(servers) > 0 {
		fmt.Println("集群已经存在，无需再次启动")
		return
	}

	// 第一次启动的情况
	fmt.Println("集群不存在，初始化集群配置并，进行启动")
	peerArray := strings.Split(r.RaftConfig.RaftCluster, ",")
	if len(peerArray) == 0 {
		return
	}

	var configuration raft.Configuration
	for _, peerInfo := range peerArray {
		peer := strings.Split(peerInfo, "/")
		id := peer[0]
		addr := peer[1]
		server := raft.Server{
			ID:      raft.ServerID(id),
			Address: raft.ServerAddress(addr),
		}
		configuration.Servers = append(configuration.Servers, server)
	}
	// 按照配置进行启动
	r.Raft.BootstrapCluster(configuration)
	fmt.Println("集群启动成功")
}

func (r *raftService) AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) error {
	future := r.Raft.AddVoter(id, address, prevIndex, timeout)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}
