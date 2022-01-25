FROM alpine

ADD ./bin/photoservice /app/bin/photoservice
COPY ./dist ./app/bin/web
RUN chmod +x /app/bin/photoservice
WORKDIR /app
RUN apk --no-cache add ca-certificates
CMD ["/app/bin/photoservice", "-d", "false"]