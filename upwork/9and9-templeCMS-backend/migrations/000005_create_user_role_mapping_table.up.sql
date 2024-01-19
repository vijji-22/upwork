CREATE TABLE IF NOT EXISTS user_role_mapping (
		id int(10) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) unsigned NOT NULL,
		role_id int(11) unsigned NOT NULL,
		reference_id int(11) unsigned,
		created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES `user`(id),
		FOREIGN KEY (role_id) REFERENCES role(id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;