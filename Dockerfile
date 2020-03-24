FROM golang:1.14 as builder
RUN apt-get update && apt-get install -y mesa-opencl-icd ocl-icd-opencl-dev

WORKDIR /tmp
RUN curl https://sh.rustup.rs -sSf > rustup.sh
RUN chmod 755 rustup.sh
RUN ./rustup.sh -y
RUN rm /tmp/rustup.sh
ENV PATH="/root/.cargo/bin:${PATH}"

RUN mkdir /app 
WORKDIR /app 

RUN mkdir -p extern extern
WORKDIR /app/extern
RUN git clone https://github.com/filecoin-project/filecoin-ffi
WORKDIR /app/extern/filecoin-ffi
RUN git checkout f20cfbe28d99beda69e5416c6829927945116428
WORKDIR /app
COPY Makefile Makefile
RUN make .filecoin-build

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o local-devnet main.go  && \
go run github.com/GeertJohan/go.rice/rice append --exec local-devnet -i ./build

FROM ubuntu
RUN apt-get update && apt-get install -y mesa-opencl-icd ocl-icd-opencl-dev
COPY --from=builder /app/local-devnet /app/local-devnet
WORKDIR /app 
EXPOSE 7777
ENTRYPOINT ["./local-devnet"]
