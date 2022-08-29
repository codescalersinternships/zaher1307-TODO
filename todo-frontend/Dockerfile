FROM node:14-alpine
WORKDIR /app/frontend

COPY . .
RUN npm install

EXPOSE 8080
ENV HOST=0.0.0.0
CMD [ "npm", "start"]
