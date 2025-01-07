# Build Stage
FROM node:20 AS build-stage

# Set working directory
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy app source code
COPY . .

# Build the application
RUN npm run build

# Deployment Stage
FROM nginx:alpine AS deploy-stage

# Copy built assets to Nginx web root
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Copy custom Nginx configuration, if needed
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port and start Nginx
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]