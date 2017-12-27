package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data := `
web:
  access_log: /var/log/http2mq/access.log
  error_log: /var/log/http2mq/error.log
  port: 8080


kafka:
  brokers: 192.168.1.182:9092
  topic: http2mq


user:
    - user1:pwd1
    - user2:pwd2
    - user3:pwd3
`

	c, err := parse([]byte(data))
	assert.Equal(t, nil, err)
	assert.Equal(t, "/var/log/http2mq/access.log", c.AccessLog)
	assert.Equal(t, "/var/log/http2mq/error.log", c.ErrorLog)
	assert.Equal(t, 8080, c.Port)

	assert.Equal(t, "192.168.1.182:9092", c.KafkaConf.Brokers)
	assert.Equal(t, "http2mq", c.KafkaConf.Topic)

	assert.Equal(t, []string{"user1:pwd1", "user2:pwd2", "user3:pwd3"}, c.UserConf)
	assert.Equal(t, AuthUser{Name: "user1", Password: "pwd1"}, c.User["user1"])
	assert.Equal(t, AuthUser{Name: "user2", Password: "pwd2"}, c.User["user2"])
	assert.Equal(t, AuthUser{Name: "user3", Password: "pwd3"}, c.User["user3"])
}
