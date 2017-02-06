-- +migrate Up
CREATE TABLE `reason` (
	`id`					int UNSIGNED NOT NULL AUTO_INCREMENT,
	`canceled_class_id`	int NOT NULL,
	`reason` 				varchar(255) NOT NULL,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE reason;
