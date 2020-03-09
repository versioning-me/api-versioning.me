
-- +migrate Up
CREATE TABLE IF NOT EXISTS `file` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `detail` text DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;

CREATE TABLE IF NOT EXISTS `version` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `file_id` bigint(20) DEFAULT NULL,
  `version_num` int(11) DEFAULT NULL,
  `ext` varchar(10) DEFAULT NULL,
  `mime` varchar(255) DEFAULT NULL,
  `size` int DEFAULT NULL,
  `url` varchar(255) DEFAULT NULL,
  `history` text DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY(`id`),
  FOREIGN KEY(`file_id`) REFERENCES `file`(`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=1 ;


-- +migrate Down
DROP TABLE IF EXISTS version;
DROP TABLE IF EXISTS uploaded_file;
