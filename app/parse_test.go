package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data := `
web:
  port: 8080
  access_log: /var/log/http2mq/access.log
  error_log: /var/log/http2mq/error.log


kafka:
  brokers: 192.168.1.182:9092
  zk: 192.168.1.169:2181,192.168.1.179:2181,192.168.1.180:2181
  zk_root: /qianbao/kafka/logtest
  topic: http2mq
  consumer_user: http2mq_unit_test

users:
    - user1:pwd1
    - user2:pwd2
    - user3:pwd3

topics:
    - ipos
    - http2mq

`

	c, err := parse([]byte(data))
	assert.Equal(t, nil, err)
	assert.Equal(t, "/var/log/http2mq/access.log", c.AccessLog)
	assert.Equal(t, "/var/log/http2mq/error.log", c.ErrorLog)
	assert.Equal(t, 8080, c.Port)

	assert.Equal(t, "192.168.1.182:9092", c.KafkaConf.Brokers)
	assert.Equal(t, "192.168.1.169:2181,192.168.1.179:2181,192.168.1.180:2181", c.KafkaConf.ZK)
	assert.Equal(t, []string{"192.168.1.169:2181", "192.168.1.179:2181", "192.168.1.180:2181"}, c.KafkaConf.ZKServers)
	assert.Equal(t, "/qianbao/kafka/logtest", c.KafkaConf.ZkRoot)
	assert.Equal(t, "http2mq", c.KafkaConf.Topic)
	assert.Equal(t, "http2mq_unit_test", c.KafkaConf.ConsumerUser)

	assert.Equal(t, []string{"user1:pwd1", "user2:pwd2", "user3:pwd3"}, c.UserConf)
	assert.Equal(t, AuthUser{Name: "user1", Password: "pwd1"}, c.User["user1"])
	assert.Equal(t, AuthUser{Name: "user2", Password: "pwd2"}, c.User["user2"])
	assert.Equal(t, AuthUser{Name: "user3", Password: "pwd3"}, c.User["user3"])

	assert.Equal(t, []string{"ipos", "http2mq"}, c.Topics)
}
