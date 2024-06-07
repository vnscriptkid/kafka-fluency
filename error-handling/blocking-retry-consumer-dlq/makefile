# produce message to the topic myTopic
produce:
	docker exec -it kafka1 kafka-console-producer \
		--topic myTopic \
		--bootstrap-server kafka1:19092

consume-dlq:
	docker exec -it kafka1 kafka-console-consumer \
		--topic myDLQ \
		--from-beginning \
		--bootstrap-server kafka1:19092 \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--property print.timestamp=true \
		--property key.separator=, \
		--max-messages 1 \
		--group myDLQGroup