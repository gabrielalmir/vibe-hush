package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabrielalmir/vibe-hush/cache"
	"github.com/gabrielalmir/vibe-hush/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type CacheServer struct {
	cache  *cache.Cache
	router *gin.Engine
	config ServerConfig
	logger *zap.Logger
}

type ServerConfig struct {
	Capacity          int
	DefaultExpiration time.Duration
	AuthToken         string
	CertFile          string
	KeyFile           string
}

type CacheItem struct {
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl,omitempty"`
}

func NewCacheServer(config ServerConfig) *CacheServer {
	logger, _ := zap.NewProduction()

	if config.AuthToken == "" {
		logger.Warn("No auth token provided, using default token")
		config.AuthToken = "default-secret-token" // Default token if none is provided
	}

	server := &CacheServer{
		cache:  cache.NewCache(config.Capacity, config.DefaultExpiration, cache.LRU{}),
		router: gin.Default(),
		config: config,
		logger: logger,
	}

	server.setupMiddlewares()
	server.setupRoutes()
	return server
}

func (s *CacheServer) setupMiddlewares() {
	s.router.Use(s.metricsMiddleware())
	s.router.Use(s.authMiddleware())
}

func (s *CacheServer) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		metrics.RequestDuration.WithLabelValues(
			c.Request.Method,
			path,
			fmt.Sprintf("%d", status),
		).Observe(duration)
	}
}

func (s *CacheServer) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			s.logger.Warn("Request without authorization header",
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token != s.config.AuthToken {
			s.logger.Warn("Invalid token provided",
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *CacheServer) setupRoutes() {
	// Metrics endpoint
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Cache endpoints
	s.router.PUT("/cache/:key", s.setItem)
	s.router.GET("/cache/:key", s.getItem)
	s.router.GET("/cache", s.getAllItems)
	s.router.DELETE("/cache/:key", s.deleteItem)
}

func (s *CacheServer) setItem(c *gin.Context) {
	key := c.Param("key")
	var item CacheItem

	if err := c.BindJSON(&item); err != nil {
		s.logger.Error("Failed to bind JSON",
			zap.String("key", key),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		metrics.CacheOperations.WithLabelValues("set", "error").Inc()
		return
	}

	s.cache.Set(key, item.Value)
	metrics.CacheOperations.WithLabelValues("set", "success").Inc()
	metrics.CacheSize.Set(float64(len(s.cache.GetAll())))

	s.logger.Info("Item set in cache",
		zap.String("key", key),
	)

	c.Status(http.StatusOK)
}

func (s *CacheServer) getItem(c *gin.Context) {
	key := c.Param("key")
	value := s.cache.Get(key)

	if value == nil {
		s.logger.Debug("Cache miss",
			zap.String("key", key),
		)
		metrics.CacheHits.WithLabelValues("miss").Inc()
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	s.logger.Debug("Cache hit",
		zap.String("key", key),
	)
	metrics.CacheHits.WithLabelValues("hit").Inc()
	c.JSON(http.StatusOK, gin.H{"value": value})
}

func (s *CacheServer) getAllItems(c *gin.Context) {
	items := s.cache.GetAll()
	metrics.CacheOperations.WithLabelValues("get_all", "success").Inc()
	c.JSON(http.StatusOK, items)
}

func (s *CacheServer) deleteItem(c *gin.Context) {
	key := c.Param("key")
	s.cache.Delete(key)

	s.logger.Info("Item deleted from cache",
		zap.String("key", key),
	)

	metrics.CacheOperations.WithLabelValues("delete", "success").Inc()
	metrics.CacheSize.Set(float64(len(s.cache.GetAll())))

	c.Status(http.StatusOK)
}

func (s *CacheServer) Run(addr string) error {
	s.logger.Info("Starting server",
		zap.String("address", addr),
		zap.Bool("tls_enabled", s.config.CertFile != "" && s.config.KeyFile != ""),
	)

	if s.config.CertFile != "" && s.config.KeyFile != "" {
		return s.router.RunTLS(addr, s.config.CertFile, s.config.KeyFile)
	}
	return s.router.Run(addr)
}
