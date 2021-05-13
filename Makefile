.PHONY: help run build db-migrate db-down db-fix

TEST_FLAGS ?=
CONFIG := $$PWD/config/api.conf
UPPER_DB_LOG=DEBUG

all:
	@echo "******************************"
	@echo "** Template API build tool  **"
	@echo "******************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run API in dev mode"
	@echo ""
	@echo "  build                 - build api into bin/ directory"
	@echo ""
	@echo "  db-create             - create dev db"
	@echo "  db-drop               - drop dev db"
	@echo "  db-reset              - reset dev db (drop, create, migrate)"
	@echo "  db-up                 - migrate dev DB to latest version"
	@echo "  db-down               - roll back dev DB to a previous version"
	@echo "  db-migrate            - create new db migration (NAME specifies migration name)"
	@echo "  db-status             - status of current dev DB version"
	@echo "  db-fix                - apply sequential ordering to migrations"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)


##
## Tools
##
tools:
	@mkdir -p ./bin
	@go get -u github.com/cosmtrek/air
	@go build -tags='no_mysql no_sqlite3' -o bin/migrate ./cmd/migrate
	@go mod tidy

##
# Database
##
db-status:
	./bin/migrate status

db-up:
	./bin/migrate up

db-down:
	./bin/migrate down

db-redo: db-down
	./bin/migrate up

db-migrate:
	@read -r -p "Enter migration op (ie. create, drop, add, remove, alter): " op;       \
	read -r -p "Enter migration type (ie. table, column): " type;                       \
	read -r -p "Enter migration target (ie. users): " target;                           \
	./bin/migrate create $${op}_$${type}_$${target} sql

db-reset:
	@./db/db.sh reset template
	@./bin/migrate up

db-fix:
	@./bin/migrate fix

db-create:
	@./db/db.sh create template

db-drop:
	@./db/db.sh drop template 

conf:
	[ -f config/api.conf ] || cp config/api.develop.conf config/api.conf

run: conf
	@(export CONFIG=${CONFIG}; air)

build:
	@mkdir -p ./bin
	GO111MODULE=on GOGC=off go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -o ./bin/api ./cmd/api/main.go
	GO111MODULE=on GOGC=off go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -o ./bin/migrate ./cmd/migrate/main.go
	
