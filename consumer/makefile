TOPIC_ORDER_CREATED=order.created
CONSUMER_GROUP=dispatch.order.created.group
BOOTSTRAP_SERVER=kafka1:19092
CONTAINER=kafka1

# Message: {"orderID":"123","item":"item 1"}
# Message: {"orderID":"123","item":"item 2"}
# Message: 1,{"orderID":"100","item":"item 3"}
# Message: 2,{"orderID":"100","item":"item 4"}
# Message: 3,{"orderID":"100","item":"item 5"}
# Message: 4,{"orderID":"100","item":"item 6"}
# Message: 4,{"orderID":"100","item":"item 8"}

produce:
	docker exec -it $(CONTAINER) kafka-console-producer \
		--topic $(TOPIC_ORDER_CREATED) \
		--bootstrap-server $(BOOTSTRAP_SERVER)

produce_key:
	docker exec -it $(CONTAINER) kafka-console-producer \
		--topic $(TOPIC_ORDER_CREATED) \
		--property parse.key=true \
		--property key.separator=, \
		--bootstrap-server $(BOOTSTRAP_SERVER)

group_members:
	docker exec -it $(CONTAINER) kafka-consumer-groups \
		--group $(CONSUMER_GROUP) \
		--describe \
		--members \
		--bootstrap-server $(BOOTSTRAP_SERVER)

group_desc:
	docker exec -it $(CONTAINER) kafka-consumer-groups \
		--describe \
		--group $(CONSUMER_GROUP) \
		--bootstrap-server $(BOOTSTRAP_SERVER)

topic_delete:
	docker exec -it $(CONTAINER) kafka-topics \
		--delete \
		--topic $(TOPIC_ORDER_CREATED) \
		--bootstrap-server $(BOOTSTRAP_SERVER)

topic_create:
	docker exec -it $(CONTAINER) kafka-topics \
		--create \
		--topic $(TOPIC_ORDER_CREATED) \
		--partitions 3 \
		--replication-factor 3 \
		--bootstrap-server $(BOOTSTRAP_SERVER)