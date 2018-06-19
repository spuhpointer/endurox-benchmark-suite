# Do recursive builds
all:
	go get -u github.com/endurox-dev/endurox-go
	$(MAKE) -C ubftab
	$(MAKE) -C exbench

clean:
	$(MAKE) -C ubftab clean
	$(MAKE) -C exbench clean

.PHONY: clean
