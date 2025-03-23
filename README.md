# shawty

A modern, lightweight **URL shortening service** built with **Go, Redis, and
HTMX**. <br> Generate **concise**, **shareable links** in seconds with a
**clean, responsive UI**.

## Features

- **Fast URL shortening** using timestamp-based encoding
- **No duplicate URLs** - reuses existing short codes for the same URLs
- **Modern, responsive UI with seamless HTMX-powered interactions**
- **Redis backend** for efficient storage and retrieval
- **Docker support** for easy deployment
- **TTL support** for links (default: 24 hours)
- **Cloud-ready** - deployed on Render with Cloud Redis

## Live Demo

_Check out the live demo at:
[https://shawty-1845.onrender.com](https://shawty-1845.onrender.com)_

> [!NOTE]
> Since Shawty is a learning project, it is deployed on a free-tier service,
> which ironically means the shortened URL might end up longer than your
> original one.

## Quick Start

### Prerequisites

- Go 1.23.2 or higher
- Redis server (local or cloud instance)

### Environment Variables

Create a `.env` file with the following variables:

```
# For local development
REDIS_HOST=localhost:6379
REDIS_PASSWORD=your_redis_password

# For production with cloud Redis
# REDIS_HOST=your-redis-instance.cloud-provider.com:port
# REDIS_PASSWORD=your_redis_password

BASE_URL=http://localhost:8080
PORT=8080
```

### Installation

```bash
# Clone the repository
git clone https://github.com/ashish0kumar/shawty.git
cd shawty

# Install dependencies
go mod download

# Run the application
go run main.go
```

Visit `http://localhost:8080` in your browser to use the URL shortener.

### Docker Deployment

```bash
# Build the Docker image
docker build -t shawty .

# Run the container with Redis
docker run -p 8080:8080 --env-file .env shawty
```

## Deployment

### Render

This project is deployed on [Render](https://render.com) as a web service:

1. Connect your GitHub repository to Render
2. Create a new Web Service
3. Choose `Docker` as the environment, and Render will automatically handle the
   build and deployment.
4. Add the environment variables:
   - `REDIS_HOST`: Your cloud Redis endpoint (eg.,
     `redis-12345.c56.us-east-1-2.ec2.cloud.redislabs.com:12345`)
   - `REDIS_PASSWORD`: Your cloud Redis password
   - `BASE_URL`: Your Render app URL (e.g., `https://shawty-1845.onrender.com`)

_Ensure Redis Cloud credentials are set in environment variables._

## Architecture

1. **URL Submission**: When a user submits a URL, the application checks if it
   has already been shortened.
2. **Shortening Algorithm**: If not already shortened, a new short code is
   generated using base64-encoded timestamps.
3. **Storage**: Both mappings (short→long and long→short) are stored in Redis
   with a configurable TTL.
4. **Redirection**: When a shortened URL is accessed, the user can be redirected
   to the original URL.

## API Endpoints

- `GET /` - Serves the main UI
- `POST /shorten` - Shortens a URL and returns HTML with the result
- `GET /r/{code}` - Redirects to the original URL

## Technical Details

### Code Structure

```
└── ashish0kumar-shawty/
    ├── Dockerfile         # Docker configuration
    ├── go.mod             # Go module definition
    ├── go.sum             # Go module checksums
    ├── main.go            # Application entry point
    ├── templates/         # HTML templates
    │   └── index.html     # Main UI template
    └── utils/             # Utility functions
        ├── shorten.go     # URL shortening logic
        └── store.go       # Redis storage operations
```

### Dependencies

- [go-redis/redis](https://github.com/go-redis/redis) - Redis client for Go
- [joho/godotenv](https://github.com/joho/godotenv) - Loading environment
  variables
- [HTMX](https://htmx.org/) - Frontend interactivity without JavaScript
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework

## TODO

- [ ] Add URL content validation and checking for malicious URLs
- [ ] Implement rate limiting for production use
- [ ] Add analytics and click tracking functionality

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

[MIT](LICENSE)
