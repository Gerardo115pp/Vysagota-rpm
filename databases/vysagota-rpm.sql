
DROP DATABASE `vysagota-rpm`;
CREATE DATABASE `vysagota-rpm`;

USE `vysagota-rpm`;

CREATE TABLE `pacients`(
    `uuid` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `password` CHAR(64) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(20) NOT NULL,
    `tutor_email` VARCHAR(255) NOT NULL,
    `birthday` DATE DEFAULT NULL,
    `appointments_count` INT UNSIGNED DEFAULT 0,
    `last_seen` DATETIME DEFAULT NULL,
    `created_at` DATETIME DEFAULT NOW(),
    `gender` enum('M','F') DEFAULT 'M'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `doctors` (
    `uuid` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `password` CHAR(64) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `registration_date` DATETIME DEFAULT NOW()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `pacientsDoctors` (
    `pacient_uuid` BIGINT UNSIGNED NOT NULL,
    `doctor_uuid` BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (`pacient_uuid`, `doctor_uuid`),
    CONSTRAINT `pd_pacient_fk` FOREIGN KEY (`pacient_uuid`) REFERENCES `pacients`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `pd_doctor_fk` FOREIGN KEY (`doctor_uuid`) REFERENCES `doctors`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionarys` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `path_to_file` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `length` INT UNSIGNED NOT NULL,
    `created_at` DATETIME DEFAULT NOW(),
    `short_name` VARCHAR(255) NOT NULL UNIQUE,
    `creator` BIGINT UNSIGNED NOT NULL,
    CONSTRAINT `q_creator_fk` FOREIGN KEY (`creator`) REFERENCES `doctors`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `stats`;
CREATE TABLE `stats` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionaryMeasures` (
    `questionary` INT UNSIGNED NOT NULL,
    `stat` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`questionary`, `stat`),
    CONSTRAINT `qm_questionary_fk` FOREIGN KEY (`questionary`) REFERENCES `questionarys`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `qm_stat_fk` FOREIGN KEY (`stat`) REFERENCES `stats`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionarysSet` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `created_at` DATETIME DEFAULT NOW(),
    `path_to_file` VARCHAR(255) NOT NULL,
    `creator` BIGINT UNSIGNED NOT NULL,
    CONSTRAINT `qs_creator_fk` FOREIGN KEY (`creator`) REFERENCES `doctors`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `questionarysSetQuestionarys` (
    `questionarySet` INT UNSIGNED NOT NULL,
    `questionary` INT UNSIGNED NOT NULL,
    PRIMARY KEY (`questionarySet`, `questionary`),
    CONSTRAINT `qsq_questionarySet_fk` FOREIGN KEY (`questionarySet`) REFERENCES `questionarysSet`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `qsq_questionary_fk` FOREIGN KEY (`questionary`) REFERENCES `questionarys`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `revisions` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `created_at` DATETIME DEFAULT NOW(),
    `starts_at` DATETIME DEFAULT NULL,
    `expires_at` DATETIME DEFAULT NULL,
    `was_finished` TINYINT(1) DEFAULT 0,
    `finished_at` DATETIME DEFAULT NULL,
    `was_canceled` TINYINT(1) DEFAULT 0,
    `show_to_pacient` TINYINT(1) DEFAULT 1,
    `questionary` INT UNSIGNED NOT NULL,
    `pacient` BIGINT UNSIGNED NOT NULL,
    `creator` BIGINT UNSIGNED NOT NULL,
    CONSTRAINT `r_creator_fk` FOREIGN KEY (`creator`) REFERENCES `doctors`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `r_questionary_fk` FOREIGN KEY (`questionary`) REFERENCES `questionarys`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `r_pacient_fk` FOREIGN KEY (`pacient`) REFERENCES `pacients`(`uuid`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Data Dictionary:
-- FieldName| Constraint | Type | Description


DROP TRIGGER IF EXISTS `clean_stats_name`;
CREATE TRIGGER `clean_stats_name` BEFORE INSERT ON `stats` FOR EACH ROW
SET NEW.`name` = UPPER(NEW.`name`);

DROP TRIGGER IF EXISTS `hash_doctor_password`;
CREATE TRIGGER `hash_doctor_password` BEFORE INSERT ON `doctors` FOR EACH ROW
SET NEW.`password` = SHA2(NEW.`password`,256);

DROP TRIGGER IF EXISTS `hash_pacient_password`;
CREATE TRIGGER `hash_pacient_password` BEFORE INSERT ON `pacients` FOR EACH ROW
SET NEW.`password` = SHA2(NEW.`password`,256);

DROP TRIGGER IF EXISTS `rehash_pacient_password`;
CREATE TRIGGER `rehash_pacient_password` BEFORE UPDATE ON `pacients` FOR EACH ROW
SET NEW.`password` = SHA2(NEW.`password`,256);

DROP TRIGGER IF EXISTS `complete_first_doctor`;
DELIMITER //
CREATE TRIGGER `complete_first_doctor` BEFORE INSERT ON `doctors` FOR EACH ROW
PRECEDES `hash_doctor_password`
BEGIN
    SET @doctor_count = (SELECT COUNT(*) FROM `doctors`);
    IF @doctor_count = 0 AND NEW.`name` = '' THEN
        SET NEW.`name` = 'vysagota-admin';
    ELSE
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Empty doctor name is only allowed for the first doctor';
    END IF;
END//
DELIMITER ;

-- insert main doctor
INSERT INTO `doctors` (`uuid`, `name`, `password`, `email`) VALUES (1, '', 'heroku', 'theronin115@gmail.com');

