WORKING_DIRS=tmp
SRC=$(shell find . -name "*.go")
BIN=tmp/$(shell basename $(CURDIR))
FMT=tmp/fmt
TEST=tmp/cover
DOC=Document.txt

.PHONY: all fmt clean cover test

all: $(WORKING_DIRS) fmt $(BIN) test $(DOC)

clean:
	rm -rf $(WORKING_DIRS)

$(WORKING_DIRS):
	mkdir -p $(WORKING_DIRS)

fmt: $(SRC)
	go fmt ./...

$(BIN): $(SRC)
	go build -o $(BIN)

test: $(BIN)
	go test -v -tags=mock -cover -coverprofile=$(TEST) ./...

$(DOC): $(SRC)
	go doc -all . > $(DOC)

cover: $(TEST)
	grep "0$$" $(TEST) || true
