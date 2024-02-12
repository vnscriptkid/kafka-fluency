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

#### zookeeper vs kraft
