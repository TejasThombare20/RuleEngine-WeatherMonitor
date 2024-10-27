#!/bin/sh

# Start Rule Engine Frontend
cd /rule_engine/client
PORT=3001 npm start &

# Start Weather Monitoring Frontend
cd /weather_monitor/client
PORT=3002 npm start &

# Start Rule Engine Backend
cd /rule_engine/server
./main &

# Start Weather Monitoring Backend
cd /weather_monitor/server
./main &

# Keep container running
wait