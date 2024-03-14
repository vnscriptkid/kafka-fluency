package processors

import (
	"fmt"

	"github.com/vnscriptkid/kafka-fluency/consumer/domains"
)

type DispatchService struct{}

func (d DispatchService) Process(orderCreated *domains.OrderCreated) {
	fmt.Printf("Dispatching order: %+v\n", orderCreated)
	// Here, you could also produce a message to the `order-dispatched` topic
}

func NewDispatchService() DispatchService {
	return DispatchService{}
}
