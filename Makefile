define build
@echo BUILD $1
@mkdir -p _bin
@go build -o _bin/$@ $1
endef

all: simplegame

simplegame:
	$(call build,"./cmd/simplegamed")
