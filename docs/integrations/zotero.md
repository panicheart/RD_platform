# Zotero Integration API Documentation

## Overview

The Zotero integration allows users to connect their Zotero library to the RDP platform, synchronize literature items, preview PDFs, and generate citations in various formats (GB/T 7714-2015, APA, MLA).

## Base URL

```
/api/v1/zotero
```

## Authentication

All endpoints require authentication via JWT Bearer token.

```
Authorization: Bearer <token>
```

---

## Connection Management

### 1. Get Connection Status

Retrieves the current user's Zotero connection configuration.

**Endpoint:** `GET /api/v1/zotero/connection`

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "01GQKH...",
    "user_id": "01GQKH...",
    "api_key": "***",
    "zotero_user_id": "12345678",
    "is_active": true,
    "last_sync_at": "2026-02-23T10:30:00Z",
    "created_at": "2026-02-23T08:00:00Z",
    "updated_at": "2026-02-23T10:30:00Z"
  }
}
```

### 2. Save Connection

Saves or updates Zotero API credentials for the current user.

**Endpoint:** `POST /api/v1/zotero/connection`

**Request Body:**
```json
{
  "api_key": "your_zotero_api_key",
  "zotero_user_id": "your_zotero_user_id"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "01GQKH...",
    "user_id": "01GQKH...",
    "api_key": "***",
    "zotero_user_id": "12345678",
    "is_active": true,
    "last_sync_at": null,
    "created_at": "2026-02-23T08:00:00Z",
    "updated_at": "2026-02-23T08:00:00Z"
  }
}
```

### 3. Delete Connection

Removes the current user's Zotero connection.

**Endpoint:** `DELETE /api/v1/zotero/connection`

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "connection deleted"
  }
}
```

### 4. Test Connection

Tests Zotero credentials without saving them.

**Endpoint:** `POST /api/v1/zotero/connection/test`

**Request Body:**
```json
{
  "api_key": "your_zotero_api_key",
  "zotero_user_id": "your_zotero_user_id"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "connected",
    "message": "Zotero connection successful"
  }
}
```

---

## Item Operations

### 5. List Items

Retrieves a paginated list of synchronized Zotero items.

**Endpoint:** `GET /api/v1/zotero/items`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| item_type | string | No | Filter by item type (book, journalArticle, etc.) |
| tag | string | No | Filter by tag |
| search | string | No | Search in title, authors, abstract |
| page | int | No | Page number (default: 1) |
| page_size | int | No | Items per page (default: 20) |

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "01GQKH...",
        "zotero_key": "A1B2C3D4",
        "title": "Research on Microwave Systems",
        "item_type": "journalArticle",
        "authors": "[{\"creatorType\":\"author\",\"firstName\":\"John\",\"lastName\":\"Smith\"}]",
        "abstract": "This paper discusses...",
        "publication": "IEEE Transactions",
        "volume": "45",
        "issue": "3",
        "pages": "123-135",
        "date": "2025-06",
        "doi": "10.1109/example.2025.1234567",
        "url": "https://doi.org/10.1109/example.2025.1234567",
        "pdf_path": "E5F6G7H8",
        "tags": "[{\"tag\":\"microwave\"},{\"tag\":\"RF\"}]",
        "synced_at": "2026-02-23T10:30:00Z",
        "created_at": "2026-02-23T08:00:00Z",
        "updated_at": "2026-02-23T10:30:00Z"
      }
    ],
    "total": 156,
    "page": 1
  }
}
```

### 6. Get Item

Retrieves a single Zotero item by ID.

**Endpoint:** `GET /api/v1/zotero/items/:id`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Item internal ID |

**Response:** Same as item in List Items response.

### 7. Delete Item

Deletes a Zotero item from the local database (does not delete from Zotero).

**Endpoint:** `DELETE /api/v1/zotero/items/:id`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Item internal ID |

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "message": "item deleted"
  }
}
```

---

## Synchronization

### 8. Sync Items

Triggers synchronization of Zotero items from the user's library.

**Endpoint:** `POST /api/v1/zotero/sync`

**Request Body:**
```json
{
  "incremental": true
}
```

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| incremental | boolean | No | If true, only sync items modified since last sync |

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "created": 5,
    "updated": 12,
    "deleted": 0,
    "errors": []
  }
}
```

---

## PDF Operations

### 9. Get PDF URL

Returns a URL for viewing the PDF attachment of a Zotero item.

**Endpoint:** `GET /api/v1/zotero/items/:id/pdf`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Item internal ID |

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "https://api.zotero.org/users/12345678/items/E5F6G7H8/file/view"
  }
}
```

---

## Citation Operations

### 10. Generate Citation

Generates a formatted citation for a Zotero item.

**Endpoint:** `POST /api/v1/zotero/items/:id/citation`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | string | Yes | Item internal ID |

**Request Body:**
```json
{
  "format": "gb7714"
}
```

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| format | string | Yes | Citation format: `gb7714`, `apa`, `mla` |

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "citation": "Smith John. Research on Microwave Systems[J]. IEEE Transactions, 45(3): 123-135. 2025-06. DOI:10.1109/example.2025.1234567.",
    "format": "gb7714"
  }
}
```

**Supported Formats:**

- **gb7714**: GB/T 7714-2015 (Chinese national standard)
  - Example: `Smith John. Title[J]. Journal, 45(3): 123-135. 2025. DOI:10.xxx.`
  
- **apa**: APA 7th Edition
  - Example: `Smith, J. (2025). Title. *Journal*, 45(3), 123-135. https://doi.org/10.xxx`
  
- **mla**: MLA 9th Edition
  - Example: `Smith, John. "Title." *Journal* 45.3 (2025): 123-135.`

---

## Collection Operations

### 11. Get Collections

Retrieves all collections/folders from the user's Zotero library.

**Endpoint:** `GET /api/v1/zotero/collections`

**Response:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "key": "A1B2C3D4",
      "version": 1234,
      "name": "Microwave Research",
      "parentCollection": ""
    },
    {
      "key": "E5F6G7H8",
      "version": 5678,
      "name": "RF Components",
      "parentCollection": "A1B2C3D4"
    }
  ]
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| 400 | Bad Request | Invalid request parameters or missing required fields |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | User does not have permission |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error or Zotero API error |

---

## Data Models

### ZoteroItem

| Field | Type | Description |
|-------|------|-------------|
| id | string | Internal ULID |
| zotero_key | string | Zotero item key |
| title | string | Item title |
| item_type | string | Type: book, journalArticle, conferencePaper, etc. |
| authors | string | JSON array of creators |
| abstract | string | Abstract text |
| publication | string | Publication title |
| volume | string | Volume number |
| issue | string | Issue number |
| pages | string | Page range |
| date | string | Publication date |
| doi | string | DOI identifier |
| url | string | URL |
| pdf_path | string | Attachment key for PDF |
| tags | string | JSON array of tags |
| synced_at | timestamp | Last sync time |

---

## Setup Instructions

### 1. Get Zotero API Key

1. Log in to [Zotero](https://www.zotero.org/)
2. Go to Settings → Feeds/API
3. Create a new API key with "Allow library access" permission
4. Note your User ID from the same page

### 2. Configure Connection

```bash
# Test connection
curl -X POST /api/v1/zotero/connection/test \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your_api_key",
    "zotero_user_id": "your_user_id"
  }'

# Save connection
curl -X POST /api/v1/zotero/connection \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your_api_key",
    "zotero_user_id": "your_user_id"
  }'

# Sync items
curl -X POST /api/v1/zotero/sync \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "incremental": false
  }'
```

---

## Implementation Notes

### Sync Behavior

- **Full Sync**: Downloads all items from Zotero (use with `incremental: false`)
- **Incremental Sync**: Only downloads items modified since last sync
- **Safety Limit**: Sync stops at 10,000 items to prevent timeouts
- **PDF Attachments**: Automatically detects and links PDF attachments

### Security

- API keys are stored encrypted in the database
- API keys are masked (*** ) in API responses
- Each user can only access their own Zotero library

### Performance

- Items are paginated (default: 20 per page)
- Sync operations may take time for large libraries
- Consider scheduling incremental syncs periodically

---

**Last Updated:** 2026-02-23  
**API Version:** v1  
**Zotero API Version:** v3
