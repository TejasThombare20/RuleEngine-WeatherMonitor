FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "start"]

# Note: Even though we EXPOSE 3000, this service will be accessible 
# on port 3001 on the host machine due to the port mapping in docker-compose.yml