BINARY=whb

BUILDFOLDER = build/bin

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1

# goxc flag
GOXCFLAG= -tasks-=validate -pv=${VERSION} -d ${BUILDFOLDER}

default: fmt
	go build -o ${BUILDFOLDER}/${BINARY} *.go
	@echo "Your binary is ready. Check "${BUILDFOLDER}/${VERSION}/${BINARY}

cross-all:
	goxc ${GOXCFLAG} 

fmt:
	@echo "Run gofmt"
	@echo "Run goimports"
	bash fmt.sh
