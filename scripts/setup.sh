#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🏭 Software Factory - Setup Script${NC}"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# 1. Check dependencies
echo -e "${BLUE}📋 Checking dependencies...${NC}"
echo ""

command_exists docker && print_status 0 "Docker found" || print_status 1 "Docker not found - please install"
command_exists docker-compose && print_status 0 "Docker Compose found" || print_status 1 "Docker Compose not found"
command_exists go && print_status 0 "Go found ($(go version))" || print_status 1 "Go not found - please install Go 1.22+"
command_exists node && print_status 0 "Node.js found ($(node --version))" || print_status 1 "Node.js not found - please install Node.js 20+"
command_exists python3 && print_status 0 "Python found ($(python3 --version))" || print_status 1 "Python not found"
command_exists git && print_status 0 "Git found" || print_status 1 "Git not found"
command_exists tmux && print_status 0 "Tmux found" || print_status 1 "Tmux not found - please install"
command_exists nvim && print_status 0 "Neovim found" || print_status 1 "Neovim not found - please install"

echo ""

# 2. Check LM Studio
echo -e "${BLUE}🤖 Checking LM Studio...${NC}"
if curl -s http://localhost:1234/v1/models >/dev/null 2>&1; then
    echo -e "${GREEN}✅ LM Studio is running${NC}"
else
    echo -e "${YELLOW}⚠️  LM Studio not detected${NC}"
    echo "Please start LM Studio and enable the local server:"
    echo "  1. Open LM Studio"
    echo "  2. Go to: Developer → Local Server"
    echo "  3. Click 'Start Server'"
    echo "  4. Load a model (recommended: deepseek-coder-v2-16b)"
    echo ""
    read -p "Press Enter when LM Studio is ready..."
fi
echo ""

# 3. Check GPU
echo -e "${BLUE}🎮 Checking GPU...${NC}"
if command_exists nvidia-smi; then
    nvidia-smi --query-gpu=name,memory.total --format=csv,noheader
    echo -e "${GREEN}✅ NVIDIA GPU detected${NC}"
else
    echo -e "${YELLOW}⚠️  No NVIDIA GPU detected. LM Studio will use CPU (slower).${NC}"
fi
echo ""

# 4. Create directories
echo -e "${BLUE}📁 Creating directories...${NC}"
mkdir -p data/{postgres,redis,chromadb,prometheus,grafana,jaeger,backups}
mkdir -p config/{grafana/dashboards,grafana/datasources}
mkdir -p logs
echo -e "${GREEN}✅ Directories created${NC}"
echo ""

# 5. Create Prometheus config
echo -e "${BLUE}⚙️  Creating configuration files...${NC}"
cat > config/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'factory-api'
    static_configs:
      - targets: ['host.docker.internal:8000']
  
  - job_name: 'factory-orchestrator'
    static_configs:
      - targets: ['host.docker.internal:8001']
EOF

# Grafana datasource
cat > config/grafana/datasources/prometheus.yml <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
    editable: true
EOF

echo -e "${GREEN}✅ Config files created${NC}"
echo ""

# 6. Start infrastructure services
echo -e "${BLUE}🐳 Starting infrastructure services...${NC}"
docker-compose up -d postgres redis prometheus grafana jaeger

echo ""
echo -e "${YELLOW}⏳ Waiting for services to be ready...${NC}"
sleep 5

# Check PostgreSQL
echo -n "Checking PostgreSQL... "
until docker exec factory-postgres pg_isready -U factory >/dev/null 2>&1; do
    sleep 2
done
echo -e "${GREEN}✓${NC}"

# Check Redis
echo -n "Checking Redis... "
until docker exec factory-redis redis-cli ping >/dev/null 2>&1; do
    sleep 2
done
echo -e "${GREEN}✓${NC}"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Infrastructure setup complete!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}🌐 Services:${NC}"
echo "  - PostgreSQL:  localhost:5432 (user: factory, db: factory)"
echo "  - Redis:       localhost:6379"
echo "  - Grafana:     http://localhost:3000 (admin/admin)"
echo "  - Jaeger:      http://localhost:16686"
echo "  - Prometheus:  http://localhost:9090"
echo ""
echo -e "${BLUE}🤖 LM Studio:${NC}"
echo "  - API:         http://localhost:1234/v1"
echo "  - Dashboard:   Check LM Studio app"
echo ""
echo -e "${BLUE}📝 Next steps:${NC}"
echo "  1. Install Node.js dependencies for Sandcastle:"
echo "     cd sandcastle-bridge && npm install"
echo "  2. Install Python dependencies:"
echo "     cd orchestrator && python -m venv venv && source venv/bin/activate && pip install -r requirements.txt"
echo "  3. Install Go dependencies:"
echo "     cd cockpit && go mod init && go mod tidy"
echo ""
echo -e "${BLUE}📚 Documentation:${NC}"
echo "  - Implementation Guide: docs/implementation.md"
echo "  - GitHub Issues: https://github.com/m00nk0d3/software-factory/issues"
echo ""
