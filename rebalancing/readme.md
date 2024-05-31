# kafka rebalancing

## scenario 1: adding 1 more consumer
- topic `test` has 3 partitions: 0,1,2
- publish random messages to topic `test` every 1 second
- start 1 consumer X
- X will be assigned partitions 0,1,2
- start 1 more consumer Y
- observe rebalancing using `kafka-consumer-groups.sh` tool
- start 1 more consumer Z
- observe rebalancing using `kafka-consumer-groups.sh` tool

## scenario 2: removing 1 consumer
- rm 1 consumer Z
- observe rebalancing using `kafka-consumer-groups.sh` tool

## scenario 3: adding 1 more partition
- add 1 more partition to topic `test` using `kafka-topics.sh` tool
- observe rebalancing using `kafka-consumer-groups.sh` tool