FROM golang:1.15.6-buster
WORKDIR /usr/local/go/src/alteroSmartTestTask
COPY . .
WORKDIR /usr/local/go/src/alteroSmartTestTask/backend/services/MS_Persistence/migrations
ENTRYPOINT ["go", "run", "run_migration.go"]
