-- docker/init/01_schema.sql
-- Library Management System Schema
-- Note: ID uses ULID (26 characters)

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) PRIMARY KEY,  -- ULID
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status TINYINT UNSIGNED NOT NULL DEFAULT 1,  -- 1:ACTIVE, 2:SUSPENDED, 3:INACTIVE
    max_loans INT NOT NULL DEFAULT 5,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
    id VARCHAR(26) PRIMARY KEY,  -- ULID
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    -- Availability is derived from the loans table
);

CREATE TABLE IF NOT EXISTS loans (
    id VARCHAR(26) PRIMARY KEY,  -- ULID
    user_id VARCHAR(26) NOT NULL,  -- ULID
    book_id VARCHAR(26) NOT NULL,  -- ULID
    borrowed_at DATETIME NOT NULL,
    due_date DATETIME NOT NULL,
    returned_at DATETIME DEFAULT NULL,  -- NULL = active loan (borrowed)
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (book_id) REFERENCES books(id),
    INDEX idx_user_active (user_id, returned_at),
    INDEX idx_book_active (book_id, returned_at)
);