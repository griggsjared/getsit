# if .env file exists, include it and export all variables to the environment.
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# default port for the server.
PORT ?= 8080

# default host for the server.
HOST ?= localhost

# default port for the templ proxy.
TEMPL_PROXY_PORT ?= 7331

# default host for the templ proxy.
TEMPL_PROXY_HOST ?= localhost

# run templ generation in watch mode to detect all .templ files and
# re-create _templ.txt files on change, then send reload event to browser.
# Default url: http://localhost:7331
dev/templ:
	templ generate \
	--watch \
	--proxy="http://${HOST}:${PORT}" \
	--proxybind="${TEMPL_PROXY_HOST}" --proxyport="${TEMPL_PROXY_PORT}" \
	--open-browser=true -v

# run tailwindcss to generate the main.css bundle in watch mode.
dev/tailwind:
	npx tailwindcss -i ./web/main.css -o ./web/public/main.css -c ./web/tailwind.config.js \
	--watch \
	--minify

# run air to detect any go file changes to re-build and re-run the server.
dev/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/web ./cmd/web/" \
	--build.bin "tmp/bin/web" \
	--build.delay "100" \
	--build.include_ext "go,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
dev/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "web/public" \
	--build.include_ext "js,css" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# start all dev/ tasks in parallel.
dev:
	make migrate/up && make -j6 dev/templ dev/server dev/sync_assets dev/tailwind

# run the base goose command to manage the database migrations.
migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DATABASE_DSN} goose -dir ./database/migrations $(ARGS)

# run the goose status command to check the current migration status.
migrate/status:
	ARGS="status" make migrate

# run the goose up command to apply all pending migrations.
migrate/up:
	ARGS="up" make migrate

migrate/up-one:
	ARGS="up-by-one" make migrate

# run the goose down command to rollback the last migration.
migrate/down:
	ARGS="down" make migrate

# run a series of commands to reset the database and re-apply all migrations.
migrate/fresh:
	ARGS="reset" make migrate && ARGS="up" make migrate/up
