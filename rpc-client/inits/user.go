package init_client

import (
	"ticketing-infra/rpc-client/kitex-gens/user/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

var NewUserClient userservice.Client

func InitUserSvClient() error {
	var err error
	NewUserClient, err = userservice.NewClient("userservice",
		client.WithHostPorts("0.0.0.0:8889"),
		client.WithTransportProtocol(transport.TTHeader),        // ← 明确 TTHeader
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler)) // ← 必须)
	if err != nil {
		return err
	}
	return nil
}
