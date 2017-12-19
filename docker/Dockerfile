FROM node:8-alpine

ADD ./hmemo.tar.gz /root/hmemo
WORKDIR /root/hmemo
RUN npm install --production

ENTRYPOINT ["node", "./bin/www"]
