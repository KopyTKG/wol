# Compiler variable
CC = g++
CFLAGS = -Wall -g -std=c++23
LDFLAGS = 

# Source files (all .cpp files in the directory)
SRCS = $(wildcard *.cpp */*.cpp)

# Object files (change .cpp to .o)
OBJS = $(SRCS:.cpp=.o)

# Executable name
EXEC = wol
BIN = ${EXEC}.bin

# All target
all: build install

# Compile source
build: deepclean memclean $(BIN) memclean

# Link the object files to create the executable
$(BIN): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $(OBJS) $(LDFLAGS)

# Compile each .cpp file to .o
%.o: %.cpp
	$(CC) $(CFLAGS) -c $< -o $@

# Clean up object files and executable
memclean:
	rm -rf */*.o *.o

deepclean: memclean
	rm -rf $(BIN)

install: 
	cp ./${BIN} /usr/bin/${EXEC}

# Phony targets
.PHONY: all build memclean deepclean install

