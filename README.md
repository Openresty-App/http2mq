# http2mq 
send http body to kafka


# install
```
go get -u github.com/Openresty-App/http2mq

go run main.go

```

# configure
```
web:
  port: 8080  #http listen port


kafka:
  brokers: 192.168.1.182:9092 # kafka brokers
  topic: http2mq # topic


user:
    - user1:pwd1 # username:password,  basic auth
    - user2:pwd2
    - user3:pwd3
```
