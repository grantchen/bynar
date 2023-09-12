# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START cloudrun_helloworld_dockerfile]
# [START run_helloworld_dockerfile]docker run -p 8080:8080

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19-buster as builder

# Copy local code to the container image.
COPY . /

WORKDIR /service/main
RUN go mod download
# # Build the binary.
RUN go build -v -o bynar_server
# # Copy the binary to the production image from the builder stage.
# COPY service/main/bynar_server /bynar_server
EXPOSE 8080
# # Run the web service on container startup.
CMD ["/service/main/bynar_server"]

# [END run_helloworld_dockerfile]
# [END cloudrun_helloworld_dockerfile]