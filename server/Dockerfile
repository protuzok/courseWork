FROM golang:1.25-alpine

WORKDIR /app

# Copy the workspace files
# Go workspaces allow us to build multiple modules together easily
COPY go.work go.work.sum ./

# Copy the server and shared modules
# We need shared because server depends on it
COPY server ./server
COPY shared ./shared

# Build the server application
# Since we are in a workspace, we can build from the module path
RUN go build -o server_app ./server/server.go

# Expose the port the server listens on
EXPOSE 1323

# Command to run the executable
CMD ["./server_app"]
