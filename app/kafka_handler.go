package app

import (
	"net/http"

	"io/ioutil"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
)

type Kafka struct {
}

func (k *Kafka) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !CheckAuth(req) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logrus.Warnf("%s: res:%s", req.RequestURI, err.Error())
		res.Write([]byte(""))
		return
	}

	Conf.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: "http2mq",
		Value: sarama.ByteEncoder(b),
	})
	//
	//Conf.AsyncProducer.Input() <- &sarama.ProducerMessage{
	//	Topic: "http2mq",
	//	Value: sarama.ByteEncoder(b),
	//}

	res.Write([]byte(""))
}

func NewKafka() http.Handler {
	return &Kafka{}
}
