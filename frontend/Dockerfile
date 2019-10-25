FROM node:12.13

WORKDIR /app

COPY package*.json ./

RUN npm install

EXPOSE 3000
EXPOSE 35729

CMD ["npm", "start"]