# Zippy Project Documentation

## Project Structure

```
zippy/
├── cmd/
│   └── zippy/           # Main application entry point
├── internal/
│   ├── backup/          # Core backup logic
│   ├── config/          # Configuration handling
│   ├── retention/       # Retention policy logic
│   └── storage/         # Storage abstraction
├── configs/             # Example configurations
├── Dockerfile           # Container image definition
└── README.md
```

## Roadmap

- [ ] .backupignore config
- [ ] Restore functionality
- [ ] Implement retention policy cleanup
- [ ] Add S3 storage backend
- [ ] Support for additional compression algorithms (zstd)
- [ ] Metrics and monitoring endpoints
- [ ] Backup encryption
- [ ] Incremental backups
- [ ] Restore functionality
