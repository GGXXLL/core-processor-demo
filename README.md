### core-process

The functions of this project are as follows:

- consume kafka messages
- handle messages
- batch processing of data

Normally, the directories (files) that users need to pay attention to and modify are as follows:

- `entity` define entity models
- `handler` add handlers
- `internal/config/register.go` register global module or dependency
