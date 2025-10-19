# start all 5 watch processes in parallel.
.PHONY: live live/templ live/server live/tailwind live/sync_assets
live: 
	make -j4 live/templ live/server live/tailwind live/sync_assets
	#
# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
live/templ:
	templ generate --watch --proxy="http://localhost:9999" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
live/server:
	COUCHDB_URL=http://admin:password@localhost:5984/harmony \
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main ./cmd/server && templ generate --notify-proxy" --build.bin "tmp/bin/main" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go,tmpl" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
live/tailwind:
	pnpm run watch

# run esbuild to generate the index.js bundle in watch mode.
#live/esbuild:
	#npx --yes esbuild js/index.ts --bundle --outdir=assets/ --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

live/kill:
	killall templ || echo "Cannot kill templ"
	killall main || echo "Cannot kill main"
	killall air || echo "Cannot kill air"

test:
	gow -s -P "Start suite" -S "Suite done" test -vet=off ./...

.PHONY: codegen
codegen:
	rm -rf internal/testing/mocks
	go tool mockery

.PHONY: workspace/make
workspace/make:
	go work init
	go work use .
	go work use ../browser
	go work use ../shaman
	go work use ../surgeon

.PHONY: workspace/clean
workspace/clean:
	rm -f go.work go.work.sum
