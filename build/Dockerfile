FROM node:lts-alpine as frontend
# Set the Current Working Directory inside the container
WORKDIR /build
# Copy package.json and package-lock.json files
COPY ./web/package.json ./web/package-lock.json ./
# Install dependencies from lock file
RUN npm ci
# Copy the source from the project web directory to the Working Directory inside the container
COPY ./web .
# Build frontend
RUN npm run build


FROM golang:alpine as backend
# Set the Current Working Directory inside the container
WORKDIR /build
# Copy go mod and sum files
COPY ./go.mod ./go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the project root directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o helmcabin ./cmd/helmcabin/main.go


FROM golang:alpine
# Add Maintainer Info
LABEL maintainer="Nick Triller"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy executable from backend stage
COPY --from=backend /build/helmcabin ./
# Copy frontend artifacts
COPY --from=frontend /build/dist ./web/dist/
# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
ENTRYPOINT ["./helmcabin"]
