FROM node:20.4-alpine3.17

WORKDIR /multimessenger_bot
COPY bot-webapp ./
COPY dev-full.crt ./bot-webapp/dev-full.crt
COPY dev-key.key ./bot-webapp/dev-key.key

RUN npm install -g serve http-server
RUN cd bot-webapp && npm run build

CMD ["http-server build -S -C dev-full.crt -K dev-key.key -p 443"]