# Scenario
- Topic with 3 partitions and 1 consumer (number of consumers <= number of partitions)
    - The consumer reads from all 3 partitions
- Publish 10 messages with key [0, 9]
- Consumer leaves the group
- 2 new consumers join the group (of the same group id)
    - Should not re-read the messages that were already read by the previous consumer
    - one consumer is assigned 2 partitions, the other is assigned 1 partition
- Publish 10 messages with key [0, 9]
- Publish 10 messages with key [0, 9]