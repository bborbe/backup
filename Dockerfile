FROM golang:1.27.4 AS build
COPY . /workspace
WORKDIR /workspace
# No explicit -mod flag on purpose. vendor/ is gitignored, so a bare clone (the
# cluster Job clones the repo and has none) would fail a -mod=vendor build
# outright; with no flag Go uses vendor/ when it is present and the module cache
# otherwise, which keeps the local arm on vendor and lets the cluster arm build.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s" -a -installsuffix cgo -o /main
CMD ["/bin/bash"]

FROM node:lts-alpine AS build-node
ARG NPM_REGISTRY=https://verdaccio.prod.nuke.benjamin-borbe.de
COPY frontend /frontend
WORKDIR /frontend
RUN npm set registry ${NPM_REGISTRY}
RUN npm install -g npm@11.8.0 --verbose
RUN	npm install --verbose --include=optional
RUN npm run lint:analyse
RUN npm run build

FROM alpine:3.24 AS alpine
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
