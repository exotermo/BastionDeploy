# Stage 1 — build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .

# API_URL vem do docker-compose via ARG. Se vazio, usa path relativa.
ARG DASHBOARD_API_URL=""
ENV VITE_API_URL=$DASHBOARD_API_URL

RUN npm run build

# Stage 2 — serve com nginx
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
