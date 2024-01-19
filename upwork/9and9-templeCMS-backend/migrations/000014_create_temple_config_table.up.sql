CREATE TABLE IF NOT EXISTS temple_config (
    id int(11) unsigned NOT NULL AUTO_INCREMENT,
    parent_id int(11) unsigned,
    group_id int(11) unsigned NOT NULL,
    node_type_id int(11) unsigned NOT NULL,
    value varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;