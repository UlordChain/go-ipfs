FROM zaqwsx_udfs-test-img

RUN udfs init -b=1024
ADD . /tmp/id
RUN mv -f /tmp/id/config /root/.udfs/config
RUN udfs id

EXPOSE 4031 4032/udp

ENV UDFS_PROF true
ENV UDFS_LOGGING_FMT nocolor

ENTRYPOINT ["/bin/bash"]
CMD ["/tmp/id/run.sh"]
