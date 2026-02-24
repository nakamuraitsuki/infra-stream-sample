# Build Stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# Runtime Stage
FROM nginx:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

# NOTE: nginx のレイヤで ENTRYPOINT と CMD が定義されているため
#       ここでは CMD を上書きするだけで十分
CMD ["nginx", "-g", "daemon off;"]