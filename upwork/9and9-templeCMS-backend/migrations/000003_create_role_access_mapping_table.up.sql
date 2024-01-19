CREATE TABLE IF NOT EXISTS role_access_mapping (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		role_id int(11) unsigned NOT NULL,
		access_id int(11) unsigned NOT NULL,
		project varchar(255),
		
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (role_id) REFERENCES role(id),
		FOREIGN KEY (access_id) REFERENCES access(id),
		UNIQUE KEY (role_id, access_id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;