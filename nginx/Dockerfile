FROM nginx

RUN rm /etc/nginx/conf.d/default.conf

COPY default.conf /etc/nginx/conf.d/default.conf

COPY index.html /usr/share/nginx/html

VOLUME /usr/share/nginx/html

VOLUME /etc/nginx