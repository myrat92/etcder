package etcd

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/myrat92/etcder/internal/engine/infrastructure/etcdv3"
)

type Operator interface {
	ListAll() ([]string, error)
	Get(key string) string
	Update(key, value string) error

	Health() error
}

var operate *Operate

type Operate struct {
	cli *etcdv3.Etcdv3
}

func NewInstance(endpoint string) {
	cli, _ := etcdv3.NewEtcdv3Client(endpoint)
	operate = &Operate{
		cli: cli,
	}
}

func Instance() *Operate {
	return operate
}

func (o *Operate) ListAll() ([]string, error) {
	return o.cli.List(context.Background(), "/")
}

func (o *Operate) Get(key string) string {
	res, err := o.cli.Get(context.Background(), key)
	if err != nil {
		slog.Warn("get key", err)
		return ""
	}
	return res
}

func (o *Operate) Health() error {
	return o.cli.Health(context.Background())
}

func (o *Operate) Update(key, value string) error {
	_, err := o.cli.Update(context.Background(), key, value)
	return err
}
