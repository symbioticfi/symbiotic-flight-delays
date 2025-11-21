FROM node:20-alpine

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN apk add --no-cache python3 make g++ \
    && corepack enable \
    && pnpm install --frozen-lockfile

COPY . .
RUN cp .env.docker .env && pnpm run build

EXPOSE 5173
CMD ["pnpm", "run", "preview", "--host", "0.0.0.0", "--port", "5173"]
