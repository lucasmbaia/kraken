package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"reflect"
	"context"
	"errors"
	"encoding/json"
	pb "go.etcd.io/etcd/mvcc/mvccpb"
)

type Response struct {
	Key	string
	Value	string
	Action	pb.Event_EventType
}

type Config struct {
	Username	string
	Password	string
	Endpoints	[]string
	Timeout		int32
}

type Client struct {
	client	*clientv3.Client
	timeout	int32
}

func NewClient(cfg Config) (c Client, err error) {
	if c.client, err = clientv3.New(clientv3.Config{
		Endpoints:	cfg.Endpoints,
		Username:	cfg.Username,
		Password:	cfg.Password,
		DialTimeout:	time.Duration(cfg.Timeout) * time.Second,
	}); err != nil {
		return
	}

	c.timeout = cfg.Timeout
	return
}

func (c *Client) Get(key string, value interface{}) (err error) {
	var (
		response	*clientv3.GetResponse
		ctx		context.Context
		cancel		context.CancelFunc
	)

	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("Expected a pointer to a variable")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.timeout) * time.Second)
	defer cancel()

	if response, err = c.client.Get(ctx, key); err != nil {
		return
	}

	for _, ev := range response.Kvs {
		if err = json.Unmarshal(ev.Value, value); err != nil {
			return
		}
	}

	return
}

func (c *Client) Put(key, value string) (err error) {
	var (
		ctx		context.Context
		cancel		context.CancelFunc
	)

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(c.timeout) * time.Second)
	defer cancel()

	_, err = c.client.Put(ctx, key, value)
	return
}

func (c *Client) Watch(key string, values chan<- Response) {
	var (
		watch		clientv3.WatchChan
	)

	watch = c.client.Watch(context.Background(), key, clientv3.WithPrefix())

	for w := range watch {
		for _, ev := range w.Events {
			values <- Response{
				Action:	ev.Type,
				Key:	string(ev.Kv.Key),
				Value:	string(ev.Kv.Value),
			}
		}
	}

	return
}
