# Monitor API Documentation

## Overview

The Monitor API provides system monitoring, metrics collection, logging, and alerting capabilities for the RDP platform.

## Base URL

```
/api/v1/monitor
```

## Authentication

All endpoints require authentication via JWT Bearer token.

## Endpoints

### System Metrics

#### Get System Metrics

Retrieves historical system metrics data.

```http
GET /api/v1/monitor/metrics/system
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_time | string | No | Start time in RFC3339 format |
| end_time | string | No | End time in RFC3339 format |
| limit | integer | No | Maximum number of records (default: 100) |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "01HJSQ8N7Q2...",
      "timestamp": "2026-02-23T10:30:00Z",
      "cpu_usage": 45.2,
      "memory_usage": 67.5,
      "memory_total": 17179869184,
      "memory_used": 11596411699,
      "disk_usage": 55.3,
      "disk_total": 107374182400,
      "disk_used": 59381349952,
      "network_in": 1234567890,
      "network_out": 987654321,
      "db_connections": 15,
      "api_requests": 1250,
      "created_at": "2026-02-23T10:30:00Z"
    }
  ]
}
```

#### Get System Metric Statistics

Returns aggregated statistics for system metrics.

```http
GET /api/v1/monitor/metrics/system/stats
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_time | string | No | Start time in RFC3339 format |
| end_time | string | No | End time in RFC3339 format |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "cpu": {
      "avg": 42.5,
      "max": 85.3
    },
    "memory": {
      "avg": 65.2,
      "max": 78.9
    },
    "disk": {
      "avg": 52.1,
      "max": 55.3
    },
    "network": {
      "in": 12345678900,
      "out": 9876543210
    }
  }
}
```

### API Metrics

#### Get API Metrics

Retrieves API call metrics.

```http
GET /api/v1/monitor/metrics/api
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_time | string | No | Start time in RFC3339 format |
| end_time | string | No | End time in RFC3339 format |
| endpoint | string | No | Filter by endpoint path |
| limit | integer | No | Maximum number of records (default: 100) |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "01HJSQ8N7Q2...",
      "timestamp": "2026-02-23T10:30:00Z",
      "endpoint": "/api/v1/projects",
      "method": "GET",
      "duration": 45,
      "status_code": 200,
      "user_id": "01HJSP...",
      "ip_address": "192.168.1.100",
      "created_at": "2026-02-23T10:30:00Z"
    }
  ]
}
```

#### Get API Metric Statistics

Returns aggregated API performance statistics.

```http
GET /api/v1/monitor/metrics/api/stats
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_requests": 12500,
    "total_errors": 45,
    "error_rate": 0.36,
    "avg_duration": 52.3,
    "endpoints": [
      {
        "endpoint": "/api/v1/projects",
        "count": 3500,
        "avg_duration": 45.2,
        "max_duration": 1250,
        "min_duration": 12,
        "error_count": 12
      }
    ]
  }
}
```

### Prometheus Metrics

Returns metrics in Prometheus exposition format.

```http
GET /api/v1/monitor/metrics/prometheus
```

**Response:**

```text
# HELP rdp_cpu_usage CPU usage percentage
# TYPE rdp_cpu_usage gauge
rdp_cpu_usage 45.20
# HELP rdp_memory_usage Memory usage percentage
# TYPE rdp_memory_usage gauge
rdp_memory_usage 67.50
```

### Logs

#### Get Log Entries

Retrieves log entries with filtering and pagination.

```http
GET /api/v1/monitor/logs
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number (default: 1) |
| page_size | integer | No | Items per page (default: 20) |
| level | string | No | Log level: DEBUG, INFO, WARN, ERROR |
| source | string | No | Log source/service name |
| module | string | No | Module/component name |
| keyword | string | No | Search keyword in message |
| start_time | string | No | Start time in RFC3339 format |
| end_time | string | No | End time in RFC3339 format |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01HJSQ8N7Q2...",
        "timestamp": "2026-02-23T10:30:00Z",
        "level": "ERROR",
        "message": "Failed to connect to database",
        "source": "rdp-api",
        "module": "database",
        "user_id": "01HJSP...",
        "request_id": "req-12345",
        "metadata": "{\"retry_count\": 3}",
        "created_at": "2026-02-23T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

#### Get Log Sources

Returns list of unique log sources.

```http
GET /api/v1/monitor/logs/sources
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": ["rdp-api", "rdp-worker", "nginx"]
}
```

### Alert Rules

#### List Alert Rules

Retrieves all alert rules.

```http
GET /api/v1/monitor/alerts/rules
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number (default: 1) |
| page_size | integer | No | Items per page (default: 20) |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01HJSQ8N7Q2...",
        "name": "High CPU Usage",
        "description": "CPU usage exceeds 80% for more than 5 minutes",
        "metric": "cpu_usage",
        "condition": ">",
        "threshold": 80.0,
        "duration": 5,
        "severity": "warning",
        "is_active": true,
        "notify_channels": "[\"in_app\", \"email\"]",
        "created_at": "2026-02-23T10:00:00Z",
        "updated_at": "2026-02-23T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

#### Create Alert Rule

Creates a new alert rule.

```http
POST /api/v1/monitor/alerts/rules
```

**Request Body:**

```json
{
  "name": "High Memory Usage",
  "description": "Memory usage exceeds 90%",
  "metric": "memory_usage",
  "condition": ">",
  "threshold": 90.0,
  "duration": 10,
  "severity": "critical",
  "is_active": true,
  "notify_channels": ["in_app", "email"]
}
```

**Response:**

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": "01HJSQ8N7Q2...",
    "name": "High Memory Usage",
    ...
  }
}
```

#### Get Alert Rule

Retrieves a specific alert rule.

```http
GET /api/v1/monitor/alerts/rules/:id
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "01HJSQ8N7Q2...",
    "name": "High CPU Usage",
    ...
  }
}
```

#### Update Alert Rule

Updates an existing alert rule.

```http
PUT /api/v1/monitor/alerts/rules/:id
```

**Request Body:**

```json
{
  "threshold": 85.0,
  "is_active": false
}
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "alert rule updated"
  }
}
```

#### Delete Alert Rule

Deletes an alert rule.

```http
DELETE /api/v1/monitor/alerts/rules/:id
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "alert rule deleted"
  }
}
```

### Alert History

#### Get Alert History

Retrieves alert history with filters.

```http
GET /api/v1/monitor/alerts/history
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | integer | No | Page number (default: 1) |
| page_size | integer | No | Items per page (default: 20) |
| status | string | No | Filter by status: firing, resolved |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01HJSQ8N7Q2...",
        "rule_id": "01HJSQ8N7Q1...",
        "rule_name": "High CPU Usage",
        "severity": "warning",
        "message": "cpu_usage is > 80.00 (current: 85.30)",
        "value": 85.30,
        "threshold": 80.00,
        "status": "firing",
        "resolved_at": null,
        "created_at": "2026-02-23T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

#### Resolve Alert

Marks an alert as resolved.

```http
PUT /api/v1/monitor/alerts/history/:id/resolve
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "alert resolved"
  }
}
```

#### Get Alert Statistics

Returns alert statistics.

```http
GET /api/v1/monitor/alerts/stats
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_rules": 5,
    "active_rules": 4,
    "firing_alerts": 1,
    "resolved_alerts": 15
  }
}
```

### Health Check

#### Get Health Status

Returns detailed health status of the system.

```http
GET /api/v1/monitor/health
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "timestamp": "2026-02-23T10:30:00Z",
    "system_info": {
      "go_version": "go1.22.0",
      "go_os": "linux",
      "go_arch": "amd64",
      "num_cpu": 8,
      "num_goroutine": 45
    },
    "checks": {
      "cpu": {
        "status": "healthy",
        "value": 45.2
      },
      "memory": {
        "status": "warning",
        "value": 82.5
      },
      "disk": {
        "status": "healthy",
        "value": 55.3
      }
    }
  }
}
```

**Status Values:**

- `healthy`: All metrics within normal range
- `warning`: One or more metrics above warning threshold
- `degraded`: One or more metrics above critical threshold
- `critical`: System unavailable

### System Information

#### Get System Information

Returns static system information.

```http
GET /api/v1/monitor/system/info
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "go_version": "go1.22.0",
    "go_os": "linux",
    "go_arch": "amd64",
    "num_cpu": 8,
    "num_goroutine": 45,
    "memory_alloc": 123456789,
    "memory_sys": 234567890,
    "memory_heap_alloc": 98765432,
    "memory_heap_sys": 123456789,
    "partitions": [
      {
        "device": "root",
        "mountpoint": "/",
        "fstype": "ext4",
        "total": 107374182400,
        "used": 59381349952,
        "free": 47992832448,
        "used_percent": 55.3
      }
    ]
  }
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

## Metrics Reference

### System Metrics

| Metric | Unit | Description |
|--------|------|-------------|
| cpu_usage | percentage | CPU usage percentage |
| memory_usage | percentage | Memory usage percentage |
| memory_total | bytes | Total memory |
| memory_used | bytes | Used memory |
| disk_usage | percentage | Disk usage percentage |
| disk_total | bytes | Total disk space |
| disk_used | bytes | Used disk space |
| network_in | bytes | Network bytes received |
| network_out | bytes | Network bytes sent |
| db_connections | count | Active database connections |
| api_requests | count | Total API requests |

### Alert Conditions

| Condition | Description |
|-----------|-------------|
| > | Greater than |
| >= | Greater than or equal |
| < | Less than |
| <= | Less than or equal |
| == | Equal to |
| != | Not equal to |

### Alert Severity Levels

| Severity | Description |
|----------|-------------|
| warning | Warning level - requires attention |
| critical | Critical level - immediate action required |

### Notification Channels

| Channel | Description |
|---------|-------------|
| in_app | In-app notification |
| email | Email notification |
| webhook | Webhook callback |

## Rate Limits

- Metrics endpoints: 60 requests per minute
- Log endpoints: 30 requests per minute
- Alert endpoints: 30 requests per minute

## Data Retention

- System metrics: 30 days
- API metrics: 7 days
- Log entries: 14 days
- Alert history: 90 days

---

*Last Updated: 2026-02-23*
