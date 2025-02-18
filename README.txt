Personal Portfolio Website
A modern, responsive portfolio website built with FastAPI, React, and PostgreSQL. This website showcases projects, skills, and professional experience with an intuitive and elegant design that works seamlessly across all devices.

🚀 Features

Interactive project showcase with detailed descriptions and live demos
Responsive design that adapts to any screen size or device
Server-side rendering for optimal performance
Database-driven content management
Modern, clean UI with smooth animations
Comprehensive API documentation
Secure authentication system
Docker containerization for easy deployment

🛠️ Technology Stack
Backend

Python 3.11+
FastAPI - Modern web framework
PostgreSQL - Database
SQLAlchemy - ORM
Alembic - Database migrations
Pydantic - Data validation
Uvicorn - ASGI server

Frontend

React 18
TypeScript
Tailwind CSS
React Router
Axios - HTTP client
React Query - Data fetching
Framer Motion - Animations

Development & Deployment

Docker & Docker Compose
Git
GitHub Actions (CI/CD)
Poetry (Python dependency management)
Node.js & npm

📋 Prerequisites
Before you begin, ensure you have the following installed:

Python 3.11 or higher
Node.js 18 or higher
PostgreSQL 14 or higher
Docker & Docker Compose (optional, but recommended)
Git

🔧 Local Development Setup

Clone the repository:
bashCopygit clone https://github.com/yourusername/portfolio-website.git
cd portfolio-website

Set up the backend:
bashCopycd backend
python -m venv venv
source venv/bin/activate  # On Windows: .\venv\Scripts\activate
pip install poetry
poetry install

Set up the database:
bashCopy# Create a PostgreSQL database
createdb portfolio_db

# Run migrations
alembic upgrade head

Set up the frontend:
bashCopycd ../frontend
npm install

Create a .env file in the backend directory:
CopyDATABASE_URL=postgresql://username:password@localhost/portfolio_db
SECRET_KEY=your-secret-key
ENVIRONMENT=development

Start the development servers:
bashCopy# Terminal 1 - Backend
cd backend
uvicorn app.main:app --reload

# Terminal 2 - Frontend
cd frontend
npm run dev

Access the application:

Frontend: http://localhost:5173
Backend API: http://localhost:8000
API Documentation: http://localhost:8000/docs



🐳 Docker Setup (Alternative)

Build and start the containers:
bashCopydocker-compose up --build

Access the application:

Frontend: http://localhost:3000
Backend API: http://localhost:8000



📁 Project Structure
project-website/
├── backend/
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py              # FastAPI application entry point
│   │   ├── models/              # SQLAlchemy models
│   │   │   ├── __init__.py
│   │   │   ├── project.py       # Project model
│   │   │   └── user.py          # User model
│   │   ├── routes/              # API endpoints
│   │   │   ├── __init__.py
│   │   │   ├── projects.py      # Project-related endpoints
│   │   │   └── users.py         # User-related endpoints
│   │   └── utils/               # Utility functions
│   ├── alembic/                 # Database migrations
│   ├── tests/                   # Backend tests
│   └── requirements.txt         # Python dependencies
├── frontend/
│   ├── src/
│   │   ├── components/          # Reusable React components
│   │   │   ├── Layout/
│   │   │   ├── ProjectCard/
│   │   │   └── Navigation/
│   │   ├── pages/              # Page components
│   │   │   ├── Home/
│   │   │   ├── Projects/
│   │   │   └── About/
│   │   ├── hooks/             # Custom React hooks
│   │   ├── utils/             # Frontend utilities
│   │   ├── types/            # TypeScript type definitions
│   │   └── App.tsx           # Root component
│   ├── public/               # Static assets
│   └── package.json         # Node.js dependencies
├── docker/                  # Docker configuration
│   ├── Dockerfile.backend
│   └── Dockerfile.frontend
├── docker-compose.yml      # Docker Compose configuration
└── README.md              # Project documentation


🚀 Deployment
Instructions for deploying to DigitalOcean or similar platforms will be added soon.
📝 License
This project is licensed under the MIT License - see the LICENSE file for details.
🤝 Contributing
Contributions, issues, and feature requests are welcome! Feel free to check the issues page.
👤 Author
Your Name

GitHub: @yourusername
Website: yourwebsite.com