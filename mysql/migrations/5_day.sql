-- +migrate Up
CREATE TABLE `day` (
	`id`			int UNSIGNED NOT NULL AUTO_INCREMENT,
	`canceled_class_id` 			int NOT NULL,
	`day` 			varchar(255) NOt NULL,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE day;
