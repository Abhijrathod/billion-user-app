package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/my-username/billion-user-app/pkg/jwtutils"
	"github.com/my-username/billion-user-app/services/media-service/internal/domain"
	"github.com/my-username/billion-user-app/services/media-service/internal/service"
)

type MediaHandler struct {
	mediaService service.MediaService
}

func NewMediaHandler(mediaService service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

// CreateMediaRequest represents a media creation request
type CreateMediaRequest struct {
	FileName     string `json:"file_name"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Metadata     string `json:"metadata"`
}

// PresignedURLRequest represents a presigned URL request
type PresignedURLRequest struct {
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	FileSize   int64  `json:"file_size"`
	ExpiresIn  int    `json:"expires_in"` // seconds
}

func (h *MediaHandler) CreateMedia(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*jwtutils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req CreateMediaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	media := &domain.Media{
		UserID:       claims.UserID,
		FileName:     req.FileName,
		FileType:     domain.MediaType(req.FileType),
		FileSize:     req.FileSize,
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		Metadata:     req.Metadata,
	}

	createdMedia, err := h.mediaService.CreateMedia(media)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create media",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdMedia)
}

func (h *MediaHandler) GetMedia(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid media ID",
		})
	}

	media, err := h.mediaService.GetMediaByID(id)
	if err != nil {
		if err == service.ErrMediaNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Media not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get media",
		})
	}

	return c.JSON(media)
}

func (h *MediaHandler) GetMyMedia(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*jwtutils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit > 100 {
		limit = 100
	}

	media, err := h.mediaService.GetMediaByUserID(claims.UserID, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get media",
		})
	}

	return c.JSON(fiber.Map{
		"media":  media,
		"offset": offset,
		"limit":  limit,
	})
}

func (h *MediaHandler) DeleteMedia(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid media ID",
		})
	}

	claims, ok := c.Locals("claims").(*jwtutils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := h.mediaService.DeleteMedia(id, claims.UserID); err != nil {
		if err == service.ErrMediaNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Media not found",
			})
		}
		if err == service.ErrUnauthorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Unauthorized to delete this media",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete media",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Media deleted successfully",
	})
}

func (h *MediaHandler) GetPresignedURL(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*jwtutils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req PresignedURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate a key for the file (user_id/filename)
	key := strconv.FormatUint(claims.UserID, 10) + "/" + req.FileName
	expiresIn := req.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600 // Default 1 hour
	}

	url, err := h.mediaService.GeneratePresignedURL("media-bucket", key, expiresIn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate presigned URL",
		})
	}

	return c.JSON(fiber.Map{
		"url":        url,
		"key":        key,
		"expires_in": expiresIn,
	})
}

// JWTMiddleware validates JWT tokens
func JWTMiddleware(jwtManager *jwtutils.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("claims", claims)
		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}

