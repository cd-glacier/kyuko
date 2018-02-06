-- +migrate Up
CREATE TABLE `canceled_class` (
	`id`			int UNSIGNED NOT NULL AUTO_INCREMENT,
	`canceled` 	int NOT NULL,
	`place` 		int NOT NULL,
	`week` 			int NOT NULL,
	`period`		int NOT NULL,
	`year` 			int NOT NULL,
	`season` 			varchar(255) NOt NULL,
	`class_name` 	varchar(255) NOT NULL,
	`instructor` 	varchar(255) NOT NULL,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE canceled_class;
