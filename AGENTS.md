# Agent Instructions

zippy is a backup tool that will run as sidecar containers in docker and kubernetes to backup data.

High level overview of how it works:

- A simple backup tool. No incremental backups. Just archives source directories and copies them into the dest directory as a single compressed file.
- Config driven Grandfather-Father-Son style backups. Should either use the destination backup file name or metadata to cleanup based on retention policy.
- Builds into docker image and run as sidecar containers in both docker and kubernetes
- Runs on cron schedule

Config:

- Use yaml for config
- Source can be one or many directories
- Destination can be a mount dir or object storage bucket

## Coding style

- Built in go
- Prefer functional patterns. Avoid side effects.
- Write defensive code and write tests to cover possible use cases
- Follow modern go code patterns and practices leaning more towards established patterns that's used in big orgs
