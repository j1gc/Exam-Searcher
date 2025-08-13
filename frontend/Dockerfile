FROM oven/bun:latest AS builder
WORKDIR /app
COPY package.json package.json
# COPY bun.lockb bun.lockb
RUN bun install
COPY . .
RUN bun run build

FROM builder AS production
COPY --from=builder /app/build ./build


EXPOSE 3000
CMD ["bun", "./build"]