# Distributed Workflow / Job Orchestration Platform

This project simulates a production-grade distributed workflow orchestration platform similar to CI/CD systems or data pipelines. The platform allows users to define workflows consisting of multiple tasks, execute them across distributed workers, store execution metadata, and receive notifications about results.

The system is intentionally designed to practice real-world backend engineering concepts:
- Go microservices
- gRPC and HTTP APIs
- Kafka event-driven architecture
- Envoy API Gateway
- MongoDB for flexible job metadata
- PostgreSQL for relational services
