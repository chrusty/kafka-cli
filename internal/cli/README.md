CLI
===

Commands and logic for the Kafka CLI.

- Built around [Cobra](github.com/spf13/cobra) which makes beautiful CLIs easy.
- Commands are broken up and grouped into different files here, with an eye on keeping it easy to expand with new commands.
- A CLI struct is provided here which ties the various dependencies together
- Cobra commands should be registered with this CLI struct
