package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// User 用户结构体
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// Response 通用响应结构体
type Response struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Error     string      `json:"error,omitempty"`
}

// APIInfo API信息结构体
type APIInfo struct {
	Service     string `json:"service"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// HealthStatus 健康状态结构体
type HealthStatus struct {
	Status    string  `json:"status"`
	Uptime    float64 `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

// 动态路由规则结构体
type DynamicRoute struct {
	ID     int    `json:"id"`         // 路由ID
	Path   string `json:"path"`       // 匹配的路径（如 /foo）
	Target string `json:"target"`     // 目标服务地址（仅做展示，不实际转发）
	Desc   string `json:"desc"`       // 路由描述
}

var (
	// 模拟用户数据
	users = []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
		{ID: 3, Name: "王五", Email: "wangwu@example.com"},
	}
	startTime = time.Now()
	// 动态路由规则存储（内存）
	dynamicRoutes   = make([]DynamicRoute, 0) // 路由规则列表
	dynamicRouteSeq = 1                       // 自增ID
)

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.New()

	// 使用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 配置CORS中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 自定义中间件：记录请求信息
	r.Use(func(c *gin.Context) {
		log.Printf("%s - %s %s", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.Path)
		log.Printf("请求头: %v", c.Request.Header)
		c.Next()
	})

	// 根路径路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			Message:   "欢迎使用Traefik演示应用！",
			Timestamp: time.Now(),
			Data: gin.H{
				"headers":  c.Request.Header,
				"hostname": c.Request.Host,
				"remoteIP": c.ClientIP(),
			},
		})
	})

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthStatus{
			Status:    "healthy",
			Uptime:    time.Since(startTime).Seconds(),
			Timestamp: time.Now(),
		})
	})

	// API信息端点
	r.GET("/api/info", func(c *gin.Context) {
		env := os.Getenv("NODE_ENV")
		if env == "" {
			env = "development"
		}

		c.JSON(http.StatusOK, APIInfo{
			Service:     "custom-app",
			Version:     "1.0.0",
			Environment: env,
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 获取用户列表
		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Data: gin.H{
					"users": users,
					"count": len(users),
				},
				Timestamp: time.Now(),
			})
		})

		// 获取特定用户
		api.GET("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			// 这里简化处理，实际应该解析ID
			if id == "1" {
				c.JSON(http.StatusOK, Response{
					Data: gin.H{
						"user": users[0],
					},
					Timestamp: time.Now(),
				})
			} else if id == "2" {
				c.JSON(http.StatusOK, Response{
					Data: gin.H{
						"user": users[1],
					},
					Timestamp: time.Now(),
				})
			} else if id == "3" {
				c.JSON(http.StatusOK, Response{
					Data: gin.H{
						"user": users[2],
					},
					Timestamp: time.Now(),
				})
			} else {
				c.JSON(http.StatusNotFound, Response{
					Error:     "用户未找到",
					Timestamp: time.Now(),
				})
			}
		})

		// 创建新用户
		api.POST("/users", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, Response{
					Error:     "请求数据格式错误",
					Timestamp: time.Now(),
				})
				return
			}

			// 验证必填字段
			if newUser.Name == "" || newUser.Email == "" {
				c.JSON(http.StatusBadRequest, Response{
					Error:     "姓名和邮箱是必需的",
					Timestamp: time.Now(),
				})
				return
			}

			// 生成新用户ID（简化处理）
			newUser.ID = len(users) + 1
			newUser.CreatedAt = time.Now()

			c.JSON(http.StatusCreated, Response{
				Message: "用户创建成功",
				Data: gin.H{
					"user": newUser,
				},
				Timestamp: time.Now(),
			})
		})
	}

	// 动态路由管理API
	r.GET("/routes", getRoutes)      // 查询所有动态路由
	r.POST("/routes", addRoute)      // 添加新动态路由
	r.DELETE("/routes/:id", deleteRoute) // 删除动态路由

	// 404处理（动态路由兜底）
	r.NoRoute(dynamicRouteHandler)

	// 启动服务器
	port := ":3000"
	log.Printf("🚀 服务器运行在 http://0.0.0.0%s", port)
	log.Printf("📊 健康检查: http://0.0.0.0%s/health", port)
	log.Printf("📋 API信息: http://0.0.0.0%s/api/info", port)
	log.Printf("👥 用户列表: http://0.0.0.0%s/api/users", port)

	if err := r.Run(port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

// 获取所有动态路由
func getRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Data:      dynamicRoutes,
		Timestamp: time.Now(),
	})
}

// 添加新动态路由
func addRoute(c *gin.Context) {
	var req struct {
		Path   string `json:"path"`   // 必填
		Target string `json:"target"` // 必填
		Desc   string `json:"desc"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Path == "" || req.Target == "" {
		c.JSON(http.StatusBadRequest, Response{
			Error:     "参数错误，path 和 target 必填",
			Timestamp: time.Now(),
		})
		return
	}
	// 检查是否已存在
	for _, r := range dynamicRoutes {
		if r.Path == req.Path {
			c.JSON(http.StatusConflict, Response{
				Error:     "该 path 已存在",
				Timestamp: time.Now(),
			})
			return
		}
	}
	newRoute := DynamicRoute{
		ID:     dynamicRouteSeq,
		Path:   req.Path,
		Target: req.Target,
		Desc:   req.Desc,
	}
	dynamicRouteSeq++
	dynamicRoutes = append(dynamicRoutes, newRoute)
	c.JSON(http.StatusCreated, Response{
		Message:   "路由添加成功",
		Data:      newRoute,
		Timestamp: time.Now(),
	})
}

// 删除动态路由
func deleteRoute(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error:     "ID格式错误",
			Timestamp: time.Now(),
		})
		return
	}
	idx := -1
	for i, r := range dynamicRoutes {
		if r.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		c.JSON(http.StatusNotFound, Response{
			Error:     "未找到该路由",
			Timestamp: time.Now(),
		})
		return
	}
	// 删除该路由
	dynamicRoutes = append(dynamicRoutes[:idx], dynamicRoutes[idx+1:]...)
	c.JSON(http.StatusOK, Response{
		Message:   "路由已删除",
		Timestamp: time.Now(),
	})
}

// 动态路由处理器（未命中路由时调用）
func dynamicRouteHandler(c *gin.Context) {
	path := c.Request.URL.Path
	for _, r := range dynamicRoutes {
		if r.Path == path {
			// 假装转发到 r.Target，实际只返回模拟响应
			c.JSON(http.StatusOK, Response{
				Message:   "已匹配动态路由，模拟转发到目标服务",
				Data: gin.H{
					"path":   r.Path,
					"target": r.Target,
					"desc":   r.Desc,
				},
				Timestamp: time.Now(),
			})
			return
		}
	}
	// 未匹配到动态路由，返回404
	c.JSON(http.StatusNotFound, Response{
		Error:     "页面未找到（含动态路由）",
		Data: gin.H{
			"path":   path,
			"method": c.Request.Method,
		},
		Timestamp: time.Now(),
	})
} 