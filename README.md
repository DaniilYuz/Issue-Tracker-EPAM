Designed and developed a gRPC-based project tracking microservice in Go, supporting full CRUD and List operations for Users, Issues, and Projects. Defined a language-agnostic API contract using Protocol Buffers and extended it with grpc-gateway annotations to expose RESTful HTTP/JSON endpoints alongside gRPC. Implemented comprehensive business logic including entity validation, issue lifecycle management with status transition rules, and conditional field requirements.

Built a layered storage architecture following SOLID and DRY principles with a clean store interface, enabling seamless switching between in-memory, PostgreSQL, and Redis-cached implementations. Integrated Redis as a caching layer with TTL support to reduce database load on read-heavy operations. Containerized the entire stack (application, PostgreSQL, Redis) using Docker and Docker Compose with persistent volumes for data durability.
Responsibilities:
Designed and implemented a gRPC service API using Protocol Buffers, defining entities, enums, messages, and service methods
Developed REST gateway using grpc-gateway to expose HTTP/JSON endpoints from gRPC definitions
Implemented business logic including entity validation, issue lifecycle management, and status transition rules
Built a storage abstraction layer following SOLID and DRY principles with interchangeable in-memory, PostgreSQL, and Redis-backed implementations
Integrated Redis caching with TTL support to optimize read-heavy database operations
Containerized the application and its dependencies (PostgreSQL, Redis) using Docker and Docker Compose with persistent volume configuration
