# ---- Base image ----
FROM node:20-slim

# ---- App directory ----
WORKDIR /app

# ---- Environment defaults ----
ENV NODE_ENV=production
ENV OTEL_TRACES_EXPORTER=otlp
ENV OTEL_EXPORTER_OTLP_PROTOCOL=grpc
ENV OTEL_EXPORTER_OTLP_ENDPOINT=grpc://tempo:4317
ENV OTEL_SERVICE_NAME=inventory-service

# ---- Install dependencies ----
COPY package*.json ./
RUN npm ci --omit=dev

# ---- Copy application code ----
COPY app.js otel.js ./

# ---- Non-root user (security best practice) ----
RUN useradd -m appuser \
  && chown -R appuser:appuser /app
USER appuser

# ---- Expose port ----
EXPOSE 8000

# ---- Start app ----
CMD ["node", "app.js"]