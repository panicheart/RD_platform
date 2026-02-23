package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"rdp-platform/rdp-api/services"
)

type ForumHandler struct {
	forumService *services.ForumService
}

func NewForumHandler(forumService *services.ForumService) *ForumHandler {
	return &ForumHandler{
		forumService: forumService,
	}
}

func (h *ForumHandler) ListBoards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := services.ListBoardsQuery{
		Category: category,
		Page:     page,
		PageSize: pageSize,
	}

	boards, total, err := h.forumService.ListBoards(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     boards,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *ForumHandler) GetBoard(c *gin.Context) {
	id := c.Param("id")

	board, err := h.forumService.GetBoardByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    board,
	})
}

func (h *ForumHandler) CreateBoard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	var req services.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	board, err := h.forumService.CreateBoard(req, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "board created",
		"data":    board,
	})
}

func (h *ForumHandler) UpdateBoard(c *gin.Context) {
	id := c.Param("id")

	var req services.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	board, err := h.forumService.UpdateBoard(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "board updated",
		"data":    board,
	})
}

func (h *ForumHandler) DeleteBoard(c *gin.Context) {
	id := c.Param("id")

	if err := h.forumService.DeleteBoard(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "board deleted",
		"data":    nil,
	})
}

func (h *ForumHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	boardID := c.Query("board_id")
	authorID := c.Query("author_id")
	search := c.Query("search")

	isPinnedStr := c.Query("is_pinned")
	var isPinned *bool
	if isPinnedStr == "true" {
		pinned := true
		isPinned = &pinned
	} else if isPinnedStr == "false" {
		pinned := false
		isPinned = &pinned
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := services.ListPostsQuery{
		BoardID:  boardID,
		AuthorID: authorID,
		Search:   search,
		IsPinned: isPinned,
		Page:     page,
		PageSize: pageSize,
	}

	posts, total, err := h.forumService.ListPosts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     posts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *ForumHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	post, err := h.forumService.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    4040,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    post,
	})
}

func (h *ForumHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	userName, _ := c.Get("user_name")
	authorName := "Anonymous"
	if userName != nil {
		authorName = userName.(string)
	}

	var req services.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	post, err := h.forumService.CreatePost(req, userID.(string), authorName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "post created",
		"data":    post,
	})
}

func (h *ForumHandler) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	var req services.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	post, err := h.forumService.UpdatePost(id, userID.(string), isAdmin, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "post updated",
		"data":    post,
	})
}

func (h *ForumHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.forumService.DeletePost(id, userID.(string), isAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "post deleted",
		"data":    nil,
	})
}

func (h *ForumHandler) PinPost(c *gin.Context) {
	id := c.Param("id")

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.forumService.PinPost(id, isAdmin); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    4030,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "post pin status toggled",
		"data":    nil,
	})
}

func (h *ForumHandler) LockPost(c *gin.Context) {
	id := c.Param("id")

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.forumService.LockPost(id, isAdmin); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    4030,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "post lock status toggled",
		"data":    nil,
	})
}

func (h *ForumHandler) ListReplies(c *gin.Context) {
	postID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	parentID := c.Query("parent_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := services.ListRepliesQuery{
		PostID:   postID,
		Page:     page,
		PageSize: pageSize,
	}

	if parentID != "" {
		query.ParentID = &parentID
	}

	replies, total, err := h.forumService.ListReplies(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     replies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *ForumHandler) CreateReply(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	userName, _ := c.Get("user_name")
	authorName := "Anonymous"
	if userName != nil {
		authorName = userName.(string)
	}

	var req services.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	reply, err := h.forumService.CreateReply(postID, req, userID.(string), authorName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "reply created",
		"data":    reply,
	})
}

func (h *ForumHandler) UpdateReply(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	var req services.UpdateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	reply, err := h.forumService.UpdateReply(id, userID.(string), isAdmin, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "reply updated",
		"data":    reply,
	})
}

func (h *ForumHandler) DeleteReply(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    4010,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	if err := h.forumService.DeleteReply(id, userID.(string), isAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "reply deleted",
		"data":    nil,
	})
}

func (h *ForumHandler) ListTags(c *gin.Context) {
	tags, err := h.forumService.ListTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tags,
	})
}

func (h *ForumHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	tag, err := h.forumService.CreateTag(req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "tag created",
		"data":    tag,
	})
}

func (h *ForumHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")

	if err := h.forumService.DeleteTag(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "tag deleted",
		"data":    nil,
	})
}

func (h *ForumHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    4000,
			"message": "search query is required",
			"data":    nil,
		})
		return
	}

	boardID := c.Query("board_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	posts, total, err := h.forumService.SearchPosts(keyword, boardID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    5000,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     posts,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
