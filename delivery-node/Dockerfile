FROM nginx:1.27-alpine

RUN apk add --no-cache iproute2

# キャッシュディレクトリの作成と権限設定を追加
RUN mkdir -p /var/cache/nginx/video_cache && \
    chown -R nginx:nginx /var/cache/nginx/video_cache

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY 30-setup-tc.sh /docker-entrypoint.d/30-setup-tc.sh
RUN chmod +x /docker-entrypoint.d/30-setup-tc.sh

EXPOSE 80