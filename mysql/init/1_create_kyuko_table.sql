use kyuko;
CREATE TABLE `kyuko_data` (
	`id`			int UNSIGNED NOT NULL AUTO_INCREMENT,
	`place` 		int NOT NULL,
	`week` 			int NOT NULL,
	`period`		int NOT NULL,
	`day` 			varchar(255) NOt NULL,
	`class_name` 	varchar(255) NOT NULL,
	`instructor` 	varchar(255) NOT NULL,
	`reason` 		varchar(255) NOT NULL,
	PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

