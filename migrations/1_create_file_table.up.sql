CREATE TABLE IF NOT EXISTS files
(
    id              VARCHAR(36) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    meta_data       JSON  DEFAULT NULL,
    owner_id        VARCHAR(100) NOT NULL,
    bucket_path     VARCHAR(255) NOT NULL,
    provider        VARCHAR(20) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at      TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

