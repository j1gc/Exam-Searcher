FROM oven/bun:latest AS builder
WORKDIR /app
COPY package.json package.json
# COPY bun.lockb bun.lockb
RUN bun install --production
COPY . .
ENV PUBLIC_BACKEND_URL=https://exambackend.jonas-floerke.de
RUN bun run build

FROM builder AS production
COPY --from=builder /app/build ./build


EXPOSE 3000
CMD ["bun", "./build"]