package main

import (
	"flag"

	httpClient "github.com/hb-go/micro-plugins/client/istio_http"
	httpServer "github.com/hb-go/micro-plugins/server/istio_http"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/registry/noop"

	apiClient "github.com/hb-go/micro/istio/http/api/client"
	"github.com/hb-go/micro/istio/http/api/handler"
	example "github.com/hb-go/micro/istio/http/api/proto/example"
)

var (
	serverAddr string
	callAddr   string
	cmdHelp    bool
)

func init() {
	flag.StringVar(&serverAddr, "server_address", "0.0.0.0:9080", "server address.")
	flag.StringVar(&callAddr, "client_call_address", ":9080", "client call options address.")
	flag.BoolVar(&cmdHelp, "h", false, "help")
	flag.Parse()
}

func main() {
	if cmdHelp {
		flag.PrintDefaults()
		return
	}

	// TODO 多client需要统一端口，或者在client中hard code
	c := httpClient.NewClient(
		client.ContentType("application/json"),
		func(o *client.Options) {
			o.CallOptions.Address = callAddr
		},
	)
	s := httpServer.NewServer(
		server.Address(serverAddr),
	)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.sample"),
		micro.Version("latest"),
		micro.Registry(noop.NewRegistry()),
		micro.Client(c),
		micro.Server(s),

		// 兼容micro cmd parse
		micro.Flags(cli.StringFlag{
			Name:   "client_call_address",
			EnvVar: "MICRO_CLIENT_CALL_ADDRESS",
			Usage:  " Invalid!!!",
		}),
	)

	service.Options().Cmd.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Initialise service
	service.Init(
		// create wrap for the Example srv client
		micro.WrapHandler(apiClient.ExampleWrapper(service)),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}