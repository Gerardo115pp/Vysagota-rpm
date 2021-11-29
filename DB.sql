DROP DATABASE IF EXISTS `interviews`;


CREATE DATABASE `interviews` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci;

USE `interviews`;


DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) UNIQUE NOT NULL,
  `hsh` VARCHAR(40) NOT NULL,
  `data_folder` VARCHAR(255) NOT NULL, 
  PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;


DROP TABLE IF EXISTS `interviews`;
CREATE TABLE `interviews` (
  `id` VARCHAR(40) NOT NULL,
  `was_finished` BOOLEAN DEFAULT FALSE,
  `created` DATETIME DEFAULT NOW(),
  `created_by` INT UNSIGNED NOT NULL,
  `applyed_to` VARCHAR(40) NOT NULL,
  `profile` INT UNSIGNED,
  `tokenized` BOOLEAN DEFAULT FALSE,
  CONSTRAINT `created_by_fk` FOREIGN KEY(`created_by`) REFERENCES users(`id`),
  PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;


DROP TABLE IF EXISTS `interviewees`;
CREATE TABLE `interviewees` (
  `id` VARCHAR(40) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `interviewer` INT UNSIGNED NOT NULL,
  CONSTRAINT `interviewer_fk` FOREIGN KEY(`interviewer`) REFERENCES users(`id`),
  PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;

ALTER TABLE `interviews` ADD CONSTRAINT `interviewee_fk` FOREIGN KEY(`applyed_to`) REFERENCES interviewees(`id`);


DROP TABLE IF EXISTS `interviewstests`;
CREATE TABLE `interviewstests` (
  `test` INT UNSIGNED NOT NULL,
  `interview` VARCHAR(40) NOT NULL,
  CONSTRAINT `interview_IT_fk` FOREIGN KEY(`interview`) REFERENCES interviews(`id`)
) ENGINE = InnoDB;


DROP TABLE IF EXISTS `tests`;
CREATE TABLE `tests` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `short_name`  VARCHAR(15) NOT NULL UNIQUE,
  `full_name` VARCHAR(255) NOT NULL,
  `length` INT UNSIGNED NOT NULL,
  `path` VARCHAR(255) NOT NULL,
  PRIMARY KEY `pk_id`(`id`)
) ENGINE = InnoDB;

ALTER TABLE `interviewstests` ADD CONSTRAINT `test_IT_fk` FOREIGN KEY(`test`) REFERENCES tests(`id`);

  DROP TABLE IF EXISTS `stats`;
  CREATE TABLE `stats` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255),
    PRIMARY KEY `pk_id`(`id`)
  ) ENGINE = InnoDB;


  DROP TABLE IF EXISTS `testsstats`;
  CREATE TABLE `testsstats` (
    `stat` INT UNSIGNED NOT NULL,
    `test` INT UNSIGNED NOT NULL,
    CONSTRAINT `fk_testsstats_stat` FOREIGN KEY(`stat`) REFERENCES stats(`id`),
    CONSTRAINT `fk_testsstats_test` FOREIGN KEY(`test`) REFERENCES tests(`id`)
  ) ENGINE = InnoDB;


DROP TABLE IF EXISTS `profiles`;
CREATE TABLE `profiles` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(40),
  `path_to_file` VARCHAR(255) NOT NULL,
  `creator` INT UNSIGNED NOT NULL,
  PRIMARY KEY `pk_id`(`id`),
  CONSTRAINT `fk_creator` FOREIGN KEY (`creator`) REFERENCES users(`id`) 
) ENGINE = InnoDB;

ALTER TABLE `interviews` ADD CONSTRAINT `fk_profile` FOREIGN KEY (`profile`) REFERENCES profiles(`id`);