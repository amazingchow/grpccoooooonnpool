SRC         := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGETS     := mock-proxy mock-server
ALL_TARGETS := $(TARGETS)

ifeq ($(race), 1)
	BUILD_FLAGS := -race
endif

ifeq ($(debug), 1)
	BUILD_FLAGS += -gcflags=all="-N -l"
endif

.PHONY: all
all: build

.PHONY: build
build: clean $(ALL_TARGETS)

$(TARGETS): $(SRC)
	@go build $(BUILD_FLAGS) $(PWD)/mocker/$@

.PHONY: clean
clean:
	@rm -f $(ALL_TARGETS)
