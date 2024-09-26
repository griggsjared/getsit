ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# default port for the server.
PORT ?= 8080

# default port for the templ proxy.
TEMPL_PROXY_PORT ?= 7331

# run templ generation in watch mode to detect all .templ files and
# re-create _templ.txt files on change, then send reload event to browser.
# Default url: http://localhost:7331
dev/templ:
	templ generate \
	--watch \
	--proxy="http://localhost:${PORT}" \
	--proxybind="127.0.0.1" --proxyport="${TEMPL_PROXY_PORT}" \
	--open-browser=true -v

# run tailwindcss to generate the styles.css bundle in watch mode.
dev/tailwind:
	npx tailwindcss -i ./web/main.css -o ./web/public/main.css -c ./web/tailwind.config.js

# run air to detect any go file changes to re-build and re-run the server.
dev/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
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
	--build.include_ext "js,css"

# start all dev/ tasks in parallel.
dev:
	make -j5 dev/templ dev/tailwind dev/server dev/sync_assets
