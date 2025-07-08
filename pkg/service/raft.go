package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/raft"
	myraft "github.com/liuzheran/lockInRaft/pkg/raft"
	"github.com/liuzheran/lockInRaft/pkg/schema"
	"github.com/liuzheran/lockInRaft/pkg/setting"
)

type RaftService interface {
	BootStrap()
	GetRaft() *raft.Raft
	IsLeaderPtr() *int64
	AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) error
	RemoveServer(id raft.ServerID, prevIndex uint64, timeout time.Duration) error
	GetRaftClusterInfo() ([]schema.RaftPeer, error)
	GetRaftLeader() (schema.RaftPeer, error)
}

type raftService struct {
	IsLeader   int64
	Raft       *raft.Raft
	MyFsm      *myraft.MyEmptyFsm
	RaftConfig *setting.RaftConfig
}

func NewRaftService(raft *raft.Raft, myFsm *myraft.MyEmptyFsm, raftConfig *setting.RaftConfig) RaftService {
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

func (r *raftService) GetRaft() *raft.Raft {
	return r.Raft
}

func (r *raftService) RemoveServer(id raft.ServerID, prevIndex uint64, timeout time.Duration) error {
	future := r.Raft.RemoveServer(id, prevIndex, timeout)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

func (r *raftService) IsLeaderPtr() *int64 {
	// 返回指针，方便在协程中修改
	return &r.IsLeader
}

func (r *raftService) GetRaftClusterInfo() ([]schema.RaftPeer, error) {
	future := r.Raft.GetConfiguration()
	if err := future.Error(); err != nil {
		fmt.Sprintf("failed to get raftconfiguration: %s", err)
		return nil, err
	}

	var nodes []schema.RaftPeer
	for _, server := range future.Configuration().Servers {
		s := strings.Split(string(server.Address), ":")
		if len(s) >= 2 {
			nodes = append(nodes, schema.RaftPeer{Ip: s[0], Port: s[1]}) // 将ip和port添加到nodes切片中
		}
	}

	addr, _ := r.Raft.LeaderWithID()
	addrStr := string(addr)
	// TODO 这里可能发生Leader切换
	leader := strings.Split(addrStr, ":")
	for i, node := range nodes {
		if node.Ip == leader[0] && node.Port == leader[1] {
			nodes[i].Role = "leader" // 标记为领导者
		} else {
			nodes[i].Role = "follower" // 标记为追随者
		}
	}
	return nodes, nil
}

func (r *raftService) GetRaftLeader() (schema.RaftPeer, error) {
	addr, _ := r.Raft.LeaderWithID()
	addrStr := string(addr)
	leader := strings.Split(addrStr, ":")
	return schema.RaftPeer{Ip: leader[0], Port: leader[1], Role: "leader"}, nil
}
