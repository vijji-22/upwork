CREATE TABLE IF NOT EXISTS temple (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    logo VARCHAR(255),
    hero_image VARCHAR(255),

    open_time TIME,
    close_time TIME,

    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(64) NOT NULL,
    state VARCHAR(64) NOT NULL,
    zip VARCHAR(6),
    country VARCHAR(32) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;