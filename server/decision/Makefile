COVERAGE_DIR := coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html
COVERAGE_MIN := 90
COVERAGE_TARGET := internal/service/decider.go

.PHONY: test coverage

test:
	go test ./...

coverage:
	mkdir -p $(COVERAGE_DIR)

	go test \
		-covermode=atomic \
		-coverpkg=./internal/service \
		-coverprofile=$(COVERAGE_FILE) \
		./test

	go tool cover -func=$(COVERAGE_FILE)

# 	awk -v file="$(COVERAGE_TARGET)" -v min="$(COVERAGE_MIN)" '\
# 		index($$1, file ":") { total += $$2; if ($$3 > 0) covered += $$2 } \
# 		END { \
# 			if (total == 0) { print "coverage target not found: " file; exit 1 } \
# 			pct = covered / total * 100; \
# 			printf "%s coverage: %.1f%% (minimum %.1f%%)\n", file, pct, min; \
# 			if (pct < min) exit 1; \
# 		}' $(COVERAGE_FILE)

	go tool cover \
		-html=$(COVERAGE_FILE) \
		-o $(COVERAGE_HTML)
