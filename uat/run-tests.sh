#!/bin/bash

# Run UAT tests against PostgreSQL database
# Focused script that only handles testing functionality

set -e

# Configuration
UAT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$UAT_DIR")"
BINARY_NAME="pg-goer"
OUTPUT_FILE="uat-test-output.md"
CONNECTION_STRING="postgresql://testuser:testpass@127.0.0.1:5555/testdb?sslmode=disable"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[UAT-RUN]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

main() {
    log "Running UAT tests against PostgreSQL database"
    
    # Verify binary exists
    cd "$PROJECT_ROOT"
    if [[ ! -f "$BINARY_NAME" ]]; then
        error "Binary $BINARY_NAME not found. Run 'make uat-up' first."
        return 1
    fi
    
    # Test basic CLI functionality
    log "Testing CLI help and version commands..."
    
    if ! ./"$BINARY_NAME" --help >/dev/null 2>&1; then
        error "Failed to run --help command"
        return 1
    fi
    
    if ! ./"$BINARY_NAME" --version >/dev/null 2>&1; then
        error "Failed to run --version command"
        return 1
    fi
    
    success "CLI commands working correctly"
    
    # Test database connectivity
    log "Testing database connectivity..."
    
    if ! ./"$BINARY_NAME" -o "$UAT_DIR/$OUTPUT_FILE" "$CONNECTION_STRING"; then
        error "pg-goer failed to generate documentation"
        error "Check database status: docker logs pg-goer-test-db"
        return 1
    fi
    
    success "Documentation generated successfully"
    
    # Validate the output
    cd "$UAT_DIR"
    log "Validating generated documentation..."
    
    if [[ ! -f "$OUTPUT_FILE" ]]; then
        error "Output file $OUTPUT_FILE not found"
        return 1
    fi
    
    # Check file size (should be substantial)
    file_size=$(wc -c < "$OUTPUT_FILE")
    if [[ $file_size -lt 1000 ]]; then
        error "Output file is too small ($file_size bytes), likely incomplete"
        return 1
    fi
    
    success "Output file exists and has reasonable size ($file_size bytes)"
    
    # Validate content
    log "Validating documentation content..."
    
    # Required sections
    required_sections=(
        "# PostgreSQL Database Documentation"
        "## Table of Contents"
        "## Database Summary"
        "## Database Relationships"
        "\`\`\`mermaid"
        "erDiagram"
        "## Tables"
        "## users"
        "## orders"
        "## order_items"
        "## categories"
        "## products"
    )
    
    missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -q "$section" "$OUTPUT_FILE"; then
            missing_sections+=("$section")
        fi
    done
    
    if [[ ${#missing_sections[@]} -gt 0 ]]; then
        error "Missing required sections in output:"
        for section in "${missing_sections[@]}"; do
            error "  - $section"
        done
        return 1
    fi
    
    success "All required sections found in documentation"
    
    # Validate relationships are documented
    log "Validating foreign key relationships..."
    
    expected_relationships=(
        "users ||--o{ orders"
        "orders ||--o{ order_items"
        "categories ||--o{ categories"
        "categories ||--o{ products"
    )
    
    missing_relationships=()
    for relationship in "${expected_relationships[@]}"; do
        if ! grep -q "$relationship" "$OUTPUT_FILE"; then
            missing_relationships+=("$relationship")
        fi
    done
    
    if [[ ${#missing_relationships[@]} -gt 0 ]]; then
        error "Missing expected relationships in Mermaid diagram:"
        for relationship in "${missing_relationships[@]}"; do
            error "  - $relationship"
        done
        return 1
    fi
    
    success "All expected relationships found in documentation"
    
    # Validate row counts are included
    log "Validating row counts..."
    
    if ! grep -q "Row Count:" "$OUTPUT_FILE"; then
        error "No row counts found in documentation"
        return 1
    fi
    
    # Check for specific table row counts (approximate)
    if ! grep -q "Row Count: 10" "$OUTPUT_FILE"; then
        warning "Expected ~10 rows in users table, check data insertion"
    fi
    
    success "Row counts found in documentation"
    
    # Test schema filtering
    log "Testing schema filtering functionality..."
    
    cd "$PROJECT_ROOT"
    if ! ./"$BINARY_NAME" -schemas public -o "$UAT_DIR/schema-filtered-output.md" "$CONNECTION_STRING"; then
        error "Failed to run with schema filtering"
        return 1
    fi
    
    if [[ ! -f "$UAT_DIR/schema-filtered-output.md" ]]; then
        error "Schema-filtered output file not created"
        return 1
    fi
    
    success "Schema filtering works correctly"
    
    # Validate the generated documentation shows proper structure
    log "Final validation of documentation structure..."
    
    cd "$UAT_DIR"
    line_count=$(wc -l < "$OUTPUT_FILE")
    if [[ $line_count -lt 50 ]]; then
        error "Documentation seems too short ($line_count lines)"
        return 1
    fi
    
    success "Documentation has proper length ($line_count lines)"
    
    # Display summary
    log "UAT Test Results Summary:"
    echo
    success "✓ CLI commands (--help, --version) working"
    success "✓ PostgreSQL connection successful"
    success "✓ Documentation generation successful"
    success "✓ All required sections present"
    success "✓ Foreign key relationships documented"
    success "✓ Row counts included"
    success "✓ Schema filtering functional"
    success "✓ Output file structure validated"
    echo
    success "🎉 All UAT tests passed! pg-goer is working correctly."
    
    # Show sample output
    log "Sample output (first 20 lines):"
    echo "----------------------------------------"
    head -20 "$OUTPUT_FILE"
    echo "----------------------------------------"
    
    log "Full documentation written to: $UAT_DIR/$OUTPUT_FILE"
}

main "$@"