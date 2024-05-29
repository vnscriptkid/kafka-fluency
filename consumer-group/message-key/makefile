produce:
	cd ./producer && go run .
	
consume:
	cd ./consumer && go run .

topic_create:
	docker exec -it kafka1 kafka-topics \
		--create \
		--topic demo-topic \
		--partitions 3 \
		--replication-factor 3 \
		--if-not-exists \
		--bootstrap-server kafka1:19092

topic_describe:
	docker exec -it kafka1 kafka-topics \
		--describe \
		--topic demo-topic \
		--bootstrap-server kafka1:19092

consumerg_describe:
	docker exec -it kafka1 kafka-consumer-groups \
		--describe \
		--group demo-group \
		--bootstrap-server kafka1:19092