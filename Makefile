SUBMODULES=.update-modules

FFI_PATH:=./extern/filecoin-ffi/
FFI_DEPS:=libfilecoin.a filecoin.pc filecoin.h
FFI_DEPS:=$(addprefix $(FFI_PATH),$(FFI_DEPS))

$(FFI_DEPS): .filecoin-build ;

.filecoin-build: $(FFI_PATH)
	$(MAKE) -C $(FFI_PATH) $(FFI_DEPS:$(FFI_PATH)%=%)
	@touch $@

.update-modules:
	git submodule update --init --recursive
	@touch $@

build: .update-modules .filecoin-build
.PHONY: build

PARAMCACHE_PATH:=/var/tmp/fil-tools/filecoin-proof-parameters
test: build
	rm -f .filecoin-build
	rm -f .update-modules
	git submodule deinit --all -f
	mkdir -p $(PARAMCACHE_PATH)
	cat build/proof-params/parameters.json | jq 'keys[]' | xargs touch
	mv -n v20* $(PARAMCACHE_PATH)
	rm v20* || true
	PARAMCACHE_PATH=$(PARAMCACHE_PATH) go test -count 1 ./...
.PHONY: test