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

# create a topic user.activity
topic_create:
	docker exec -it kafka1 kafka-topics \
		--create \
		--topic user.activity \
		--partitions 3 \
		--replication-factor 3 \
		--if-not-exists \
		--bootstrap-server kafka1:19092

# list all topics
topic_list:
	docker exec -it kafka1 kafka-topics \
		--list \
		--bootstrap-server kafka1:19092

# describe a topic user.activity
topic_describe:
	docker exec -it kafka1 kafka-topics \
		--describe \
		--topic user.activity \
		--bootstrap-server kafka1:19092

# delete a topic user.activity
topic_delete:
	docker exec -it kafka1 kafka-topics \
		--delete \
		--topic user.activity \
		--bootstrap-server kafka1:19092

# increase the number of partitions of a topic by 1
topic_more_partition:
	docker exec -it kafka1 kafka-topics \
		--alter \
		--topic user.activity \
		--partitions 4 \
		--bootstrap-server kafka1:19092

# produce message to the topic user.activity
produce:
	docker exec -it kafka1 kafka-console-producer \
		--topic user.activity \
		--bootstrap-server kafka1:19092

# {"userId": "user1", "event": "login", "timestamp": "2024-02-12T10:00:00Z"}
# {"userId": "user2", "event": "postMessage", "content": "This is a message", "timestamp": "2024-02-12T10:05:00Z"}
# {"userId": "user3", "event": "logout", "timestamp": "2024-02-12T10:10:00Z"}

# consume messages from the topic user.activity
# each terminal is one consumer, that belongs to different consumer group
consume:
	docker exec -it kafka1 kafka-console-consumer \
		--topic user.activity \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

consume_groupa:
	docker exec -it kafka1 kafka-console-consumer \
		--topic user.activity \
		--group group.a \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

consume_groupb:
	docker exec -it kafka1 kafka-console-consumer \
		--topic user.activity \
		--group group.b \
		--from-beginning \
		--property print.key=true \
		--property print.partition=true \
		--property print.offset=true \
		--bootstrap-server kafka1:19092

# list all consumer groups
consumer_groups:
	docker exec -it kafka1 kafka-consumer-groups \
		--list \
		--bootstrap-server kafka1:19092

# describe a consumer group
consumer_groupa:
	docker exec -it kafka1 kafka-consumer-groups \
		--describe \
		--group group.a \
		--bootstrap-server kafka1:19092

# check health of consumer group
consumer_groupa_health:
	docker exec -it kafka1 kafka-consumer-groups \
		--group group.a \
		--describe \
		--state \
		--bootstrap-server kafka1:19092

# check members of consumer group
consumer_groupa_members:
	docker exec -it kafka1 kafka-consumer-groups \
		--group group.a \
		--describe \
		--members \
		--bootstrap-server kafka1:19092

# version of kafka
version:
	docker exec -it kafka1 kafka-topics --version

# produce with key: [key,content]
produce_with_key:
	docker exec -it kafka1 kafka-console-producer \
		--topic user.activity \
		--property parse.key=true \
		--property key.separator=, \
		--bootstrap-server kafka1:19092

# viewing all existing consumer groups
consumer_groups:
	docker exec -it kafka1 kafka-consumer-groups \
		--list \
		--bootstrap-server kafka1:19092
	
# kafka1
kafka1:
	docker exec -it kafka1 bash

