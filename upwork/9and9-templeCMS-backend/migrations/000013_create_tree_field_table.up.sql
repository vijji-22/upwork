CREATE TABLE IF NOT EXISTS tree_fields (
    id int(11) unsigned NOT NULL AUTO_INCREMENT,
    parent_id int(11) unsigned,
    config_id int(11) unsigned NOT NULL,

    name varchar(255) NOT NULL,
    type varchar(255) NOT NULL,
    required boolean NOT NULL DEFAULT 0,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (name, parent_id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;