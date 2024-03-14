# ============= Compilation Stage ================
FROM golang:1.20.12-bullseye AS builder

ARG CAMINO_VERSION

RUN mkdir -p $GOPATH/src/github.com/chain4travel
WORKDIR $GOPATH/src/github.com/chain4travel

RUN git clone -b $CAMINO_VERSION --single-branch https://github.com/chain4travel/caminogo.git

# Copy caminoethvm repo into desired location
COPY . caminoethvm

# Set the workdir to CaminoGo and update caminoethvm dependency to local version
WORKDIR $GOPATH/src/github.com/chain4travel/caminogo
# Run go mod download here to improve caching of CaminoGo specific depednencies
RUN go mod download
# Replace the coreth dependency
RUN go mod edit -replace github.com/ava-labs/coreth=../caminoethvm
RUN go mod download && go mod tidy -compat=1.20

# Build the CaminoGo binary with local version of caminoethvm.
RUN ./scripts/build_camino.sh
# Create the plugins directory in the standard location so the build directory will be recognized
# as valid.
RUN mkdir build/plugins

# ============= Cleanup Stage ================
FROM debian:11-slim AS execution

# Maintain compatibility with previous images
RUN mkdir -p /caminogo/build
WORKDIR /caminogo/build

# Copy the executables into the container
COPY --from=builder /go/src/github.com/chain4travel/caminogo/build .

CMD [ "./caminogo" ]
