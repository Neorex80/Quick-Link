# QuickLink - URL Shortener 🔗

A modern, fast, and secure URL shortener built with Go, featuring a beautiful web interface and robust API.

![QuickLink Demo](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![Docker](https://img.shields.io/badge/Docker-Ready-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## ✨ Features

- 🎨 **Modern UI**: Beautiful gradient design with glass-morphism effects
- ⚡ **Fast & Secure**: Built with Go for optimal performance
- 🔒 **Input Validation**: Comprehensive URL validation and sanitization
- 🐳 **Docker Ready**: Easy deployment with Docker and Docker Compose
- 📱 **Responsive**: Works perfectly on desktop and mobile devices
- 🛡️ **Security**: Request size limits, collision detection, and error handling
- 🔄 **Thread-Safe**: Concurrent operations with mutex protection
- 📊 **Logging**: Detailed server logs for monitoring

## 🚀 Quick Start

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/Neorex80/quicklink-url-shortener.git
cd quicklink-url-shortener

# Run with Docker Compose
docker-compose up -d

# Or build and run manually
docker build -t quicklink .
docker run -p 8080:8080 quicklink
```

### Option 2: Go Direct

```bash
# Clone the repository
git clone https://github.com/Neorex80/quicklink-url-shortener.git
cd quicklink-url-shortener

# Install dependencies
go mod tidy

# Run the application
go run main.go
```

The application will be available at `http://localhost:8080`

## 🖥️ Web Interface

QuickLink features a modern, intuitive web interface:

- **Clean Design**: Professional gradient background with smooth animations
- **Easy to Use**: Simply paste your long URL and click "Shorten URL"
- **Instant Results**: Get your shortened URL immediately with success feedback
- **Mobile Friendly**: Responsive design that works on all devices

## 🔧 API Usage

### Shorten a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/very/long/url/path"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://www.example.com/very/long/url/path"
}
```

### Use a Short URL

Simply visit the short URL in your browser:
```
http://localhost:8080/abc123
```

You'll be automatically redirected to the original URL.

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

```yaml
version: '3.8'
services:
  quicklink:
    build: .
    ports:
      - "8080:8080"
    restart: unless-stopped
```

### Manual Docker Commands

```bash
# Build the image
docker build -t quicklink .

# Run the container
docker run -d -p 8080:8080 --name quicklink-app quicklink

# View logs
docker logs quicklink-app

# Stop the container
docker stop quicklink-app
```

## 🏗️ Architecture

### Core Components

- **HTTP Server**: Handles requests on port 8080
- **URL Storage**: Thread-safe in-memory storage
- **Short Code Generator**: Collision-resistant 6-character codes
- **Validation Layer**: Comprehensive input validation

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/` | Web interface |
| `POST` | `/shorten` | Shorten a URL |
| `GET` | `/{code}` | Redirect to original URL |

### Security Features

- ✅ Request body size limiting (1MB max)
- ✅ Content-type validation
- ✅ URL format validation (HTTP/HTTPS only)
- ✅ Short code format validation
- ✅ Collision detection for generated codes
- ✅ Structured error responses

## 📁 Project Structure

```
quicklink-url-shortener/
├── main.go              # Main application
├── go.mod              # Go module definition
├── Dockerfile          # Docker configuration
├── docker-compose.yml  # Docker Compose setup
├── .dockerignore       # Docker ignore rules
├── README.md           # This file
├── test_api.sh         # Linux/Mac test script
└── test_api.bat        # Windows test script
```

## 🧪 Testing

### Automated Testing

```bash
# Linux/Mac
./test_api.sh

# Windows
test_api.bat
```

### Manual Testing

```bash
# Valid URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Invalid URL (will return error)
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "not-a-valid-url"}'
```

## 🔮 Future Enhancements

- [ ] **Persistent Storage**: Database integration (PostgreSQL/MongoDB)
- [ ] **Custom Short Codes**: User-defined short URLs
- [ ] **Analytics**: Click tracking and statistics
- [ ] **URL Expiration**: Time-based URL expiry
- [ ] **Rate Limiting**: API rate limiting
- [ ] **Admin Dashboard**: Management interface
- [ ] **Bulk Operations**: Batch URL shortening
- [ ] **QR Code Generation**: QR codes for short URLs

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Neorex80**
- GitHub: [@Neorex80](https://github.com/Neorex80)

## 🙏 Acknowledgments

- Built with ❤️ using Go
- Inspired by modern web design principles
- Thanks to the Go community for excellent documentation

---

⭐ **Star this repository if you found it helpful!**
