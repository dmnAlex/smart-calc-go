CC=g++
CFLAGS=-c -fPIC
LDFLAGS=-shared -lstdc++ -lm

SRCFILES=calcmodel.cc calcmodelwrapper.cc
OBJFILES=$(SRCFILES:.cc=.o)
DYNAMICLIB=libcalcmodel.so

.PHONY: all dynamic clean

all: dynamic

dynamic: $(OBJFILES)
	$(CC) $(LDFLAGS) $^ -o $(DYNAMICLIB)

%.o: %.cc
	$(CC) $(CFLAGS) $^ -o $@

clean:
	rm -rf *.o $(DYNAMICLIB)
