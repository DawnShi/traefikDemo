#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Traefik演示应用 - Python版本
使用Flask框架构建RESTful API
"""

from flask import Flask, request, jsonify
from datetime import datetime
import time
import os
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 创建Flask应用
app = Flask(__name__)

# 模拟用户数据
users = [
    {"id": 1, "name": "张三", "email": "zhangsan@example.com"},
    {"id": 2, "name": "李四", "email": "lisi@example.com"},
    {"id": 3, "name": "王五", "email": "wangwu@example.com"}
]

# 应用启动时间
start_time = time.time()

@app.before_request
def log_request_info():
    """记录请求信息的中间件"""
    logger.info(f"{request.method} {request.path}")
    logger.info(f"请求头: {dict(request.headers)}")

@app.route('/')
def index():
    """根路径 - 返回欢迎信息"""
    return jsonify({
        "message": "欢迎使用Traefik演示应用！",
        "timestamp": datetime.now().isoformat(),
        "headers": dict(request.headers),
        "hostname": request.host,
        "remote_ip": request.remote_addr
    })

@app.route('/health')
def health():
    """健康检查端点"""
    uptime = time.time() - start_time
    return jsonify({
        "status": "healthy",
        "uptime": uptime,
        "timestamp": datetime.now().isoformat()
    })

@app.route('/api/info')
def api_info():
    """API信息端点"""
    env = os.getenv('NODE_ENV', 'development')
    return jsonify({
        "service": "custom-app",
        "version": "1.0.0",
        "environment": env
    })

@app.route('/api/users', methods=['GET'])
def get_users():
    """获取用户列表"""
    return jsonify({
        "users": users,
        "count": len(users),
        "timestamp": datetime.now().isoformat()
    })

@app.route('/api/users/<int:user_id>', methods=['GET'])
def get_user(user_id):
    """获取特定用户"""
    user = next((u for u in users if u["id"] == user_id), None)
    
    if user:
        return jsonify({
            "user": user,
            "timestamp": datetime.now().isoformat()
        })
    else:
        return jsonify({
            "error": "用户未找到",
            "user_id": user_id,
            "timestamp": datetime.now().isoformat()
        }), 404

@app.route('/api/users', methods=['POST'])
def create_user():
    """创建新用户"""
    data = request.get_json()
    
    if not data:
        return jsonify({
            "error": "请求数据格式错误",
            "timestamp": datetime.now().isoformat()
        }), 400
    
    name = data.get('name')
    email = data.get('email')
    
    if not name or not email:
        return jsonify({
            "error": "姓名和邮箱是必需的",
            "timestamp": datetime.now().isoformat()
        }), 400
    
    # 生成新用户ID
    new_id = max(u["id"] for u in users) + 1 if users else 1
    
    new_user = {
        "id": new_id,
        "name": name,
        "email": email,
        "created_at": datetime.now().isoformat()
    }
    
    users.append(new_user)
    
    return jsonify({
        "message": "用户创建成功",
        "user": new_user,
        "timestamp": datetime.now().isoformat()
    }), 201

@app.errorhandler(404)
def not_found(error):
    """404错误处理"""
    return jsonify({
        "error": "页面未找到",
        "path": request.path,
        "method": request.method,
        "timestamp": datetime.now().isoformat()
    }), 404

@app.errorhandler(500)
def internal_error(error):
    """500错误处理"""
    return jsonify({
        "error": "服务器内部错误",
        "timestamp": datetime.now().isoformat()
    }), 500

if __name__ == '__main__':
    port = 3000
    logger.info(f"🚀 服务器运行在 http://0.0.0.0:{port}")
    logger.info(f"📊 健康检查: http://0.0.0.0:{port}/health")
    logger.info(f"📋 API信息: http://0.0.0.0:{port}/api/info")
    logger.info(f"👥 用户列表: http://0.0.0.0:{port}/api/users")
    
    # 启动Flask应用
    app.run(host='0.0.0.0', port=port, debug=False) 