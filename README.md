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
- `admin consume <topic>`: Consume messages from a specific topic (optionally with a consumer-group ID). Messages go to STDOUT, logs to STDERR.


The following commands are under development:

- `admin config params`: Print various config parameters
- `admin groups delete <group>`: Delete a group


Config
------

Instead of having to provide your config as parameters every time, just set some env-vars:

- `KAFKA_BOOTSTRAPSERVERS`: The Kafka brokers to connect to ("**localhost:9092**")
- `KAFKA_PASSWORD`: The SASL password to authenticate with _(optional)_
- `KAFKA_USERNAME`: The SASL username to authenticate with _(optional)_
- `KAFKA_SASLMECHANISM`: The mechanism for SASL auth ["SCRAM-SHA-256", "**SCRAM-SHA-512**"]
- `KAFKA_SECURITYPROTOCOL`: The security protocol ["SASL_SSL", "SASL_PLAINTEXT", "SSL", "**PLAINTEXT**"]
