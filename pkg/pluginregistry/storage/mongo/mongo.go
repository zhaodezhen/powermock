/**
 * @file mongo.go
 * @author dezhenzhao
 * @brief mongo存储：支持添加、删除、查询列表
 * @version 0.1
 * @date 2023-03-14
 * @copyright Copyright (c) 2021 The powermock Authors. All rights reserved.
**/

package mongo

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bilibili-base/powermock/pkg/pluginregistry"
	"github.com/bilibili-base/powermock/pkg/util/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Config defines the config structure
type Config struct {
	Enable      bool
	Host        string
	Port        string
	User        string
	Password    string
	DB          string
	TimeOut     int
	MaxPoolConn int
}

func (c *Config) Validate() error {
	//TODO implement me
	panic("implement me")
}

func (c *Config) RegisterFlagsWithPrefix(prefix string, f *pflag.FlagSet) {
	//TODO implement me
	panic("implement me")
}

func (c *Config) IsEnabled() bool {
	return c.Enable
}

type Plugin struct {
	cfg        *Config
	client     *mongo.Database
	registerer prometheus.Registerer
	logger.Logger

	announcement chan struct{}
}

func (p *Plugin) initMongo() error {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?w=majority", p.cfg.User, p.cfg.Password, p.cfg.Host, p.cfg.Port)
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.cfg.TimeOut))
	defer cancel()
	// 通过传进来的uri连接相关的配置
	o := options.Client().ApplyURI(uri)
	// 设置最大连接数 - 默认是100 ，不设置就是最大 max 64
	o.SetMaxPoolSize(uint64(p.cfg.MaxPoolConn))
	// 发起链接
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return err
	}
	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		fmt.Println("ConnectToDB", err)
		return err
	}
	p.client = client.Database(p.cfg.DB)

	return nil
}

// New 初始化 mongo
func New(cfg *Config, logger logger.Logger, registerer prometheus.Registerer) (pluginregistry.StoragePlugin, error) {
	s := &Plugin{
		cfg:          cfg,
		registerer:   registerer,
		Logger:       logger.NewLogger("mongoPlugin"),
		announcement: make(chan struct{}, 1),
	}
	if err := s.initMongo(); err != nil {
		return nil, err
	}
	s.LogInfo(nil, "start to init mongo(host: %s, port: %s)", cfg.Host, cfg.Port)

	return s, nil
}

// Name is used to return the plugin name
func (p *Plugin) Name() string {
	return "mongo"
}

func (p *Plugin) Set(ctx context.Context, key string, val string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Plugin) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

type MockData struct {
	Path      string
	Method    []string
	UniqueKey string
	Mock      string
}

func getAllMock(db *mongo.Database) []*MockData {
	//findOptions := options.Find()
	//findOptions.SetLimit(10)

	cur, err := db.Collection("mock").Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println(err)
	}
	var results []*MockData
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem MockData
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println(err)
	}
	return results
}

//List 查询出所有的，返回map结构 key=path
func (p *Plugin) List(ctx context.Context) (map[string]string, error) {
	rst := getAllMock(p.client)

	data := make(map[string]string, len(rst))
	for _, item := range rst {
		b, err := base64.StdEncoding.DecodeString(item.Mock)
		if err != nil {
			p.LogError(nil, "base64.StdEncoding.DecodeString, err %s", err.Error())
		}

		data[item.Path] = string(b)
	}
	return data, nil
}

// GetAnnouncement 如果更新的数据，给里面添加一条消息，会自动load 数据
func (p *Plugin) GetAnnouncement() chan struct{} {
	return p.announcement

}

func (p *Plugin) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	return p.watch(ctx, cancelFunc)
}

//TODO  监控，心跳包，mock数据更新
func (p *Plugin) watch(ctx context.Context, cancelFunc context.CancelFunc) error {
	//revision, err := s.getRevision(ctx)
	//if err != nil && err != redis.Nil {
	//	return err
	//}
	//s.LogInfo(nil, "start to watch redis revisions, current: %d", revision)
	//util.StartServiceAsync(ctx, cancelFunc, s.Logger, func() error {
	//	ticker := time.NewTicker(time.Second)
	//	defer ticker.Stop()
	//	for {
	//		select {
	//		case <-ticker.C:
	//			got, err := s.getRevision(ctx)
	//			if err != nil && err != redis.Nil {
	//				s.LogError(nil, "failed to get revision key: %s", err)
	//				continue
	//			}
	//			if revision != got {
	//				s.LogInfo(nil, "event found (known %d vs got %d)", revision, got)
	//				revision = got
	//				select {
	//				case s.announcement <- struct{}{}:
	//				default:
	//				}
	//			}
	//		case <-ctx.Done():
	//			s.LogWarn(nil, "redis stop watching...")
	//			return nil
	//		}
	//	}
	//}, func() error {
	//	return nil
	//})
	return nil
}
