FROM node:6-alpine

ADD ./ /root/{my-node-project}
WORKDIR /root/{my-node-project}
RUN npm install

ENTRYPOINT ["node", "app.js"]
