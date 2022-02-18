FROM alpine

ADD ./bin/photoservice /app/bin/photoservice

COPY /home/runner/work/images-search/images-search/dist ./app/bin/web
RUN chmod +x /app/bin/photoservice
WORKDIR /app
RUN apk --no-cache add ca-certificates
CMD ["/app/bin/photoservice", "-d=true", "-log=false"]