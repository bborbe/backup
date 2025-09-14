FROM golang:1.25.1 AS build
COPY . /workspace
WORKDIR /workspace
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -mod=vendor -ldflags "-s" -a -installsuffix cgo -o /main
CMD ["/bin/bash"]

FROM node:lts-alpine AS build-node
COPY frontend /frontend
WORKDIR /frontend
# RUN npm set registry https://registry.npmmirror.com
RUN npm set registry https://verdaccio.quant.benjamin-borbe.de
RUN npm install -g npm@11.6.0 --verbose
RUN	npm install --verbose
RUN npm run lint:analyse
RUN npm run build

FROM alpine:3.22 AS alpine
RUN apk --no-cache add \
	ca-certificates \
	rsync \
	openssh-client \
	tzdata \
	&& rm -rf /var/cache/apk/*
COPY --from=build-node /frontend/dist /frontend/dist
COPY --from=build /main /main
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
ENTRYPOINT ["/main"]
