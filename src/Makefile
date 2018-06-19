# Do recursive builds
all:
	go get -u github.com/endurox-dev/endurox-go
	$(MAKE) -C ubftab
	$(MAKE) -C libbench

clean:
	$(MAKE) -C ubftab clean
	$(MAKE) -C libbench clean

.PHONY: clean
