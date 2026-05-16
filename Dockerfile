FROM alpine:3.19

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# create directory and set ownership to appuser:appgroup
RUN mkdir -p /home/appuser/solar-monitor \
 && chown -R appuser:appgroup /home/appuser/solar-monitor \
 && chmod 755 /home/appuser/solar-monitor

# Copy binary and set permissions
COPY bin/solar-monitor_linux-arm64 /usr/local/bin/solar-monitor
RUN chown appuser:appgroup /usr/local/bin/solar-monitor && chmod 755 /usr/local/bin/solar-monitor

USER appuser

EXPOSE 2112

ENTRYPOINT ["/usr/local/bin/solar-monitor"]
