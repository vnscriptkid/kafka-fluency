# Config

## `auto.offset.reset`
- applicable when no last committed offset is found for a consumer group (e.g. new consumer group, or the committed offset is deleted due to retention policy)
- Options:
  - `earliest`: start reading from the earliest offset
  - `latest`: start reading from the latest offset
  - `none`: throw an exception if no offset is found for the consumer group
- Default: `latest`
- Best practice: set to `earliest`