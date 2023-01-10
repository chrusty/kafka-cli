# kafka-cli

A powerful and easy-to-use CLI tool for Kafka

The default cli-tools for Kafka are very powerful, but are also rather inconsistent and sometimes baffling. This project aims to provide a simple CLI tool for Kafka, starting with the most common operations.

Requests and suggestions are welcome!


Commands
--------

So far the following commands are supported:

- `admin config metadata`: Print various metadata about the Kafka cluster and brokers
- `admin groups list`: List groups
- `admin groups describe <group>`: Describe a specific group
- `admin topics list`: List topics
- `admin topics describe <topic>`: Describe the config for a specific topic
- `admin topics delete <topic>`: Delete a topic
- `admin consume <topic>`: Consume messages from a specific topic (optionally with a consumer-group ID)


The following commands are under development:

- `admin config params`: Print various config parameters
- `admin groups delete <group>`: Delete a group
