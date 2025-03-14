# Server configuration
PORT=8080
BASE_URL=http://localhost:8080
ENVIRONMENT=development

# Database configuration
DB_PATH=data/urls.db

# Security
JWT_SECRET=replace-this-with-a-long-secure-random-string-in-production

# CORS configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# Rate limiting
API_RATE_LIMIT=100
REDIRECT_RATE_LIMIT=1000
