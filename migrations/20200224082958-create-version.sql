
-- +migrate Up
CREATE TABLE IF NOT EXISTS `version` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS `uploaded_file` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
-- `version_id` bigint(20) DEFAULT NULL,
  `version_name` varchar(255) DEFAULT NULL,
  `uuid` varchar(255) DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `ext` varchar(10) DEFAULT NULL,
  `mime` varchar(255) DEFAULT NULL,
  `size` int,
  `url` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY(`id`)
--  FOREIGN KEY(`version_id`) REFERENCES `version`(`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;



-- +migrate Down
DROP TABLE IF EXISTS version;
DROP TABLE IF EXISTS uploaded_file;
