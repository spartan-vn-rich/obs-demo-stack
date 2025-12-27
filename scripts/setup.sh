#!/bin/bash
# scripts/setup.sh

# Colors for pretty output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==> Switching to Minikube Docker Env...${NC}"
eval $(minikube docker-env)

build_image() {
    SERVICE=$1
    echo -e "${GREEN}Building $SERVICE...${NC}"
    docker build -t $SERVICE:latest ./src/$SERVICE
}

echo -e "${BLUE}==> Building Microservices...${NC}"
build_image "checkout-api"
build_image "inventory-service"
build_image "shipping-worker"

echo -e "${BLUE}==> Verifying Images in Minikube...${NC}"
docker images | grep "latest" | grep -E "checkout-api|inventory-service|shipping-worker"