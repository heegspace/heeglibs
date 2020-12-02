package heeglibs

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	d *KafKa
)

func TestMain(m *testing.M) {
	conf := &KafkaConfig {
		Topic: "test",
		Group: "",
		Brokers: []string {
			"127.0.0.1:9092",
		}
	}
	d = NewKafka(conf)
	// go func() {
	// 	fmt.Println("start consumer")
	// 	d.Consume(func(err error) {
	// 		fmt.Println(err)
	// 	}, func(noti interface{}) {
	// 		fmt.Println(noti)
	// 	}, func(msg *KafkaMsg) {
	// 		fmt.Println(msg)
	// 	})
	// }()

	os.Exit(m.Run())
}

func TestProduct(t *testing.T) {
	var (
		c   = context.Background()
		msg = []byte("Kafka msg")
	)
	err := d.Product(c, "key", 0, msg)
	err = d.Product(c, "key", 0, msg)
	err = d.Product(c, "key", 0, msg)
	err = d.Product(c, "key", 0, msg)
	err = d.Product(c, "key", 0, msg)

	assert.Nil(t, err)
}
