package deserializers

import (
	"encoding/json"

	"github.com/vnscriptkid/kafka-fluency/consumer/domains"
)

func DeserializeOrderCreated(data []byte) (*domains.OrderCreated, error) {
	var orderCreated domains.OrderCreated
	err := json.Unmarshal(data, &orderCreated)
	if err != nil {
		return nil, err
	}
	return &orderCreated, nil
}
