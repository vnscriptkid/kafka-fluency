## KAFKA

#### Dashboard
- http://localhost:9000
- Add cluster:
  - Name: local
  - Cluster zookeeper hosts: zoo1:2181
  
#### Load Balancing (which partition to send message to)
- Without key: 
  - Sticky
- With key: 
  - Consistent Hashing (same key to same partition -> same consumer)

#### QnA
- Why replicas?
  - Fault tolerance (resilience)
  - High availability
- Why partitions?
  - Parallelism -> high throughput -> Scalability
- Why need message key?
  - Maintain Order of messages
- How to replay an event that has already been acked?
- How to query one event from kafka?
- How to know which msg is on disk, which msg is on RAM?

#### zookeeper vs kraft

#### Todos
- Migrate from zookeeper to kraft

#### Caveats
- `my.first.topic` -> `my_first_topic`

#### kraft setup
- KAFKA_CLUSTER_ID="$(/bin/kafka-storage random-uuid)"
- echo $KAFKA_CLUSTER_ID
- /bin/kafka-storage format -t $KAFKA_CLUSTER_ID -c /etc/kafka/kraft/server.properties
- /bin/kafka-server-start /etc/kafka/kraft/server.properties

#### Assigning partitions to consumers
- Rebalance when 
  - new consumer joins or existing consumer leaves
  - number of partitions changes
