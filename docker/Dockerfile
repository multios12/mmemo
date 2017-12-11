FROM node:6-alpine

ADD ./ /root/hmemo
WORKDIR /root/hmemo
RUN npm install

ENTRYPOINT ["node", "./bin/www"]
