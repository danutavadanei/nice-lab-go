
-- +migrate Up
CREATE TABLE `labs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) DEFAULT NULL,
  `name` varchar(255)  NOT NULL,
  `type` varchar(255) NOT NULL,
  `hostname` varchar(1024) NOT NULL,
  `instance_id` varchar(1024) NOT NULL,
  `available` tinyint DEFAULT 1 NOT NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE `labs`;
