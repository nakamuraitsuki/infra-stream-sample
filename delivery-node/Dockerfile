FROM nginx:1.27-alpine

RUN apk add --no-cache iproute2
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY 30-setup-tc.sh /docker-entrypoint.d/30-setup-tc.sh
RUN chmod +x /docker-entrypoint.d/30-setup-tc.sh

EXPOSE 80