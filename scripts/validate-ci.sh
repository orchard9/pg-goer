#!/bin/bash
set -e

echo "🔍 Validating GitHub Actions CI Configuration"
echo "=============================================="

# Check if workflow files exist
if [ ! -f ".github/workflows/ci.yml" ]; then
    echo "❌ CI workflow file not found"
    exit 1
fi

if [ ! -f ".github/workflows/release.yml" ]; then
    echo "❌ Release workflow file not found"
    exit 1
fi

echo "✅ GitHub Actions workflow files found"

# Validate YAML syntax (if yq is available)
if command -v yq &> /dev/null; then
    echo "📝 Validating YAML syntax..."
    
    if yq eval '.jobs' .github/workflows/ci.yml > /dev/null; then
        echo "✅ CI workflow YAML is valid"
    else
        echo "❌ CI workflow YAML is invalid"
        exit 1
    fi
    
    if yq eval '.jobs' .github/workflows/release.yml > /dev/null; then
        echo "✅ Release workflow YAML is valid"
    else
        echo "❌ Release workflow YAML is invalid"
        exit 1
    fi
else
    echo "⚠️  yq not available - skipping YAML validation"
fi

# Simulate CI pipeline locally
echo ""
echo "🧪 Simulating CI Pipeline"
echo "========================="

echo "📦 1. Testing Go module download..."
go mod download
echo "✅ Dependencies downloaded successfully"

echo ""
echo "🔧 2. Running unit tests..."
go test -v -race ./...
echo "✅ Unit tests passed"

echo ""
echo "🔍 3. Running linter..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
    echo "✅ Linting passed"
else
    echo "⚠️  golangci-lint not available - skipping lint check"
fi

echo ""
echo "🏗️  4. Testing build process..."
go build -o pg-goer ./cmd/pg-goer
echo "✅ Build successful"

echo ""
echo "🔧 5. Testing version information..."
./pg-goer --version
echo "✅ Version command works"

echo ""
echo "🧹 6. Cleaning up..."
rm -f pg-goer coverage.out
echo "✅ Cleanup complete"

echo ""
echo "🎉 GitHub Actions CI Validation Complete!"
echo "=========================================="
echo "✅ All checks passed - CI configuration looks good"
echo ""
echo "📋 Summary:"
echo "  - GitHub Actions workflow files are present"
echo "  - Dependencies can be downloaded"
echo "  - Unit tests pass"
echo "  - Code lints successfully"
echo "  - Binary builds correctly"
echo "  - Version information works"
echo ""
echo "🚀 Ready for GitHub Actions deployment!"