#!/bin/bash
# Comprehensive RDP System Test Script
# Tests all features against requirements

echo "=========================================="
echo "RDP 系统全面测试报告"
echo "测试时间: $(date)"
echo "=========================================="
echo ""

# Configuration
BASE_URL="http://localhost:8080/api/v1"
FRONTEND_URL="http://localhost:3002"
TOKEN=""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0
WARNINGS=0

# Function to test API
test_api() {
    local description="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local auth="$5"
    local expected_code="$6"
    
    echo -n "Testing: $description ... "
    
    if [ "$auth" = "true" ]; then
        if [ "$method" = "GET" ]; then
            response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" -H "Authorization: Bearer $TOKEN" 2>/dev/null)
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                -d "$data" 2>/dev/null)
        fi
    else
        if [ "$method" = "GET" ]; then
            response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint" 2>/dev/null)
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data" 2>/dev/null)
        fi
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected $expected_code, got $http_code)"
        echo "  Response: $body"
        ((FAILED++))
        return 1
    fi
}

echo "=========================================="
echo "Phase 1: 基础功能测试"
echo "=========================================="
echo ""

# 1. Health Check
test_api "健康检查" "GET" "/health" "" "false" "200"

echo ""
echo "=========================================="
echo "Phase 2: 用户认证测试"
echo "=========================================="
echo ""

# 2. Login
echo -n "Testing: 用户登录 ... "
login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"Admin@123"}' 2>/dev/null)

if echo "$login_response" | jq -e '.data.access_token' > /dev/null 2>&1; then
    TOKEN=$(echo "$login_response" | jq -r '.data.access_token')
    echo -e "${GREEN}✓ PASS${NC} (Token generated)"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC} (Login failed)"
    echo "  Response: $login_response"
    ((FAILED++))
    exit 1
fi

# 3. Get current user
test_api "获取当前用户信息" "GET" "/users/me" "" "true" "200"

echo ""
echo "=========================================="
echo "Phase 3: 门户功能测试"
echo "=========================================="
echo ""

# Test portal pages
echo -n "Testing: 门户首页可访问性 ... "
portal_status=$(curl -s -o /dev/null -w "%{http_code}" "$FRONTEND_URL/portal" 2>/dev/null)
if [ "$portal_status" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $portal_status)"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ WARNING${NC} (HTTP $portal_status)"
    ((WARNINGS++))
fi

echo -n "Testing: 登录页面可访问性 ... "
login_status=$(curl -s -o /dev/null -w "%{http_code}" "$FRONTEND_URL/login" 2>/dev/null)
if [ "$login_status" = "200" ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $login_status)"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ WARNING${NC} (HTTP $login_status)"
    ((WARNINGS++))
fi

echo ""
echo "=========================================="
echo "Phase 4: 项目管理功能测试"
echo "=========================================="
echo ""

test_api "项目列表查询" "GET" "/projects" "" "true" "200"
test_api "项目详情查询" "GET" "/projects/01HPVJ8ZM00000000000000501" "" "true" "200"

echo ""
echo "=========================================="
echo "Phase 5: 技术论坛功能测试"
echo "=========================================="
echo ""

test_api "论坛板块列表" "GET" "/boards" "" "true" "200"
test_api "帖子列表查询" "GET" "/posts" "" "true" "200"
test_api "帖子详情查询" "GET" "/posts/01HPVJ8ZM00000000000000301" "" "true" "200"
test_api "帖子回复列表" "GET" "/posts/01HPVJ8ZM00000000000000301/replies" "" "true" "200"

echo ""
echo "=========================================="
echo "Phase 6: 知识库功能测试"
echo "=========================================="
echo ""

test_api "知识库列表查询" "GET" "/knowledge" "" "true" "200"
test_api "知识分类树查询" "GET" "/categories/tree" "" "true" "200"
test_api "标签列表查询" "GET" "/tags" "" "true" "200"
test_api "知识搜索功能" "GET" "/knowledge?q=Go" "" "true" "200"

echo ""
echo "=========================================="
echo "Phase 7: 数据统计测试"
echo "=========================================="
echo ""

# Get data counts
echo "Data Statistics:"
echo -n "  项目数量: "
project_count=$(curl -s "$BASE_URL/projects" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r '.data.total // 0')
echo "$project_count"

echo -n "  知识文档数量: "
knowledge_count=$(curl -s "$BASE_URL/knowledge" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r '.data.total // 0')
echo "$knowledge_count"

echo -n "  论坛板块数量: "
board_count=$(curl -s "$BASE_URL/boards" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r '.data | length // 0')
echo "$board_count"

echo -n "  帖子数量: "
post_count=$(curl -s "$BASE_URL/posts" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r '.data | length // 0')
echo "$post_count"

echo -n "  分类数量: "
category_count=$(curl -s "$BASE_URL/categories/tree" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r '.data | length // 0')
echo "$category_count"

echo ""
echo "=========================================="
echo "测试摘要"
echo "=========================================="
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"
echo -e "${YELLOW}警告: $WARNINGS${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ 所有关键测试通过！系统运行正常。${NC}"
    exit_code=0
else
    echo -e "${RED}✗ 存在失败的测试，请检查系统状态。${NC}"
    exit_code=1
fi

echo ""
echo "=========================================="
echo "测试完成时间: $(date)"
echo "=========================================="

exit $exit_code
