.
├── bin
│   └── devpanel
├── cmd
│   ├── admin
│   │   └── main.go
│   ├── api
│   │   └── main.go
│   ├── devpanel
│   │   ├── main.go
│   │   └── main.go.bak
│   ├── migration
│   │   └── main.go
│   ├── server
│   │   └── main.go
│   └── worker
│       └── main.go
├── config
│   ├── app.yaml
│   ├── config.go
│   ├── development.yaml
│   ├── devpanel.yaml
│   ├── env.go
│   ├── production.yaml
│   ├── projects.yaml
│   └── testing.yaml
├── deployments
│   ├── docker
│   │   ├── api
│   │   │   └── Dockerfile
│   │   ├── docker-compose.dev.yml
│   │   ├── docker-compose.yml
│   │   └── worker
│   │       └── Dockerfile
│   ├── nginx
│   │   ├── api.conf
│   │   └── ssl
│   └── systemd
│       ├── api.service
│       └── worker.service
├── docs
│   ├── api
│   │   ├── endpoints
│   │   │   ├── messaging.md
│   │   │   └── urlshortener.md
│   │   ├── openapi.json
│   │   └── swagger.yaml
│   └── architecture.md 
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── bootstrap.go
│   │   ├── config
│   │   │   └── config.go
│   │   ├── context
│   │   │   └── context.go
│   │   └── server
│   │       ├── middleware
│   │       │   ├── auth.go
│   │       │   ├── cors.go
│   │       │   ├── logging.go
│   │       │   ├── recovery.go
│   │       │   └── spa.go
│   │       ├── router.go
│   │       └── server.go
│   ├── auth
│   │   └── auth.go
│   ├── common
│   │   ├── auth
│   │   │   ├── auth.go
│   │   │   ├── jwt.go
│   │   │   └── password.go
│   │   ├── cache
│   │   │   ├── cache.go
│   │   │   └── redis.go
│   │   ├── compression
│   │   │   └── compression.go
│   │   ├── database
│   │   │   ├── db.go
│   │   │   └── transaction.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── health
│   │   │   └── health.go
│   │   ├── logger
│   │   │   ├── http.go
│   │   │   └── logger.go
│   │   ├── metrics
│   │   │   ├── grafana.go
│   │   │   └── prometheus.go
│   │   ├── ratelimit
│   │   │   └── ratelimit.go
│   │   ├── storage
│   │   │   └── storage.go
│   │   ├── tracing
│   │   │   ├── jaeger.go
│   │   │   └── opentelemetry.go
│   │   ├── utils
│   │   │   ├── token.go
│   │   │   └── url_validator.go
│   │   └── validator
│   │       └── validator.go
│   ├── config
│   │   └── config.go
│   ├── core
│   │   ├── base_service.go
│   │   ├── config
│   │   │   └── config.go
│   │   ├── db
│   │   │   └── db.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   └── service_manager.go
│   ├── devpanel
│   │   ├── logs.go
│   │   ├── metrics.go
│   │   ├── project
│   │   │   └── project.go
│   │   ├── repository
│   │   │   └── repository.go
│   │   ├── server
│   │   │   └── server.go
│   │   └── service.go
│   ├── domain
│   │   ├── entity
│   │   │   ├── audit.go
│   │   │   └── user.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── messaging.go
│   │   ├── models.go
│   │   └── permission.go
│   ├── events
│   │   └── events.go
│   ├── gateway
│   │   └── gateway.go
│   ├── messaging
│   │   ├── attachments
│   │   │   └── service.go
│   │   ├── delivery
│   │   │   ├── http
│   │   │   │   ├── handlers.go
│   │   │   │   ├── middleware.go
│   │   │   │   ├── read_receipt_handler.go
│   │   │   │   └── routes.go
│   │   │   └── websocket
│   │   │       ├── client.go
│   │   │       ├── connection_manager.go
│   │   │       ├── hub.go
│   │   │       ├── presence.go
│   │   │       └── types.go
│   │   ├── domain
│   │   │   ├── channel.go
│   │   │   ├── message.go
│   │   │   ├── reaction.go
│   │   │   ├── read_receipts.go
│   │   │   └── repository.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── events
│   │   │   ├── events.go
│   │   │   ├── event_types.go
│   │   │   └── handlers.go
│   │   ├── interfaces
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   ├── repository
│   │   │   ├── cache
│   │   │   │   ├── channel.go
│   │   │   │   ├── common.go
│   │   │   │   └── message.go
│   │   │   ├── errors.go
│   │   │   ├── factory.go
│   │   │   ├── gorm_repository.go
│   │   │   ├── init.go
│   │   │   ├── mock
│   │   │   │   └── repository.go
│   │   │   ├── postgres
│   │   │   │   ├── attachment_repository.go
│   │   │   │   ├── channel_repository.go
│   │   │   │   ├── embed_repository.go
│   │   │   │   ├── message_repository.go
│   │   │   │   ├── reaction_repository.go
│   │   │   │   └── read_receipt_repository.go
│   │   │   └── repository.go
│   │   ├── service
│   │   │   ├── attachment_service.go
│   │   │   ├── channel_service.go
│   │   │   ├── interface.go
│   │   │   ├── messaging_service.go
│   │   │   ├── reaction_service.go
│   │   │   ├── read_receipt_service.go
│   │   │   └── service.go
│   │   ├── service.go
│   │   ├── usecase
│   │   │   ├── delete_message.go
│   │   │   ├── mark_as_read.go
│   │   │   ├── pin_message.go
│   │   │   ├── search_messages.go
│   │   │   ├── send_message.go
│   │   │   └── upload_attachment.go
│   │   └── websocket
│   │       └── websocket.go
│   ├── urlshortener
│   │   ├── delivery
│   │   │   └── http
│   │   │       ├── handlers.go
│   │   │       ├── middleware.go
│   │   │       └── routes.go
│   │   ├── domain
│   │   │   ├── repository.go
│   │   │   ├── stats.go
│   │   │   └── url.go
│   │   ├── repository
│   │   │   ├── cache
│   │   │   │   └── url.go
│   │   │   ├── gorm_repository.go
│   │   │   ├── mock
│   │   │   │   └── repository.go
│   │   │   ├── postgres
│   │   │   │   ├── stats.go
│   │   │   │   └── url.go
│   │   │   └── repository.go
│   │   ├── service
│   │   │   ├── service.go
│   │   │   ├── service_imp.go
│   │   │   ├── stats.go
│   │   │   └── url.go
│   │   ├── service.go
│   │   └── usecase
│   │       ├── resolve_url.go
│   │       ├── shorten_url.go
│   │       └── track_click.go
│   └── worker
│       ├── queue
│       │   ├── kafka.go
│       │   └── rabbitmq.go
│       └── tasks
│           ├── messaging_tasks.go
│           └── scheduled.go
├── middleware
│   └── auth
│       └── auth.go
├── migrations
│   ├── common
│   │   ├── 000001_create_users_table.down.sql
│   │   └── 000001_create_users_table.up.sql
│   ├── messaging
│   │   ├── 000001_create_channels_table.down.sql
│   │   └── 000001_create_channels_table.up.sql
│   └── urlshortener
│       ├── 000001_create_urls_table.down.sql
│       └── 000001_create_urls_table.up.sql
├── pkg
│   ├── httputil
│   ├── pagination
│   └── validator
├── PROJECT_STRUCTURE.md
└── scripts
    ├── backup_db.sh
    ├── lint.sh
    ├── restore_db.sh
    ├── run.sh
    ├── seed.sh
    └── setup.sh

92 directories, 176 files
