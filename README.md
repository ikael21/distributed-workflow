# Distributed Workflow / Job Orchestration Platform

This project simulates a production-grade distributed workflow orchestration platform similar to CI/CD systems or data pipelines. The platform allows users to define workflows consisting of multiple tasks, execute them across distributed workers, store execution metadata, and receive notifications about results.

The system is intentionally designed to practice real-world backend engineering concepts:

- Go microservices
- gRPC and HTTP APIs
- Kafka event-driven architecture
- Envoy API Gateway
- MongoDB for flexible job metadata
- PostgreSQL for relational services

## Project Structure

```bash
distributed-workflow/
├── infra/
│   ├── envoy/               # Envoy API Gateway configuration
│   └── compose.yml          # Docker Compose configuration
│
├── pkg/                     # Shared packages
│
├── services/
│   ├── artifact/            # Artifact storage service
│   ├── iam/                 # Identity and Access Management service
│   ├── metadata/            # Metadata storage service
│   ├── notification/        # Notification service
│   ├── scheduler/           # Workflow/job scheduler service
│   ├── worker/              # Distributed worker service
│   └── workflow/            # Workflow definition and management service
│
├── Makefile                 # Build automation and development commands
└── README.md                # Project documentation
```
