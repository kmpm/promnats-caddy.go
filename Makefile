
GOOS=$(shell go env GOOS)



ifeq ($(GOOS),windows) 
	BINEXT = .exe
else
	BINEXT =
endif


CADDYBIN=caddy$(BINEXT)

.PHONY: build
build: $(CADDYBIN)
	

.PHONY: run
run: 
	xcaddy run --watch --config testdata/Caddyfile


.PHONY: cleancaddy
clean: 
	del $(CADDYBIN)
	

$(CADDYBIN):
	xcaddy build --with github.com/kmpm/promnats.go/plugin/caddy

