# Copied from Mathias Bynens PHP URL Shortener: https://github.com/mathiasbynens/php-url-shortener
SET NAMES utf8;

# Why you should use `utf8mb4` instead of `utf8`: http://mathiasbynens.be/notes/mysql-utf8mb4
DROP TABLE IF EXISTS `shortr`;
CREATE TABLE `shortr` (
	`slug` varchar(14) collate utf8mb4_unicode_ci NOT NULL,
	`url` varchar(620) collate utf8mb4_unicode_ci NOT NULL,
	`date` datetime NOT NULL,
	`hits` bigint(20) NOT NULL default '0',
	`ip` varchar(40) NOT NULL,
	PRIMARY KEY (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Used for the URL shortener';

INSERT INTO `shortr` VALUES ('test', 'http://domain.tld', NOW(), 1, '127.0.0.1');

# Why you should use `utf8mb4` instead of `utf8`: http://mathiasbynens.be/notes/mysql-utf8mb4
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
	`id` smallint NOT NULL AUTO_INCREMENT,
  `user` varchar(64) NOT NULL collate utf8mb4_unicode_ci,
  `password` varchar(64) NOT NULL collate utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Authenticate shortr user';

INSERT INTO `users` VALUES ('1', 'admin', '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918');
