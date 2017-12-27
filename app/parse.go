package app

import (
	"io/ioutil"

	"github.com/Shopify/sarama"

	"strings"

	"fmt"

	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Conf *Http2MQ

type WebConf struct {
	AccessLog string `yaml:"access_log"`
	ErrorLog  string `yaml:"error_log"`
	Port      int    `yaml:"port"`
}

type KafkaConf struct {
	Brokers string `yaml:"brokers"`
	Topic   string `yaml:"topic"`

	SyncProducer  sarama.SyncProducer  `yaml:"-"`
	AsyncProducer sarama.AsyncProducer `yaml:"-"`
}

type AuthUser struct {
	Name     string
	Password string
}

type Http2MQ struct {
	WebConf   `yaml:"web"`
	KafkaConf `yaml:"kafka"`
	UserConf  []string            `yaml:"user"`
	User      map[string]AuthUser `yaml:"-"`
}

func InitConf(file string) (*Http2MQ, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c, err := parse(buf)
	if err != nil {
		return nil, err
	}

	brokers := strings.Split(c.KafkaConf.Brokers, ",")
	c.KafkaConf.SyncProducer = newSyncProducer(brokers)
	c.KafkaConf.AsyncProducer = newAsyncProducer(brokers)

	Conf = c

	return c, nil
}

func newSyncProducer(brokers []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	//tlsConfig := createTlsConfiguration()
	//if tlsConfig != nil {
	//	config.Net.TLS.Config = tlsConfig
	//	config.Net.TLS.Enable = true
	//}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logrus.Warnf("init AsyncProducer err:%s", err.Error())
	}

	return producer
}

func newAsyncProducer(brokers []string) sarama.AsyncProducer {

	config := sarama.NewConfig()
	//tlsConfig := createTlsConfiguration()
	//if tlsConfig != nil {
	//	config.Net.TLS.Enable = true
	//	config.Net.TLS.Config = tlsConfig
	//}
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		logrus.Warnf("init AsyncProducer err:%s", err.Error())
	}

	return producer
}

func parse(d []byte) (*Http2MQ, error) {
	c := &Http2MQ{}

	if err := yaml.Unmarshal(d, c); err != nil {
		return nil, err
	}

	c.User = make(map[string]AuthUser)
	for _, v := range c.UserConf {
		ds := strings.Split(v, ":")
		if len(ds) != 2 {
			return nil, fmt.Errorf("user must be name:password, error in :%s", v)
		}
		c.User[ds[0]] = AuthUser{
			Name:     ds[0],
			Password: ds[1],
		}
	}

	return c, nil
}
