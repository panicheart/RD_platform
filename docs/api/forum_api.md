# Forum API Documentation

## Overview

Forum API provides endpoints for managing discussion boards, posts, replies, and tags.

Base URL: `/api/v1`

## Authentication

All endpoints require authentication via JWT Bearer token.

## Boards

### List Boards

```http
GET /boards
```

Query Parameters:
- `category` (string, optional): Filter by category
- `page` (int, optional): Page number (default: 1)
- `page_size` (int, optional): Items per page (default: 20, max: 100)

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "name": "技术讨论",
        "description": "技术相关讨论板块",
        "category": "tech",
        "icon": "code",
        "sort_order": 1,
        "topic_count": 100,
        "post_count": 500,
        "last_post_at": "2026-02-23T10:30:00Z",
        "is_active": true,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-02-23T10:30:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### Get Board

```http
GET /boards/:id
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
    "name": "技术讨论",
    "description": "技术相关讨论板块",
    "category": "tech",
    "icon": "code",
    "sort_order": 1,
    "topic_count": 100,
    "post_count": 500,
    "last_post_at": "2026-02-23T10:30:00Z",
    "is_active": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-02-23T10:30:00Z"
  }
}
```

### Create Board

```http
POST /boards
```

Required Role: admin, team_leader

Request Body:
```json
{
  "name": "技术讨论",
  "description": "技术相关讨论板块",
  "category": "tech",
  "icon": "code",
  "sort_order": 1
}
```

Response:
```json
{
  "code": 0,
  "message": "board created",
  "data": {
    "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
    "name": "技术讨论",
    "description": "技术相关讨论板块",
    "category": "tech",
    "icon": "code",
    "sort_order": 1,
    "is_active": true,
    "created_at": "2026-02-23T10:30:00Z",
    "updated_at": "2026-02-23T10:30:00Z"
  }
}
```

### Update Board

```http
PUT /boards/:id
```

Required Role: admin, team_leader

Request Body:
```json
{
  "name": "技术讨论区",
  "description": "更新后的描述",
  "category": "tech",
  "icon": "code",
  "sort_order": 2,
  "is_active": true
}
```

Response:
```json
{
  "code": 0,
  "message": "board updated",
  "data": { ... }
}
```

### Delete Board

```http
DELETE /boards/:id
```

Required Role: admin

Response:
```json
{
  "code": 0,
  "message": "board deleted",
  "data": null
}
```

## Posts

### List Posts

```http
GET /posts
```

Query Parameters:
- `board_id` (string, optional): Filter by board
- `author_id` (string, optional): Filter by author
- `search` (string, optional): Search in title and content
- `is_pinned` (boolean, optional): Filter by pinned status
- `page` (int, optional): Page number (default: 1)
- `page_size` (int, optional): Items per page (default: 20, max: 100)

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "board_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "title": "如何优化Go代码性能",
        "content": "内容...",
        "author_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "author_name": "张三",
        "view_count": 100,
        "reply_count": 20,
        "is_pinned": false,
        "is_locked": false,
        "is_best_answer": false,
        "tags": "[\"go\",\"performance\"]",
        "last_reply_at": "2026-02-23T10:30:00Z",
        "created_at": "2026-02-23T10:00:00Z",
        "updated_at": "2026-02-23T10:30:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### Get Post

```http
GET /posts/:id
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
    "board_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
    "title": "如何优化Go代码性能",
    "content": "内容...",
    "author_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
    "author_name": "张三",
    "view_count": 101,
    "reply_count": 20,
    "is_pinned": false,
    "is_locked": false,
    "is_best_answer": false,
    "tags": "[\"go\",\"performance\"]",
    "knowledge_id": null,
    "last_reply_at": "2026-02-23T10:30:00Z",
    "created_at": "2026-02-23T10:00:00Z",
    "updated_at": "2026-02-23T10:30:00Z"
  }
}
```

### Create Post

```http
POST /posts
```

Request Body:
```json
{
  "title": "如何优化Go代码性能",
  "content": "详细内容支持 **Markdown** 格式\n\n```go\nfunc main() {\n  fmt.Println(\"Hello\")\n}\n```",
  "board_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
  "tags": ["go", "performance"],
  "knowledge_id": null
}
```

Response:
```json
{
  "code": 0,
  "message": "post created",
  "data": { ... }
}
```

### Update Post

```http
PUT /posts/:id
```

Request Body:
```json
{
  "title": "更新后的标题",
  "content": "更新后的内容",
  "tags": ["go", "performance", "optimization"]
}
```

Response:
```json
{
  "code": 0,
  "message": "post updated",
  "data": { ... }
}
```

### Delete Post

```http
DELETE /posts/:id
```

Response:
```json
{
  "code": 0,
  "message": "post deleted",
  "data": null
}
```

### Pin Post

```http
POST /posts/:id/pin
```

Required Role: admin, team_leader

Response:
```json
{
  "code": 0,
  "message": "post pin status toggled",
  "data": null
}
```

### Lock Post

```http
POST /posts/:id/lock
```

Required Role: admin, team_leader

Response:
```json
{
  "code": 0,
  "message": "post lock status toggled",
  "data": null
}
```

## Replies

### List Replies

```http
GET /posts/:postId/replies
```

Query Parameters:
- `parent_id` (string, optional): Filter by parent reply (for nested replies)
- `page` (int, optional): Page number (default: 1)
- `page_size` (int, optional): Items per page (default: 20, max: 100)

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "post_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "parent_id": null,
        "content": "回复内容",
        "author_id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
        "author_name": "李四",
        "is_best_answer": false,
        "mentions": "[\"zhangsan\"]",
        "created_at": "2026-02-23T10:30:00Z",
        "updated_at": "2026-02-23T10:30:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 20
  }
}
```

### Create Reply

```http
POST /posts/:postId/replies
```

Request Body:
```json
{
  "content": "@zhangsan 这是一个很好的问题...",
  "parent_id": null
}
```

Response:
```json
{
  "code": 0,
  "message": "reply created",
  "data": { ... }
}
```

### Update Reply

```http
PUT /replies/:id
```

Request Body:
```json
{
  "content": "更新后的回复内容"
}
```

Response:
```json
{
  "code": 0,
  "message": "reply updated",
  "data": { ... }
}
```

### Delete Reply

```http
DELETE /replies/:id
```

Response:
```json
{
  "code": 0,
  "message": "reply deleted",
  "data": null
}
```

## Tags

### List Tags

```http
GET /forum/tags
```

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "01J8K2M3N4P5Q6R7S8T9U0V1W",
      "name": "go",
      "color": "#1890ff",
      "count": 50,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

### Create Tag

```http
POST /forum/tags
```

Required Role: admin

Request Body:
```json
{
  "name": "javascript",
  "color": "#f7df1e"
}
```

Response:
```json
{
  "code": 0,
  "message": "tag created",
  "data": { ... }
}
```

### Delete Tag

```http
DELETE /forum/tags/:id
```

Required Role: admin

Response:
```json
{
  "code": 0,
  "message": "tag deleted",
  "data": null
}
```

## Search

### Search Posts

```http
GET /forum/search?q=keyword&board_id=&page=1&page_size=20
```

Query Parameters:
- `q` (string, required): Search query
- `board_id` (string, optional): Filter by board
- `page` (int, optional): Page number (default: 1)
- `page_size` (int, optional): Items per page (default: 20, max: 100)

Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [ ... ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 4000 | Bad Request - Invalid parameters |
| 4001 | Bad Request - Business logic error |
| 4010 | Unauthorized - Authentication required |
| 4030 | Forbidden - Insufficient permissions |
| 4040 | Not Found - Resource not found |
| 5000 | Internal Server Error |

## Features

### @Mentions

Support @username mentions in posts and replies. Mentioned users will receive notifications.

Example:
```markdown
@zhangsan @lisi 请大家看一下这个问题
```

### Markdown Support

Posts and replies support Markdown formatting:
- Headers (# ## ###)
- Bold (**text**)
- Italic (*text*)
- Code blocks (```language)
- Links [text](url)
- Lists (- or 1.)

### Nested Replies

Replies support nested threading:
- Top-level replies have `parent_id: null`
- Nested replies reference their parent reply ID

### Auto-counters

The system automatically maintains counters:
- `topic_count`: Number of topics in a board
- `post_count`: Total posts (topics + replies)
- `reply_count`: Number of replies to a post
- `view_count`: Number of times a post is viewed
