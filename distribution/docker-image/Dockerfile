# 阶段1：前端构建
FROM node:18.20.6-alpine AS frontend-builder
WORKDIR /build

# 安装依赖
COPY frontend .
RUN npm install
RUN npm run build

# 阶段2：后端构建
FROM golang:1.24-alpine AS backend-builder
ARG VERSION
WORKDIR /build
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache gcc musl-dev git
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off
COPY backend .
# COPY --from=frontend-builder /build/dist ./web
RUN go install
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-X sun-panel/global.RUNCODE=release -X sun-panel/global.VERSION=${VERSION}" -o sun-panel main.go

# 阶段3：最终镜像
FROM alpine:3.19
WORKDIR /app
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache tzdata && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone
COPY --from=backend-builder /build/sun-panel .
COPY --from=backend-builder /build/lang ./lang
COPY --from=backend-builder /build/conf.ini .
COPY --from=frontend-builder /backend/web ./web
EXPOSE 3002
CMD ["./sun-panel"]
