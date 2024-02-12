up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

# simulate one broker goes down
down_broker:
	docker stop kafka1

# simulate one broker goes up
up_broker:
	docker start kafka1

# create a topic userActivity
topic_create:
	docker exec -it kafka1 kafka-topics \
		--create \
		--topic userActivity \
		--partitions 3 \
		--replication-factor 3 \
		--if-not-exists \
		--bootstrap-server kafka1:19092

# list all topics
topic_list:
	docker exec -it kafka1 kafka-topics \
		--list \
		--bootstrap-server kafka1:19092

# delete a topic userActivity
topic_delete:
	docker exec -it kafka1 kafka-topics \
		--delete \
		--topic userActivity \
		--bootstrap-server kafka1:19092

# produce message to the topic userActivity
produce:
	docker exec -it kafka1 kafka-console-producer \
		--topic userActivity \
		--bootstrap-server kafka1:19092

# {"userId": "user1", "event": "login", "timestamp": "2024-02-12T10:00:00Z"}
# {"userId": "user2", "event": "postMessage", "content": "This is a message", "timestamp": "2024-02-12T10:05:00Z"}
# {"userId": "user3", "event": "logout", "timestamp": "2024-02-12T10:10:00Z"}

# consume messages from the topic userActivity
# each terminal is one consumer, that belongs to different consumer group
consume:
	docker exec -it kafka1 kafka-console-consumer \
		--topic userActivity \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

consume_groupa:
	docker exec -it kafka1 kafka-console-consumer \
		--topic userActivity \
		--group a \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

consume_groupb:
	docker exec -it kafka1 kafka-console-consumer \
		--topic userActivity \
		--group b \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

# version of kafka
version:
	docker exec -it kafka1 kafka-topics --version

# produce with key: [key,content]
produce_with_key:
	docker exec -it kafka1 kafka-console-producer \
		--topic userActivity \
		--property parse.key=true \
		--property key.separator=, \
		--bootstrap-server kafka1:19092



	