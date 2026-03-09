FROM nginx:1.17.1

RUN rm -f /etc/nginx/nginx.conf \
    && rm -f /etc/nginx/conf.d/default.conf
COPY .docker/nginx.80.conf /etc/nginx/conf.d/nginx.80.conf
COPY .docker/nginx.conf /etc/nginx/nginx.conf
LABEL maintainer="shaohq@foxmail.com"

EXPOSE 80

RUN useradd -m -s /bin/bash www
RUN mkdir -p /var/cache/nginx/client_temp && chown -R www:www /var/cache/nginx && chmod 777 /var/run
USER www

COPY --chown=www ./.vuepress/dist /usr/share/nginx/html

ENTRYPOINT ["nginx", "-g", "daemon off;"]
