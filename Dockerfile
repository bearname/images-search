FROM debian

ADD ./bin/photoservice /app/bin/photoservice
COPY ./dist ./app/bin/web
RUN chmod +x /app/bin/photoservice
WORKDIR /app

CMD ["/app/bin/photoservice"]