# build frontend
FROM node:18.20.6-alpine AS web_image
WORKDIR /build
COPY . /build
RUN npm install && npm run build

# build backend
FROM golang:1.24-alpine as server_image
WORKDIR /build
COPY ./service .
RUN apk add --no-cache bash curl gcc git musl-dev && \
    go env -w GO111MODULE=on && \
    export PATH=$PATH:/go/bin && \
    go build -o sun-panel --ldflags="-X sun-panel/global.RUNCODE=release" main.go

# final image
FROM alpine
WORKDIR /app
COPY --from=web_image /build/dist /app/web
COPY --from=server_image /build/sun-panel /app/sun-panel
EXPOSE 3002
RUN apk add --no-cache bash ca-certificates su-exec tzdata && \
    chmod +x ./sun-panel
CMD ./sun-panel