#!/bin/bash

# Wait for PostgreSQL to be ready for UAT testing
# Focused script that only handles database readiness

set -e

# Configuration
CONNECTION_STRING="postgresql://testuser:testpass@127.0.0.1:5555/testdb?sslmode=disable"
MAX_ATTEMPTS=30
RETRY_INTERVAL=2

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[UAT-WAIT]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

main() {
    log "Waiting for PostgreSQL to be ready..."
    
    for attempt in $(seq 1 $MAX_ATTEMPTS); do
        log "Connection attempt $attempt/$MAX_ATTEMPTS"
        
        # Check if Docker container is healthy
        if docker-compose exec -T postgres pg_isready -U testuser -d testdb >/dev/null 2>&1; then
            success "PostgreSQL container is healthy"
            
            # Additional check: try to connect with our tool's connection logic
            if timeout 5 docker-compose exec -T postgres psql -U testuser -d testdb -c "SELECT 1;" >/dev/null 2>&1; then
                success "PostgreSQL is ready for connections"
                success "Database startup completed in $((attempt * RETRY_INTERVAL)) seconds"
                return 0
            fi
        fi
        
        if [ $attempt -eq $MAX_ATTEMPTS ]; then
            error "PostgreSQL failed to start within $((MAX_ATTEMPTS * RETRY_INTERVAL)) seconds"
            error "Check Docker logs: docker logs pg-goer-test-db"
            return 1
        fi
        
        sleep $RETRY_INTERVAL
    done
}

main "$@"