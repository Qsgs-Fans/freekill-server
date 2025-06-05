-- 为了讨好goctl得给主键定义单开一行

-- 用户基本信息表
CREATE TABLE IF NOT EXISTS users (
  id BIGINT NOT NULL AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  password CHAR(64) NOT NULL COMMENT 'SHA-256哈希密码',
  salt CHAR(8) NOT NULL,
  avatar VARCHAR(255) NOT NULL,
  -- is_banned BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户登录历史表（包含设备信息）
CREATE TABLE IF NOT EXISTS user_login_history (
  id BIGINT NOT NULL AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  ip_address VARBINARY(16) NOT NULL COMMENT '4字节IPv4或16字节IPv6',
  device_uuid VARCHAR(64) NOT NULL,
  login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  login_success BOOLEAN DEFAULT TRUE,
  -- FOREIGN KEY (user_id) REFERENCES users(id),
  PRIMARY KEY (id),
  INDEX idx_user_id (user_id),
  INDEX idx_ip (ip_address),
  INDEX idx_device_uuid (device_uuid),
  INDEX idx_login_time (login_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 封禁IP表（二进制存储IP）
CREATE TABLE IF NOT EXISTS banned_ips (
  id BIGINT NOT NULL AUTO_INCREMENT,
  ip_address VARBINARY(16) NOT NULL UNIQUE COMMENT '4字节IPv4或16字节IPv6',
  ban_reason VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '封禁过期时间，NULL表示永久',
  PRIMARY KEY (id),
  INDEX idx_ip (ip_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 封禁设备表
CREATE TABLE IF NOT EXISTS banned_devices (
  id BIGINT NOT NULL AUTO_INCREMENT,
  device_uuid VARCHAR(64) NOT NULL UNIQUE,
  ban_reason VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '封禁过期时间，NULL表示永久',
  PRIMARY KEY (id),
  INDEX idx_device_uuid (device_uuid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户名黑名单表
CREATE TABLE IF NOT EXISTS username_blacklist (
  id BIGINT NOT NULL AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  ban_reason VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '封禁过期时间，NULL表示永久',
  PRIMARY KEY (id),
  INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户名白名单表
CREATE TABLE IF NOT EXISTS username_whitelist (
  id BIGINT NOT NULL AUTO_INCREMENT,
  username VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
