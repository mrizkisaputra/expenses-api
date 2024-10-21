# ==============================================================================
# Golang commands

run:
	go run ./cmd/api/main.go

build:
	go build -o expense-api ./cmd/api/main.go

test:
	go test

# ==============================================================================
# Docker compose commands

local:
	echo "Starting local environment"
	docker-compose -f docker-compose-local.yaml up --build -d

develop:
	echo "Starting docker environment"
	docker compose -f docker-compose-dev.yaml up --build -d

# ==============================================================================
# Go migrate postgresql

migrate_up:
	migrate -database=postgres://postgres:password@localhost:5444/db_20102024?sslmode=disable -path=migrations up

migrate_down:
	migrate -database=postgres://postgres:postgres@localhost:5444/db_20102024?sslmode=disable -path=migrations down

# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES) 	#perintah Docker untuk menghentikan semua container yang id nya ada di dalam variable files
	docker rm $(FILES) 		#perintah ini akan menghapus semua container yang id nya ada di dalam variabel files

# membersihkan resource Docker yang tidak digunakan seperti container berhenti,
# image lama, volume, dan network yang tidak terpakai
clean:
	docker system prune -f