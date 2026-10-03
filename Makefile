GO ?= go
BIN := bin/guides

.PHONY: all sync build serve watch dev css check clean

all: build

$(BIN): $(shell find cmd internal web -type f) go.mod
	$(GO) build -o $(BIN) ./cmd/guides

sync: $(BIN)            ## copy the latest guides from the wiki into docs/
	$(BIN) sync

build: $(BIN)           ## render docs/ into dist/
	$(BIN) build

serve: $(BIN)           ## sync, build, serve on :8080, regenerate on every change
	$(BIN) serve -addr :8080

watch: $(BIN)           ## regenerate dist/ on every change, no web server
	$(BIN) watch

dev: $(BIN)             ## serve, reading templates/CSS from disk; run `npm run css:watch` alongside
	$(BIN) serve -addr :8080 -dev

css:                    ## rebuild web/static/app.css with Tailwind (needs Node)
	npm run css

check: build            ## verify every link and anchor in dist/
	python3 scripts/checklinks.py dist https://frontendlabs.xyz

clean:
	rm -rf dist .preview bin
