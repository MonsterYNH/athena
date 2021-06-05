package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/resolver"
)

type EtcdRegistry struct {
	client     *clientv3.Client
	leaseID    clientv3.LeaseID
	clientConn resolver.ClientConn
	config     *RegistryConfig
}

func (etcd *EtcdRegistry) Scheme() string {
	return "etcd"
}

func (etcd *EtcdRegistry) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	etcd.clientConn = cc
	return etcd, etcd.updateAddress()
}

func (etcd *EtcdRegistry) ResolveNow(rn resolver.ResolveNowOptions) {
	log.Println("ResolveNow")
	fmt.Println(rn)
}

//解析器关闭时调用
func (etcd *EtcdRegistry) Close() {
	log.Println("Close")
}

func (etcd *EtcdRegistry) updateAddress() error {
	if etcd.client == nil {
		client, err := clientv3.New(clientv3.Config{
			Endpoints: []string{"localhost:2379"},
		})

		if err != nil {
			return err
		}

		etcd.client = client
	}
	var addrList []resolver.Address

	for _, service := range etcd.config.DependServices {
		resp, err := etcd.client.Get(context.Background(), fmt.Sprintf("%s://%s", etcd.Scheme(), service), clientv3.WithPrefix())
		if err != nil {
			return err
		}

		for index := range resp.Kvs {
			keyEntry, err := url.Parse(string(resp.Kvs[index].Key))
			if err != nil {
				return err
			}
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(keyEntry.Path, "/")})
			log.Println(keyEntry.Path)
		}
	}

	etcd.clientConn.NewAddress(addrList)

	return nil
}

func (etcd *EtcdRegistry) watchAddress(cc resolver.ClientConn) {
	watch := etcd.client.Watch(context.Background(), fmt.Sprintf("%s://", etcd.Scheme()), clientv3.WithPrefix())
	for watchResponse := range watch {
		log.Println("etcd watch:", watchResponse.CompactRevision, etcd.updateAddress())
	}
}

func (etcd *EtcdRegistry) Regist() error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.config.Entrypoints,
		DialTimeout: etcd.config.DialTimeout,
	})

	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s://%s/%s", etcd.Scheme(), etcd.config.Name, etcd.config.IP)
	val, _ := json.Marshal(etcd.config.RouteInfos)

	ctx := context.Background()

	resp, err := client.Grant(ctx, etcd.config.TTL)
	if err != nil {
		return err
	}

	if _, err := client.Put(ctx, key, string(val), clientv3.WithLease(resp.ID)); err != nil {
		return err
	}

	etcd.client = client
	etcd.leaseID = resp.ID

	ch, err := client.KeepAlive(ctx, resp.ID)
	if err != nil {
		return err
	}

	go etcd.listen(ch)

	return nil

}

func (etcd *EtcdRegistry) listen(ch <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case <-etcd.client.Ctx().Done():
			log.Println("ERROR: etcd registry exit with context")
			return
		case _, ok := <-ch:
			if !ok {
				log.Println("ERROR: etcd registry keep alive exist")
				return
			}
		}
	}
}

func (etcd *EtcdRegistry) Stop() error {
	_, err := etcd.client.Revoke(context.Background(), etcd.leaseID)
	return err
}

func (etcd *EtcdRegistry) SetConfig(config *RegistryConfig) error {
	etcd.config = config
	return nil
}
