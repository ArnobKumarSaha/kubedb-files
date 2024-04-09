FROM ubuntu:20.04
ENV DEBIAN_FRONTEND noninteractive
WORKDIR /opt

COPY objectCount /opt/objectCount

RUN chmod a+x /opt/objectCount

CMD ["/opt/objectCount"]