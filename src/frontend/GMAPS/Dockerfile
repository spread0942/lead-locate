# 1. building phase
FROM node:22-slim as build

WORKDIR /app

COPY --chown=node:node package*.json ./
RUN npm install

COPY  --chown=node:node . .
RUN npm run build

# 2. production phase
FROM node:22-slim as production
WORKDIR /app

COPY  --chown=node:node --from=build /app/package*.json ./
RUN npm install --omit=dev

COPY  --chown=node:node --from=build /app/dist ./dist

RUN npm install -g serve

USER node
EXPOSE 80

CMD ["serve", "-s", "dist", "-l", "80"]