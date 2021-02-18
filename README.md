[![Build Status](https://travis-ci.com/fr123k/golang-template.svg?branch=main)](https://travis-ci.com/fr123k/golang-template)

# golang-template

## Continuous Build

The travis build 

### Setup

```

```

## Targets

### Build

The following command will build the golang binary and run the unit tests.
The result of this build step is the standalone binary in the `./build/` folder. 

```
make build
```

Example Output:
```
  go build -o build/ cmd/main.go
  go test -v --cover ./...
  ?   	github.com/fr123k/golang-template/cmd	[no test files]
  === RUN   TestHelloWorld
  Hello World
  --- PASS: TestHelloWorld (0.00s)
  === RUN   TestHello
  --- PASS: TestHello (0.00s)
  PASS
  coverage: 100.0% of statements
  ok  	github.com/fr123k/golang-template/pkg/utility	(cached)	coverage: 100.0% of statements
```

### Run

The following make target will first build and then execute the golang binary.
```
  make run
```

Example Output:
```
  go build -o build/main cmd/main.go
  go test -v --cover ./...
  ?   	github.com/fr123k/golang-template/cmd	[no test files]
  === RUN   TestHelloWorld
  Hello World
  --- PASS: TestHelloWorld (0.00s)
  === RUN   TestHello
  --- PASS: TestHello (0.00s)
  PASS
  coverage: 100.0% of statements
  ok  	github.com/fr123k/golang-template/pkg/utility	(cached)	coverage: 100.0% of statements
  ./build/main
  Hello World
```

### Clean

The following make target will remove the `./build/` folder.
**No confirmation needed**
```
  make clean
```

Example Output:
```
  rm -rfv ./build
  ./build/main
  ./build
```

# Changelog

* setup travis build


# Todos

* setup travis build
