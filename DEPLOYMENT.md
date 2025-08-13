# 🚀 Deployment Guide

This guide covers various deployment options for QuickLink URL Shortener.

## 📋 Prerequisites

- Git
- Docker & Docker Compose (for containerized deployment)
- Go 1.21+ (for direct deployment)

## 🐳 Docker Deployment (Recommended)

### Quick Start with Docker Compose

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Neorex80/quicklink-url-shortener.git
   cd quicklink-url-shortener
   ```

2. **Start the application:**
   ```bash
   docker-compose up -d
   ```

3. **Access the application:**
   - Open your browser and go to `http://localhost:8080`

4. **View logs:**
   ```bash
   docker-compose logs -f quicklink
   ```

5. **Stop the application:**
   ```bash
   docker-compose down
   ```

### Manual Docker Deployment

1. **Build the Docker image:**
   ```bash
   docker build -t quicklink .
   ```

2. **Run the container:**
   ```bash
   docker run -d \
     --name quicklink-app \
     -p 8080:8080 \
     --restart unless-stopped \
     quicklink
   ```

3. **Manage the container:**
   ```bash
   # View logs
   docker logs quicklink-app
   
   # Stop container
   docker stop quicklink-app
   
   # Remove container
   docker rm quicklink-app
   ```

## ☁️ Cloud Deployment

### Deploy to Heroku

1. **Install Heroku CLI and login:**
   ```bash
   heroku login
   ```

2. **Create a new Heroku app:**
   ```bash
   heroku create your-app-name
   ```

3. **Set the buildpack:**
   ```bash
   heroku buildpacks:set heroku/go
   ```

4. **Deploy:**
   ```bash
   git push heroku main
   ```

### Deploy to Railway

1. **Connect your GitHub repository to Railway**
2. **Railway will automatically detect the Go application**
3. **Set environment variables if needed**
4. **Deploy with one click**

### Deploy to Render

1. **Connect your GitHub repository to Render**
2. **Create a new Web Service**
3. **Use the following settings:**
   - Build Command: `go build -o main .`
   - Start Command: `./main`
   - Port: `8080`

### Deploy to Google Cloud Run

1. **Build and push to Google Container Registry:**
   ```bash
   gcloud builds submit --tag gcr.io/PROJECT-ID/quicklink
   ```

2. **Deploy to Cloud Run:**
   ```bash
   gcloud run deploy --image gcr.io/PROJECT-ID/quicklink --platform managed
   ```

### Deploy to AWS ECS

1. **Build and push to ECR:**
   ```bash
   aws ecr get-login-password --region region | docker login --username AWS --password-stdin aws_account_id.dkr.ecr.region.amazonaws.com
   docker build -t quicklink .
   docker tag quicklink:latest aws_account_id.dkr.ecr.region.amazonaws.com/quicklink:latest
   docker push aws_account_id.dkr.ecr.region.amazonaws.com/quicklink:latest
   ```

2. **Create ECS task definition and service**

## 🖥️ VPS/Server Deployment

### Using Docker on VPS

1. **Install Docker on your VPS:**
   ```bash
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   ```

2. **Clone and deploy:**
   ```bash
   git clone https://github.com/Neorex80/quicklink-url-shortener.git
   cd quicklink-url-shortener
   docker-compose up -d
   ```

3. **Set up reverse proxy (Nginx):**
   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

### Direct Go Deployment

1. **Install Go on your server:**
   ```bash
   wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

2. **Clone and build:**
   ```bash
   git clone https://github.com/Neorex80/quicklink-url-shortener.git
   cd quicklink-url-shortener
   go mod tidy
   go build -o quicklink main.go
   ```

3. **Create systemd service:**
   ```ini
   [Unit]
   Description=QuickLink URL Shortener
   After=network.target
   
   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/path/to/quicklink-url-shortener
   ExecStart=/path/to/quicklink-url-shortener/quicklink
   Restart=always
   RestartSec=5
   
   [Install]
   WantedBy=multi-user.target
   ```

4. **Enable and start service:**
   ```bash
   sudo systemctl enable quicklink
   sudo systemctl start quicklink
   ```

## 🔧 Environment Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `HOST` | Server host | `0.0.0.0` |

### Docker Environment

```yaml
version: '3.8'
services:
  quicklink:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
    restart: unless-stopped
```

## 🔒 Security Considerations

1. **Use HTTPS in production**
2. **Set up proper firewall rules**
3. **Regular security updates**
4. **Monitor application logs**
5. **Use environment variables for sensitive data**

## 📊 Monitoring

### Health Check Endpoint

The application provides a health check at the root endpoint (`/`).

### Docker Health Check

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1
```

### Monitoring with Docker Compose

```yaml
version: '3.8'
services:
  quicklink:
    build: .
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## 🔄 Updates and Maintenance

### Updating the Application

1. **Pull latest changes:**
   ```bash
   git pull origin main
   ```

2. **Rebuild and restart:**
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

### Backup Considerations

Since the application uses in-memory storage, consider:
- Implementing persistent storage for production
- Regular data backups if using database
- Container volume backups

## 🆘 Troubleshooting

### Common Issues

1. **Port already in use:**
   ```bash
   # Find process using port 8080
   lsof -i :8080
   # Kill the process
   kill -9 <PID>
   ```

2. **Docker build fails:**
   ```bash
   # Clean Docker cache
   docker system prune -a
   ```

3. **Application not accessible:**
   - Check firewall settings
   - Verify port mapping
   - Check application logs

### Logs

```bash
# Docker logs
docker logs quicklink-app

# Docker Compose logs
docker-compose logs -f quicklink

# System logs (if using systemd)
journalctl -u quicklink -f
```

## 📞 Support

If you encounter any issues:

1. Check the [GitHub Issues](https://github.com/Neorex80/quicklink-url-shortener/issues)
2. Create a new issue with detailed information
3. Contact [@Neorex80](https://github.com/Neorex80)

---

Happy deploying! 🚀
