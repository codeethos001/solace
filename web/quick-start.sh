#!/bin/bash
# Quick Start Script for Solace Web Dashboard

set -e

echo "🚀 Solace Web Dashboard - Quick Start"
echo "======================================"
echo ""

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 16+ first."
    exit 1
fi

echo "✅ Node.js $(node --version) detected"
echo ""

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

echo "✅ npm $(npm --version) detected"
echo ""

# Navigate to web directory
cd "$(dirname "$0")"
echo "📁 Current directory: $(pwd)"
echo ""

# Install dependencies
echo "📦 Installing frontend dependencies..."
npm install
echo "✅ Dependencies installed"
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env.local ]; then
    echo "⚙️  Creating .env.local..."
    cat > .env.local << EOF
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=Solace Dashboard
EOF
    echo "✅ .env.local created"
    echo ""
fi

# Display next steps
echo "🎉 Setup complete!"
echo ""
echo "📋 Next steps:"
echo "  1. Ensure the Go backend is running:"
echo "     cd /home/ciph3r/Programs/Go/solace"
echo "     make build-linux"
echo "     ./bin/solace-linux --api --port 8080"
echo ""
echo "  2. Start the development server:"
echo "     npm run dev"
echo ""
echo "  3. Open your browser:"
echo "     http://localhost:3000"
echo ""
echo "  4. Login with default credentials:"
echo "     Username: admin"
echo "     Password: password"
echo ""
echo "📚 For more information, see README.md"
