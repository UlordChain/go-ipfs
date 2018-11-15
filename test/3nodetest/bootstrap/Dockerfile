FROM zaqwsx_udfs-test-img

RUN udfs init -b=1024
ADD . /tmp/id
RUN mv -f /tmp/id/config /root/.udfs/config
RUN udfs id

ENV UDFS_PROF true
ENV UDFS_LOGGING_FMT nocolor

EXPOSE 4011 4012/udp
