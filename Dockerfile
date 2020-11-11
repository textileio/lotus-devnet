FROM golang:1.14 as builder
RUN apt-get update && apt-get install -y mesa-opencl-icd clang ocl-icd-opencl-dev jq hwloc libhwloc-dev

WORKDIR /tmp
RUN curl https://sh.rustup.rs -sSf > rustup.sh
RUN chmod 755 rustup.sh
RUN ./rustup.sh -y
RUN rm /tmp/rustup.sh
ENV PATH="/root/.cargo/bin:${PATH}"

RUN mkdir /app 
WORKDIR /app 

COPY . .
RUN mkdir -p extern/filecoin-ffi
RUN make clean build
RUN GOOS=linux go build -o local-devnet main.go  && \
go run github.com/GeertJohan/go.rice/rice append --exec local-devnet -i ./build

FROM golang:1.14
RUN apt-get update && apt-get install -y mesa-opencl-icd ocl-icd-opencl-dev clang hwloc libhwloc-dev
COPY --from=builder /app/local-devnet /app/local-devnet
WORKDIR /app 
EXPOSE 7777
ENTRYPOINT ["./local-devnet"]
