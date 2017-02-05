FROM progrium/busybox

RUN opkg-install ca-certificates
RUN for cert in `ls -1 /etc/ssl/certs/*.crt | grep -v /etc/ssl/certs/ca-certificates.crt`; do cat "$cert" >> /etc/ssl/certs/ca-certificates.crt; done

WORKDIR /usr/bin
ADD ./build/notification-service /usr/bin/

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

ENTRYPOINT ["notification-service"]