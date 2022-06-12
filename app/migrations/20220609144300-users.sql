
-- +migrate Up
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) DEFAULT NULL,
  `name` varchar(255)  NOT NULL,
  `username` varchar(20)  NOT NULL,
  `email` varchar(255)  NOT NULL,
  `type` varchar(255) NOT NULL DEFAULT 'student',
  `password` varchar(255)  NOT NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE `users`;
