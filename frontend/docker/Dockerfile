ARG NODE_VERSION=25-alpine3.22
ARG NGINX_VERSION=alpine3.22

# =====================
# Stage 1: Build
# =====================
FROM node:${NODE_VERSION} AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# =====================
# Stage 2: Run Nginx
# =====================
FROM nginx:${NGINX_VERSION}

# copy custom nginx config
COPY docker/nginx.conf /etc/nginx/nginx.conf

# copy the static build output from the build stage to Nginx's default HTML serving directory
COPY --chown=nginx:nginx --from=builder /app/dist/*/browser /usr/share/nginx/html

# use a built-in non-root user
USER nginx

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]