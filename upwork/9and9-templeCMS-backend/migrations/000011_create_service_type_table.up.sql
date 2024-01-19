CREATE TABLE IF NOT EXISTS service_type (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    name varchar(64) NOT NULL,
    short_description varchar(255) NOT NULL,
    long_description varchar(255) NOT NULL,
    image varchar(255) NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;