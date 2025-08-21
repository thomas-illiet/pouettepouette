package supervisor

import (
	"common/log"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"supervisor/pkg/config"
	"supervisor/pkg/editor"
	"supervisor/pkg/service"
	"supervisor/pkg/service/system"
	"supervisor/pkg/service/utility"
	"sync"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	Version = ""
)

// Run serves as main entrypoint to the supervisor.
func Run() {
	exitCode := 0
	defer handleExit(&exitCode)
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		log.WithField("cause", r).WithField("stack", string(debug.Stack())).Error("panicked")
		if exitCode == 0 {
			exitCode = 1
		}
	}()

	// Load supervisor configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("configuration error")
	}

	// Check if the program is called with "run" as an argument to start the supervisor.
	if len(os.Args) < 2 || os.Args[1] != "run" {
		fmt.Println("supervisor makes sure your workspace/Editor keeps running smoothly.\n" +
			"You don't have to call this thing, Opencoder calls it for you.")
		return
	}

	// Set git credential helper configuration
	//configureGit(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	// Start editor
	var ideWG sync.WaitGroup
	var ideReady = editor.NewEditorReadyState()
	ideWG.Add(1)
	go editor.StartAndWatchEditor(ctx, cfg, &ideWG, ideReady)

	//
	var wg sync.WaitGroup
	wg.Add(1)
	services := []service.RegisterableService{
		&system.SystemService{Cfg: cfg},
		&utility.UtilityService{},
	}
	services = append(services)
	go startGrpcEndpoint(ctx, cfg, &wg, services)

	// to shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigChan:
	}

	log.Info("received SIGTERM (or shutdown) - tearing down")

	cancel()
	ideWG.Wait()
	wg.Wait()
}

func startGrpcEndpoint(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup, services []service.RegisterableService) {
	defer wg.Done()
	defer log.Debug("startGrpcEndpoint shutdown")

	//
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.APIEndpointPort))
	if err != nil {
		log.WithError(err).Fatal("cannot start health endpoint")
	}

	//
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	//if cfg.DebugEnable {
	//	unaryInterceptors = append(unaryInterceptors, grpc_logrus.UnaryServerInterceptor(log.Log))
	//	streamInterceptors = append(streamInterceptors, grpc_logrus.StreamServerInterceptor(log.Log))
	//}

	// Add gprc recover, must be last, to be executed first after the rpc handler,
	// we want upstream interceptors to have a meaningful response to work with
	unaryInterceptors = append(unaryInterceptors, grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p interface{}) error {
			log.WithField("stack", string(debug.Stack())).Errorf("[PANIC] %s", p)
			return status.Errorf(codes.Internal, "%s", p)
		},
	)))
	streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p interface{}) error {
			log.WithField("stack", string(debug.Stack())).Errorf("[PANIC] %s", p)
			return status.Errorf(codes.Internal, "%s", p)
		},
	)))

	//
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	)

	grpcServer := grpc.NewServer(opts...)
	for _, reg := range services {
		if reg, ok := reg.(service.RegisterableGRPCService); ok {
			reg.RegisterGRPC(grpcServer)
		}
	}

	// for debuging
	reflection.Register(grpcServer)

	err = grpcServer.Serve(l)
	if err != nil {
		log.WithError(err).Fatal("cannot start grpc server")
	}

	//
	<-ctx.Done()

	//
	log.Info("shutting down grpc endpoint")
	grpcServer.GracefulStop()
}

func handleExit(ec *int) {
	exitCode := *ec
	log.WithField("exitCode", exitCode).Debug("supervisor exit")
	os.Exit(exitCode)
}
