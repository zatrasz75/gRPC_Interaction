package app

import (
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"zatrasz75/gRPC_Interaction/users/configs"
	"zatrasz75/gRPC_Interaction/users/internal/handlers"
	"zatrasz75/gRPC_Interaction/users/internal/repository"
	"zatrasz75/gRPC_Interaction/users/internal/tokens"
	"zatrasz75/gRPC_Interaction/users/pkg/gRPCLogger"
	"zatrasz75/gRPC_Interaction/users/pkg/gRPCserver"
	usersPb "zatrasz75/gRPC_Interaction/users/pkg/grpc/users"
	"zatrasz75/gRPC_Interaction/users/pkg/logger"
	"zatrasz75/gRPC_Interaction/users/pkg/postgres"
)

func Run(cfg *configs.Config, l logger.LoggersInterface) {
	pg, err := postgres.New(cfg.PostgreSQL.ConnStr, l, postgres.OptionSet(cfg.PostgreSQL.PoolMax, cfg.PostgreSQL.ConnAttempts, cfg.PostgreSQL.ConnTimeout))
	if err != nil {
		l.Fatal("ошибка запуска - postgres.New:", err)
	}
	defer pg.Close()

	err = pg.Migrate(l)
	if err != nil {
		l.Fatal("ошибка миграции", err)
	}
	repo := repository.New(pg)
	storeHandler := handlers.New(repo, l, cfg)

	lg := gRPCLogger.NewGRPCLogger()
	authInterceptor := tokens.AuthInterceptor(cfg.Token.SecretKeyHere)
	loggingInterceptor := gRPCLogger.LogUnaryServerInterceptor(lg)

	// Передаем список interceptor'ов
	interceptors := []grpc.UnaryServerInterceptor{
		loggingInterceptor,
		authInterceptor,
	}

	options := gRPCserver.OptionSet(cfg.GRPC.AddrHost, cfg.GRPC.AddrPort)
	lis := gRPCserver.New(interceptors, options)

	registerHandlers(lis.GrpcServer, storeHandler)

	go func() {
		err = lis.Start()
		if err != nil {
			l.Error("Остановка сервера:", err)
		} else {
			<-lis.Notify()
		}
	}()

	l.Info("Запуск gRPC сервера на " + cfg.GRPC.AddrHost + ":" + cfg.GRPC.AddrPort)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("принят сигнал прерывания прерывание %s , Остановка сервера gRPC", s.String())
	case err := <-lis.Notify():
		l.Error("получена ошибка сигнала прерывания сервера", err)
	}
}

func registerHandlers(s *grpc.Server, handler *handlers.Store) {
	usersPb.RegisterProfileServer(s, handler)
}
