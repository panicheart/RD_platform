#!/bin/bash

check_api() {
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/health)
    if [[ $response -eq 200 ]]; then
        echo "API: OK"
        return 0
    else
        echo "API: FAILED (HTTP $response)"
        return 1
    fi
}

check_nginx() {
    response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/health)
    if [[ $response -eq 200 ]]; then
        echo "Nginx: OK"
        return 0
    else
        echo "Nginx: FAILED (HTTP $response)"
        return 1
    fi
}

check_postgresql() {
    if pg_isready -h localhost -U rdp -d rdp &>/dev/null; then
        echo "PostgreSQL: OK"
        return 0
    else
        echo "PostgreSQL: FAILED"
        return 1
    fi
}

check_disk() {
    usage=$(df -h /opt/rdp | awk 'NR==2 {print $5}' | sed 's/%//')
    if [[ $usage -gt 90 ]]; then
        echo "Disk: WARNING (${usage}% used)"
        return 1
    else
        echo "Disk: OK (${usage}% used)"
        return 0
    fi
}

echo "RDP Platform Health Check"
echo "========================="
echo ""

failed=0

check_api || failed=$((failed + 1))
check_nginx || failed=$((failed + 1))
check_postgresql || failed=$((failed + 1))
check_disk || failed=$((failed + 1))

echo ""
if [[ $failed -eq 0 ]]; then
    echo "All checks passed!"
    exit 0
else
    echo "$failed check(s) failed!"
    exit 1
fi
