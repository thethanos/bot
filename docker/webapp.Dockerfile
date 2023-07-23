FROM node:20.4-alpine3.17

COPY bot-webapp bot-webapp
COPY dev-full.crt ./bot-webapp/dev-full.crt
COPY dev-key.key ./bot-webapp/dev-key.key

ENV NODE_OPTIONS=--openssl-legacy-provider
RUN npm install -g serve http-server
RUN cd bot-webapp && npm install && npm run build

WORKDIR /bot-webapp/build
CMD ["http-server", "--proxy", "https://bot-dev-domain.com/index.html?", "--cors", "-S", "-C", "../dev-full.crt", "-K", "../dev-key.key", "-p", "443"]