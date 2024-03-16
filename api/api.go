package api

import (
	"context"
	"fmt"
	"github.com/Filecoin-Titan/titan/api"
	"github.com/Filecoin-Titan/titan/api/client"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-explorer/config"
	"github.com/gnasnik/titan-explorer/core/cleanup"
	"github.com/gnasnik/titan-explorer/core/statistics"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var (
	DefaultAreaId            = "Asia-China-Guangdong-Shenzhen"
	SchedulerConfigKeyPrefix = "TITAN::SCHEDULERCFG"
)

type Server struct {
	cfg             config.Config
	router          *gin.Engine
	etcdClient      *statistics.EtcdClient
	statistic       *statistics.Statistic
	statisticCloser func()
}

func NewServer(cfg config.Config) (*Server, error) {
	gin.SetMode(cfg.Mode)
	router := gin.Default()

	//router.Use(Cors())

	// logging request body
	router.Use(RequestLoggerMiddleware())

	RegisterRouters(router, cfg)

	etcdClient, err := statistics.NewEtcdClient([]string{cfg.EtcdAddress})
	if err != nil {
		log.Errorf("New etcdClient Failed: %v", err)
		return nil, err
	}

	s := &Server{
		cfg:        cfg,
		router:     router,
		statistic:  statistics.New(cfg.Statistic, etcdClient),
		etcdClient: etcdClient,
	}

	go cleanup.Run(context.Background())

	return s, nil
}

func (s *Server) Run() {
	s.statistic.Run()
	err := s.router.Run(s.cfg.ApiListen)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Close() {
	s.statistic.Stop()
}

func getSchedulerClient(ctx context.Context, areaId string) (api.Scheduler, error) {
	schedulers, err := statistics.GetSchedulerConfigs(ctx, fmt.Sprintf("%s::%s", SchedulerConfigKeyPrefix, areaId))
	if err == redis.Nil && areaId != DefaultAreaId {
		return getSchedulerClient(ctx, DefaultAreaId)
	}

	if err != nil || len(schedulers) == 0 {
		log.Errorf("no scheduler found")
		return nil, errors.New("no scheduler found")
	}

	schedulerApiUrl := schedulers[0].SchedulerURL
	schedulerApiToken := schedulers[0].AccessToken
	SchedulerURL := strings.Replace(schedulerApiUrl, "https", "http", 1)
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+schedulerApiToken)
	schedulerClient, _, err := client.NewScheduler(ctx, SchedulerURL, headers)
	if err != nil {
		log.Errorf("create scheduler rpc client: %v", err)
	}

	return schedulerClient, nil
}
