package raft

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type MyEmptyFsm struct{}

func (m *MyEmptyFsm) Apply(l *raft.Log) interface{} {
	return nil
}

func (m *MyEmptyFsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (m *MyEmptyFsm) Restore(rc io.ReadCloser) error {
	return nil
}

func NewRaft(raftAddr string, raftId string, raftDir string) (*raft.Raft, *MyEmptyFsm, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftId)

	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, nil, err
	}

	// Transport是raft集群内部节点之间的信息通道，节点之间需要通过该通道来同步log、选举leader等
	transport, err := raft.NewTCPTransport(raftAddr, addr, 2, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, nil, err
	}

	// 用于raft节点恢复的快照，可以使用默认的 NewFileSnapshotStore
	snapshots, err := raft.NewFileSnapshotStore(raftDir, 2, os.Stderr)
	if err != nil {
		return nil, nil, err
	}
	// 存储raft集群传递的log
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, nil, err
	}
	// stable用于存储保存Raft选举信息，如角色、term等信息
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.db"))
	if err != nil {
		return nil, nil, err
	}

	myFsm := &MyEmptyFsm{}

	raft, err := raft.NewRaft(config, myFsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, nil, err
	}
	return raft, myFsm, nil
}
