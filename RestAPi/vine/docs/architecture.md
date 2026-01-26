
vine-automation/
│
├── cmd/                                    # Application entry points
│   ├── server/                             # HTTP API server
│   │   └── main.go
│   ├── worker/                             # Background job worker
│   │   └── main.go
│   └── migrate/                            # Database migration CLI
│       └── main.go
│
├── internal/                               # Private application code
│   │
│   ├── app/                                # Application initialization layer
│   │   ├── server.go                       # Server composition & DI
│   │   ├── worker.go                       # Worker composition & DI
│   │   └── router.go                       # HTTP route composition
│   │
│   ├── booking/                            # Booking domain (bounded context)
│   │   ├── domain/                         # Business entities & logic
│   │   │   ├── booking.go                  # Booking entity
│   │   │   ├── prisoner.go                 # Prisoner entity
│   │   │   ├── charge.go                   # Charge entity
│   │   │   ├── arrest.go                   # Arrest entity
│   │   │   ├── facility.go                 # Facility entity
│   │   │   ├── release.go                  # Release entity
│   │   │   ├── mugshot.go                  # Mugshot entity
│   │   │   ├── repository.go               # Repository interfaces
│   │   │   └── errors.go                   # Domain-specific errors
│   │   │
│   │   ├── repository/                     # Data access implementations
│   │   │   ├── booking.go                  # Booking repository
│   │   │   ├── prisoner.go                 # Prisoner repository
│   │   │   ├── charge.go                   # Charge repository
│   │   │   ├── arrest.go                   # Arrest repository
│   │   │   ├── facility.go                 # Facility repository
│   │   │   ├── release.go                  # Release repository
│   │   │   └── mugshot.go                  # Mugshot repository
│   │   │
│   │   ├── service/                        # Business logic layer
│   │   │   ├── processor.go                # Main booking processor
│   │   │   ├── file_generator.go           # VINE file generation
│   │   │   ├── mugshot_handler.go          # Mugshot processing
│   │   │   └── validator.go                # Business validation
│   │   │
│   │   └── handler/                        # HTTP handlers
│   │       ├── http.go                     # HTTP routes & handlers
│   │       ├── dto.go                      # Request/Response DTOs
│   │       └── mapper.go                   # DTO ↔ Domain mapping
│   │
│   ├── scheduler/                          # Scheduler domain
│   │   ├── domain/
│   │   │   ├── job.go                      # Job entity
│   │   │   ├── config.go                   # Scheduler config entity
│   │   │   ├── execution.go                # Job execution entity
│   │   │   ├── repository.go               # Repository interfaces
│   │   │   └── errors.go
│   │   │
│   │   ├── repository/
│   │   │   ├── config.go                   # Scheduler config repo
│   │   │   └── execution.go                # Job execution history repo
│   │   │
│   │   ├── service/
│   │   │   ├── scheduler.go                # Scheduler service
│   │   │   ├── cron.go                     # Cron management
│   │   │   └── executor.go                 # Job executor
│   │   │
│   │   └── handler/
│   │       ├── http.go                     # Scheduler API endpoints
│   │       └── dto.go
│   │
│   ├── transfer/                           # File transfer domain
│   │   ├── domain/
│   │   │   ├── transfer.go                 # Transfer entity
│   │   │   ├── ftp_config.go               # FTP config entity
│   │   │   ├── repository.go
│   │   │   └── errors.go
│   │   │
│   │   ├── repository/
│   │   │   └── config.go                   # FTP config repository
│   │   │
│   │   ├── service/
│   │   │   ├── transfer.go                 # Main transfer service
│   │   │   ├── ftp.go                      # FTP client
│   │   │   ├── sftp.go                     # SFTP client
│   │   │   └── file_manager.go             # File operations
│   │   │
│   │   └── handler/
│   │       ├── http.go                     # Transfer API endpoints
│   │       └── dto.go
│   │
│   ├── sysconfig/                          # System configuration domain
│   │   ├── domain/
│   │   │   ├── config.go                   # System config entity
│   │   │   ├── lookup.go                   # System lookup entity
│   │   │   ├── repository.go
│   │   │   └── errors.go
│   │   │
│   │   ├── repository/
│   │   │   ├── config.go                   # System config repo
│   │   │   └── lookup.go                   # System lookup repo
│   │   │
│   │   ├── service/
│   │   │   └── config.go                   # Config management service
│   │   │
│   │   └── handler/
│   │       ├── http.go                     # Config API endpoints
│   │       └── dto.go
│   │
│   └── platform/                           # Shared platform code
│       │
│       ├── config/                         # Application configuration
│       │   ├── config.go                   # Main config struct
│       │   ├── loader.go                   # Config loading logic
│       │   └── validator.go                # Config validation
│       │
│       ├── database/                       # Database utilities
│       │   ├── sqlserver.go                # SQL Server connection
│       │   ├── transaction.go              # Transaction helper
│       │   ├── health.go                   # DB health check
│       │   └── types.go                    # Custom SQL types
│       │
│       ├── logger/                         # Structured logging
│       │   ├── logger.go                   # Logger interface & impl
│       │   ├── context.go                  # Context-aware logging
│       │   └── fields.go                   # Log field helpers
│       │
│       ├── server/                         # HTTP server utilities
│       │   ├── server.go                   # Server struct & methods
│       │   ├── middleware/                 # HTTP middleware
│       │   │   ├── logger.go               # Request logging
│       │   │   ├── recovery.go             # Panic recovery
│       │   │   ├── cors.go                 # CORS handling
│       │   │   ├── metrics.go              # Prometheus metrics
│       │   │   ├── request_id.go           # Request ID generation
│       │   │   └── timeout.go              # Request timeout
│       │   │
│       │   └── response/                   # HTTP response helpers
│       │       ├── json.go                 # JSON responses
│       │       ├── error.go                # Error responses
│       │       └── pagination.go           # Pagination helpers
│       │
│       ├── crypto/                         # Cryptography utilities
│       │   ├── encryption.go               # AES encryption
│       │   └── password.go                 # Password hashing
│       │
│       ├── metrics/                        # Observability
│       │   ├── metrics.go                  # Prometheus metrics
│       │   └── collector.go                # Custom collectors
│       │
│       └── validator/                      # Validation utilities
│           └── validator.go                # Request validation
│
├── pkg/                                    # Public reusable libraries
│   │
│   ├── vine/                               # VINE protocol library
│   │   ├── encoder.go                      # VINE file format encoder
│   │   ├── decoder.go                      # VINE file format decoder
│   │   ├── validator.go                    # VINE data validation
│   │   └── constants.go                    # VINE constants
│   │
│   ├── sqlutil/                            # SQL utilities
│   │   ├── null.go                         # Null type handling
│   │   ├── scanner.go                      # Custom scanners
│   │   └── time.go                         # Time handling
│   │
│   └── fileutil/                           # File utilities
│       ├── copy.go                         # File copy operations
│       ├── cleanup.go                      # File cleanup
│       └── path.go                         # Path utilities
│
├── migrations/                             # Database migrations
│   ├── 000001_initial_schema.up.sql
│   ├── 000001_initial_schema.down.sql
│   ├── 000002_scheduler_tables.up.sql
│   ├── 000002_scheduler_tables.down.sql
│   ├── 000003_job_history.up.sql
│   └── 000003_job_history.down.sql
│
├── api/                                    # API specifications
│   └── openapi/
│       ├── swagger.yaml                    # OpenAPI 3.0 spec
│       └── README.md
│
├── scripts/                                # Development & build scripts
│   ├── build.sh                            # Build script
│   ├── test.sh                             # Test script
│   ├── lint.sh                             # Linting script
│   ├── migrate-up.sh                       # Run migrations up
│   ├── migrate-down.sh                     # Run migrations down
│   └── dev.sh                              # Start dev environment
│
├── deployments/                            # Deployment configurations
│   │
│   ├── docker/
│   │   ├── Dockerfile                      # Production dockerfile
│   │   ├── Dockerfile.dev                  # Development dockerfile
│   │   └── docker-compose.yml              # Local dev stack
│   │
│   ├── kubernetes/                         # K8s manifests
│   │   ├── namespace.yaml
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   ├── secret.yaml
│   │   ├── ingress.yaml
│   │   └── hpa.yaml                        # Horizontal Pod Autoscaler
│   │
│   └── terraform/                          # Infrastructure as Code
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
│
├── test/                                   # Additional test files
│   │
│   ├── integration/                        # Integration tests
│   │   ├── booking_test.go
│   │   ├── scheduler_test.go
│   │   └── transfer_test.go
│   │
│   ├── e2e/                                # End-to-end tests
│   │   └── api_test.go
│   │
│   ├── fixtures/                           # Test data
│   │   ├── bookings.json
│   │   ├── prisoners.json
│   │   └── charges.json
│   │
│   └── testutil/                           # Test utilities
│       ├── database.go                     # Test DB setup
│       ├── mock.go                         # Mock helpers
│       └── assert.go                       # Custom assertions
│
├── docs/                                   # Documentation
│   ├── README.md                           # Main documentation
│   ├── architecture.md                     # Architecture decisions
│   ├── api.md                              # API documentation
│   ├── deployment.md                       # Deployment guide
│   ├── development.md                      # Development guide
│   └── migration-from-spring.md            # Migration guide
│
├── configs/                                # Configuration files
│   ├── config.yaml                         # Default config
│   ├── config.dev.yaml                     # Dev environment
│   ├── config.staging.yaml                 # Staging environment
│   └── config.prod.yaml                    # Production environment
│
├── .github/                                # GitHub specific files
│   ├── workflows/
│   │   ├── ci.yml                          # CI pipeline
│   │   ├── release.yml                     # Release pipeline
│   │   └── security.yml                    # Security scanning
│   │
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── ISSUE_TEMPLATE/
│       ├── bug_report.md
│       └── feature_request.md
│
├── tools/                                  # Go tools
│   └── tools.go                            # Tool dependencies
│
├── vendor/                                 # Vendored dependencies (optional)
│
├── go.mod                                  # Go module definition
├── go.sum                                  # Go module checksums
├── Makefile                                # Build automation
├── README.md                               # Project README
├── LICENSE                                 # License file
├── .gitignore                              # Git ignore rules
├── .golangci.yml                           # Golangci-lint config
├── .air.toml                               # Air hot reload config
├── .env.example                            # Example environment vars
└── .editorconfig                           # Editor configuration