FROM zaqwsx_udfs-test-img

RUN udfs init -b=1024
ADD . /tmp/test
RUN mv -f /tmp/test/config /root/.udfs/config
RUN udfs id
RUN chmod +x /tmp/test/run.sh

EXPOSE 4021 4022/udp

ENV UDFS_PROF true
ENV UDFS_LOGGING_FMT nocolor

ENTRYPOINT ["/bin/bash"]
CMD ["/tmp/test/run.sh"]
