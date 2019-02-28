package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rs/cors"
	"github.com/spf13/viper"

	_service "github.com/stevepartridge/service"

	pb "github.com/stevepartridge/service/example/protos"
	"github.com/stevepartridge/service/swagger"
)

var (
	serviceName = "example"
	version     = "0.0.0"
	builtAt     = ""
	build       = "0"
	githash     = ""

	defaultHost = "example.local"
	defaultPort = 8000
	port        int

	enableInsecure = false

	service     *_service.Service
	grpcService *exampleService

	swag *swagger.Swagger
)

func main() {

	prepare()

	if builtAt == "" {
		builtAt = time.Now().Format(time.RFC3339Nano)
	}

	var err error

	port = viper.GetInt("port")
	if port == 0 {
		port = defaultPort
	}

	service, err = _service.New(port)
	ifError(err)

	if viper.GetBool("ENABLE_CORS") {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			Debug:          true,
		})
		service.AddHttpMiddlware(c.Handler)
	}

	serve()

}

func serve() {

	var err error

	enableInsecure = viper.GetBool("ENABLE_INSECURE")

	if enableInsecure {
		log.Println("Insecure is enabled.  Not recommended when in production.")
		service.EnableInsecure()
	}

	rootCA := getRootCA()
	if rootCA != nil {
		err = service.AppendCertsFromPEM(rootCA)
		ifError(err)
	}

	cert, err := getCert()
	ifError(err)
	key, err := getKey()
	ifError(err)

	err = service.AddKeyPair(cert, key)
	ifError(err)

	// service.Grpc.AddUnaryInterceptors(
	// RequestInterceptor(),
	// TelemetryInterceptor(),
	// )

	grpcService = &exampleService{}

	pb.RegisterExampleServer(service.GrpcServer(), grpcService)

	err = service.AddGatewayHandler(pb.RegisterExampleHandlerFromEndpoint)
	if err != nil {
		fmt.Printf("serve: %v\n", err)
		return
	}

	swag, err = swagger.New(service.Router)
	if err != nil {
		fmt.Printf("Error setting up swagger %s \n", err.Error())
	}

	swag.Title = serviceName
	swag.Version = version
	swag.Schemes = []string{"https"}
	// swag.Path = "/docs"
	data, err := Asset("example/protos/service.swagger.json")
	if err != nil {
		fmt.Printf("Error loading swagger json: %s \n", err.Error())
	}
	swag.JSONData = data

	swag.Serve()

	log.Printf("Listening on port %d", service.Port)

	err = service.Serve()
	if err != nil {
		log.Fatalf("Serve failed")
	}

}

func prepare() {

	viper.SetConfigName("service")
	viper.AddConfigPath("/etc/service/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.service") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// panic(fmt.Errorf("Fatal error config file: %s \n", err))
		fmt.Printf("%s \n  Assuming environment holds the knowledge. \n", err.Error())
	}

	// go viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("Config file changed:", e.Name)

	// 	// initialize()
	// })

	// initialize()

}

func ifError(err error) bool {
	if err != nil {
		log.Printf("Err: %\n", err.Error())
		return true
	}
	return false
}
