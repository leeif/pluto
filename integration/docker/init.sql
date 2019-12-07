-- MySQL dump 10.13  Distrib 5.7.27, for Linux (x86_64)
--
-- Host: localhost    Database: pluto_server
-- ------------------------------------------------------
-- Server version	5.7.27

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `applications`
--

DROP TABLE IF EXISTS `applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `applications` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `webhook` varchar(255) NOT NULL,
  `default_role` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `applictaions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `applications`
--

LOCK TABLES `applications` WRITE;
/*!40000 ALTER TABLE `applications` DISABLE KEYS */;
INSERT INTO `applications` VALUES (1,'2019-12-07 19:56:50','2019-12-07 19:56:50',NULL,'pluto','',NULL),(2,'2019-12-07 21:47:35','2019-12-07 21:47:48',NULL,'test','',2);
/*!40000 ALTER TABLE `applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `device_apps`
--

DROP TABLE IF EXISTS `device_apps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `device_apps` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `device_id` varchar(255) NOT NULL,
  `app_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_device_apps_deleted_at` (`deleted_at`),
  KEY `device_id_app_id` (`device_id`,`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `device_apps`
--

LOCK TABLES `device_apps` WRITE;
/*!40000 ALTER TABLE `device_apps` DISABLE KEYS */;
INSERT INTO `device_apps` VALUES (1,'2019-12-07 21:47:18','2019-12-07 21:47:18',NULL,'xxxx','pluto'),(5,'2019-12-07 22:02:12','2019-12-07 22:02:12',NULL,'xxxxxx','test');
/*!40000 ALTER TABLE `device_apps` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migrations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `migrations`
--

LOCK TABLES `migrations` WRITE;
/*!40000 ALTER TABLE `migrations` DISABLE KEYS */;
INSERT INTO `migrations` VALUES (1,'2019-12-07 19:56:38','create_users_table'),(2,'2019-12-07 19:56:38','create_refresh_tokens_table'),(3,'2019-12-07 19:56:38','create_salts_table'),(4,'2019-12-07 19:56:38','create_applictaions_table'),(5,'2019-12-07 19:56:39','create_device_apps_table'),(6,'2019-12-07 19:56:39','create_history_operations_table'),(7,'2019-12-07 19:56:39','drop_history_operations_table'),(8,'2019-12-07 19:56:39','create_rbac_user_application_roles_table'),(9,'2019-12-07 19:56:39','create_rbac_roles_table'),(10,'2019-12-07 19:56:39','create_rbac_role_scopes_table'),(11,'2019-12-07 19:56:39','create_rbac_scopes_table'),(12,'2019-12-07 19:56:39','add_default_role_column_in_applications_table'),(13,'2019-12-07 19:56:39','add_scopes_column_in_refresh_token_table');
/*!40000 ALTER TABLE `migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_role_scopes`
--

DROP TABLE IF EXISTS `rbac_role_scopes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbac_role_scopes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `role_id` int(10) unsigned NOT NULL,
  `scope_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `role_id_scope_id` (`role_id`,`scope_id`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_role_scopes`
--

LOCK TABLES `rbac_role_scopes` WRITE;
/*!40000 ALTER TABLE `rbac_role_scopes` DISABLE KEYS */;
INSERT INTO `rbac_role_scopes` VALUES (1,'2019-12-07 19:56:51','2019-12-07 19:56:51',NULL,1,1),(2,'2019-12-07 21:49:38','2019-12-07 21:49:38',NULL,2,2);
/*!40000 ALTER TABLE `rbac_role_scopes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_roles`
--

DROP TABLE IF EXISTS `rbac_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbac_roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(20) NOT NULL,
  `app_id` int(10) unsigned NOT NULL,
  `default_scope` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `app_id` (`app_id`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_roles`
--

LOCK TABLES `rbac_roles` WRITE;
/*!40000 ALTER TABLE `rbac_roles` DISABLE KEYS */;
INSERT INTO `rbac_roles` VALUES (1,'2019-12-07 19:56:50','2019-12-07 22:02:03',NULL,'admin',1,1),(2,'2019-12-07 21:47:45','2019-12-07 22:01:56',NULL,'user',2,2);
/*!40000 ALTER TABLE `rbac_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_scopes`
--

DROP TABLE IF EXISTS `rbac_scopes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbac_scopes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(20) NOT NULL,
  `app_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `app_id` (`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_scopes`
--

LOCK TABLES `rbac_scopes` WRITE;
/*!40000 ALTER TABLE `rbac_scopes` DISABLE KEYS */;
INSERT INTO `rbac_scopes` VALUES (1,'2019-12-07 19:56:50','2019-12-07 19:56:50',NULL,'pluto.admin',1),(2,'2019-12-07 21:49:28','2019-12-07 21:49:28',NULL,'test.user_test',2);
/*!40000 ALTER TABLE `rbac_scopes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_user_application_roles`
--

DROP TABLE IF EXISTS `rbac_user_application_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbac_user_application_roles` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `app_id` int(10) unsigned NOT NULL,
  `role_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id_app_id` (`user_id`,`app_id`),
  KEY `idx_user_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_user_application_roles`
--

LOCK TABLES `rbac_user_application_roles` WRITE;
/*!40000 ALTER TABLE `rbac_user_application_roles` DISABLE KEYS */;
INSERT INTO `rbac_user_application_roles` VALUES (1,'2019-12-07 19:56:51','2019-12-07 22:02:03',NULL,1,1,1);
/*!40000 ALTER TABLE `rbac_user_application_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refresh_tokens`
--

DROP TABLE IF EXISTS `refresh_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `refresh_tokens` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `refresh_token` varchar(255) NOT NULL,
  `device_app_id` int(10) unsigned NOT NULL,
  `scopes` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `refresh_tokens_deleted_at` (`deleted_at`),
  KEY `user_id_refresh_token` (`user_id`,`refresh_token`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
INSERT INTO `refresh_tokens` VALUES (1,'2019-12-07 21:47:18','2019-12-07 21:47:18',NULL,1,'832875c74f4f88ce2f350a7bbe093c2f21520619',1,'pluto.admin'),(2,'2019-12-07 22:02:12','2019-12-07 22:02:12',NULL,1,'17f48e57b4d53d24b55a953c166bebabf95ecf25',5,'test.user_test');
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `salts`
--

DROP TABLE IF EXISTS `salts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `salts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `salt` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `salts_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `salts`
--

LOCK TABLES `salts` WRITE;
/*!40000 ALTER TABLE `salts` DISABLE KEYS */;
INSERT INTO `salts` VALUES (1,'2019-12-07 19:56:51','2019-12-07 19:56:51',NULL,1,'WjJWbGEyeDVaamt5TmpFd1FHZHRZV2xzTG1OdmJRbTrS_lPubpqTvw==');
/*!40000 ALTER TABLE `salts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `mail` varchar(255) DEFAULT NULL,
  `name` varchar(60) NOT NULL,
  `gender` varchar(10) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `birthday` timestamp NULL DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `verified` tinyint(1) DEFAULT NULL,
  `login_type` varchar(10) NOT NULL,
  `identify_token` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `login_type_identify_token` (`login_type`,`identify_token`),
  KEY `users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2019-12-07 19:56:51','2019-12-07 19:56:51',NULL,'geeklyf92610@gmail.com','yifan.li',NULL,'nZCzLRuC7CCZEFUq9o5SN6kZY6b0a4oB_cgoNW9Qhqw=',NULL,'https://pluto-staging.oss-cn-hongkong.aliyuncs.com/avatar/15037733c65880a1.png',1,'mail','Z2Vla2x5ZjkyNjEwQGdtYWlsLmNvbQ');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-12-07 13:04:02
