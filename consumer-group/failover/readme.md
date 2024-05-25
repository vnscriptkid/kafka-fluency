# Scenario
- 1 topics with 3 partitions: 1, 2, 3
- 1 consumer group with 2 consumers: consumer-1, consumer-2
    - consumer-1 is assigned to partition 1, 2
    - consumer-2 is assigned to partition 3
- 10 messages are produced to the topic (across all partitions)
- 2 consumers start consuming messages
- after 2 messages are consumed by consumer-1, consumer-1 is stopped
- failover happens, rebalance is triggered: consumer-2 is assigned to partition 1, 2
- consumer-2 continues consuming messages from partition 1, 2, 3

# Key: failover -> rebalance
- fault-tolerant