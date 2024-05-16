package transport_client

import (
	"fmt"
	"jdb/jrpc/rpc_common/entities"
	"jdb/jrpc/rpc_core/load_balancer"
	"jdb/jrpc/rpc_core/serializer"
	"jdb/jrpc/rpc_core/services/discovery"
	"jdb/jrpc/rpc_core/transport/utils"
	"net"
)

type SocketClient struct {
	ServiceDiscovery services_discovery.ServiceDiscovery
	Serializer       serializer.CommonSerializer
	AimIp            string
}

func NewDefaultSocketClient() *SocketClient {
	return NewSocketClient(DEFAULT_SERIALIZER, &load_balancer.RandomLoadBalancer{}, "")
}

func NewDefaultSocketClientWithAimIp(aimIp string) *SocketClient {
	return NewSocketClient(DEFAULT_SERIALIZER, &load_balancer.RandomLoadBalancer{}, aimIp)
}

func NewSocketClient(serializerId int, loadBalancer load_balancer.LoadBalancer, aimIp string) *SocketClient {
	return &SocketClient{
		Serializer:       serializer.GetByCode(serializerId),
		ServiceDiscovery: services_discovery.NewZkServiceDiscovery(loadBalancer),
		AimIp:            aimIp,
	}
}

func (client *SocketClient) SendRequest(req entities.RPCdata) (*entities.RPCdata, error) {
	var addr *net.TCPAddr
	var err error
	if client.AimIp == "" {
		fmt.Println("走的是ServiceDiscovery")
		addr, err = client.ServiceDiscovery.LookupService(req.Name)
		fmt.Println("结果为addr", addr)
	} else {
		fmt.Println("请求的client.AimIp为", client.AimIp)
		addr, err = net.ResolveTCPAddr("tcp", client.AimIp)
		fmt.Println("走的是ResolveTCPAddr")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lookup service: %w", err)
	}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}
	err = transport_utils.NewObjectWriter(conn).WriteObject(&req, client.Serializer)
	if err != nil {
		return nil, err
	}
	resp, err := transport_utils.NewObjectReader(conn).ReadObject()
	if err != nil {
		return nil, err
	}
	return resp, nil
}