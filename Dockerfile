FROM alpine:3

COPY ./do-ddns /do-ddns
COPY ./crontab.txt /crontab.txt
COPY ./entry.sh /entry.sh

RUN chmod 755 /do-ddns /entry.sh

RUN apk add libc6-compat

RUN /usr/bin/crontab /crontab.txt

CMD [ "/entry.sh" ]
