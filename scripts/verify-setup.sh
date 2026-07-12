#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  🏭 Software Factory - Setup Verification Script     ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Track overall status
ALL_GOOD=true

# Function to check command
check_command() {
    local cmd=$1
    local name=$2
    local required_version=$3
    
    if command -v "$cmd" >/dev/null 2>&1; then
        local version=$($cmd --version 2>&1 | head -1)
        echo -e "${GREEN}✅ $name found${NC} - $version"
        return 0
    else
        echo -e "${RED}❌ $name NOT found${NC}"
        ALL_GOOD=false
        return 1
    fi
}

# 1. Check dependencies
echo -e "${BLUE}📋 Checking Dependencies...${NC}"
echo ""

check_command "docker" "Docker" "20.10+"
check_command "docker-compose" "Docker Compose" "2.0+"
check_command "go" "Go" "1.22+"
check_command "node" "Node.js" "20+"
check_command "python3" "Python3" "3.11+"
check_command "git" "Git" "-"
check_command "tmux" "Tmux" "-"
check_command "nvim" "Neovim" "0.9+"

echo ""

# 2. Check GPU
echo -e "${BLUE}🎮 Checking GPU...${NC}"
if command -v nvidia-smi >/dev/null 2>&1; then
    GPU_INFO=$(nvidia-smi --query-gpu=name,memory.total --format=csv,noheader 2>/dev/null)
    if [ -n "$GPU_INFO" ]; then
        echo -e "${GREEN}✅ NVIDIA GPU detected:${NC}"
        echo "   $GPU_INFO"
    else
        echo -e "${YELLOW}⚠️  nvidia-smi found but no GPU info${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  No NVIDIA GPU detected (LM Studio will use CPU)${NC}"
fi
echo ""

# 3. Check LM Studio
echo -e "${BLUE}🤖 Checking LM Studio...${NC}"
LM_STUDIO_RESPONSE=$(curl -s http://localhost:1234/v1/models 2>&1)
if echo "$LM_STUDIO_RESPONSE" | grep -q '"data"'; then
    echo -e "${GREEN}✅ LM Studio API is running${NC}"
    
    # Check for required models
    echo ""
    echo -e "${BLUE}   Required Models:${NC}"
    
    if echo "$LM_STUDIO_RESPONSE" | grep -q "deepseek.*coder.*v2.*16b"; then
        echo -e "   ${GREEN}✅ deepseek-coder-v2-16b${NC}"
    else
        echo -e "   ${RED}❌ deepseek-coder-v2-16b NOT loaded${NC}"
        ALL_GOOD=false
    fi
    
    if echo "$LM_STUDIO_RESPONSE" | grep -q "qwen.*2.5.*coder.*7b"; then
        echo -e "   ${GREEN}✅ qwen2.5-coder-7b${NC}"
    else
        echo -e "   ${RED}❌ qwen2.5-coder-7b NOT loaded${NC}"
        ALL_GOOD=false
    fi
    
    if echo "$LM_STUDIO_RESPONSE" | grep -q "nomic.*embed"; then
        echo -e "   ${GREEN}✅ nomic-embed-text${NC}"
    else
        echo -e "   ${RED}❌ nomic-embed-text NOT loaded${NC}"
        ALL_GOOD=false
    fi
    
    echo ""
    echo -e "${BLUE}   Currently Loaded Models:${NC}"
    echo "$LM_STUDIO_RESPONSE" | grep -oP '"id":\s*"\K[^"]+' | sed 's/^/   - /'
    
else
    echo -e "${RED}❌ LM Studio NOT running or not responding${NC}"
    echo ""
    echo "   Please start LM Studio:"
    echo "   1. Open LM Studio"
    echo "   2. Go to: Developer → Local Server"
    echo "   3. Click 'Start Server'"
    echo "   4. Load required models"
    ALL_GOOD=false
fi
echo ""

# 4. Check Docker services
echo -e "${BLUE}🐳 Checking Docker Services...${NC}"

# Check if services are defined
if [ ! -f "docker-compose.yml" ]; then
    echo -e "${RED}❌ docker-compose.yml not found${NC}"
    ALL_GOOD=false
else
    # Check running containers
    RUNNING_CONTAINERS=$(docker ps --format "{{.Names}}" 2>/dev/null)
    
    check_docker_service() {
        local service=$1
        local port=$2
        
        if echo "$RUNNING_CONTAINERS" | grep -q "$service"; then
            echo -e "   ${GREEN}✅ $service${NC} (running on port $port)"
        else
            echo -e "   ${YELLOW}⚠️  $service${NC} (not running)"
        fi
    }
    
    check_docker_service "factory-postgres" "5432"
    check_docker_service "factory-redis" "6379"
    check_docker_service "factory-prometheus" "9090"
    check_docker_service "factory-grafana" "3000"
    check_docker_service "factory-jaeger" "16686"
fi

echo ""

# 5. Test connections
echo -e "${BLUE}🔌 Testing Service Connections...${NC}"

# Test PostgreSQL
if docker exec factory-postgres pg_isready -U factory >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ PostgreSQL${NC} - Ready"
else
    echo -e "   ${YELLOW}⚠️  PostgreSQL${NC} - Not ready (start with: docker-compose up -d postgres)"
fi

# Test Redis
if docker exec factory-redis redis-cli ping >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Redis${NC} - Ready"
else
    echo -e "   ${YELLOW}⚠️  Redis${NC} - Not ready (start with: docker-compose up -d redis)"
fi

# Test Prometheus
if curl -s http://localhost:9090/-/healthy >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Prometheus${NC} - Ready"
else
    echo -e "   ${YELLOW}⚠️  Prometheus${NC} - Not ready (start with: docker-compose up -d prometheus)"
fi

# Test Grafana
if curl -s http://localhost:3000/api/health >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Grafana${NC} - Ready"
else
    echo -e "   ${YELLOW}⚠️  Grafana${NC} - Not ready (start with: docker-compose up -d grafana)"
fi

# Test Jaeger
if curl -s http://localhost:16686/ >/dev/null 2>&1; then
    echo -e "   ${GREEN}✅ Jaeger${NC} - Ready"
else
    echo -e "   ${YELLOW}⚠️  Jaeger${NC} - Not ready (start with: docker-compose up -d jaeger)"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

if [ "$ALL_GOOD" = true ]; then
    echo -e "${GREEN}✅ All Phase 0 requirements met! Ready for Phase 1.${NC}"
    echo ""
    echo -e "${BLUE}🚀 Next Steps:${NC}"
    echo "   1. Review docs/implementation.md for Phase 1"
    echo "   2. Start building the Cockpit TUI"
    echo "   3. Implement Sandcastle bridge"
else
    echo -e "${YELLOW}⚠️  Some requirements not met. Please address the issues above.${NC}"
    echo ""
    echo -e "${BLUE}📝 To fix missing models:${NC}"
    echo "   1. Open LM Studio"
    echo "   2. Go to 'Discover' tab"
    echo "   3. Search for and download:"
    echo "      - deepseek-coder-v2-16b (Q4_K_M quantization)"
    echo "      - qwen2.5-coder-7b (Q4_K_M quantization)"
    echo "      - nomic-embed-text-v1.5"
    echo "   4. Go to 'Local Server' and load the models"
    echo ""
    echo -e "${BLUE}📝 To start Docker services:${NC}"
    echo "   cd $(pwd) && docker-compose up -d"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
