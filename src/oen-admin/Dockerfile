FROM node:20-alpine AS builder

WORKDIR /build

# Copy package files
COPY package.json package-lock.json* ./
RUN npm ci

# Copy source code
COPY . .

# Build the frontend
RUN npm run build

# Final stage - Nginx
FROM nginx:alpine

# Copy built assets
COPY --from=builder /build/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
