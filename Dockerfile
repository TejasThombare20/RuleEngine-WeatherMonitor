    # Build stage for Rule Engine Frontend
    FROM node:18-alpine AS rule-engine-frontend
    WORKDIR /rule_engine/client
    COPY rule_engine/client/package*.json ./
    RUN npm install
    COPY rule_engine/client ./
    RUN npm run build

    # Build stage for Weather Monitoring Frontend
    FROM node:18-alpine AS weather-frontend
    WORKDIR /weather_monitoring/client
    COPY weather_monitoring/client/package*.json ./
    RUN npm install
    COPY weather_monitoring/client ./
    RUN npm run build

    # Build stage for Rule Engine Backend
    # Updated to use Go 1.22 which is the latest stable version
    FROM golang:1.23.2-alpine AS rule-engine-backend
    WORKDIR /rule_engine/server
    COPY rule_engine/server/go.mod rule_engine/server/go.sum ./
    RUN go mod download
    COPY rule_engine/server ./
    RUN go build -o main

    # Build stage for Weather Monitoring Backend
    FROM golang:1.23.2-alpine AS weather-backend
    WORKDIR /weather_monitoring/server
    COPY weather_monitoring/server/go.mod weather_monitoring/server/go.sum ./
    RUN go mod download
    COPY weather_monitoring/server ./
    RUN go build -o main

    # Final stage
    FROM alpine:latest

    # Install Node.js and npm
    RUN apk add --update nodejs npm

    # Create directories matching your structure
    RUN mkdir -p /rule_engine/client /rule_engine/server \
        /weather_monitoring/client /weather_monitoring/server

    # Copy Rule Engine Frontend
    COPY --from=rule-engine-frontend /rule_engine/client/.next /rule_engine/client/.next
    COPY --from=rule-engine-frontend /rule_engine/client/node_modules /rule_engine/client/node_modules
    COPY --from=rule-engine-frontend /rule_engine/client/package.json /rule_engine/client/

    # Copy Weather Monitoring Frontend
    COPY --from=weather-frontend /weather_monitoring/client/.next /weather_monitoring/client/.next
    COPY --from=weather-frontend /weather_monitoring/client/node_modules /weather_monitoring/client/node_modules
    COPY --from=weather-frontend /weather_monitoring/client/package.json /weather_monitoring/client/

    # Copy backends
    COPY --from=rule-engine-backend /rule_engine/server/main /rule_engine/server/
    COPY --from=weather-backend /weather_monitoring/server/main /weather_monitoring/server/

    # Copy start script
    COPY start.sh ./
    RUN chmod +x start.sh

    EXPOSE 3001 3002 8000 9000

    CMD ["./start.sh"]