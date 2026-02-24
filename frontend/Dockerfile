# Build Stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Runtime Stage
FROM nginx:1.27-alpine

RUN apk add --no-cache iproute2
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
COPY 30-setup-tc.sh /docker-entrypoint.d/30-setup-tc.sh
RUN chmod +x /docker-entrypoint.d/30-setup-tc.sh

EXPOSE 80
