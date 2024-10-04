#
# Copyright (c) 2022 Jiangxing Intelligence
# Copyright (c) 20223 Intel Corporation
# Copyright (c) 2024 IOTech Ltd
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

ARG BASE=golang:1.23-alpine3.20
FROM ${BASE} AS builder

ARG ALPINE_PKG_BASE="make git openssh-client"
ARG ALPINE_PKG_EXTRA=""
ARG ADD_BUILD_TAGS=""

# set the working directory
WORKDIR /device-uart

# Install our build time packages.
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .

# To run tests in the build container:
#   docker build --build-arg 'MAKE=build test' .
# This is handy of you do your Docker business on a Mac
ARG MAKE="make -e ADD_BUILD_TAGS=$ADD_BUILD_TAGS build"
RUN ${MAKE}

FROM alpine:3.20

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2021: Jiangxing Intelligence'

# dumb-init needed for injected secure bootstrapping entry point script when run in secure mode.
RUN apk add --update --no-cache dumb-init
# Ensure using latest versions of all installed packages to avoid any recent CVEs
RUN apk --no-cache upgrade

WORKDIR /
COPY --from=builder /device-uart/Attribution.txt /
COPY --from=builder /device-uart/LICENSE /
COPY --from=builder /device-uart/cmd/ /

EXPOSE 59995

ENTRYPOINT ["/device-uart"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry"]
