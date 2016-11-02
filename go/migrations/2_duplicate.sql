-- +migrate Up
ALTER TABLE `kyuko_data` ADD UNIQUE INDEX `class` (`day` ASC, `class_name` ASC, `instructor` ASC);

-- +migrate Down
ALTER TABLE `kyuko_data` DROP INDEX `class`;
