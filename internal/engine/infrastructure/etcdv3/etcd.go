package etcdv3

import (
	"context"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcdv3 struct {
	client *clientv3.Client
}

const (
	Timeout time.Duration = 5 * time.Second
)

func NewEtcdv3Client(endpoints string) (*Etcdv3, error) {
	endpoint := strings.Split(endpoints, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &Etcdv3{client: cli}, err
}

func (e *Etcdv3) List(ctx context.Context, key string) ([]string, error) {
	ctx, _ = context.WithTimeout(ctx, Timeout)
	getResp, err := e.client.KV.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var res []string
	for _, kvs := range getResp.Kvs {
		res = append(res, string(kvs.Key))
	}

	return res, nil
}

func (e *Etcdv3) Get(ctx context.Context, key string) (string, error) {
	ctx, _ = context.WithTimeout(ctx, Timeout)
	getResp, err := e.client.KV.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(getResp.Kvs) == 0 {
		return "", NewKeyNotFoundError(key, 0)
	}

	return string(getResp.Kvs[0].Value), nil
}

func (e *Etcdv3) Health(ctx context.Context) error {
	ctx, _ = context.WithTimeout(ctx, Timeout)
	_, err := e.client.Status(ctx, e.client.Endpoints()[0])
	return err
}

func (e *Etcdv3) Update(ctx context.Context, key, value string) (*clientv3.PutResponse, error) {
	ctx, _ = context.WithTimeout(ctx, Timeout)
	return e.client.KV.Put(ctx, key, value)
}
