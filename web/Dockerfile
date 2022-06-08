ARG NODE_VERSION=14

FROM node:${NODE_VERSION} as base

ENV APP_PATH="/var/app" \
    VITE_API_BASE_URL=http://127.0.0.1:8080

COPY web/ ${APP_PATH}

WORKDIR ${APP_PATH}

RUN npm ci

FROM base as build

RUN npm run build

FROM nginx:alpine as production

COPY --from=build /var/app/dist /usr/share/nginx/html
