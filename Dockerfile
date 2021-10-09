# 基础镜像
FROM golang:1.17.0

MAINTAINER cunoe

# 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 操作目录
WORKDIR /go/src/sst

# 复制源文件至操作目录
COPY . .

# 编译
RUN go build -installsuffix cgo .

# 暴露端口
EXPOSE 8080

CMD ["./Solar-Space-Server"]