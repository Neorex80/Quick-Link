# QuickLink 🔗

> A fast, secure URL shortener built with Go featuring a beautiful web interface and QR code generation.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![CI/CD](https://img.shields.io/github/actions/workflow/status/Neorex80/Quick-Link/ci.yml?style=for-the-badge&logo=github)](https://github.com/Neorex80/Quick-Link/actions)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://hub.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Stars](https://img.shields.io/github/stars/Neorex80/Quick-Link?style=for-the-badge&logo=github)](https://github.com/Neorex80/Quick-Link/stargazers)

<div align="center">
  <img src="https://github.com/user-attachments/assets/0433b494-56a2-4dda-b53b-c112f6450dd0" alt="QuickLink Demo" width="600"/>
</div>

## ✨ Features

- 🎨 **Modern UI** - Beautiful gradient design with glass-morphism effects
- ⚡ **Lightning Fast** - Built with Go for optimal performance
- 🔒 **Secure** - Input validation, rate limiting, and collision detection
- 📱 **Responsive** - Works perfectly on all devices
- 🔄 **QR Codes** - Auto-generated QR codes for easy sharing
- 🐳 **Docker Ready** - One-command deployment
- 🛡️ **Thread-Safe** - Concurrent operations with mutex protection

## 🚀 One-Click Deployment

### Deploy to Railway
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/quicklink?referralCode=bonus)

### Deploy to Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/Neorex80/Quick-Link)

### Deploy to Heroku
[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/Neorex80/Quick-Link)

### Deploy to Vercel
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/Neorex80/Quick-Link)

## 🏃‍♂️ Quick Start

### Docker (Recommended)
```bash
git clone https://github.com/Neorex80/Quick-Link.git
cd Quick-Link
docker-compose up -d
```

### Go Direct
```bash
git clone https://github.com/Neorex80/Quick-Link.git
cd Quick-Link
go mod tidy && go run main.go
```

Visit `http://localhost:8080` 🎉

## 🔧 API Usage

**Shorten URL:**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "custom_code": "my-link"}'
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/my-link",
  "original_url": "https://example.com"
}
```

## 🛠️ Tech Stack

- **Backend:** Go 1.21+
- **QR Codes:** go-qrcode
- **Frontend:** Vanilla JS + Modern CSS
- **Deployment:** Docker, Docker Compose

## 📊 Project Stats

![GitHub repo size](https://img.shields.io/github/repo-size/Neorex80/Quick-Link?style=flat-square)
![GitHub last commit](https://img.shields.io/github/last-commit/Neorex80/Quick-Link?style=flat-square)
![GitHub issues](https://img.shields.io/github/issues/Neorex80/Quick-Link?style=flat-square)
![GitHub pull requests](https://img.shields.io/github/issues-pr/Neorex80/Quick-Link?style=flat-square)

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Neorex80**
- GitHub: [@Neorex80](https://github.com/Neorex80)
- Project: [Quick-Link](https://github.com/Neorex80/Quick-Link)

---

<div align="center">
  <strong>⭐ Star this repository if you found it helpful!</strong>
</div>
