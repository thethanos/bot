FROM node:20.4-alpine3.17

COPY webapp webapp
COPY dev-full.crt ./webapp/dev-full.crt
COPY dev-key.key ./webapp/dev-key.key

ENV NODE_OPTIONS=--openssl-legacy-provider
RUN npm install -g serve http-server
RUN cd webapp && npm install && npm run build

WORKDIR /webapp/build
CMD ["http-server", "--proxy", "https://bot-dev-domain.com/index.html?", "--cors", "-S", "-C", "../dev-full.crt", "-K", "../dev-key.key", "-p", "443"]