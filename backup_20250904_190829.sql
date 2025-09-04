-- MySQL dump 10.13  Distrib 8.0.30, for macos12 (x86_64)
--
-- Host: localhost    Database: rentpro_admin
-- ------------------------------------------------------
-- Server version	8.0.30

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `sys_buildings`
--

DROP TABLE IF EXISTS `sys_buildings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_buildings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `developer` varchar(100) DEFAULT NULL,
  `detailed_address` varchar(500) NOT NULL,
  `city` varchar(50) NOT NULL,
  `district` varchar(50) NOT NULL,
  `business_area` varchar(100) DEFAULT NULL,
  `sub_district` varchar(50) DEFAULT NULL,
  `property_type` varchar(50) DEFAULT NULL,
  `property_company` varchar(100) DEFAULT NULL,
  `description` text,
  `sale_count` bigint DEFAULT '0',
  `rent_count` bigint DEFAULT '0',
  `sale_deals_count` bigint DEFAULT '0',
  `rent_deals_count` bigint DEFAULT '0',
  `status` varchar(20) NOT NULL DEFAULT 'active',
  `is_hot` tinyint(1) DEFAULT '0',
  `created_by` varchar(50) DEFAULT NULL,
  `updated_by` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_sale_count` (`sale_count`),
  KEY `idx_rent_count` (`rent_count`),
  KEY `idx_status` (`status`),
  KEY `idx_is_hot` (`is_hot`),
  KEY `idx_sys_buildings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_buildings`
--

LOCK TABLES `sys_buildings` WRITE;
/*!40000 ALTER TABLE `sys_buildings` DISABLE KEYS */;
INSERT INTO `sys_buildings` VALUES (1,'北京滨江一号','滨江地产','北京市朝阳区张江高科技园区博云路2号','北京市','朝阳区','望京商圈',NULL,'住宅','滨江物业','滨江一号是滨江地产打造的高品质住宅项目，位于朝阳区核心位置，配套设施完善。',15,25,5,10,'active',1,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(2,'北京城市之光','城市发展','北京市朝阳区建国路88号','北京市','朝阳区','CBD商圈',NULL,'商业','城市物业','城市之光是集购物、餐饮、娱乐于一体的综合性商业中心，地处北京市中心繁华地段。',0,8,0,3,'active',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(3,'北京科技园大厦','科技地产','北京市海淀区中关村大街1号','北京市','海淀区','中关村商圈',NULL,'办公','科技物业','科技园大厦是专为高科技企业打造的现代化办公大楼，配备高速网络和智能办公系统。',3,42,1,15,'active',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(4,'北京湖景花园','湖景地产','北京市海淀区文三路478号','北京市','海淀区','中关村商圈',NULL,'住宅','湖景物业','湖景花园依湖而建，环境优美，空气清新，是理想的居住选择。',22,18,8,5,'active',1,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(5,'北京金融中心','金融地产','北京市西城区金融大街1号','北京市','西城区','金融街商圈',NULL,'商业','金融物业','金融中心是北京市标志性建筑，吸引了众多金融机构和企业入驻。',0,12,0,7,'pending',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(6,'北京绿城雅居','绿城集团','北京市朝阳区东三环中路3888号','北京市','朝阳区','国贸商圈',NULL,'住宅','绿城物业','绿城雅居是绿城集团打造的高端住宅项目，绿化率高，居住环境舒适。',18,30,6,12,'active',1,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(7,'北京万达广场','万达集团','北京市朝阳区建国路93号','北京市','朝阳区','CBD商圈',NULL,'商业','万达物业','万达广场是集购物、餐饮、娱乐、办公于一体的大型城市综合体。',0,15,0,8,'active',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(8,'北京SOHO现代城','SOHO中国','北京市朝阳区东三环中路39号','北京市','朝阳区','国贸商圈',NULL,'办公','第一太平戴维斯','SOHO现代城是CBD核心区的标志性写字楼项目，交通便利，配套齐全。',5,35,2,20,'active',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(9,'北京保利国际','保利地产','北京市朝阳区朝阳公园南路1088号','北京市','朝阳区','朝阳公园商圈',NULL,'住宅/商业','保利物业','保利国际是集高端住宅、甲级写字楼、大型购物中心于一体的综合体项目。',28,22,10,6,'active',1,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(10,'北京华润中心','华润置地','北京市东城区东直门南大街5001号','北京市','东城区','东直门商圈',NULL,'商业','华润物业','华润中心是北京的地标性建筑，包含购物中心、写字楼和公寓等多种业态。',0,20,0,12,'active',0,'admin','admin','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL);
/*!40000 ALTER TABLE `sys_buildings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_business_areas`
--

DROP TABLE IF EXISTS `sys_business_areas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_business_areas` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `name` varchar(100) NOT NULL,
  `district_id` bigint unsigned NOT NULL,
  `city_code` varchar(20) NOT NULL,
  `sort` bigint DEFAULT '0',
  `status` varchar(20) DEFAULT 'active',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_business_areas_code` (`code`),
  KEY `idx_sys_business_areas_district_id` (`district_id`),
  KEY `idx_sys_business_areas_city_code` (`city_code`),
  CONSTRAINT `fk_sys_districts_business_areas` FOREIGN KEY (`district_id`) REFERENCES `sys_districts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_business_areas`
--

LOCK TABLES `sys_business_areas` WRITE;
/*!40000 ALTER TABLE `sys_business_areas` DISABLE KEYS */;
INSERT INTO `sys_business_areas` VALUES (1,'BJ001001','国贸商圈',1,'BJ',1,'active'),(2,'BJ001002','三里屯商圈',1,'BJ',2,'active'),(3,'BJ001003','望京商圈',1,'BJ',3,'active'),(4,'BJ001004','亚运村商圈',1,'BJ',4,'active'),(5,'BJ001005','CBD商圈',1,'BJ',5,'active'),(6,'BJ002001','中关村商圈',2,'BJ',1,'active'),(7,'BJ002002','五道口商圈',2,'BJ',2,'active'),(8,'BJ002003','西二旗商圈',2,'BJ',3,'active'),(9,'BJ002004','上地商圈',2,'BJ',4,'active'),(10,'BJ002005','万柳商圈',2,'BJ',5,'active'),(11,'BJ003001','金融街商圈',3,'BJ',1,'active'),(12,'BJ003002','西单商圈',3,'BJ',2,'active'),(13,'BJ003003','什刹海商圈',3,'BJ',3,'active'),(14,'BJ003004','德胜门商圈',3,'BJ',4,'active'),(15,'BJ004001','王府井商圈',4,'BJ',1,'active'),(16,'BJ004002','东单商圈',4,'BJ',2,'active'),(17,'BJ004003','前门商圈',4,'BJ',3,'active'),(18,'BJ004004','崇文门商圈',4,'BJ',4,'active'),(19,'BJ005001','丽泽商圈',5,'BJ',1,'active'),(20,'BJ005002','丰台科技园商圈',5,'BJ',2,'active'),(21,'BJ005003','方庄商圈',5,'BJ',3,'active'),(22,'BJ006001','石景山万达商圈',6,'BJ',1,'active'),(23,'BJ006002','八角商圈',6,'BJ',2,'active'),(24,'BJ007001','通州万达商圈',7,'BJ',1,'active'),(25,'BJ007002','运河商务区商圈',7,'BJ',2,'active'),(26,'BJ008001','回龙观商圈',8,'BJ',1,'active'),(27,'BJ008002','天通苑商圈',8,'BJ',2,'active'),(28,'BJ009001','亦庄商圈',9,'BJ',1,'active'),(29,'BJ009002','大兴新城商圈',9,'BJ',2,'active'),(30,'BJ010001','良乡商圈',10,'BJ',1,'active'),(31,'BJ010002','燕山商圈',10,'BJ',2,'active');
/*!40000 ALTER TABLE `sys_business_areas` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dept`
--

DROP TABLE IF EXISTS `sys_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dept` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned DEFAULT '0',
  `dept_path` varchar(255) DEFAULT NULL,
  `dept_name` varchar(128) NOT NULL,
  `sort` bigint DEFAULT '1',
  `leader` varchar(128) DEFAULT NULL,
  `phone` varchar(32) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `status` varchar(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sys_dept_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dept`
--

LOCK TABLES `sys_dept` WRITE;
/*!40000 ALTER TABLE `sys_dept` DISABLE KEYS */;
INSERT INTO `sys_dept` VALUES (1,0,'0,1','RentPro科技',1,'系统管理员','15888888888','admin@rentpro.com','0','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(2,1,'0,1,2','技术部',1,'技术总监','15666666666','tech@rentpro.com','0','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(3,1,'0,1,3','运营部',2,'运营总监','15777777777','ops@rentpro.com','0','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL);
/*!40000 ALTER TABLE `sys_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_districts`
--

DROP TABLE IF EXISTS `sys_districts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_districts` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `city_code` varchar(20) NOT NULL,
  `sort` bigint DEFAULT '0',
  `status` varchar(20) DEFAULT 'active',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sys_districts_code` (`code`),
  KEY `idx_sys_districts_city_code` (`city_code`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_districts`
--

LOCK TABLES `sys_districts` WRITE;
/*!40000 ALTER TABLE `sys_districts` DISABLE KEYS */;
INSERT INTO `sys_districts` VALUES (1,'BJ001','朝阳区','BJ',1,'active'),(2,'BJ002','海淀区','BJ',2,'active'),(3,'BJ003','西城区','BJ',3,'active'),(4,'BJ004','东城区','BJ',4,'active'),(5,'BJ005','丰台区','BJ',5,'active'),(6,'BJ006','石景山区','BJ',6,'active'),(7,'BJ007','通州区','BJ',7,'active'),(8,'BJ008','昌平区','BJ',8,'active'),(9,'BJ009','大兴区','BJ',9,'active'),(10,'BJ010','房山区','BJ',10,'active');
/*!40000 ALTER TABLE `sys_districts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_house_types`
--

DROP TABLE IF EXISTS `sys_house_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_house_types` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `code` varchar(50) NOT NULL,
  `description` text,
  `building_id` bigint unsigned NOT NULL,
  `area` decimal(8,2) NOT NULL,
  `rooms` bigint NOT NULL DEFAULT '1',
  `halls` bigint NOT NULL DEFAULT '1',
  `bathrooms` bigint NOT NULL DEFAULT '1',
  `balconies` bigint DEFAULT '0',
  `floor_height` decimal(4,2) DEFAULT NULL,
  `orientation` varchar(50) DEFAULT NULL,
  `view` varchar(100) DEFAULT NULL,
  `sale_price` decimal(12,2) DEFAULT '0.00',
  `rent_price` decimal(8,2) DEFAULT '0.00',
  `sale_price_per` decimal(8,2) DEFAULT '0.00',
  `rent_price_per` decimal(6,2) DEFAULT '0.00',
  `total_stock` bigint NOT NULL DEFAULT '0',
  `sale_stock` bigint NOT NULL DEFAULT '0',
  `rent_stock` bigint NOT NULL DEFAULT '0',
  `reserved_stock` bigint NOT NULL DEFAULT '0',
  `status` varchar(20) NOT NULL DEFAULT 'active',
  `sale_status` varchar(20) DEFAULT 'available',
  `rent_status` varchar(20) DEFAULT 'available',
  `is_hot` tinyint(1) DEFAULT '0',
  `main_image` varchar(500) DEFAULT NULL,
  `image_urls` json DEFAULT NULL,
  `tags` json DEFAULT NULL,
  `created_by` varchar(50) DEFAULT NULL,
  `updated_by` varchar(50) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_name` (`name`),
  KEY `idx_building_id` (`building_id`),
  KEY `idx_area` (`area`,`sale_price`),
  KEY `idx_rent_price` (`rent_price`),
  KEY `idx_status` (`status`),
  KEY `idx_is_hot` (`is_hot`),
  KEY `idx_sys_house_types_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_sys_house_types_building` FOREIGN KEY (`building_id`) REFERENCES `sys_buildings` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_house_types`
--

LOCK TABLES `sys_house_types` WRITE;
/*!40000 ALTER TABLE `sys_house_types` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_house_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu`
--

DROP TABLE IF EXISTS `sys_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `title` varchar(128) DEFAULT NULL,
  `icon` varchar(128) DEFAULT NULL,
  `path` varchar(128) DEFAULT NULL,
  `redirect` varchar(128) DEFAULT NULL,
  `component` varchar(128) DEFAULT NULL,
  `permission` varchar(255) DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT '0',
  `type` varchar(1) DEFAULT 'M',
  `sort` bigint DEFAULT '1',
  `visible` varchar(1) DEFAULT '0',
  `is_frame` varchar(1) DEFAULT '1',
  `is_cache` varchar(1) DEFAULT '0',
  `menu_type` varchar(1) DEFAULT '',
  `status` varchar(1) DEFAULT '0',
  `perms` varchar(100) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sys_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu`
--

LOCK TABLES `sys_menu` WRITE;
/*!40000 ALTER TABLE `sys_menu` DISABLE KEYS */;
INSERT INTO `sys_menu` VALUES (1,'System','系统管理','Setting','/system','','Layout','system:view',0,'M',1,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(2,'Rental','租赁管理','OfficeBuilding','/rental','','Layout','rental:view',0,'M',2,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(11,'User','用户管理','User','/system/user','','system/user/index','system:user:view',1,'C',1,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(12,'Role','角色管理','UserFilled','/system/role','','system/role/index','system:role:view',1,'C',2,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(13,'Menu','菜单管理','Menu','/system/menu','','system/menu/index','system:menu:view',1,'C',3,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(21,'Building','楼盘管理','House','/rental/building','','rental/building/building-management','rental:building:view',2,'C',1,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(22,'House','房屋管理','House','/rental/house','','rental/house/index','rental:house:view',2,'C',2,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(23,'Tenant','租户管理','User','/rental/tenant','','rental/tenant/index','rental:tenant:view',2,'C',3,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(24,'Agent','经纪人管理','UserFilled','/rental/agent','','rental/agent/index','rental:agent:view',2,'C',4,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(25,'Landlord','房东管理','UserFilled','/rental/landlord','','rental/landlord/index','rental:landlord:view',2,'C',5,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL),(26,'Contract','合同管理','Document','/rental/contract','','rental/contract/index','rental:contract:view',2,'C',6,'0','1','0','','0','','2025-09-02 21:43:18.000','2025-09-02 21:43:18.000',NULL);
/*!40000 ALTER TABLE `sys_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_migration`
--

DROP TABLE IF EXISTS `sys_migration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_migration` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `version` varchar(191) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'completed',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_migration_version` (`version`),
  KEY `idx_version` (`version`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_migration`
--

LOCK TABLES `sys_migration` WRITE;
/*!40000 ALTER TABLE `sys_migration` DISABLE KEYS */;
INSERT INTO `sys_migration` VALUES (1,'1756303530910','创建所有系统表和业务表','completed','2025-09-02 01:03:05.732','2025-09-02 01:03:05.732');
/*!40000 ALTER TABLE `sys_migration` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_post`
--

DROP TABLE IF EXISTS `sys_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_post` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `post_code` varchar(64) NOT NULL,
  `post_name` varchar(128) NOT NULL,
  `sort` bigint DEFAULT '1',
  `status` varchar(1) DEFAULT '0',
  `remark` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_post_post_code` (`post_code`),
  KEY `idx_post_code` (`post_code`),
  KEY `idx_sys_post_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_post`
--

LOCK TABLES `sys_post` WRITE;
/*!40000 ALTER TABLE `sys_post` DISABLE KEYS */;
INSERT INTO `sys_post` VALUES (1,'ceo','董事长',1,'0','董事长','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(2,'se','项目经理',2,'0','项目经理','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(3,'hr','人力资源',3,'0','人力资源','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL),(4,'user','普通员工',4,'0','普通员工','2025-09-02 01:03:06.000','2025-09-02 01:03:06.000',NULL);
/*!40000 ALTER TABLE `sys_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `key` varchar(128) NOT NULL,
  `status` bigint DEFAULT '1',
  `sort` bigint DEFAULT '1',
  `flag` varchar(128) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  `admin` tinyint(1) DEFAULT '0',
  `data_scope` varchar(128) DEFAULT '1',
  `params` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_role_name` (`name`),
  UNIQUE KEY `uni_sys_role_key` (`key`),
  KEY `idx_name` (`name`),
  KEY `idx_key` (`key`),
  KEY `idx_sys_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'超级管理员','admin',1,1,'','超级管理员',1,'1','','2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL),(2,'普通用户','common',1,2,'','普通用户',0,'2','','2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL),(3,'租客','tenant',1,3,'','租客用户',0,'5','','2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL),(4,'房东','landlord',1,4,'','房东用户',0,'4','','2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL);
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_menu`
--

DROP TABLE IF EXISTS `sys_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_menu` (
  `sys_menu_id` bigint unsigned NOT NULL,
  `sys_role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`sys_menu_id`,`sys_role_id`),
  KEY `fk_sys_role_menu_sys_role` (`sys_role_id`),
  CONSTRAINT `fk_sys_role_menu_sys_menu` FOREIGN KEY (`sys_menu_id`) REFERENCES `sys_menu` (`id`),
  CONSTRAINT `fk_sys_role_menu_sys_role` FOREIGN KEY (`sys_role_id`) REFERENCES `sys_role` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_menu`
--

LOCK TABLES `sys_role_menu` WRITE;
/*!40000 ALTER TABLE `sys_role_menu` DISABLE KEYS */;
INSERT INTO `sys_role_menu` VALUES (1,1),(2,1),(11,1),(12,1),(13,1),(21,1),(22,1),(23,1),(24,1),(25,1),(26,1),(2,2),(21,2),(22,2),(23,2),(24,2),(25,2),(26,2);
/*!40000 ALTER TABLE `sys_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password` varchar(128) NOT NULL,
  `nick_name` varchar(128) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `phone` varchar(32) DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `is_admin` tinyint(1) DEFAULT '0',
  `remark` varchar(255) DEFAULT NULL,
  `dept_id` bigint unsigned DEFAULT '0',
  `post_id` bigint unsigned DEFAULT '0',
  `role_id` bigint unsigned DEFAULT '0',
  `salt` varchar(255) DEFAULT NULL,
  `last_login_ip` varchar(128) DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_sys_user_username` (`username`),
  KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_sys_user_deleted_at` (`deleted_at`),
  KEY `fk_sys_role_users` (`role_id`),
  KEY `fk_sys_dept_users` (`dept_id`),
  KEY `fk_sys_post_users` (`post_id`),
  CONSTRAINT `fk_sys_dept_users` FOREIGN KEY (`dept_id`) REFERENCES `sys_dept` (`id`),
  CONSTRAINT `fk_sys_post_users` FOREIGN KEY (`post_id`) REFERENCES `sys_post` (`id`),
  CONSTRAINT `fk_sys_role_users` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,'admin','123456','超级管理员','','admin@rentpro.com','15888888888',1,1,'超级管理员账号',1,1,1,'','',NULL,'2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL),(2,'test','123456','测试用户','','test@rentpro.com','15666666666',1,0,'测试用户账号',2,4,2,'','',NULL,'2025-09-02 01:03:07.000','2025-09-02 01:03:07.000',NULL);
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'rentpro_admin'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-04 19:08:30
