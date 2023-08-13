MAKE=make

.PHONY: run

run:
	$(MAKE) -C model/lib clean dynamic
	go run .
