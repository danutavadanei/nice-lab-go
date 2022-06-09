
-- +migrate Up
CREATE TABLE `auth_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(255) NOT NULL,
  `expire_at` timestamp NOT NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8;

-- +migrate Down
DROP TABLE `auth_tokens`;
