FROM node:16.9.0-alpine AS builder

WORKDIR /opt/web
# RUN npm install -g npm@6.14.4
COPY package.json package-lock.json ./
RUN npm install

ENV PATH="./node_modules/.bin:$PATH"

COPY . ./
RUN npm run build

# EXPOSE  10888

# CMD ["npm", "start"]

FROM nginx:1.16.0-alpine

COPY --from=builder /opt/web/dist /usr/share/nginx/html

COPY ./default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]