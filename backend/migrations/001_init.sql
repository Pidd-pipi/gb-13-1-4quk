-- campusbooks 数据库表结构（由 GORM AutoMigrate 自动生成，此文件用于文档参考与手工初始化）
-- 启动时后端 internal/database/database.go 中的 Migrate 会按 model 定义自动建表/迁移。

-- users 用户表
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  student_no VARCHAR(32) NOT NULL,
  email VARCHAR(128) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(64) DEFAULT '',
  department VARCHAR(64) DEFAULT '',
  campus VARCHAR(64) DEFAULT '',
  contact VARCHAR(64) DEFAULT '',
  avatar_url VARCHAR(255) DEFAULT '',
  role VARCHAR(16) NOT NULL DEFAULT 'student',
  email_verified TINYINT(1) NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'active',
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  UNIQUE KEY uni_users_student_no (student_no),
  UNIQUE KEY uni_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- books 书籍表
CREATE TABLE IF NOT EXISTS books (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  seller_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(128) NOT NULL,
  author VARCHAR(64) DEFAULT '',
  isbn VARCHAR(32) DEFAULT '',
  course_name VARCHAR(128) DEFAULT '',
  original_price DECIMAL(10,2) DEFAULT 0,
  price DECIMAL(10,2) NOT NULL,
  `condition` VARCHAR(16) NOT NULL,
  subject_category VARCHAR(16) NOT NULL,
  trade_type VARCHAR(16) NOT NULL,
  campus VARCHAR(64) DEFAULT '',
  description TEXT,
  images JSON,
  status VARCHAR(16) NOT NULL DEFAULT 'on_sale',
  reserved_by BIGINT UNSIGNED DEFAULT 0,
  reserved_at DATETIME(3) NULL,
  borrowable TINYINT(1) NOT NULL DEFAULT 0,
  borrow_duration INT NOT NULL DEFAULT 0,
  view_count INT DEFAULT 0,
  favorite_count INT DEFAULT 0,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  KEY idx_books_seller (seller_id),
  KEY idx_books_title (title),
  KEY idx_books_status (status),
  KEY idx_books_subject (subject_category),
  KEY idx_books_condition (`condition`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- wishes 求购信息表
CREATE TABLE IF NOT EXISTS wishes (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  book_title VARCHAR(128) NOT NULL,
  author VARCHAR(64) DEFAULT '',
  isbn VARCHAR(32) DEFAULT '',
  expected_price DECIMAL(10,2) DEFAULT 0,
  condition_requirement VARCHAR(16) DEFAULT '',
  subject_category VARCHAR(16) NOT NULL,
  description TEXT,
  status VARCHAR(16) NOT NULL DEFAULT 'open',
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  KEY idx_wishes_user (user_id),
  KEY idx_wishes_status (status),
  KEY idx_wishes_subject (subject_category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- conversations 交易会话表
CREATE TABLE IF NOT EXISTS conversations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  book_id BIGINT UNSIGNED DEFAULT 0,
  wish_id BIGINT UNSIGNED DEFAULT 0,
  buyer_id BIGINT UNSIGNED NOT NULL,
  seller_id BIGINT UNSIGNED NOT NULL,
  last_message VARCHAR(2000) DEFAULT '',
  last_message_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  KEY idx_convs_buyer (buyer_id),
  KEY idx_convs_seller (seller_id),
  KEY idx_convs_book (book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- messages 消息表
CREATE TABLE IF NOT EXISTS messages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  conversation_id BIGINT UNSIGNED NOT NULL,
  sender_id BIGINT UNSIGNED NOT NULL,
  content VARCHAR(2000) DEFAULT '',
  image_url VARCHAR(255) DEFAULT '',
  is_read TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME(3) NULL,
  KEY idx_msgs_conv (conversation_id),
  KEY idx_msgs_sender (sender_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- evaluations 交易评价表
CREATE TABLE IF NOT EXISTS evaluations (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  from_user_id BIGINT UNSIGNED NOT NULL,
  to_user_id BIGINT UNSIGNED NOT NULL,
  book_id BIGINT UNSIGNED NOT NULL,
  type VARCHAR(16) NOT NULL,
  content VARCHAR(500) DEFAULT '',
  created_at DATETIME(3) NULL,
  UNIQUE KEY uni_eval_pair (from_user_id, to_user_id, book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- favorites 收藏表
CREATE TABLE IF NOT EXISTS favorites (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  book_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME(3) NULL,
  UNIQUE KEY uni_fav_user_book (user_id, book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- browse_histories 浏览历史表
CREATE TABLE IF NOT EXISTS browse_histories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  book_id BIGINT UNSIGNED NOT NULL,
  viewed_at DATETIME(3) NOT NULL,
  KEY idx_hist_user_time (user_id, viewed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- borrows 短借申请表（pending -> approved -> returning -> returned，另有 rejected）
CREATE TABLE IF NOT EXISTS borrows (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  book_id BIGINT UNSIGNED NOT NULL,
  lender_id BIGINT UNSIGNED NOT NULL,
  borrower_id BIGINT UNSIGNED NOT NULL,
  duration INT NOT NULL DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  due_at DATETIME(3) NULL,
  approved_at DATETIME(3) NULL,
  returned_at DATETIME(3) NULL,
  confirmed_at DATETIME(3) NULL,
  reject_reason VARCHAR(255) DEFAULT '',
  reminded_at DATETIME(3) NULL,
  created_at DATETIME(3) NULL,
  updated_at DATETIME(3) NULL,
  KEY idx_borrows_book (book_id),
  KEY idx_borrows_lender (lender_id),
  KEY idx_borrows_borrower (borrower_id),
  KEY idx_borrows_status (status),
  KEY idx_borrows_due (due_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- audit_logs 操作审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED DEFAULT 0,
  action VARCHAR(64) NOT NULL,
  resource_type VARCHAR(32) DEFAULT '',
  resource_id BIGINT UNSIGNED DEFAULT 0,
  detail VARCHAR(500) DEFAULT '',
  ip VARCHAR(64) DEFAULT '',
  request_id VARCHAR(64) DEFAULT '',
  created_at DATETIME(3) NULL,
  KEY idx_audit_user (user_id),
  KEY idx_audit_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
