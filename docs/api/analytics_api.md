# Analytics API Documentation

## Overview

The Analytics API provides statistical data and reporting capabilities for the RDP platform. It includes project statistics, user activity metrics, shelf utilization data, and knowledge base analytics.

## Base URL

```
/api/v1/analytics
```

## Authentication

All endpoints require authentication via JWT Bearer token.

## Endpoints

### Dashboard

#### Get Dashboard Overview

Returns an overview of key metrics for the dashboard.

```
GET /api/v1/analytics/dashboard
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_projects": 150,
    "active_projects": 45,
    "completed_projects": 98,
    "delayed_projects": 7,
    "total_users": 85,
    "active_users": 72,
    "total_knowledge": 320,
    "total_products": 45,
    "avg_project_progress": 67.5
  }
}
```

#### Get Dashboard Widgets

Returns all widget data for the dashboard including time series, charts, and heatmaps.

```
GET /api/v1/analytics/dashboard/widgets?start_date=2026-01-01&end_date=2026-02-23
```

**Query Parameters:**

| Parameter  | Type   | Required | Description                      |
|------------|--------|----------|----------------------------------|
| start_date | string | No       | Start date (YYYY-MM-DD), default: 30 days ago |
| end_date   | string | No       | End date (YYYY-MM-DD), default: today         |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "stat_cards": [
      {
        "title": "总项目数",
        "value": 150,
        "change": 12.5,
        "change_direction": "up",
        "icon": "project",
        "color": "#1890ff"
      }
    ],
    "time_series": [
      {
        "name": "project_trend",
        "label": "项目趋势",
        "data": [
          {
            "timestamp": 1704067200000,
            "value": 45,
            "label": "2026-01"
          }
        ]
      }
    ],
    "pie_charts": [...],
    "bar_charts": [...],
    "heatmap": {...}
  }
}
```

### Project Statistics

#### Get Project Statistics

Returns detailed project statistics for a given date range.

```
GET /api/v1/analytics/projects?start_date=2026-01-01&end_date=2026-02-23
```

**Query Parameters:**

| Parameter  | Type   | Required | Description                      |
|------------|--------|----------|----------------------------------|
| start_date | string | No       | Start date (YYYY-MM-DD)          |
| end_date   | string | No       | End date (YYYY-MM-DD)            |

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_projects": 50,
    "status_distribution": {
      "active": 20,
      "completed": 25,
      "draft": 5
    },
    "category_distribution": {
      "R&D": 30,
      "Product": 20
    },
    "monthly_trend": [
      {
        "month": "2026-01",
        "created": 10,
        "completed": 8,
        "active": 15
      }
    ],
    "top_projects": [
      {
        "id": "01HMD...",
        "code": "RDP-PD-20260101-001",
        "name": "微波组件研发",
        "status": "active",
        "progress": 85,
        "leader_id": "01HMD..."
      }
    ]
  }
}
```

### User Statistics

#### Get User Statistics

Returns user activity and contribution statistics.

```
GET /api/v1/analytics/users?start_date=2026-01-01&end_date=2026-02-23
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_users": 85,
    "active_users": 72,
    "new_users": 5,
    "top_contributors": [
      {
        "user_id": "01HMD...",
        "display_name": "张三",
        "project_count": 12,
        "knowledge_count": 8,
        "contribution": 200.0
      }
    ],
    "department_stats": [
      {
        "department": "产品研发部",
        "user_count": 25,
        "project_count": 18
      }
    ],
    "monthly_activity": [
      {
        "month": "2026-01",
        "active_users": 68,
        "new_users": 3
      }
    ]
  }
}
```

### Shelf Statistics

#### Get Shelf Statistics

Returns product shelf and technology utilization statistics.

```
GET /api/v1/analytics/shelf?start_date=2026-01-01&end_date=2026-02-23
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_products": 45,
    "published_products": 38,
    "total_technologies": 25,
    "adoption_rate": 84.4,
    "reuse_rate": 62.5,
    "category_stats": [
      {
        "category": "微波组件",
        "count": 15,
        "usage": 45
      }
    ],
    "top_products": [
      {
        "product_id": "01HMD...",
        "product_name": "高频滤波器",
        "cart_count": 12
      }
    ]
  }
}
```

### Knowledge Statistics

#### Get Knowledge Statistics

Returns knowledge base content and engagement statistics.

```
GET /api/v1/analytics/knowledge?start_date=2026-01-01&end_date=2026-02-23
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_knowledge": 320,
    "published_count": 280,
    "draft_count": 40,
    "total_views": 15420,
    "category_stats": [
      {
        "category_id": "01HMD...",
        "category_name": "技术文档",
        "count": 120,
        "views": 8500
      }
    ],
    "top_knowledge": [
      {
        "id": "01HMD...",
        "title": "微波电路设计指南",
        "author_id": "01HMD...",
        "view_count": 450
      }
    ],
    "tag_stats": [
      {
        "tag_id": "01HMD...",
        "tag_name": "设计规范",
        "count": 45,
        "color": "#1890ff"
      }
    ],
    "monthly_trend": [
      {
        "month": "2026-01",
        "created": 25,
        "published": 20
      }
    ]
  }
}
```

### Dashboard Configurations

#### List Dashboard Configurations

```
GET /api/v1/analytics/dashboards
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "01HMD...",
      "name": "默认仪表盘",
      "description": "系统默认仪表盘配置",
      "layout": "{...}",
      "is_default": true,
      "created_by": "01HMD...",
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

#### Get Dashboard Configuration

```
GET /api/v1/analytics/dashboards/:id
```

#### Create Dashboard Configuration

```
POST /api/v1/analytics/dashboards
```

**Request Body:**

```json
{
  "name": "自定义仪表盘",
  "description": "我的自定义仪表盘",
  "layout": "{\"widgets\": [...]}",
  "is_default": false
}
```

#### Update Dashboard Configuration

```
PUT /api/v1/analytics/dashboards/:id
```

**Request Body:** Same as Create

#### Delete Dashboard Configuration

```
DELETE /api/v1/analytics/dashboards/:id
```

#### Set Default Dashboard

```
PUT /api/v1/analytics/dashboards/:id/default
```

### Export

#### Export Statistics

Exports statistics data in various formats.

```
GET /api/v1/analytics/export?type=projects&format=json&start_date=2026-01-01&end_date=2026-02-23
```

**Query Parameters:**

| Parameter  | Type   | Required | Description                                  |
|------------|--------|----------|----------------------------------------------|
| type       | string | Yes      | Export type: projects, users, shelf, knowledge |
| format     | string | Yes      | Export format: json, csv, excel              |
| start_date | string | No       | Start date (YYYY-MM-DD)                      |
| end_date   | string | No       | End date (YYYY-MM-DD)                        |

### Snapshots

#### Generate Snapshot

Generates a statistics snapshot for the current date.

```
POST /api/v1/analytics/snapshot
```

**Response:**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "snapshot generated successfully",
    "date": "2026-02-23"
  }
}
```

## Error Responses

All errors follow the standard API response format:

```json
{
  "code": 400,
  "message": "error description",
  "data": null
}
```

### Common Error Codes

| Code | Description                    |
|------|--------------------------------|
| 400  | Bad Request - Invalid parameters |
| 401  | Unauthorized - Missing or invalid token |
| 403  | Forbidden - Insufficient permissions |
| 404  | Not Found - Resource not found |
| 500  | Internal Server Error          |

## Data Types

### Date Formats

- All dates use ISO 8601 format: `YYYY-MM-DD`
- Timestamps are in milliseconds since Unix epoch

### Response Format

All responses follow the standard format:

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

## Rate Limiting

Analytics endpoints have a rate limit of 100 requests per minute per user.

## Caching

Dashboard overview data is cached for 5 minutes. Other statistics are calculated in real-time.
