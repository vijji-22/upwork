CREATE TABLE IF NOT EXISTS service (
    id int(10) unsigned NOT NULL AUTO_INCREMENT,
    name varchar(255) NOT NULL,
    short_description varchar(255) NOT NULL,
    long_description varchar(255) NOT NULL,
    image varchar(255) NOT NULL,
    max_amount int(11) NOT NULL,
    min_amount int(11) NOT NULL,

    temple_id int(10) unsigned NOT NULL,
    service_type_id int(10) unsigned NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name),
    FOREIGN KEY (temple_id) REFERENCES temple(id),
    FOREIGN KEY (service_type_id) REFERENCES service_type(id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;