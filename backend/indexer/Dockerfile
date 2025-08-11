# Stage 1: Base build
FROM rust:latest AS base
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 2: Production image
FROM debian:bookworm-slim AS prod
WORKDIR /app

# Install any runtime dependencies your Rust binary needs
RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the build stage
COPY --from=base /app/target/release/indexer /app/indexer
COPY --from=base /app/db.sqlite3 /app/db.sqlite3
COPY --from=base /app/embeddings.postcard /app/embeddings.postcard

# Set binary as the container entrypoint
ENTRYPOINT ["./indexer"]