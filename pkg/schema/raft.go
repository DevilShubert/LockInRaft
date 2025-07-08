package schema

type NodeParam struct {
	// 结构体标签：当从 JSON 解码时，也会自动把 id 字段的值填充到 Id，peer_addr 填充到 PeerAddr，反之填充时也会这样
	Id       string `json:"id"`
	PeerAddr string `json:"peer_addr"`
}

type RaftPeer struct {
	Ip   string
	Port string
	Role string
}

type GetClusterInfoParam struct {
}
