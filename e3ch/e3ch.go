package e3ch

import (
	"crypto/tls"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/soyking/e3ch"
	"github.com/soyking/e3w/conf"
)

var tlsConfig *tls.Config

func NewE3chClient(config *conf.Config) (*client.EtcdHRCHYClient, error) {
	var err error

	tlsInfo := transport.TLSInfo{
		CertFile:           config.CertFile,
		KeyFile:            config.KeyFile,
		TrustedCAFile:      config.CAFile,
		InsecureSkipVerify: true,
	}
	tlsConfig, err = tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	clt, err := clientv3.New(clientv3.Config{
		Endpoints: config.EtcdEndPoints,
		Username:  config.EtcdUsername,
		Password:  config.EtcdPassword,
		TLS:       tlsConfig,
	})
	if err != nil {
		return nil, err
	}

	client, err := client.New(clt, config.EtcdRootKey, config.DirValue)
	if err != nil {
		return nil, err
	}
	return client, client.FormatRootKey()
}

func CloneE3chClient(username, password string, client *client.EtcdHRCHYClient) (*client.EtcdHRCHYClient, error) {
	clt, err := clientv3.New(clientv3.Config{
		Endpoints: client.EtcdClient().Endpoints(),
		Username:  username,
		Password:  password,
		TLS:       tlsConfig,
	})
	if err != nil {
		return nil, err
	}
	return client.Clone(clt), nil
}
