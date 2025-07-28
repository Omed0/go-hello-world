#!/bin/bash

# GO Hello World - Build Script
# This script builds both the API server and CLI application

echo "🚀 Building GO Hello World Task Management System..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f bin/server bin/task-cli

# Build API Server
echo "🔧 Building API Server..."
go build -ldflags="-w -s" -o bin/server .
if [ $? -eq 0 ]; then
    echo "✅ API Server built successfully -> bin/server"
else
    echo "❌ Failed to build API Server"
    exit 1
fi

# Build CLI Application (Server-Connected)
echo "🎨 Building Server-Connected CLI Application..."
go build -ldflags="-w -s" -o bin/task-cli ./cmd/cli
if [ $? -eq 0 ]; then
    echo "✅ CLI Application built successfully -> bin/task-cli"
else
    echo "❌ Failed to build CLI Application"
    exit 1
fi

echo ""
echo "🎉 Build completed successfully!"
echo ""

# Starter section
echo "🚀 Starting applications..."
echo ""
echo "Choose an option:"
echo "1) Start API Server"
echo "2) Start CLI Application"
echo "3) Start both (server in background)"
echo "4) Exit"
echo ""
read -p "Enter your choice (1-4): " choice

case $choice in
    1)
        echo "🌐 Starting API Server..."
        ./bin/server
        ;;
    2)
        echo "💻 Starting CLI Application..."
        ./bin/task-cli
        ;;
    3)
        echo "🌐 Starting API Server in background..."
        ./bin/server &
        SERVER_PID=$!
        echo "Server started with PID: $SERVER_PID"
        echo "💻 Starting CLI Application..."
        ./bin/task-cli
        echo "Stopping server..."
        kill $SERVER_PID 2>/dev/null
        ;;
    4)
        echo "👋 Goodbye!"
        ;;
    *)
        echo "❌ Invalid choice"
        ;;
esac

echo ""
echo "📚 Usage:"
echo "  API Server: ./bin/server"
echo "  CLI App:    ./bin/task-cli"
echo ""
echo "📖 Documentation: See README.md for detailed usage instructions"
