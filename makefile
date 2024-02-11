up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

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

# produce message to the topic userActivity
produce:
	docker exec -it kafka1 kafka-console-producer \
		--topic userActivity \
		--bootstrap-server kafka1:19092

# {"userId": "user1", "event": "login", "timestamp": "2024-02-12T10:00:00Z"}
# {"userId": "user2", "event": "postMessage", "content": "This is a message", "timestamp": "2024-02-12T10:05:00Z"}
# {"userId": "user3", "event": "logout", "timestamp": "2024-02-12T10:10:00Z"}

# consume messages from the topic userActivity
consume:
	docker exec -it kafka1 kafka-console-consumer \
		--topic userActivity \
		--from-beginning \
		--bootstrap-server kafka1:19092