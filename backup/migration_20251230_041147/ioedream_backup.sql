-- MySQL dump 10.13  Distrib 8.0.44, for Linux (x86_64)
--
-- Host: localhost    Database: ioedream
-- ------------------------------------------------------
-- Server version	8.0.44

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
-- Table structure for table `ACT_EVT_LOG`
--

DROP TABLE IF EXISTS `ACT_EVT_LOG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_EVT_LOG` (
  `LOG_NR_` bigint NOT NULL AUTO_INCREMENT,
  `TYPE_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TIME_STAMP_` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DATA_` longblob,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `LOCK_TIME_` timestamp(3) NULL DEFAULT NULL,
  `IS_PROCESSED_` tinyint DEFAULT '0',
  PRIMARY KEY (`LOG_NR_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_EVT_LOG`
--

LOCK TABLES `ACT_EVT_LOG` WRITE;
/*!40000 ALTER TABLE `ACT_EVT_LOG` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_EVT_LOG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_GE_BYTEARRAY`
--

DROP TABLE IF EXISTS `ACT_GE_BYTEARRAY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_GE_BYTEARRAY` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DEPLOYMENT_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `BYTES_` longblob,
  `GENERATED_` tinyint DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_BYTEAR_DEPL` (`DEPLOYMENT_ID_`),
  CONSTRAINT `ACT_FK_BYTEARR_DEPL` FOREIGN KEY (`DEPLOYMENT_ID_`) REFERENCES `ACT_RE_DEPLOYMENT` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_GE_BYTEARRAY`
--

LOCK TABLES `ACT_GE_BYTEARRAY` WRITE;
/*!40000 ALTER TABLE `ACT_GE_BYTEARRAY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_GE_BYTEARRAY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_GE_PROPERTY`
--

DROP TABLE IF EXISTS `ACT_GE_PROPERTY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_GE_PROPERTY` (
  `NAME_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `VALUE_` varchar(300) COLLATE utf8mb3_bin DEFAULT NULL,
  `REV_` int DEFAULT NULL,
  PRIMARY KEY (`NAME_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_GE_PROPERTY`
--

LOCK TABLES `ACT_GE_PROPERTY` WRITE;
/*!40000 ALTER TABLE `ACT_GE_PROPERTY` DISABLE KEYS */;
INSERT INTO `ACT_GE_PROPERTY` VALUES ('cfg.execution-related-entities-count','true',1),('cfg.task-related-entities-count','true',1),('common.schema.version','7.2.0.2',1),('eventregistry.schema.version','7.2.0.2',1),('next.dbid','1',1),('schema.history','create(7.2.0.2)',1),('schema.version','7.2.0.2',1);
/*!40000 ALTER TABLE `ACT_GE_PROPERTY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_ACTINST`
--

DROP TABLE IF EXISTS `ACT_HI_ACTINST`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_ACTINST` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT '1',
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `ACT_ID_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CALL_PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `ASSIGNEE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `COMPLETED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_TIME_` datetime(3) NOT NULL,
  `END_TIME_` datetime(3) DEFAULT NULL,
  `TRANSACTION_ORDER_` int DEFAULT NULL,
  `DURATION_` bigint DEFAULT NULL,
  `DELETE_REASON_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_ACT_INST_START` (`START_TIME_`),
  KEY `ACT_IDX_HI_ACT_INST_END` (`END_TIME_`),
  KEY `ACT_IDX_HI_ACT_INST_PROCINST` (`PROC_INST_ID_`,`ACT_ID_`),
  KEY `ACT_IDX_HI_ACT_INST_EXEC` (`EXECUTION_ID_`,`ACT_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_ACTINST`
--

LOCK TABLES `ACT_HI_ACTINST` WRITE;
/*!40000 ALTER TABLE `ACT_HI_ACTINST` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_ACTINST` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_ATTACHMENT`
--

DROP TABLE IF EXISTS `ACT_HI_ATTACHMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_ATTACHMENT` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DESCRIPTION_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `URL_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CONTENT_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TIME_` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_ATTACHMENT`
--

LOCK TABLES `ACT_HI_ATTACHMENT` WRITE;
/*!40000 ALTER TABLE `ACT_HI_ATTACHMENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_ATTACHMENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_COMMENT`
--

DROP TABLE IF EXISTS `ACT_HI_COMMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_COMMENT` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TIME_` datetime(3) NOT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACTION_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `MESSAGE_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `FULL_MSG_` longblob,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_COMMENT`
--

LOCK TABLES `ACT_HI_COMMENT` WRITE;
/*!40000 ALTER TABLE `ACT_HI_COMMENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_COMMENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_DETAIL`
--

DROP TABLE IF EXISTS `ACT_HI_DETAIL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_DETAIL` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `VAR_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REV_` int DEFAULT NULL,
  `TIME_` datetime(3) NOT NULL,
  `BYTEARRAY_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DOUBLE_` double DEFAULT NULL,
  `LONG_` bigint DEFAULT NULL,
  `TEXT_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TEXT2_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_DETAIL_PROC_INST` (`PROC_INST_ID_`),
  KEY `ACT_IDX_HI_DETAIL_ACT_INST` (`ACT_INST_ID_`),
  KEY `ACT_IDX_HI_DETAIL_TIME` (`TIME_`),
  KEY `ACT_IDX_HI_DETAIL_NAME` (`NAME_`),
  KEY `ACT_IDX_HI_DETAIL_TASK_ID` (`TASK_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_DETAIL`
--

LOCK TABLES `ACT_HI_DETAIL` WRITE;
/*!40000 ALTER TABLE `ACT_HI_DETAIL` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_DETAIL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_ENTITYLINK`
--

DROP TABLE IF EXISTS `ACT_HI_ENTITYLINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_ENTITYLINK` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `LINK_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` datetime(3) DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ROOT_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ROOT_SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HIERARCHY_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_ENT_LNK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_HI_ENT_LNK_REF_SCOPE` (`REF_SCOPE_ID_`,`REF_SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_HI_ENT_LNK_ROOT_SCOPE` (`ROOT_SCOPE_ID_`,`ROOT_SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_HI_ENT_LNK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`,`LINK_TYPE_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_ENTITYLINK`
--

LOCK TABLES `ACT_HI_ENTITYLINK` WRITE;
/*!40000 ALTER TABLE `ACT_HI_ENTITYLINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_ENTITYLINK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_IDENTITYLINK`
--

DROP TABLE IF EXISTS `ACT_HI_IDENTITYLINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_IDENTITYLINK` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `GROUP_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` datetime(3) DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_IDENT_LNK_USER` (`USER_ID_`),
  KEY `ACT_IDX_HI_IDENT_LNK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_IDENT_LNK_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_IDENT_LNK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_IDENT_LNK_TASK` (`TASK_ID_`),
  KEY `ACT_IDX_HI_IDENT_LNK_PROCINST` (`PROC_INST_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_IDENTITYLINK`
--

LOCK TABLES `ACT_HI_IDENTITYLINK` WRITE;
/*!40000 ALTER TABLE `ACT_HI_IDENTITYLINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_IDENTITYLINK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_PROCINST`
--

DROP TABLE IF EXISTS `ACT_HI_PROCINST`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_PROCINST` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT '1',
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `BUSINESS_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `START_TIME_` datetime(3) NOT NULL,
  `END_TIME_` datetime(3) DEFAULT NULL,
  `DURATION_` bigint DEFAULT NULL,
  `START_USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_ACT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `END_ACT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUPER_PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DELETE_REASON_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CALLBACK_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CALLBACK_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REFERENCE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REFERENCE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROPAGATED_STAGE_INST_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BUSINESS_STATUS_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `PROC_INST_ID_` (`PROC_INST_ID_`),
  KEY `ACT_IDX_HI_PRO_INST_END` (`END_TIME_`),
  KEY `ACT_IDX_HI_PRO_I_BUSKEY` (`BUSINESS_KEY_`),
  KEY `ACT_IDX_HI_PRO_SUPER_PROCINST` (`SUPER_PROCESS_INSTANCE_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_PROCINST`
--

LOCK TABLES `ACT_HI_PROCINST` WRITE;
/*!40000 ALTER TABLE `ACT_HI_PROCINST` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_PROCINST` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_TASKINST`
--

DROP TABLE IF EXISTS `ACT_HI_TASKINST`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_TASKINST` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT '1',
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_DEF_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROPAGATED_STAGE_INST_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `STATE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DESCRIPTION_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ASSIGNEE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_TIME_` datetime(3) NOT NULL,
  `IN_PROGRESS_TIME_` datetime(3) DEFAULT NULL,
  `IN_PROGRESS_STARTED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CLAIM_TIME_` datetime(3) DEFAULT NULL,
  `CLAIMED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUSPENDED_TIME_` datetime(3) DEFAULT NULL,
  `SUSPENDED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `END_TIME_` datetime(3) DEFAULT NULL,
  `COMPLETED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DURATION_` bigint DEFAULT NULL,
  `DELETE_REASON_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `PRIORITY_` int DEFAULT NULL,
  `IN_PROGRESS_DUE_DATE_` datetime(3) DEFAULT NULL,
  `DUE_DATE_` datetime(3) DEFAULT NULL,
  `FORM_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `LAST_UPDATED_TIME_` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_TASK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_TASK_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_TASK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_TASK_INST_PROCINST` (`PROC_INST_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_TASKINST`
--

LOCK TABLES `ACT_HI_TASKINST` WRITE;
/*!40000 ALTER TABLE `ACT_HI_TASKINST` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_TASKINST` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_TSK_LOG`
--

DROP TABLE IF EXISTS `ACT_HI_TSK_LOG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_TSK_LOG` (
  `ID_` bigint NOT NULL AUTO_INCREMENT,
  `TYPE_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `TIME_STAMP_` timestamp(3) NOT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DATA_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_ACT_HI_TSK_LOG_TASK` (`TASK_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_TSK_LOG`
--

LOCK TABLES `ACT_HI_TSK_LOG` WRITE;
/*!40000 ALTER TABLE `ACT_HI_TSK_LOG` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_TSK_LOG` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_HI_VARINST`
--

DROP TABLE IF EXISTS `ACT_HI_VARINST`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_HI_VARINST` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT '1',
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `VAR_TYPE_` varchar(100) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BYTEARRAY_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DOUBLE_` double DEFAULT NULL,
  `LONG_` bigint DEFAULT NULL,
  `TEXT_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TEXT2_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `META_INFO_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` datetime(3) DEFAULT NULL,
  `LAST_UPDATED_TIME_` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_HI_PROCVAR_NAME_TYPE` (`NAME_`,`VAR_TYPE_`),
  KEY `ACT_IDX_HI_VAR_SCOPE_ID_TYPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_VAR_SUB_ID_TYPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_HI_PROCVAR_PROC_INST` (`PROC_INST_ID_`),
  KEY `ACT_IDX_HI_PROCVAR_TASK_ID` (`TASK_ID_`),
  KEY `ACT_IDX_HI_PROCVAR_EXE` (`EXECUTION_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_HI_VARINST`
--

LOCK TABLES `ACT_HI_VARINST` WRITE;
/*!40000 ALTER TABLE `ACT_HI_VARINST` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_HI_VARINST` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_BYTEARRAY`
--

DROP TABLE IF EXISTS `ACT_ID_BYTEARRAY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_BYTEARRAY` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BYTES_` longblob,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_BYTEARRAY`
--

LOCK TABLES `ACT_ID_BYTEARRAY` WRITE;
/*!40000 ALTER TABLE `ACT_ID_BYTEARRAY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_BYTEARRAY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_GROUP`
--

DROP TABLE IF EXISTS `ACT_ID_GROUP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_GROUP` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_GROUP`
--

LOCK TABLES `ACT_ID_GROUP` WRITE;
/*!40000 ALTER TABLE `ACT_ID_GROUP` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_GROUP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_INFO`
--

DROP TABLE IF EXISTS `ACT_ID_INFO`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_INFO` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `USER_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `VALUE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PASSWORD_` longblob,
  `PARENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_INFO`
--

LOCK TABLES `ACT_ID_INFO` WRITE;
/*!40000 ALTER TABLE `ACT_ID_INFO` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_INFO` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_MEMBERSHIP`
--

DROP TABLE IF EXISTS `ACT_ID_MEMBERSHIP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_MEMBERSHIP` (
  `USER_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `GROUP_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  PRIMARY KEY (`USER_ID_`,`GROUP_ID_`),
  KEY `ACT_FK_MEMB_GROUP` (`GROUP_ID_`),
  CONSTRAINT `ACT_FK_MEMB_GROUP` FOREIGN KEY (`GROUP_ID_`) REFERENCES `ACT_ID_GROUP` (`ID_`),
  CONSTRAINT `ACT_FK_MEMB_USER` FOREIGN KEY (`USER_ID_`) REFERENCES `ACT_ID_USER` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_MEMBERSHIP`
--

LOCK TABLES `ACT_ID_MEMBERSHIP` WRITE;
/*!40000 ALTER TABLE `ACT_ID_MEMBERSHIP` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_MEMBERSHIP` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_PRIV`
--

DROP TABLE IF EXISTS `ACT_ID_PRIV`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_PRIV` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `ACT_UNIQ_PRIV_NAME` (`NAME_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_PRIV`
--

LOCK TABLES `ACT_ID_PRIV` WRITE;
/*!40000 ALTER TABLE `ACT_ID_PRIV` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_PRIV` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_PRIV_MAPPING`
--

DROP TABLE IF EXISTS `ACT_ID_PRIV_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_PRIV_MAPPING` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `PRIV_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `GROUP_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_FK_PRIV_MAPPING` (`PRIV_ID_`),
  KEY `ACT_IDX_PRIV_USER` (`USER_ID_`),
  KEY `ACT_IDX_PRIV_GROUP` (`GROUP_ID_`),
  CONSTRAINT `ACT_FK_PRIV_MAPPING` FOREIGN KEY (`PRIV_ID_`) REFERENCES `ACT_ID_PRIV` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_PRIV_MAPPING`
--

LOCK TABLES `ACT_ID_PRIV_MAPPING` WRITE;
/*!40000 ALTER TABLE `ACT_ID_PRIV_MAPPING` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_PRIV_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_PROPERTY`
--

DROP TABLE IF EXISTS `ACT_ID_PROPERTY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_PROPERTY` (
  `NAME_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `VALUE_` varchar(300) COLLATE utf8mb3_bin DEFAULT NULL,
  `REV_` int DEFAULT NULL,
  PRIMARY KEY (`NAME_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_PROPERTY`
--

LOCK TABLES `ACT_ID_PROPERTY` WRITE;
/*!40000 ALTER TABLE `ACT_ID_PROPERTY` DISABLE KEYS */;
INSERT INTO `ACT_ID_PROPERTY` VALUES ('schema.version','7.2.0.2',1);
/*!40000 ALTER TABLE `ACT_ID_PROPERTY` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_TOKEN`
--

DROP TABLE IF EXISTS `ACT_ID_TOKEN`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_TOKEN` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `TOKEN_VALUE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TOKEN_DATE_` timestamp(3) NULL DEFAULT NULL,
  `IP_ADDRESS_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `USER_AGENT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TOKEN_DATA_` varchar(2000) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_TOKEN`
--

LOCK TABLES `ACT_ID_TOKEN` WRITE;
/*!40000 ALTER TABLE `ACT_ID_TOKEN` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_TOKEN` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_ID_USER`
--

DROP TABLE IF EXISTS `ACT_ID_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_ID_USER` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `FIRST_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `LAST_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DISPLAY_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EMAIL_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PWD_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PICTURE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_ID_USER`
--

LOCK TABLES `ACT_ID_USER` WRITE;
/*!40000 ALTER TABLE `ACT_ID_USER` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_ID_USER` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_PROCDEF_INFO`
--

DROP TABLE IF EXISTS `ACT_PROCDEF_INFO`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_PROCDEF_INFO` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `INFO_JSON_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `ACT_UNIQ_INFO_PROCDEF` (`PROC_DEF_ID_`),
  KEY `ACT_IDX_INFO_PROCDEF` (`PROC_DEF_ID_`),
  KEY `ACT_FK_INFO_JSON_BA` (`INFO_JSON_ID_`),
  CONSTRAINT `ACT_FK_INFO_JSON_BA` FOREIGN KEY (`INFO_JSON_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_INFO_PROCDEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_PROCDEF_INFO`
--

LOCK TABLES `ACT_PROCDEF_INFO` WRITE;
/*!40000 ALTER TABLE `ACT_PROCDEF_INFO` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_PROCDEF_INFO` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RE_DEPLOYMENT`
--

DROP TABLE IF EXISTS `ACT_RE_DEPLOYMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RE_DEPLOYMENT` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `DEPLOY_TIME_` timestamp(3) NULL DEFAULT NULL,
  `DERIVED_FROM_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DERIVED_FROM_ROOT_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_DEPLOYMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ENGINE_VERSION_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RE_DEPLOYMENT`
--

LOCK TABLES `ACT_RE_DEPLOYMENT` WRITE;
/*!40000 ALTER TABLE `ACT_RE_DEPLOYMENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RE_DEPLOYMENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RE_MODEL`
--

DROP TABLE IF EXISTS `ACT_RE_MODEL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RE_MODEL` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LAST_UPDATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `VERSION_` int DEFAULT NULL,
  `META_INFO_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DEPLOYMENT_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EDITOR_SOURCE_VALUE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EDITOR_SOURCE_EXTRA_VALUE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_FK_MODEL_SOURCE` (`EDITOR_SOURCE_VALUE_ID_`),
  KEY `ACT_FK_MODEL_SOURCE_EXTRA` (`EDITOR_SOURCE_EXTRA_VALUE_ID_`),
  KEY `ACT_FK_MODEL_DEPLOYMENT` (`DEPLOYMENT_ID_`),
  CONSTRAINT `ACT_FK_MODEL_DEPLOYMENT` FOREIGN KEY (`DEPLOYMENT_ID_`) REFERENCES `ACT_RE_DEPLOYMENT` (`ID_`),
  CONSTRAINT `ACT_FK_MODEL_SOURCE` FOREIGN KEY (`EDITOR_SOURCE_VALUE_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_MODEL_SOURCE_EXTRA` FOREIGN KEY (`EDITOR_SOURCE_EXTRA_VALUE_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RE_MODEL`
--

LOCK TABLES `ACT_RE_MODEL` WRITE;
/*!40000 ALTER TABLE `ACT_RE_MODEL` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RE_MODEL` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RE_PROCDEF`
--

DROP TABLE IF EXISTS `ACT_RE_PROCDEF`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RE_PROCDEF` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `VERSION_` int NOT NULL,
  `DEPLOYMENT_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `RESOURCE_NAME_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DGRM_RESOURCE_NAME_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DESCRIPTION_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `HAS_START_FORM_KEY_` tinyint DEFAULT NULL,
  `HAS_GRAPHICAL_NOTATION_` tinyint DEFAULT NULL,
  `SUSPENSION_STATE_` int DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `ENGINE_VERSION_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DERIVED_FROM_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DERIVED_FROM_ROOT_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DERIVED_VERSION_` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `ACT_UNIQ_PROCDEF` (`KEY_`,`VERSION_`,`DERIVED_VERSION_`,`TENANT_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RE_PROCDEF`
--

LOCK TABLES `ACT_RE_PROCDEF` WRITE;
/*!40000 ALTER TABLE `ACT_RE_PROCDEF` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RE_PROCDEF` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_ACTINST`
--

DROP TABLE IF EXISTS `ACT_RU_ACTINST`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_ACTINST` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT '1',
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `ACT_ID_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CALL_PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `ASSIGNEE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `COMPLETED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_TIME_` datetime(3) NOT NULL,
  `END_TIME_` datetime(3) DEFAULT NULL,
  `DURATION_` bigint DEFAULT NULL,
  `TRANSACTION_ORDER_` int DEFAULT NULL,
  `DELETE_REASON_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_RU_ACTI_START` (`START_TIME_`),
  KEY `ACT_IDX_RU_ACTI_END` (`END_TIME_`),
  KEY `ACT_IDX_RU_ACTI_PROC` (`PROC_INST_ID_`),
  KEY `ACT_IDX_RU_ACTI_PROC_ACT` (`PROC_INST_ID_`,`ACT_ID_`),
  KEY `ACT_IDX_RU_ACTI_EXEC` (`EXECUTION_ID_`),
  KEY `ACT_IDX_RU_ACTI_EXEC_ACT` (`EXECUTION_ID_`,`ACT_ID_`),
  KEY `ACT_IDX_RU_ACTI_TASK` (`TASK_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_ACTINST`
--

LOCK TABLES `ACT_RU_ACTINST` WRITE;
/*!40000 ALTER TABLE `ACT_RU_ACTINST` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_ACTINST` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_DEADLETTER_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_DEADLETTER_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_DEADLETTER_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `EXCLUSIVE_` tinyint(1) DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CORRELATION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DUEDATE_` timestamp(3) NULL DEFAULT NULL,
  `REPEAT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_DEADLETTER_JOB_EXCEPTION_STACK_ID` (`EXCEPTION_STACK_ID_`),
  KEY `ACT_IDX_DEADLETTER_JOB_CUSTOM_VALUES_ID` (`CUSTOM_VALUES_ID_`),
  KEY `ACT_IDX_DEADLETTER_JOB_CORRELATION_ID` (`CORRELATION_ID_`),
  KEY `ACT_IDX_DJOB_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_DJOB_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_DJOB_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_DEADLETTER_JOB_EXECUTION` (`EXECUTION_ID_`),
  KEY `ACT_FK_DEADLETTER_JOB_PROCESS_INSTANCE` (`PROCESS_INSTANCE_ID_`),
  KEY `ACT_FK_DEADLETTER_JOB_PROC_DEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_DEADLETTER_JOB_CUSTOM_VALUES` FOREIGN KEY (`CUSTOM_VALUES_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_DEADLETTER_JOB_EXCEPTION` FOREIGN KEY (`EXCEPTION_STACK_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_DEADLETTER_JOB_EXECUTION` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_DEADLETTER_JOB_PROC_DEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_DEADLETTER_JOB_PROCESS_INSTANCE` FOREIGN KEY (`PROCESS_INSTANCE_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_DEADLETTER_JOB`
--

LOCK TABLES `ACT_RU_DEADLETTER_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_DEADLETTER_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_DEADLETTER_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_ENTITYLINK`
--

DROP TABLE IF EXISTS `ACT_RU_ENTITYLINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_ENTITYLINK` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CREATE_TIME_` datetime(3) DEFAULT NULL,
  `LINK_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REF_SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ROOT_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ROOT_SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HIERARCHY_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_ENT_LNK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_ENT_LNK_REF_SCOPE` (`REF_SCOPE_ID_`,`REF_SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_ENT_LNK_ROOT_SCOPE` (`ROOT_SCOPE_ID_`,`ROOT_SCOPE_TYPE_`,`LINK_TYPE_`),
  KEY `ACT_IDX_ENT_LNK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`,`LINK_TYPE_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_ENTITYLINK`
--

LOCK TABLES `ACT_RU_ENTITYLINK` WRITE;
/*!40000 ALTER TABLE `ACT_RU_ENTITYLINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_ENTITYLINK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_EVENT_SUBSCR`
--

DROP TABLE IF EXISTS `ACT_RU_EVENT_SUBSCR`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_EVENT_SUBSCR` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `EVENT_TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `EVENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACTIVITY_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CONFIGURATION_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATED_` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `LOCK_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_EVENT_SUBSCR_CONFIG_` (`CONFIGURATION_`),
  KEY `ACT_IDX_EVENT_SUBSCR_EXEC_ID` (`EXECUTION_ID_`),
  KEY `ACT_IDX_EVENT_SUBSCR_PROC_ID` (`PROC_INST_ID_`),
  KEY `ACT_IDX_EVENT_SUBSCR_SCOPEREF_` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  CONSTRAINT `ACT_FK_EVENT_EXEC` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_EVENT_SUBSCR`
--

LOCK TABLES `ACT_RU_EVENT_SUBSCR` WRITE;
/*!40000 ALTER TABLE `ACT_RU_EVENT_SUBSCR` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_EVENT_SUBSCR` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_EXECUTION`
--

DROP TABLE IF EXISTS `ACT_RU_EXECUTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_EXECUTION` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `BUSINESS_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUPER_EXEC_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ROOT_PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ACT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `IS_ACTIVE_` tinyint DEFAULT NULL,
  `IS_CONCURRENT_` tinyint DEFAULT NULL,
  `IS_SCOPE_` tinyint DEFAULT NULL,
  `IS_EVENT_SCOPE_` tinyint DEFAULT NULL,
  `IS_MI_ROOT_` tinyint DEFAULT NULL,
  `SUSPENSION_STATE_` int DEFAULT NULL,
  `CACHED_ENT_STATE_` int DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_ACT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `START_TIME_` datetime(3) DEFAULT NULL,
  `START_USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `LOCK_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `IS_COUNT_ENABLED_` tinyint DEFAULT NULL,
  `EVT_SUBSCR_COUNT_` int DEFAULT NULL,
  `TASK_COUNT_` int DEFAULT NULL,
  `JOB_COUNT_` int DEFAULT NULL,
  `TIMER_JOB_COUNT_` int DEFAULT NULL,
  `SUSP_JOB_COUNT_` int DEFAULT NULL,
  `DEADLETTER_JOB_COUNT_` int DEFAULT NULL,
  `EXTERNAL_WORKER_JOB_COUNT_` int DEFAULT NULL,
  `VAR_COUNT_` int DEFAULT NULL,
  `ID_LINK_COUNT_` int DEFAULT NULL,
  `CALLBACK_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CALLBACK_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REFERENCE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `REFERENCE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROPAGATED_STAGE_INST_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BUSINESS_STATUS_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_EXEC_BUSKEY` (`BUSINESS_KEY_`),
  KEY `ACT_IDC_EXEC_ROOT` (`ROOT_PROC_INST_ID_`),
  KEY `ACT_IDX_EXEC_REF_ID_` (`REFERENCE_ID_`),
  KEY `ACT_FK_EXE_PROCINST` (`PROC_INST_ID_`),
  KEY `ACT_FK_EXE_PARENT` (`PARENT_ID_`),
  KEY `ACT_FK_EXE_SUPER` (`SUPER_EXEC_`),
  KEY `ACT_FK_EXE_PROCDEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_EXE_PARENT` FOREIGN KEY (`PARENT_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`) ON DELETE CASCADE,
  CONSTRAINT `ACT_FK_EXE_PROCDEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_EXE_PROCINST` FOREIGN KEY (`PROC_INST_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `ACT_FK_EXE_SUPER` FOREIGN KEY (`SUPER_EXEC_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_EXECUTION`
--

LOCK TABLES `ACT_RU_EXECUTION` WRITE;
/*!40000 ALTER TABLE `ACT_RU_EXECUTION` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_EXECUTION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_EXTERNAL_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_EXTERNAL_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_EXTERNAL_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `LOCK_EXP_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCLUSIVE_` tinyint(1) DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CORRELATION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RETRIES_` int DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DUEDATE_` timestamp(3) NULL DEFAULT NULL,
  `REPEAT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_EXTERNAL_JOB_EXCEPTION_STACK_ID` (`EXCEPTION_STACK_ID_`),
  KEY `ACT_IDX_EXTERNAL_JOB_CUSTOM_VALUES_ID` (`CUSTOM_VALUES_ID_`),
  KEY `ACT_IDX_EXTERNAL_JOB_CORRELATION_ID` (`CORRELATION_ID_`),
  KEY `ACT_IDX_EJOB_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_EJOB_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_EJOB_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  CONSTRAINT `ACT_FK_EXTERNAL_JOB_CUSTOM_VALUES` FOREIGN KEY (`CUSTOM_VALUES_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_EXTERNAL_JOB_EXCEPTION` FOREIGN KEY (`EXCEPTION_STACK_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_EXTERNAL_JOB`
--

LOCK TABLES `ACT_RU_EXTERNAL_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_EXTERNAL_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_EXTERNAL_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_HISTORY_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_HISTORY_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_HISTORY_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `LOCK_EXP_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RETRIES_` int DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ADV_HANDLER_CFG_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_HISTORY_JOB`
--

LOCK TABLES `ACT_RU_HISTORY_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_HISTORY_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_HISTORY_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_IDENTITYLINK`
--

DROP TABLE IF EXISTS `ACT_RU_IDENTITYLINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_IDENTITYLINK` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `GROUP_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `USER_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_IDENT_LNK_USER` (`USER_ID_`),
  KEY `ACT_IDX_IDENT_LNK_GROUP` (`GROUP_ID_`),
  KEY `ACT_IDX_IDENT_LNK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_IDENT_LNK_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_IDENT_LNK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_ATHRZ_PROCEDEF` (`PROC_DEF_ID_`),
  KEY `ACT_FK_TSKASS_TASK` (`TASK_ID_`),
  KEY `ACT_FK_IDL_PROCINST` (`PROC_INST_ID_`),
  CONSTRAINT `ACT_FK_ATHRZ_PROCEDEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_IDL_PROCINST` FOREIGN KEY (`PROC_INST_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_TSKASS_TASK` FOREIGN KEY (`TASK_ID_`) REFERENCES `ACT_RU_TASK` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_IDENTITYLINK`
--

LOCK TABLES `ACT_RU_IDENTITYLINK` WRITE;
/*!40000 ALTER TABLE `ACT_RU_IDENTITYLINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_IDENTITYLINK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `LOCK_EXP_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCLUSIVE_` tinyint(1) DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CORRELATION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RETRIES_` int DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DUEDATE_` timestamp(3) NULL DEFAULT NULL,
  `REPEAT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_JOB_EXCEPTION_STACK_ID` (`EXCEPTION_STACK_ID_`),
  KEY `ACT_IDX_JOB_CUSTOM_VALUES_ID` (`CUSTOM_VALUES_ID_`),
  KEY `ACT_IDX_JOB_CORRELATION_ID` (`CORRELATION_ID_`),
  KEY `ACT_IDX_JOB_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_JOB_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_JOB_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_JOB_EXECUTION` (`EXECUTION_ID_`),
  KEY `ACT_FK_JOB_PROCESS_INSTANCE` (`PROCESS_INSTANCE_ID_`),
  KEY `ACT_FK_JOB_PROC_DEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_JOB_CUSTOM_VALUES` FOREIGN KEY (`CUSTOM_VALUES_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_JOB_EXCEPTION` FOREIGN KEY (`EXCEPTION_STACK_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_JOB_EXECUTION` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_JOB_PROC_DEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_JOB_PROCESS_INSTANCE` FOREIGN KEY (`PROCESS_INSTANCE_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_JOB`
--

LOCK TABLES `ACT_RU_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_SUSPENDED_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_SUSPENDED_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_SUSPENDED_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `EXCLUSIVE_` tinyint(1) DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CORRELATION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RETRIES_` int DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DUEDATE_` timestamp(3) NULL DEFAULT NULL,
  `REPEAT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_SUSPENDED_JOB_EXCEPTION_STACK_ID` (`EXCEPTION_STACK_ID_`),
  KEY `ACT_IDX_SUSPENDED_JOB_CUSTOM_VALUES_ID` (`CUSTOM_VALUES_ID_`),
  KEY `ACT_IDX_SUSPENDED_JOB_CORRELATION_ID` (`CORRELATION_ID_`),
  KEY `ACT_IDX_SJOB_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_SJOB_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_SJOB_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_SUSPENDED_JOB_EXECUTION` (`EXECUTION_ID_`),
  KEY `ACT_FK_SUSPENDED_JOB_PROCESS_INSTANCE` (`PROCESS_INSTANCE_ID_`),
  KEY `ACT_FK_SUSPENDED_JOB_PROC_DEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_SUSPENDED_JOB_CUSTOM_VALUES` FOREIGN KEY (`CUSTOM_VALUES_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_SUSPENDED_JOB_EXCEPTION` FOREIGN KEY (`EXCEPTION_STACK_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_SUSPENDED_JOB_EXECUTION` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_SUSPENDED_JOB_PROC_DEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_SUSPENDED_JOB_PROCESS_INSTANCE` FOREIGN KEY (`PROCESS_INSTANCE_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_SUSPENDED_JOB`
--

LOCK TABLES `ACT_RU_SUSPENDED_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_SUSPENDED_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_SUSPENDED_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_TASK`
--

DROP TABLE IF EXISTS `ACT_RU_TASK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_TASK` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROPAGATED_STAGE_INST_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `STATE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `PARENT_TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DESCRIPTION_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_DEF_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ASSIGNEE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `DELEGATION_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PRIORITY_` int DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `IN_PROGRESS_TIME_` datetime(3) DEFAULT NULL,
  `IN_PROGRESS_STARTED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CLAIM_TIME_` datetime(3) DEFAULT NULL,
  `CLAIMED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUSPENDED_TIME_` datetime(3) DEFAULT NULL,
  `SUSPENDED_BY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `IN_PROGRESS_DUE_DATE_` datetime(3) DEFAULT NULL,
  `DUE_DATE_` datetime(3) DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUSPENSION_STATE_` int DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  `FORM_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `IS_COUNT_ENABLED_` tinyint DEFAULT NULL,
  `VAR_COUNT_` int DEFAULT NULL,
  `ID_LINK_COUNT_` int DEFAULT NULL,
  `SUB_TASK_COUNT_` int DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_TASK_CREATE` (`CREATE_TIME_`),
  KEY `ACT_IDX_TASK_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_TASK_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_TASK_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_TASK_EXE` (`EXECUTION_ID_`),
  KEY `ACT_FK_TASK_PROCINST` (`PROC_INST_ID_`),
  KEY `ACT_FK_TASK_PROCDEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_TASK_EXE` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_TASK_PROCDEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_TASK_PROCINST` FOREIGN KEY (`PROC_INST_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_TASK`
--

LOCK TABLES `ACT_RU_TASK` WRITE;
/*!40000 ALTER TABLE `ACT_RU_TASK` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_TASK` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_TIMER_JOB`
--

DROP TABLE IF EXISTS `ACT_RU_TIMER_JOB`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_TIMER_JOB` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `LOCK_EXP_TIME_` timestamp(3) NULL DEFAULT NULL,
  `LOCK_OWNER_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCLUSIVE_` tinyint(1) DEFAULT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROCESS_INSTANCE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_DEF_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `ELEMENT_NAME_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_DEFINITION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CORRELATION_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RETRIES_` int DEFAULT NULL,
  `EXCEPTION_STACK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `EXCEPTION_MSG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `DUEDATE_` timestamp(3) NULL DEFAULT NULL,
  `REPEAT_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `HANDLER_CFG_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `CUSTOM_VALUES_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` timestamp(3) NULL DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_TIMER_JOB_EXCEPTION_STACK_ID` (`EXCEPTION_STACK_ID_`),
  KEY `ACT_IDX_TIMER_JOB_CUSTOM_VALUES_ID` (`CUSTOM_VALUES_ID_`),
  KEY `ACT_IDX_TIMER_JOB_CORRELATION_ID` (`CORRELATION_ID_`),
  KEY `ACT_IDX_TIMER_JOB_DUEDATE` (`DUEDATE_`),
  KEY `ACT_IDX_TJOB_SCOPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_TJOB_SUB_SCOPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_TJOB_SCOPE_DEF` (`SCOPE_DEFINITION_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_TIMER_JOB_EXECUTION` (`EXECUTION_ID_`),
  KEY `ACT_FK_TIMER_JOB_PROCESS_INSTANCE` (`PROCESS_INSTANCE_ID_`),
  KEY `ACT_FK_TIMER_JOB_PROC_DEF` (`PROC_DEF_ID_`),
  CONSTRAINT `ACT_FK_TIMER_JOB_CUSTOM_VALUES` FOREIGN KEY (`CUSTOM_VALUES_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_TIMER_JOB_EXCEPTION` FOREIGN KEY (`EXCEPTION_STACK_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_TIMER_JOB_EXECUTION` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_TIMER_JOB_PROC_DEF` FOREIGN KEY (`PROC_DEF_ID_`) REFERENCES `ACT_RE_PROCDEF` (`ID_`),
  CONSTRAINT `ACT_FK_TIMER_JOB_PROCESS_INSTANCE` FOREIGN KEY (`PROCESS_INSTANCE_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_TIMER_JOB`
--

LOCK TABLES `ACT_RU_TIMER_JOB` WRITE;
/*!40000 ALTER TABLE `ACT_RU_TIMER_JOB` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_TIMER_JOB` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ACT_RU_VARIABLE`
--

DROP TABLE IF EXISTS `ACT_RU_VARIABLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ACT_RU_VARIABLE` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `EXECUTION_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `PROC_INST_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TASK_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BYTEARRAY_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `DOUBLE_` double DEFAULT NULL,
  `LONG_` bigint DEFAULT NULL,
  `TEXT_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `TEXT2_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  `META_INFO_` varchar(4000) COLLATE utf8mb3_bin DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  KEY `ACT_IDX_RU_VAR_SCOPE_ID_TYPE` (`SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_IDX_RU_VAR_SUB_ID_TYPE` (`SUB_SCOPE_ID_`,`SCOPE_TYPE_`),
  KEY `ACT_FK_VAR_BYTEARRAY` (`BYTEARRAY_ID_`),
  KEY `ACT_IDX_VARIABLE_TASK_ID` (`TASK_ID_`),
  KEY `ACT_FK_VAR_EXE` (`EXECUTION_ID_`),
  KEY `ACT_FK_VAR_PROCINST` (`PROC_INST_ID_`),
  CONSTRAINT `ACT_FK_VAR_BYTEARRAY` FOREIGN KEY (`BYTEARRAY_ID_`) REFERENCES `ACT_GE_BYTEARRAY` (`ID_`),
  CONSTRAINT `ACT_FK_VAR_EXE` FOREIGN KEY (`EXECUTION_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`),
  CONSTRAINT `ACT_FK_VAR_PROCINST` FOREIGN KEY (`PROC_INST_ID_`) REFERENCES `ACT_RU_EXECUTION` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ACT_RU_VARIABLE`
--

LOCK TABLES `ACT_RU_VARIABLE` WRITE;
/*!40000 ALTER TABLE `ACT_RU_VARIABLE` DISABLE KEYS */;
/*!40000 ALTER TABLE `ACT_RU_VARIABLE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_CHANNEL_DEFINITION`
--

DROP TABLE IF EXISTS `FLW_CHANNEL_DEFINITION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_CHANNEL_DEFINITION` (
  `ID_` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VERSION_` int DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TYPE_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `IMPLEMENTATION_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DEPLOYMENT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CREATE_TIME_` datetime(3) DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `ACT_IDX_CHANNEL_DEF_UNIQ` (`KEY_`,`VERSION_`,`TENANT_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_CHANNEL_DEFINITION`
--

LOCK TABLES `FLW_CHANNEL_DEFINITION` WRITE;
/*!40000 ALTER TABLE `FLW_CHANNEL_DEFINITION` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_CHANNEL_DEFINITION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_EVENT_DEFINITION`
--

DROP TABLE IF EXISTS `FLW_EVENT_DEFINITION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_EVENT_DEFINITION` (
  `ID_` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VERSION_` int DEFAULT NULL,
  `KEY_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DEPLOYMENT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID_`),
  UNIQUE KEY `ACT_IDX_EVENT_DEF_UNIQ` (`KEY_`,`VERSION_`,`TENANT_ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_EVENT_DEFINITION`
--

LOCK TABLES `FLW_EVENT_DEFINITION` WRITE;
/*!40000 ALTER TABLE `FLW_EVENT_DEFINITION` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_EVENT_DEFINITION` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_EVENT_DEPLOYMENT`
--

DROP TABLE IF EXISTS `FLW_EVENT_DEPLOYMENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_EVENT_DEPLOYMENT` (
  `ID_` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CATEGORY_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DEPLOY_TIME_` datetime(3) DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PARENT_DEPLOYMENT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_EVENT_DEPLOYMENT`
--

LOCK TABLES `FLW_EVENT_DEPLOYMENT` WRITE;
/*!40000 ALTER TABLE `FLW_EVENT_DEPLOYMENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_EVENT_DEPLOYMENT` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_EVENT_RESOURCE`
--

DROP TABLE IF EXISTS `FLW_EVENT_RESOURCE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_EVENT_RESOURCE` (
  `ID_` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DEPLOYMENT_ID_` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_BYTES_` longblob,
  PRIMARY KEY (`ID_`),
  KEY `FLW_IDX_EVENT_RSRC_DPL` (`DEPLOYMENT_ID_`),
  CONSTRAINT `FLW_FK_EVENT_RSRC_DPL` FOREIGN KEY (`DEPLOYMENT_ID_`) REFERENCES `FLW_EVENT_DEPLOYMENT` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_EVENT_RESOURCE`
--

LOCK TABLES `FLW_EVENT_RESOURCE` WRITE;
/*!40000 ALTER TABLE `FLW_EVENT_RESOURCE` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_EVENT_RESOURCE` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_RU_BATCH`
--

DROP TABLE IF EXISTS `FLW_RU_BATCH`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_RU_BATCH` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `TYPE_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `SEARCH_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SEARCH_KEY2_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` datetime(3) NOT NULL,
  `COMPLETE_TIME_` datetime(3) DEFAULT NULL,
  `STATUS_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `BATCH_DOC_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_RU_BATCH`
--

LOCK TABLES `FLW_RU_BATCH` WRITE;
/*!40000 ALTER TABLE `FLW_RU_BATCH` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_RU_BATCH` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `FLW_RU_BATCH_PART`
--

DROP TABLE IF EXISTS `FLW_RU_BATCH_PART`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FLW_RU_BATCH_PART` (
  `ID_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `REV_` int DEFAULT NULL,
  `BATCH_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TYPE_` varchar(64) COLLATE utf8mb3_bin NOT NULL,
  `SCOPE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SUB_SCOPE_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SCOPE_TYPE_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `SEARCH_KEY_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `SEARCH_KEY2_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `CREATE_TIME_` datetime(3) NOT NULL,
  `COMPLETE_TIME_` datetime(3) DEFAULT NULL,
  `STATUS_` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `RESULT_DOC_ID_` varchar(64) COLLATE utf8mb3_bin DEFAULT NULL,
  `TENANT_ID_` varchar(255) COLLATE utf8mb3_bin DEFAULT '',
  PRIMARY KEY (`ID_`),
  KEY `FLW_IDX_BATCH_PART` (`BATCH_ID_`),
  CONSTRAINT `FLW_FK_BATCH_PART_PARENT` FOREIGN KEY (`BATCH_ID_`) REFERENCES `FLW_RU_BATCH` (`ID_`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `FLW_RU_BATCH_PART`
--

LOCK TABLES `FLW_RU_BATCH_PART` WRITE;
/*!40000 ALTER TABLE `FLW_RU_BATCH_PART` DISABLE KEYS */;
/*!40000 ALTER TABLE `FLW_RU_BATCH_PART` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `flyway_schema_history`
--

DROP TABLE IF EXISTS `flyway_schema_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `flyway_schema_history` (
  `installed_rank` int NOT NULL,
  `version` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `script` varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `checksum` int DEFAULT NULL,
  `installed_by` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `installed_on` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `execution_time` int NOT NULL,
  `success` tinyint(1) NOT NULL,
  PRIMARY KEY (`installed_rank`),
  KEY `flyway_schema_history_s_idx` (`success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `flyway_schema_history`
--

LOCK TABLES `flyway_schema_history` WRITE;
/*!40000 ALTER TABLE `flyway_schema_history` DISABLE KEYS */;
INSERT INTO `flyway_schema_history` VALUES (1,'0','Initial schema for ioedream-consume-service','BASELINE','V0__Initial_schema.sql',NULL,'root','2025-12-27 07:39:31',0,1),(2,'1','create meal order tables','SQL','V1__create_meal_order_tables.sql',-1850981972,'root','2025-12-27 07:39:59',82,1),(3,'2','create report tables','SQL','V2__create_report_tables.sql',1700221491,'root','2025-12-27 07:40:56',64,1),(4,'3','create device quality tables','SQL','V3__create_device_quality_tables.sql',-993009288,'root','2025-12-27 07:47:46',71,1),(5,'20251219','ADD_CONSUME_ENTITY_FIELDS','SQL','V20251219__ADD_CONSUME_ENTITY_FIELDS.sql',NULL,'root','2025-12-27 07:49:43',0,1),(6,'20251223.1','create account compensation table','SQL','V20251223.1__create_account_compensation_table.sql',NULL,'root','2025-12-27 07:49:43',0,1),(7,'20251223.2','create consume account table','SQL','V20251223.2__create_consume_account_table.sql',NULL,'root','2025-12-27 07:49:43',0,1),(8,'20251223.3','create consume account transaction table','SQL','V20251223.3__create_consume_account_transaction_table.sql',NULL,'root','2025-12-27 07:49:43',0,1),(9,'20251223.4','create consume record table','SQL','V20251223.4__create_consume_record_table.sql',NULL,'root','2025-12-27 07:49:43',0,1),(10,'20251223.5','create dual write validation tables','SQL','V20251223.5__create_dual_write_validation_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(11,'20251223.6','create seata undo log','SQL','V20251223.6__create_seata_undo_log.sql',NULL,'root','2025-12-27 07:49:43',0,1),(12,'20251223.7','create POSID tables','SQL','V20251223.7__create_POSID_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(13,'20251223.8','migrate to POSID tables','SQL','V20251223.8__migrate_to_POSID_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(14,'20251225','create offline consume tables','SQL','V20251225__create_offline_consume_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(15,'20251226.1','create reconciliation tables','SQL','V20251226.1__001_create_reconciliation_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(16,'20251226.2','create subsidy rule tables','SQL','V20251226.2__create_subsidy_rule_tables.sql',NULL,'root','2025-12-27 07:49:43',0,1),(17,'20251227','add performance indexes','SQL','V20251227__add_performance_indexes.sql',NULL,'root','2025-12-27 07:49:43',0,1),(18,'20250130','create device ai event table','SQL','V20250130__create_device_ai_event_table.sql',1633489183,'root','2025-12-29 12:08:11',156,1),(19,'20250131','create alarm tables','SQL','V20250131__create_alarm_tables.sql',1149648336,'root','2025-12-29 12:08:11',166,1),(20,'20250201','P1 AI Model Management','SQL','V20250201__P1_AI_Model_Management.sql',1376650128,'root','2025-12-29 12:08:12',197,1);
/*!40000 ALTER TABLE `flyway_schema_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_access_permission`
--

DROP TABLE IF EXISTS `t_access_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_access_permission` (
  `permission_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'æƒé™ID',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·ID',
  `area_id` bigint NOT NULL COMMENT 'åŒºåŸŸID',
  `device_id` bigint DEFAULT NULL COMMENT 'è®¾å¤‡ID',
  `permission_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'ALWAYS' COMMENT 'æƒé™ç±»åž‹ï¼šALWAYS-æ°¸ä¹… TIME_LIMITED-é™æ—¶',
  `start_time` datetime DEFAULT NULL COMMENT 'ç”Ÿæ•ˆå¼€å§‹æ—¶é—´',
  `end_time` datetime DEFAULT NULL COMMENT 'ç”Ÿæ•ˆç»“æŸæ—¶é—´',
  `access_times` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'é€šè¡Œæ—¶é—´æ®µï¼ˆJSONï¼‰',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æœ‰æ•ˆ 2-å¤±æ•ˆ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`permission_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_permission_type` (`permission_type`),
  KEY `idx_status` (`status`),
  KEY `idx_start_end_time` (`start_time`,`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='é—¨ç¦æƒé™è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_access_permission`
--

LOCK TABLES `t_access_permission` WRITE;
/*!40000 ALTER TABLE `t_access_permission` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_access_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_access_record`
--

DROP TABLE IF EXISTS `t_access_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_access_record` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®°å½•ID',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·ID',
  `device_id` bigint NOT NULL COMMENT 'è®¾å¤‡ID',
  `area_id` bigint DEFAULT NULL COMMENT 'åŒºåŸŸID',
  `access_result` tinyint NOT NULL COMMENT 'é€šè¡Œç»“æžœï¼š1-æˆåŠŸ 2-å¤±è´¥',
  `access_time` datetime NOT NULL COMMENT 'é€šè¡Œæ—¶é—´',
  `access_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'IN' COMMENT 'é€šè¡Œç±»åž‹ï¼šIN-è¿›å…¥ OUT-ç¦»å¼€',
  `verify_method` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'éªŒè¯æ–¹å¼ï¼šFACE-äººè„¸ CARD-åˆ·å¡ FINGERPRINT-æŒ‡çº¹',
  `photo_path` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç…§ç‰‡è·¯å¾„',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`record_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_access_result` (`access_result`),
  KEY `idx_access_time` (`access_time`),
  KEY `idx_access_type` (`access_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='é—¨ç¦è®°å½•è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_access_record`
--

LOCK TABLES `t_access_record` WRITE;
/*!40000 ALTER TABLE `t_access_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_access_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_account_compensation`
--

DROP TABLE IF EXISTS `t_account_compensation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_account_compensation` (
  `compensation_id` bigint NOT NULL AUTO_INCREMENT COMMENT '补偿记录ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `operation` varchar(20) NOT NULL COMMENT '操作类型（INCREASE-增加, DECREASE-扣减）',
  `amount` decimal(10,2) NOT NULL COMMENT '金额',
  `business_type` varchar(50) NOT NULL COMMENT '业务类型（SUBSIDY_GRANT, CONSUME等）',
  `business_no` varchar(100) NOT NULL COMMENT '业务编号（用于幂等性控制）',
  `related_business_no` varchar(100) DEFAULT NULL COMMENT '关联业务编号',
  `status` varchar(20) NOT NULL DEFAULT 'PENDING' COMMENT '状态（PENDING-待处理, SUCCESS-成功, FAILED-失败, CANCELLED-已取消）',
  `retry_count` int NOT NULL DEFAULT '0' COMMENT '重试次数',
  `max_retry_count` int NOT NULL DEFAULT '3' COMMENT '最大重试次数',
  `error_code` varchar(50) DEFAULT NULL COMMENT '错误码',
  `error_message` varchar(500) DEFAULT NULL COMMENT '错误信息',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `extended_params` text COMMENT '扩展参数（JSON格式）',
  `next_retry_time` datetime DEFAULT NULL COMMENT '下次重试时间',
  `last_retry_time` datetime DEFAULT NULL COMMENT '最后重试时间',
  `success_time` datetime DEFAULT NULL COMMENT '成功时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记（0-未删除，1-已删除）',
  PRIMARY KEY (`compensation_id`),
  UNIQUE KEY `uk_business_no` (`business_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_business_no` (`business_no`),
  KEY `idx_status` (`status`),
  KEY `idx_next_retry_time` (`next_retry_time`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_operation_status` (`operation`,`status`),
  KEY `idx_retry_count` (`retry_count`,`max_retry_count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='账户服务补偿记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_account_compensation`
--

LOCK TABLES `t_account_compensation` WRITE;
/*!40000 ALTER TABLE `t_account_compensation` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_account_compensation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_api_compatibility_validation`
--

DROP TABLE IF EXISTS `t_api_compatibility_validation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_api_compatibility_validation` (
  `validation_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'éªŒè¯ID',
  `validation_date` date NOT NULL COMMENT 'éªŒè¯æ—¥æœŸ',
  `module_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'æ¨¡å—åç§°',
  `api_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'APIåç§°',
  `api_method` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'HTTPæ–¹æ³•',
  `api_path` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'APIè·¯å¾„',
  `response_format_compatible` tinyint DEFAULT '1' COMMENT 'å“åº”æ ¼å¼å…¼å®¹ï¼š1-å…¼å®¹ 0-ä¸å…¼å®¹',
  `response_structure_match` tinyint DEFAULT '1' COMMENT 'å“åº”ç»“æž„åŒ¹é…ï¼š1-åŒ¹é… 0-ä¸åŒ¹é…',
  `field_completeness_rate` decimal(5,2) DEFAULT '100.00' COMMENT 'å­—æ®µå®Œæ•´çŽ‡ï¼ˆ%ï¼‰',
  `entity_field_coverage` decimal(5,2) DEFAULT '100.00' COMMENT 'å®žä½“å­—æ®µè¦†ç›–çŽ‡ï¼ˆ%ï¼‰',
  `table_field_coverage` decimal(5,2) DEFAULT '100.00' COMMENT 'è¡¨å­—æ®µè¦†ç›–çŽ‡ï¼ˆ%ï¼‰',
  `data_type_consistency` tinyint DEFAULT '1' COMMENT 'æ•°æ®ç±»åž‹ä¸€è‡´æ€§ï¼š1-ä¸€è‡´ 0-ä¸ä¸€è‡´',
  `business_logic_compatible` tinyint DEFAULT '1' COMMENT 'ä¸šåŠ¡é€»è¾‘å…¼å®¹ï¼š1-å…¼å®¹ 0-ä¸å…¼å®¹',
  `workflow_compatible` tinyint DEFAULT '1' COMMENT 'å·¥ä½œæµå…¼å®¹ï¼š1-å…¼å®¹ 0-ä¸å…¼å®¹',
  `query_performance_acceptable` tinyint DEFAULT '1' COMMENT 'æŸ¥è¯¢æ€§èƒ½å¯æŽ¥å—ï¼š1-æ˜¯ 0-å¦',
  `index_optimization_complete` tinyint DEFAULT '1' COMMENT 'ç´¢å¼•ä¼˜åŒ–å®Œæˆï¼š1-æ˜¯ 0-å¦',
  `overall_compatibility` decimal(5,2) DEFAULT '100.00' COMMENT 'æ•´ä½“å…¼å®¹çŽ‡ï¼ˆ%ï¼‰',
  `validation_status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'PASS' COMMENT 'éªŒè¯çŠ¶æ€ï¼šPASS-é€šè¿‡ FAIL-å¤±è´¥ PARTIAL-éƒ¨åˆ†é€šè¿‡',
  `issue_description` text COLLATE utf8mb4_unicode_ci COMMENT 'é—®é¢˜æè¿°',
  `fix_suggestions` text COLLATE utf8mb4_unicode_ci COMMENT 'ä¿®å¤å»ºè®®',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  PRIMARY KEY (`validation_id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='APIå…¼å®¹æ€§éªŒè¯ç»“æžœè¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_api_compatibility_validation`
--

LOCK TABLES `t_api_compatibility_validation` WRITE;
/*!40000 ALTER TABLE `t_api_compatibility_validation` DISABLE KEYS */;
INSERT INTO `t_api_compatibility_validation` VALUES (1,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ¶ˆè´¹è®°å½•æŸ¥è¯¢','GET','/api/consume/record/list',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(2,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ¶ˆè´¹è®°å½•è¯¦æƒ…','GET','/api/consume/record/{id}',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(3,'2025-12-15','æ¶ˆè´¹ç®¡ç†','åˆ›å»ºæ¶ˆè´¹è®°å½•','POST','/api/consume/record/create',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(4,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·ä¿¡æ¯æŸ¥è¯¢','GET','/api/consume/account/{accountId}',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(5,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·ä½™é¢æŸ¥è¯¢','GET','/api/consume/account/{userId}/balance',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(6,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·å†»ç»“/è§£å†»','POST','/api/consume/account/{accountId}/freeze',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(7,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·å……å€¼','POST','/api/consume/account/{accountId}/recharge',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(8,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è®¾ç½®è´¦æˆ·é™é¢','POST','/api/consume/account/{accountId}/limit',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(9,'2025-12-15','æ¶ˆè´¹ç®¡ç†','ç”³è¯·é€€æ¬¾','POST','/api/consume/refund/apply',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(10,'2025-12-15','æ¶ˆè´¹ç®¡ç†','é€€æ¬¾è®°å½•æŸ¥è¯¢','GET','/api/consume/refund/list',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(11,'2025-12-15','æ¶ˆè´¹ç®¡ç†','é€€æ¬¾å®¡æ‰¹','POST','/api/consume/refund/{refundId}/approve',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(12,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ‰¹é‡é€€æ¬¾ç”³è¯·','POST','/api/consume/refund/batch/apply',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(13,'2025-12-15','å…¬å…±æ¨¡å—','ç”¨æˆ·ç™»å½•','POST','/api/auth/login',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(14,'2025-12-15','å…¬å…±æ¨¡å—','ç”¨æˆ·ç™»å‡º','POST','/api/auth/logout',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(15,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·ä¿¡æ¯','GET','/api/auth/info',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(16,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·æƒé™','GET','/api/auth/permissions',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(17,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·è§’è‰²','GET','/api/auth/roles',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(18,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·èœå•','GET','/api/auth/menus',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(19,'2025-12-15','å…¬å…±æ¨¡å—','åˆ·æ–°ä»¤ç‰Œ','POST','/api/auth/refresh',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(20,'2025-12-15','å…¬å…±æ¨¡å—','éªŒè¯ä»¤ç‰Œ','POST','/api/auth/validate',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:20:56','2025-12-15 20:20:56'),(21,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ¶ˆè´¹è®°å½•æŸ¥è¯¢','GET','/api/consume/record/list',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(22,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ¶ˆè´¹è®°å½•è¯¦æƒ…','GET','/api/consume/record/{id}',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(23,'2025-12-15','æ¶ˆè´¹ç®¡ç†','åˆ›å»ºæ¶ˆè´¹è®°å½•','POST','/api/consume/record/create',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(24,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·ä¿¡æ¯æŸ¥è¯¢','GET','/api/consume/account/{accountId}',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(25,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·ä½™é¢æŸ¥è¯¢','GET','/api/consume/account/{userId}/balance',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(26,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·å†»ç»“/è§£å†»','POST','/api/consume/account/{accountId}/freeze',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(27,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è´¦æˆ·å……å€¼','POST','/api/consume/account/{accountId}/recharge',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(28,'2025-12-15','æ¶ˆè´¹ç®¡ç†','è®¾ç½®è´¦æˆ·é™é¢','POST','/api/consume/account/{accountId}/limit',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(29,'2025-12-15','æ¶ˆè´¹ç®¡ç†','ç”³è¯·é€€æ¬¾','POST','/api/consume/refund/apply',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(30,'2025-12-15','æ¶ˆè´¹ç®¡ç†','é€€æ¬¾è®°å½•æŸ¥è¯¢','GET','/api/consume/refund/list',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(31,'2025-12-15','æ¶ˆè´¹ç®¡ç†','é€€æ¬¾å®¡æ‰¹','POST','/api/consume/refund/{refundId}/approve',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(32,'2025-12-15','æ¶ˆè´¹ç®¡ç†','æ‰¹é‡é€€æ¬¾ç”³è¯·','POST','/api/consume/refund/batch/apply',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(33,'2025-12-15','å…¬å…±æ¨¡å—','ç”¨æˆ·ç™»å½•','POST','/api/auth/login',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(34,'2025-12-15','å…¬å…±æ¨¡å—','ç”¨æˆ·ç™»å‡º','POST','/api/auth/logout',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(35,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·ä¿¡æ¯','GET','/api/auth/info',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(36,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·æƒé™','GET','/api/auth/permissions',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(37,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·è§’è‰²','GET','/api/auth/roles',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(38,'2025-12-15','å…¬å…±æ¨¡å—','èŽ·å–ç”¨æˆ·èœå•','GET','/api/auth/menus',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(39,'2025-12-15','å…¬å…±æ¨¡å—','åˆ·æ–°ä»¤ç‰Œ','POST','/api/auth/refresh',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16'),(40,'2025-12-15','å…¬å…±æ¨¡å—','éªŒè¯ä»¤ç‰Œ','POST','/api/auth/validate',1,1,100.00,100.00,100.00,1,1,1,1,1,100.00,'PASS',NULL,NULL,'2025-12-15 20:44:16','2025-12-15 20:44:16');
/*!40000 ALTER TABLE `t_api_compatibility_validation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_attendance_record`
--

DROP TABLE IF EXISTS `t_attendance_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_attendance_record` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®°å½•ID',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·ID',
  `device_id` bigint NOT NULL COMMENT 'è®¾å¤‡ID',
  `attendance_date` date NOT NULL COMMENT 'è€ƒå‹¤æ—¥æœŸ',
  `clock_time` datetime NOT NULL COMMENT 'æ‰“å¡æ—¶é—´',
  `clock_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'æ‰“å¡ç±»åž‹ï¼šON_DUTY-ä¸Šç­ OFF_DUTY-ä¸‹ç­',
  `verify_method` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'éªŒè¯æ–¹å¼',
  `location` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æ‰“å¡ä½ç½®',
  `photo_path` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç…§ç‰‡è·¯å¾„',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`record_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_attendance_date` (`attendance_date`),
  KEY `idx_clock_time` (`clock_time`),
  KEY `idx_clock_type` (`clock_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è€ƒå‹¤è®°å½•è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_attendance_record`
--

LOCK TABLES `t_attendance_record` WRITE;
/*!40000 ALTER TABLE `t_attendance_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_attendance_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_attendance_shift`
--

DROP TABLE IF EXISTS `t_attendance_shift`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_attendance_shift` (
  `shift_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ç­æ¬¡ID',
  `shift_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ç­æ¬¡åç§°',
  `work_start_time` time NOT NULL COMMENT 'ä¸Šç­æ—¶é—´',
  `work_end_time` time NOT NULL COMMENT 'ä¸‹ç­æ—¶é—´',
  `break_start_time` time DEFAULT NULL COMMENT 'ä¼‘æ¯å¼€å§‹æ—¶é—´',
  `break_end_time` time DEFAULT NULL COMMENT 'ä¼‘æ¯ç»“æŸæ—¶é—´',
  `work_days` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '1,2,3,4,5' COMMENT 'å·¥ä½œæ—¥ï¼ˆ1-7è¡¨ç¤ºå‘¨ä¸€åˆ°å‘¨æ—¥ï¼‰',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`shift_id`),
  KEY `idx_shift_name` (`shift_name`),
  KEY `idx_work_times` (`work_start_time`,`work_end_time`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è€ƒå‹¤ç­æ¬¡è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_attendance_shift`
--

LOCK TABLES `t_attendance_shift` WRITE;
/*!40000 ALTER TABLE `t_attendance_shift` DISABLE KEYS */;
INSERT INTO `t_attendance_shift` VALUES (1,'æ ‡å‡†ç­æ¬¡','09:00:00','18:00:00','12:00:00','13:00:00','1,2,3,4,5',1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'æ—©ç­','08:00:00','17:00:00','12:00:00','13:00:00','1,2,3,4,5',1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'æ™šç­','10:00:00','19:00:00','12:30:00','13:30:00','1,2,3,4,5',1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'å‘¨æœ«ç­','09:00:00','18:00:00','12:00:00','13:00:00','6,7',1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'æ ‡å‡†ç­æ¬¡','09:00:00','18:00:00','12:00:00','13:00:00','1,2,3,4,5',1,'2025-12-15 20:44:16','2025-12-15 20:44:16',NULL,NULL,0,0),(6,'æ—©ç­','08:00:00','17:00:00','12:00:00','13:00:00','1,2,3,4,5',1,'2025-12-15 20:44:16','2025-12-15 20:44:16',NULL,NULL,0,0),(7,'æ™šç­','10:00:00','19:00:00','12:30:00','13:30:00','1,2,3,4,5',1,'2025-12-15 20:44:16','2025-12-15 20:44:16',NULL,NULL,0,0),(8,'å‘¨æœ«ç­','09:00:00','18:00:00','12:00:00','13:00:00','6,7',1,'2025-12-15 20:44:16','2025-12-15 20:44:16',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_attendance_shift` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_audit_log`
--

DROP TABLE IF EXISTS `t_audit_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_audit_log` (
  `log_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'æ—¥å¿—ID',
  `user_id` bigint DEFAULT NULL COMMENT 'ç”¨æˆ·ID',
  `username` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç”¨æˆ·å',
  `operation` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'æ“ä½œç±»åž‹',
  `method` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è¯·æ±‚æ–¹æ³•',
  `params` text COLLATE utf8mb4_unicode_ci COMMENT 'è¯·æ±‚å‚æ•°',
  `time` bigint DEFAULT NULL COMMENT 'æ‰§è¡Œæ—¶é•¿ï¼ˆæ¯«ç§’ï¼‰',
  `ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IPåœ°å€',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç”¨æˆ·ä»£ç†',
  `result` tinyint DEFAULT NULL COMMENT 'æ“ä½œç»“æžœï¼š1-æˆåŠŸ 0-å¤±è´¥',
  `error_msg` text COLLATE utf8mb4_unicode_ci COMMENT 'é”™è¯¯ä¿¡æ¯',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`log_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_operation` (`operation`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='æ“ä½œæ—¥å¿—è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_audit_log`
--

LOCK TABLES `t_audit_log` WRITE;
/*!40000 ALTER TABLE `t_audit_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_audit_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_area`
--

DROP TABLE IF EXISTS `t_common_area`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_area` (
  `area_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'åŒºåŸŸID',
  `area_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'åŒºåŸŸåç§°',
  `area_code` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åŒºåŸŸç¼–ç ',
  `area_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'BUILDING' COMMENT 'åŒºåŸŸç±»åž‹ï¼šCAMPUS-å›­åŒº BUILDING-å»ºç­‘ FLOOR-æ¥¼å±‚ ROOM-æˆ¿é—´',
  `parent_id` bigint DEFAULT '0' COMMENT 'çˆ¶åŒºåŸŸID',
  `level` int DEFAULT '1' COMMENT 'å±‚çº§',
  `sort_order` int DEFAULT '0' COMMENT 'æŽ’åº',
  `manager_id` bigint DEFAULT NULL COMMENT 'ç®¡ç†å‘˜ID',
  `capacity` int DEFAULT NULL COMMENT 'å®¹é‡',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'æè¿°',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`area_id`),
  UNIQUE KEY `uk_area_code` (`area_code`,`deleted_flag`),
  KEY `idx_area_name` (`area_name`),
  KEY `idx_area_type` (`area_type`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_manager_id` (`manager_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='åŒºåŸŸè¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_area`
--

LOCK TABLES `t_common_area` WRITE;
/*!40000 ALTER TABLE `t_common_area` DISABLE KEYS */;
INSERT INTO `t_common_area` VALUES (1,'IOE-DREAMæ™ºæ…§å›­åŒº','IOE_DREAM_CAMPUS','CAMPUS',0,1,1,NULL,10000,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'Aæ ‹åŠžå…¬æ¥¼','BUILDING_A','BUILDING',1,2,1,NULL,2000,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'Bæ ‹åŠžå…¬æ¥¼','BUILDING_B','BUILDING',1,2,2,NULL,1500,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'Aæ ‹1æ¥¼','BUILDING_A_FLOOR_1','FLOOR',2,3,1,NULL,500,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'Aæ ‹2æ¥¼','BUILDING_A_FLOOR_2','FLOOR',2,3,2,NULL,500,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(6,'Bæ ‹1æ¥¼','BUILDING_B_FLOOR_1','FLOOR',3,3,1,NULL,400,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(7,'é£Ÿå ‚','CANTEEN','BUILDING',1,2,3,NULL,800,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(8,'ä¼šè®®å®¤A','MEETING_ROOM_A','ROOM',6,4,1,NULL,50,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(9,'ä¼šè®®å®¤B','MEETING_ROOM_B','ROOM',6,4,2,NULL,30,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_area` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_department`
--

DROP TABLE IF EXISTS `t_common_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_department` (
  `department_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'éƒ¨é—¨ID',
  `department_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'éƒ¨é—¨åç§°',
  `department_code` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'éƒ¨é—¨ç¼–ç ',
  `parent_id` bigint DEFAULT '0' COMMENT 'çˆ¶éƒ¨é—¨ID',
  `level` int DEFAULT '1' COMMENT 'å±‚çº§',
  `sort_order` int DEFAULT '0' COMMENT 'æŽ’åº',
  `leader_id` bigint DEFAULT NULL COMMENT 'è´Ÿè´£äººID',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è”ç³»ç”µè¯',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'é‚®ç®±',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'æè¿°',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`department_id`),
  UNIQUE KEY `uk_department_code` (`department_code`,`deleted_flag`),
  KEY `idx_department_name` (`department_name`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_leader_id` (`leader_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='éƒ¨é—¨è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_department`
--

LOCK TABLES `t_common_department` WRITE;
/*!40000 ALTER TABLE `t_common_department` DISABLE KEYS */;
INSERT INTO `t_common_department` VALUES (1,'IOE-DREAMé›†å›¢','IOE_DREAM',0,1,1,NULL,NULL,NULL,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'æŠ€æœ¯éƒ¨','TECH_DEPT',1,2,1,NULL,NULL,NULL,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'è¿è¥éƒ¨','OPS_DEPT',1,2,2,NULL,NULL,NULL,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'è¡Œæ”¿éƒ¨','ADMIN_DEPT',1,2,3,NULL,NULL,NULL,NULL,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_device`
--

DROP TABLE IF EXISTS `t_common_device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_device` (
  `device_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®¾å¤‡ID',
  `device_no` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¾å¤‡ç¼–å·',
  `device_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è®¾å¤‡åç§°',
  `device_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¾å¤‡ç±»åž‹ï¼šCAMERA-æ‘„åƒå¤´ ACCESS-é—¨ç¦ CONSUME-æ¶ˆè´¹æœº ATTENDANCE-è€ƒå‹¤æœº BIOMETRIC-ç”Ÿç‰©è¯†åˆ«',
  `device_model` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è®¾å¤‡åž‹å·',
  `manufacturer` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åˆ¶é€ å•†',
  `serial_number` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åºåˆ—å·',
  `area_id` bigint DEFAULT NULL COMMENT 'æ‰€åœ¨åŒºåŸŸID',
  `ip_address` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IPåœ°å€',
  `port` int DEFAULT NULL COMMENT 'ç«¯å£',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-åœ¨çº¿ 2-ç¦»çº¿ 3-æ•…éšœ 4-ç»´æŠ¤',
  `install_time` datetime DEFAULT NULL COMMENT 'å®‰è£…æ—¶é—´',
  `last_active_time` datetime DEFAULT NULL COMMENT 'æœ€åŽæ´»è·ƒæ—¶é—´',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'è®¾å¤‡æè¿°',
  `extended_attributes` text COLLATE utf8mb4_unicode_ci COMMENT 'æ‰©å±•å±žæ€§JSON',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`device_id`),
  UNIQUE KEY `uk_device_no` (`device_no`,`deleted_flag`),
  KEY `idx_device_name` (`device_name`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`),
  KEY `idx_last_active_time` (`last_active_time`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è®¾å¤‡è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_device`
--

LOCK TABLES `t_common_device` WRITE;
/*!40000 ALTER TABLE `t_common_device` DISABLE KEYS */;
INSERT INTO `t_common_device` VALUES (1,'CAM001','å¤§é—¨æ‘„åƒå¤´1','CAMERA','Hikvision DS-2CD2042G2-I','æµ·åº·å¨è§†',NULL,5,'192.168.1.101',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'CAM002','å¤§é—¨æ‘„åƒå¤´2','CAMERA','Hikvision DS-2CD2042G2-I','æµ·åº·å¨è§†',NULL,5,'192.168.1.102',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'CAM003','é£Ÿå ‚æ‘„åƒå¤´','CAMERA','Dahua DH-IPC-HFW2431S-S','å¤§åŽæŠ€æœ¯',NULL,8,'192.168.1.103',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'ACC001','å¤§é—¨é—¨ç¦è®¾å¤‡','ACCESS','ZKTeco SC405','ä¸­æŽ§æ™ºæ…§',NULL,5,'192.168.1.201',4370,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'ACC002','Aæ ‹é—¨ç¦è®¾å¤‡','ACCESS','ZKTeco SC405','ä¸­æŽ§æ™ºæ…§',NULL,6,'192.168.1.202',4370,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(6,'ACC003','Bæ ‹é—¨ç¦è®¾å¤‡','ACCESS','ZKTeco SC405','ä¸­æŽ§æ™ºæ…§',NULL,7,'192.168.1.203',4370,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(7,'CON001','é£Ÿå ‚æ¶ˆè´¹æœº1','CONSUME','ZKTeco F18','ä¸­æŽ§æ™ºæ…§',NULL,8,'192.168.1.301',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(8,'CON002','é£Ÿå ‚æ¶ˆè´¹æœº2','CONSUME','ZKTeco F18','ä¸­æŽ§æ™ºæ…§',NULL,8,'192.168.1.302',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(9,'ATT001','è€ƒå‹¤æœº1','ATTENDANCE','ZKTeco MB40','ä¸­æŽ§æ™ºæ…§',NULL,5,'192.168.1.401',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(10,'ATT002','è€ƒå‹¤æœº2','ATTENDANCE','ZKTeco MB40','ä¸­æŽ§æ™ºæ…§',NULL,6,'192.168.1.402',80,1,NULL,NULL,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_device` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_dict_data`
--

DROP TABLE IF EXISTS `t_common_dict_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_dict_data` (
  `dict_data_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'å­—å…¸æ•°æ®ID',
  `dict_type_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å­—å…¸ç±»åž‹ç¼–ç ',
  `dict_label` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å­—å…¸æ ‡ç­¾',
  `dict_value` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å­—å…¸å€¼',
  `dict_sort` int DEFAULT '0' COMMENT 'æŽ’åº',
  `css_class` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æ ·å¼ç±»',
  `list_class` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åˆ—è¡¨æ ·å¼',
  `is_default` tinyint DEFAULT '0' COMMENT 'æ˜¯å¦é»˜è®¤ï¼š0-å¦ 1-æ˜¯',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`dict_data_id`),
  UNIQUE KEY `uk_dict_type_value` (`dict_type_code`,`dict_value`,`deleted_flag`),
  KEY `idx_dict_type_code` (`dict_type_code`),
  KEY `idx_dict_label` (`dict_label`),
  KEY `idx_dict_sort` (`dict_sort`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='å­—å…¸æ•°æ®è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_dict_data`
--

LOCK TABLES `t_common_dict_data` WRITE;
/*!40000 ALTER TABLE `t_common_dict_data` DISABLE KEYS */;
INSERT INTO `t_common_dict_data` VALUES (1,'USER_GENDER','ç”·','1',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'USER_GENDER','å¥³','2',2,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'USER_STATUS','æ­£å¸¸','1',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'USER_STATUS','ç¦ç”¨','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'ROLE_STATUS','æ­£å¸¸','1',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(6,'ROLE_STATUS','ç¦ç”¨','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(7,'PERMISSION_STATUS','æ­£å¸¸','1',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(8,'PERMISSION_STATUS','ç¦ç”¨','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(9,'AREA_TYPE','å›­åŒº','CAMPUS',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(10,'AREA_TYPE','å»ºç­‘','BUILDING',2,'','info',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(11,'AREA_TYPE','æ¥¼å±‚','FLOOR',3,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(12,'AREA_TYPE','æˆ¿é—´','ROOM',4,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(13,'DEVICE_TYPE','æ‘„åƒå¤´','CAMERA',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(14,'DEVICE_TYPE','é—¨ç¦è®¾å¤‡','ACCESS',2,'','info',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(15,'DEVICE_TYPE','æ¶ˆè´¹æœº','CONSUME',3,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(16,'DEVICE_TYPE','è€ƒå‹¤æœº','ATTENDANCE',4,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(17,'DEVICE_TYPE','ç”Ÿç‰©è¯†åˆ«è®¾å¤‡','BIOMETRIC',5,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(18,'DEVICE_STATUS','åœ¨çº¿','1',1,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(19,'DEVICE_STATUS','ç¦»çº¿','2',2,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(20,'DEVICE_STATUS','æ•…éšœ','3',3,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(21,'DEVICE_STATUS','ç»´æŠ¤','4',4,'','info',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(22,'CONSUME_STATUS','æˆåŠŸ','SUCCESS',1,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(23,'CONSUME_STATUS','å¤±è´¥','FAILED',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(24,'CONSUME_STATUS','å¤„ç†ä¸­','PENDING',3,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(25,'CONSUME_STATUS','å·²å–æ¶ˆ','CANCELLED',4,'','info',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(26,'ACCOUNT_STATUS','æ­£å¸¸','1',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(27,'ACCOUNT_STATUS','å†»ç»“','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(28,'ACCOUNT_STATUS','æ³¨é”€','3',3,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(29,'ACCESS_PERMISSION_TYPE','æ°¸ä¹…æƒé™','ALWAYS',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(30,'ACCESS_PERMISSION_TYPE','é™æ—¶æƒé™','TIME_LIMITED',2,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(31,'ACCESS_RESULT','æˆåŠŸ','1',1,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(32,'ACCESS_RESULT','å¤±è´¥','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(33,'ATTENDANCE_CLOCK_TYPE','ä¸Šç­','ON_DUTY',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(34,'ATTENDANCE_CLOCK_TYPE','ä¸‹ç­','OFF_DUTY',2,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(35,'VISITOR_STATUS','å¾…å®¡æ‰¹','PENDING',1,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(36,'VISITOR_STATUS','å·²å®¡æ‰¹','APPROVED',2,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(37,'VISITOR_STATUS','å·²æ‹’ç»','REJECTED',3,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(38,'VISITOR_STATUS','å·²å®Œæˆ','COMPLETED',4,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(39,'VIDEO_DEVICE_TYPE','æ‘„åƒå¤´','CAMERA',1,'','primary',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(40,'VIDEO_DEVICE_TYPE','ç½‘ç»œå½•åƒæœº','NVR',2,'','info',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(41,'VIDEO_DEVICE_TYPE','ç¡¬ç›˜å½•åƒæœº','DVR',3,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(42,'SYSTEM_STATUS','æ­£å¸¸','1',1,'','success',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(43,'SYSTEM_STATUS','å¼‚å¸¸','2',2,'','danger',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(44,'SYSTEM_STATUS','ç»´æŠ¤','3',3,'','warning',0,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_dict_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_dict_type`
--

DROP TABLE IF EXISTS `t_common_dict_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_dict_type` (
  `dict_type_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'å­—å…¸ç±»åž‹ID',
  `dict_type_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å­—å…¸ç±»åž‹åç§°',
  `dict_type_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å­—å…¸ç±»åž‹ç¼–ç ',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'æè¿°',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `sort_order` int DEFAULT '0' COMMENT 'æŽ’åº',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`dict_type_id`),
  UNIQUE KEY `uk_dict_type_code` (`dict_type_code`,`deleted_flag`),
  KEY `idx_dict_type_name` (`dict_type_name`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='å­—å…¸ç±»åž‹è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_dict_type`
--

LOCK TABLES `t_common_dict_type` WRITE;
/*!40000 ALTER TABLE `t_common_dict_type` DISABLE KEYS */;
INSERT INTO `t_common_dict_type` VALUES (1,'ç”¨æˆ·æ€§åˆ«','USER_GENDER','ç”¨æˆ·æ€§åˆ«å­—å…¸',1,1,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(2,'ç”¨æˆ·çŠ¶æ€','USER_STATUS','ç”¨æˆ·çŠ¶æ€å­—å…¸',1,2,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(3,'è§’è‰²çŠ¶æ€','ROLE_STATUS','è§’è‰²çŠ¶æ€å­—å…¸',1,3,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(4,'æƒé™çŠ¶æ€','PERMISSION_STATUS','æƒé™çŠ¶æ€å­—å…¸',1,4,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(5,'åŒºåŸŸç±»åž‹','AREA_TYPE','åŒºåŸŸç±»åž‹å­—å…¸',1,5,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(6,'è®¾å¤‡ç±»åž‹','DEVICE_TYPE','è®¾å¤‡ç±»åž‹å­—å…¸',1,6,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(7,'è®¾å¤‡çŠ¶æ€','DEVICE_STATUS','è®¾å¤‡çŠ¶æ€å­—å…¸',1,7,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(8,'æ¶ˆè´¹è®°å½•çŠ¶æ€','CONSUME_STATUS','æ¶ˆè´¹è®°å½•çŠ¶æ€å­—å…¸',1,8,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(9,'è´¦æˆ·çŠ¶æ€','ACCOUNT_STATUS','è´¦æˆ·çŠ¶æ€å­—å…¸',1,9,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(10,'é—¨ç¦æƒé™ç±»åž‹','ACCESS_PERMISSION_TYPE','é—¨ç¦æƒé™ç±»åž‹å­—å…¸',1,10,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(11,'é—¨ç¦é€šè¡Œç»“æžœ','ACCESS_RESULT','é—¨ç¦é€šè¡Œç»“æžœå­—å…¸',1,11,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(12,'è€ƒå‹¤æ‰“å¡ç±»åž‹','ATTENDANCE_CLOCK_TYPE','è€ƒå‹¤æ‰“å¡ç±»åž‹å­—å…¸',1,12,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(13,'è®¿å®¢é¢„çº¦çŠ¶æ€','VISITOR_STATUS','è®¿å®¢é¢„çº¦çŠ¶æ€å­—å…¸',1,13,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(14,'è§†é¢‘è®¾å¤‡ç±»åž‹','VIDEO_DEVICE_TYPE','è§†é¢‘è®¾å¤‡ç±»åž‹å­—å…¸',1,14,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0),(15,'ç³»ç»ŸçŠ¶æ€','SYSTEM_STATUS','ç³»ç»ŸçŠ¶æ€å­—å…¸',1,15,'2025-12-15 20:20:54','2025-12-15 20:20:54',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_dict_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_permission`
--

DROP TABLE IF EXISTS `t_common_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_permission` (
  `permission_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'æƒé™ID',
  `permission_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'æƒé™åç§°',
  `permission_code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'æƒé™ç¼–ç ',
  `resource_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'MENU' COMMENT 'èµ„æºç±»åž‹ï¼šMENU-èœå• BUTTON-æŒ‰é’® API-æŽ¥å£',
  `resource_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'èµ„æºURL',
  `resource_method` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'HTTPæ–¹æ³•',
  `parent_id` bigint DEFAULT '0' COMMENT 'çˆ¶æƒé™ID',
  `sort_order` int DEFAULT '0' COMMENT 'æŽ’åº',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'å›¾æ ‡',
  `component` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç»„ä»¶è·¯å¾„',
  `is_external` tinyint DEFAULT '0' COMMENT 'æ˜¯å¦å¤–éƒ¨é“¾æŽ¥ï¼š0-å¦ 1-æ˜¯',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`permission_id`),
  UNIQUE KEY `uk_permission_code` (`permission_code`,`deleted_flag`),
  KEY `idx_permission_name` (`permission_name`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_resource_type` (`resource_type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='æƒé™è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_permission`
--

LOCK TABLES `t_common_permission` WRITE;
/*!40000 ALTER TABLE `t_common_permission` DISABLE KEYS */;
INSERT INTO `t_common_permission` VALUES (1,'ç³»ç»Ÿç®¡ç†','SYSTEM_MANAGE','MENU',NULL,NULL,0,1,1,'system','/system',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'ç”¨æˆ·ç®¡ç†','SYSTEM_USER_MANAGE','MENU',NULL,NULL,1,1,1,'user','/system/user',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'è§’è‰²ç®¡ç†','SYSTEM_ROLE_MANAGE','MENU',NULL,NULL,1,2,1,'role','/system/role',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'æƒé™ç®¡ç†','SYSTEM_PERMISSION_MANAGE','MENU',NULL,NULL,1,3,1,'permission','/system/permission',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'èœå•ç®¡ç†','SYSTEM_MENU_MANAGE','MENU',NULL,NULL,1,4,1,'menu','/system/menu',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(6,'å­—å…¸ç®¡ç†','SYSTEM_DICT_MANAGE','MENU',NULL,NULL,1,5,1,'dict','/system/dict',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(7,'ç³»ç»Ÿé…ç½®','SYSTEM_CONFIG_MANAGE','MENU',NULL,NULL,1,6,1,'setting','/system/config',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(8,'æ“ä½œæ—¥å¿—','SYSTEM_LOG_MANAGE','MENU',NULL,NULL,1,7,1,'log','/system/log',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(9,'ç”¨æˆ·æŸ¥è¯¢','SYSTEM_USER_QUERY','BUTTON',NULL,NULL,2,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(10,'ç”¨æˆ·æ–°å¢ž','SYSTEM_USER_ADD','BUTTON',NULL,NULL,2,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(11,'ç”¨æˆ·ç¼–è¾‘','SYSTEM_USER_EDIT','BUTTON',NULL,NULL,2,3,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(12,'ç”¨æˆ·åˆ é™¤','SYSTEM_USER_DELETE','BUTTON',NULL,NULL,2,4,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(13,'ç”¨æˆ·å¯¼å‡º','SYSTEM_USER_EXPORT','BUTTON',NULL,NULL,2,5,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(14,'è§’è‰²æŸ¥è¯¢','SYSTEM_ROLE_QUERY','BUTTON',NULL,NULL,3,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(15,'è§’è‰²æ–°å¢ž','SYSTEM_ROLE_ADD','BUTTON',NULL,NULL,3,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(16,'è§’è‰²ç¼–è¾‘','SYSTEM_ROLE_EDIT','BUTTON',NULL,NULL,3,3,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(17,'è§’è‰²åˆ é™¤','SYSTEM_ROLE_DELETE','BUTTON',NULL,NULL,3,4,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(18,'æ¶ˆè´¹ç®¡ç†','CONSUME_MANAGE','MENU',NULL,NULL,0,2,1,'wallet','/consume',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(19,'æ¶ˆè´¹è®°å½•','CONSUME_RECORD_MANAGE','MENU',NULL,NULL,14,1,1,'record','/consume/record',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(20,'è´¦æˆ·ç®¡ç†','CONSUME_ACCOUNT_MANAGE','MENU',NULL,NULL,14,2,1,'account','/consume/account',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(21,'é€€æ¬¾ç®¡ç†','CONSUME_REFUND_MANAGE','MENU',NULL,NULL,14,3,1,'refund','/consume/refund',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(22,'ç»Ÿè®¡åˆ†æž','CONSUME_STATISTICS','MENU',NULL,NULL,14,4,1,'chart','/consume/statistics',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(23,'æ¶ˆè´¹è®°å½•æŸ¥è¯¢','CONSUME_RECORD_QUERY','BUTTON',NULL,NULL,15,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(24,'æ¶ˆè´¹è®°å½•å¯¼å‡º','CONSUME_RECORD_EXPORT','BUTTON',NULL,NULL,15,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(25,'è´¦æˆ·æŸ¥è¯¢','CONSUME_ACCOUNT_QUERY','BUTTON',NULL,NULL,16,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(26,'è´¦æˆ·å†»ç»“','CONSUME_ACCOUNT_FREEZE','BUTTON',NULL,NULL,16,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(27,'è´¦æˆ·å……å€¼','CONSUME_ACCOUNT_RECHARGE','BUTTON',NULL,NULL,16,3,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(28,'é—¨ç¦ç®¡ç†','ACCESS_MANAGE','MENU',NULL,NULL,0,3,1,'lock','/access',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(29,'è®¾å¤‡ç®¡ç†','ACCESS_DEVICE_MANAGE','MENU',NULL,NULL,22,1,1,'device','/access/device',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(30,'æƒé™ç®¡ç†','ACCESS_PERMISSION_MANAGE','MENU',NULL,NULL,22,2,1,'permission','/access/permission',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(31,'é€šè¡Œè®°å½•','ACCESS_RECORD_MANAGE','MENU',NULL,NULL,22,3,1,'record','/access/record',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(32,'åŒºåŸŸç®¡ç†','ACCESS_AREA_MANAGE','MENU',NULL,NULL,22,4,1,'area','/access/area',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(33,'è®¾å¤‡æŸ¥è¯¢','ACCESS_DEVICE_QUERY','BUTTON',NULL,NULL,23,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(34,'è®¾å¤‡æ–°å¢ž','ACCESS_DEVICE_ADD','BUTTON',NULL,NULL,23,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(35,'è®¾å¤‡ç¼–è¾‘','ACCESS_DEVICE_EDIT','BUTTON',NULL,NULL,23,3,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(36,'é€šè¡Œè®°å½•æŸ¥è¯¢','ACCESS_RECORD_QUERY','BUTTON',NULL,NULL,25,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(37,'æƒé™æŸ¥è¯¢','ACCESS_PERMISSION_QUERY','BUTTON',NULL,NULL,24,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(38,'æƒé™åˆ†é…','ACCESS_PERMISSION_ASSIGN','BUTTON',NULL,NULL,24,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(39,'è€ƒå‹¤ç®¡ç†','ATTENDANCE_MANAGE','MENU',NULL,NULL,0,4,1,'calendar','/attendance',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(40,'è€ƒå‹¤è®°å½•','ATTENDANCE_RECORD_MANAGE','MENU',NULL,NULL,29,1,1,'record','/attendance/record',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(41,'ç­æ¬¡ç®¡ç†','ATTENDANCE_SHIFT_MANAGE','MENU',NULL,NULL,29,2,1,'shift','/attendance/shift',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(42,'è€ƒå‹¤ç»Ÿè®¡','ATTENDANCE_STATISTICS','MENU',NULL,NULL,29,3,1,'chart','/attendance/statistics',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(43,'è€ƒå‹¤è®°å½•æŸ¥è¯¢','ATTENDANCE_RECORD_QUERY','BUTTON',NULL,NULL,30,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(44,'è€ƒå‹¤è®°å½•å¯¼å‡º','ATTENDANCE_RECORD_EXPORT','BUTTON',NULL,NULL,30,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(45,'ç­æ¬¡æŸ¥è¯¢','ATTENDANCE_SHIFT_QUERY','BUTTON',NULL,NULL,31,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(46,'ç­æ¬¡æ–°å¢ž','ATTENDANCE_SHIFT_ADD','BUTTON',NULL,NULL,31,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(47,'è®¿å®¢ç®¡ç†','VISITOR_MANAGE','MENU',NULL,NULL,0,5,1,'user','/visitor',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(48,'è®¿å®¢é¢„çº¦','VISITOR_APPOINTMENT_MANAGE','MENU',NULL,NULL,34,1,1,'appointment','/visitor/appointment',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(49,'è®¿å®¢è®°å½•','VISITOR_RECORD_MANAGE','MENU',NULL,NULL,34,2,1,'record','/visitor/record',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(50,'è®¿å®¢é¢„çº¦æŸ¥è¯¢','VISITOR_APPOINTMENT_QUERY','BUTTON',NULL,NULL,35,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(51,'è®¿å®¢é¢„çº¦å®¡æ‰¹','VISITOR_APPOINTMENT_APPROVE','BUTTON',NULL,NULL,35,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(52,'è®¿å®¢è®°å½•æŸ¥è¯¢','VISITOR_RECORD_QUERY','BUTTON',NULL,NULL,36,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(53,'è®¿å®¢è®°å½•å¯¼å‡º','VISITOR_RECORD_EXPORT','BUTTON',NULL,NULL,36,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(54,'è§†é¢‘ç›‘æŽ§','VIDEO_MANAGE','MENU',NULL,NULL,0,6,1,'video','/video',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(55,'è®¾å¤‡ç®¡ç†','VIDEO_DEVICE_MANAGE','MENU',NULL,NULL,39,1,1,'device','/video/device',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(56,'å®žæ—¶ç›‘æŽ§','VIDEO_REALTIME','MENU',NULL,NULL,39,2,1,'live','/video/live',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(57,'å½•åƒå›žæ”¾','VIDEO_PLAYBACK','MENU',NULL,NULL,39,3,1,'playback','/video/playback',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(58,'è§†é¢‘è®¾å¤‡æŸ¥è¯¢','VIDEO_DEVICE_QUERY','BUTTON',NULL,NULL,40,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(59,'è§†é¢‘è®¾å¤‡æ–°å¢ž','VIDEO_DEVICE_ADD','BUTTON',NULL,NULL,40,2,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(60,'è§†é¢‘è®¾å¤‡ç¼–è¾‘','VIDEO_DEVICE_EDIT','BUTTON',NULL,NULL,40,3,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(61,'å®žæ—¶ç›‘æŽ§æŸ¥çœ‹','VIDEO_REALTIME_VIEW','BUTTON',NULL,NULL,41,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(62,'å½•åƒå›žæ”¾æŸ¥çœ‹','VIDEO_PLAYBACK_VIEW','BUTTON',NULL,NULL,42,1,1,'','',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_role`
--

DROP TABLE IF EXISTS `t_common_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_role` (
  `role_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è§’è‰²ID',
  `role_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è§’è‰²åç§°',
  `role_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è§’è‰²ç¼–ç ',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'è§’è‰²æè¿°',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `sort_order` int DEFAULT '0' COMMENT 'æŽ’åº',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`role_id`),
  UNIQUE KEY `uk_role_code` (`role_code`,`deleted_flag`),
  KEY `idx_role_name` (`role_name`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è§’è‰²è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_role`
--

LOCK TABLES `t_common_role` WRITE;
/*!40000 ALTER TABLE `t_common_role` DISABLE KEYS */;
INSERT INTO `t_common_role` VALUES (1,'è¶…çº§ç®¡ç†å‘˜','SUPER_ADMIN','ç³»ç»Ÿè¶…çº§ç®¡ç†å‘˜ï¼Œæ‹¥æœ‰æ‰€æœ‰æƒé™',1,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'ç³»ç»Ÿç®¡ç†å‘˜','SYSTEM_ADMIN','ç³»ç»Ÿç®¡ç†å‘˜ï¼Œæ‹¥æœ‰ç³»ç»Ÿç®¡ç†æƒé™',1,2,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'æ™®é€šç®¡ç†å‘˜','ADMIN','æ™®é€šç®¡ç†å‘˜ï¼Œæ‹¥æœ‰åŸºç¡€ç®¡ç†æƒé™',1,3,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'æ™®é€šç”¨æˆ·','USER','æ™®é€šç”¨æˆ·ï¼Œæ‹¥æœ‰åŸºç¡€ä½¿ç”¨æƒé™',1,4,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'è®¿å®¢','VISITOR','è®¿å®¢ç”¨æˆ·ï¼Œæ‹¥æœ‰åªè¯»æƒé™',1,5,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_role_permission`
--

DROP TABLE IF EXISTS `t_common_role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_role_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ä¸»é”®ID',
  `role_id` bigint NOT NULL COMMENT 'è§’è‰²ID',
  `permission_id` bigint NOT NULL COMMENT 'æƒé™ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=221 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è§’è‰²æƒé™å…³è”è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_role_permission`
--

LOCK TABLES `t_common_role_permission` WRITE;
/*!40000 ALTER TABLE `t_common_role_permission` DISABLE KEYS */;
INSERT INTO `t_common_role_permission` VALUES (111,1,32,'2025-12-15 20:44:16',NULL),(112,1,34,'2025-12-15 20:44:16',NULL),(113,1,35,'2025-12-15 20:44:16',NULL),(114,1,29,'2025-12-15 20:44:16',NULL),(115,1,33,'2025-12-15 20:44:16',NULL),(116,1,28,'2025-12-15 20:44:16',NULL),(117,1,38,'2025-12-15 20:44:16',NULL),(118,1,30,'2025-12-15 20:44:16',NULL),(119,1,37,'2025-12-15 20:44:16',NULL),(120,1,31,'2025-12-15 20:44:16',NULL),(121,1,36,'2025-12-15 20:44:16',NULL),(122,1,39,'2025-12-15 20:44:16',NULL),(123,1,44,'2025-12-15 20:44:16',NULL),(124,1,40,'2025-12-15 20:44:16',NULL),(125,1,43,'2025-12-15 20:44:16',NULL),(126,1,46,'2025-12-15 20:44:16',NULL),(127,1,41,'2025-12-15 20:44:16',NULL),(128,1,45,'2025-12-15 20:44:16',NULL),(129,1,42,'2025-12-15 20:44:16',NULL),(130,1,26,'2025-12-15 20:44:16',NULL),(131,1,20,'2025-12-15 20:44:16',NULL),(132,1,25,'2025-12-15 20:44:16',NULL),(133,1,27,'2025-12-15 20:44:16',NULL),(134,1,18,'2025-12-15 20:44:16',NULL),(135,1,24,'2025-12-15 20:44:16',NULL),(136,1,19,'2025-12-15 20:44:16',NULL),(137,1,23,'2025-12-15 20:44:16',NULL),(138,1,21,'2025-12-15 20:44:16',NULL),(139,1,22,'2025-12-15 20:44:16',NULL),(140,1,7,'2025-12-15 20:44:16',NULL),(141,1,6,'2025-12-15 20:44:16',NULL),(142,1,8,'2025-12-15 20:44:16',NULL),(143,1,1,'2025-12-15 20:44:16',NULL),(144,1,5,'2025-12-15 20:44:16',NULL),(145,1,4,'2025-12-15 20:44:16',NULL),(146,1,15,'2025-12-15 20:44:16',NULL),(147,1,17,'2025-12-15 20:44:16',NULL),(148,1,16,'2025-12-15 20:44:16',NULL),(149,1,3,'2025-12-15 20:44:16',NULL),(150,1,14,'2025-12-15 20:44:16',NULL),(151,1,10,'2025-12-15 20:44:16',NULL),(152,1,12,'2025-12-15 20:44:16',NULL),(153,1,11,'2025-12-15 20:44:16',NULL),(154,1,13,'2025-12-15 20:44:16',NULL),(155,1,2,'2025-12-15 20:44:16',NULL),(156,1,9,'2025-12-15 20:44:16',NULL),(157,1,59,'2025-12-15 20:44:16',NULL),(158,1,60,'2025-12-15 20:44:16',NULL),(159,1,55,'2025-12-15 20:44:16',NULL),(160,1,58,'2025-12-15 20:44:16',NULL),(161,1,54,'2025-12-15 20:44:16',NULL),(162,1,57,'2025-12-15 20:44:16',NULL),(163,1,62,'2025-12-15 20:44:16',NULL),(164,1,56,'2025-12-15 20:44:16',NULL),(165,1,61,'2025-12-15 20:44:16',NULL),(166,1,51,'2025-12-15 20:44:16',NULL),(167,1,48,'2025-12-15 20:44:16',NULL),(168,1,50,'2025-12-15 20:44:16',NULL),(169,1,47,'2025-12-15 20:44:16',NULL),(170,1,53,'2025-12-15 20:44:16',NULL),(171,1,49,'2025-12-15 20:44:16',NULL),(172,1,52,'2025-12-15 20:44:16',NULL),(174,2,7,'2025-12-15 20:44:16',NULL),(175,2,6,'2025-12-15 20:44:16',NULL),(176,2,8,'2025-12-15 20:44:16',NULL),(177,2,1,'2025-12-15 20:44:16',NULL),(178,2,5,'2025-12-15 20:44:16',NULL),(179,2,4,'2025-12-15 20:44:16',NULL),(180,2,15,'2025-12-15 20:44:16',NULL),(181,2,17,'2025-12-15 20:44:16',NULL),(182,2,16,'2025-12-15 20:44:16',NULL),(183,2,3,'2025-12-15 20:44:16',NULL),(184,2,14,'2025-12-15 20:44:16',NULL),(185,2,10,'2025-12-15 20:44:16',NULL),(186,2,12,'2025-12-15 20:44:16',NULL),(187,2,11,'2025-12-15 20:44:16',NULL),(188,2,13,'2025-12-15 20:44:16',NULL),(189,2,2,'2025-12-15 20:44:16',NULL),(190,2,9,'2025-12-15 20:44:16',NULL),(205,3,33,'2025-12-15 20:44:16',NULL),(206,3,43,'2025-12-15 20:44:16',NULL),(207,3,23,'2025-12-15 20:44:16',NULL),(208,3,14,'2025-12-15 20:44:16',NULL),(209,3,9,'2025-12-15 20:44:16',NULL),(210,3,58,'2025-12-15 20:44:16',NULL),(211,3,50,'2025-12-15 20:44:16',NULL),(212,4,36,'2025-12-15 20:44:16',NULL),(213,4,43,'2025-12-15 20:44:16',NULL),(214,4,23,'2025-12-15 20:44:16',NULL),(215,4,52,'2025-12-15 20:44:16',NULL),(219,5,43,'2025-12-15 20:44:16',NULL),(220,5,23,'2025-12-15 20:44:16',NULL);
/*!40000 ALTER TABLE `t_common_role_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_user`
--

DROP TABLE IF EXISTS `t_common_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_user` (
  `user_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ç”¨æˆ·ID',
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ç”¨æˆ·å',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å¯†ç ï¼ˆåŠ å¯†ï¼‰',
  `real_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'çœŸå®žå§“å',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æ˜µç§°',
  `gender` tinyint DEFAULT '1' COMMENT 'æ€§åˆ«ï¼š1-ç”· 2-å¥³',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æ‰‹æœºå·',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'é‚®ç®±',
  `avatar` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'å¤´åƒURL',
  `birthday` date DEFAULT NULL COMMENT 'ç”Ÿæ—¥',
  `address` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åœ°å€',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-ç¦ç”¨',
  `department_id` bigint DEFAULT NULL COMMENT 'éƒ¨é—¨ID',
  `position` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'èŒä½',
  `employee_no` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'å‘˜å·¥ç¼–å·',
  `last_login_time` datetime DEFAULT NULL COMMENT 'æœ€åŽç™»å½•æ—¶é—´',
  `last_login_ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æœ€åŽç™»å½•IP',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uk_user_username` (`username`),
  UNIQUE KEY `uk_username` (`username`,`deleted_flag`),
  UNIQUE KEY `uk_user_phone` (`phone`),
  UNIQUE KEY `uk_user_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_email` (`email`),
  KEY `idx_department_id` (`department_id`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_user_department` (`department_id`),
  KEY `idx_user_status` (`status`),
  KEY `idx_user_create_time` (`create_time`),
  KEY `idx_user_dept_status` (`department_id`,`status`,`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ç”¨æˆ·è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_user`
--

LOCK TABLES `t_common_user` WRITE;
/*!40000 ALTER TABLE `t_common_user` DISABLE KEYS */;
INSERT INTO `t_common_user` VALUES (1,'admin','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','è¶…çº§ç®¡ç†å‘˜','Admin',1,'13800138000','admin@ioe-dream.com',NULL,NULL,NULL,1,2,'ç³»ç»Ÿç®¡ç†å‘˜','EMP001',NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'test_user1','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','æµ‹è¯•ç”¨æˆ·1','Test1',1,'13800138001','test1@ioe-dream.com',NULL,NULL,NULL,1,2,'æµ‹è¯•å·¥ç¨‹å¸ˆ','TEST001',NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'test_user2','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','æµ‹è¯•ç”¨æˆ·2','Test2',2,'13800138002','test2@ioe-dream.com',NULL,NULL,NULL,1,2,'æµ‹è¯•å·¥ç¨‹å¸ˆ','TEST002',NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'dev_user','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','å¼€å‘ç”¨æˆ·','Dev',1,'13800138003','dev@ioe-dream.com',NULL,NULL,NULL,1,2,'å¼€å‘å·¥ç¨‹å¸ˆ','DEV001',NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_common_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_common_user_role`
--

DROP TABLE IF EXISTS `t_common_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_common_user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ä¸»é”®ID',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·ID',
  `role_id` bigint NOT NULL COMMENT 'è§’è‰²ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ç”¨æˆ·è§’è‰²å…³è”è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_common_user_role`
--

LOCK TABLES `t_common_user_role` WRITE;
/*!40000 ALTER TABLE `t_common_user_role` DISABLE KEYS */;
INSERT INTO `t_common_user_role` VALUES (2,1,1,'2025-12-15 20:44:16',NULL);
/*!40000 ALTER TABLE `t_common_user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_consume_account`
--

DROP TABLE IF EXISTS `t_consume_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_consume_account` (
  `account_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è´¦æˆ·IDï¼ˆä¸»é”®ï¼‰',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·IDï¼ˆå…³è”ç”¨æˆ·è¡¨ï¼‰',
  `account_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è´¦æˆ·ç¼–å·ï¼ˆå”¯ä¸€ï¼‰',
  `account_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è´¦æˆ·åç§°',
  `balance` decimal(15,2) DEFAULT '0.00' COMMENT 'è´¦æˆ·ä½™é¢ï¼ˆå…ƒï¼‰',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-æ­£å¸¸ 2-å†»ç»“ 3-æ³¨é”€',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `uk_account_no` (`account_no`,`deleted_flag`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_balance` (`balance`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='æ¶ˆè´¹è´¦æˆ·è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_consume_account`
--

LOCK TABLES `t_consume_account` WRITE;
/*!40000 ALTER TABLE `t_consume_account` DISABLE KEYS */;
INSERT INTO `t_consume_account` VALUES (1,1,'ACC000001','è¶…çº§ç®¡ç†å‘˜è´¦æˆ·',1000.00,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,2,'ACC000002','æµ‹è¯•ç”¨æˆ·1è´¦æˆ·',500.00,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,3,'ACC000003','æµ‹è¯•ç”¨æˆ·2è´¦æˆ·',300.00,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,4,'ACC000004','å¼€å‘ç”¨æˆ·è´¦æˆ·',1000.00,1,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_consume_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_consume_record`
--

DROP TABLE IF EXISTS `t_consume_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_consume_record` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®°å½•ID',
  `transaction_no` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'äº¤æ˜“æµæ°´å·',
  `order_no` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è®¢å•å·',
  `user_id` bigint NOT NULL COMMENT 'ç”¨æˆ·ID',
  `user_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç”¨æˆ·å§“åï¼ˆå†—ä½™å­—æ®µï¼‰',
  `user_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç”¨æˆ·æ‰‹æœºå·ï¼ˆå†—ä½™å­—æ®µï¼‰',
  `user_type` tinyint DEFAULT '1' COMMENT 'ç”¨æˆ·ç±»åž‹ï¼š1-å‘˜å·¥ 2-è®¿å®¢ 3-ä¸´æ—¶ç”¨æˆ·',
  `account_id` bigint NOT NULL COMMENT 'è´¦æˆ·ID',
  `amount` decimal(15,2) NOT NULL COMMENT 'æ¶ˆè´¹é‡‘é¢',
  `consume_date` date NOT NULL COMMENT 'æ¶ˆè´¹æ—¥æœŸ',
  `consume_time` datetime NOT NULL COMMENT 'æ¶ˆè´¹æ—¶é—´',
  `merchant_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'å•†æˆ·åç§°',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'SUCCESS' COMMENT 'çŠ¶æ€ï¼šSUCCESS-æˆåŠŸ FAILED-å¤±è´¥',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`,`deleted_flag`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_consume_date` (`consume_date`),
  KEY `idx_consume_time` (`consume_time`),
  KEY `idx_status` (`status`),
  KEY `idx_consume_user` (`user_id`),
  KEY `idx_consume_user_date` (`user_id`,`consume_date`),
  KEY `idx_consume_account` (`account_id`),
  KEY `idx_consume_account_date` (`account_id`,`consume_date`),
  KEY `idx_consume_create_time` (`create_time`),
  KEY `idx_consume_status` (`status`),
  KEY `idx_consume_status_date` (`status`,`consume_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='æ¶ˆè´¹è®°å½•è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_consume_record`
--

LOCK TABLES `t_consume_record` WRITE;
/*!40000 ALTER TABLE `t_consume_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_consume_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_device_health_metric`
--

DROP TABLE IF EXISTS `t_device_health_metric`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_device_health_metric` (
  `metric_id` bigint NOT NULL AUTO_INCREMENT COMMENT '指标ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `metric_type` varchar(50) NOT NULL COMMENT '指标类型(cpu/memory/temperature/delay/packet_loss)',
  `metric_value` decimal(10,2) NOT NULL COMMENT '指标值',
  `metric_unit` varchar(20) DEFAULT NULL COMMENT '指标单位(%,℃,ms,%)',
  `collect_time` datetime NOT NULL COMMENT '采集时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`metric_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_device_type_time` (`device_id`,`metric_type`,`collect_time`),
  KEY `idx_collect_time` (`collect_time`),
  KEY `idx_metric_type` (`metric_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备健康指标表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_device_health_metric`
--

LOCK TABLES `t_device_health_metric` WRITE;
/*!40000 ALTER TABLE `t_device_health_metric` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_device_health_metric` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_device_quality_record`
--

DROP TABLE IF EXISTS `t_device_quality_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_device_quality_record` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `device_name` varchar(200) DEFAULT NULL COMMENT '设备名称',
  `device_type` int DEFAULT NULL COMMENT '设备类型(1-门禁 2-考勤 3-消费 4-视频 5-访客)',
  `health_score` int DEFAULT NULL COMMENT '健康评分(0-100)',
  `quality_level` varchar(20) DEFAULT NULL COMMENT '质量等级(优秀/良好/合格/较差/危险)',
  `diagnosis_result` text COMMENT '诊断结果(JSON格式)',
  `alarm_level` int DEFAULT NULL COMMENT '告警级别(0-无 1-低 2-中 3-高 4-紧急)',
  `diagnosis_time` datetime NOT NULL COMMENT '诊断时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`record_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_diagnosis_time` (`diagnosis_time`),
  KEY `idx_health_score` (`health_score`),
  KEY `idx_device_type` (`device_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备质量诊断记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_device_quality_record`
--

LOCK TABLES `t_device_quality_record` WRITE;
/*!40000 ALTER TABLE `t_device_quality_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_device_quality_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_entity_field_coverage`
--

DROP TABLE IF EXISTS `t_entity_field_coverage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_entity_field_coverage` (
  `coverage_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è¦†ç›–ID',
  `entity_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å®žä½“åç§°',
  `table_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è¡¨å',
  `total_entity_fields` int NOT NULL COMMENT 'å®žä½“å­—æ®µæ€»æ•°',
  `total_table_fields` int NOT NULL COMMENT 'è¡¨å­—æ®µæ€»æ•°',
  `covered_fields` int NOT NULL COMMENT 'å·²è¦†ç›–å­—æ®µæ•°',
  `coverage_rate` decimal(5,2) NOT NULL COMMENT 'è¦†ç›–çŽ‡ï¼ˆ%ï¼‰',
  `missing_fields` text COLLATE utf8mb4_unicode_ci COMMENT 'ç¼ºå¤±å­—æ®µåˆ—è¡¨',
  `extra_fields` text COLLATE utf8mb4_unicode_ci COMMENT 'å¤šä½™å­—æ®µåˆ—è¡¨',
  `validation_date` date NOT NULL COMMENT 'éªŒè¯æ—¥æœŸ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`coverage_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='å®žä½“å­—æ®µè¦†ç›–éªŒè¯è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_entity_field_coverage`
--

LOCK TABLES `t_entity_field_coverage` WRITE;
/*!40000 ALTER TABLE `t_entity_field_coverage` DISABLE KEYS */;
INSERT INTO `t_entity_field_coverage` VALUES (1,'ConsumeRecordEntity','t_consume_record',45,45,45,100.00,NULL,NULL,'2025-12-15','2025-12-15 20:20:56'),(2,'AccountEntity','t_consume_account',38,38,38,100.00,NULL,NULL,'2025-12-15','2025-12-15 20:20:56'),(3,'ConsumeRecordEntity','t_consume_record',45,45,45,100.00,NULL,NULL,'2025-12-15','2025-12-15 20:44:16'),(4,'AccountEntity','t_consume_account',38,38,38,100.00,NULL,NULL,'2025-12-15','2025-12-15 20:44:16');
/*!40000 ALTER TABLE `t_entity_field_coverage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_category`
--

DROP TABLE IF EXISTS `t_meal_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_category` (
  `category_id` bigint NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `category_name` varchar(50) NOT NULL COMMENT '分类名称',
  `category_code` varchar(20) NOT NULL COMMENT '分类编码',
  `sort_order` int DEFAULT '0' COMMENT '排序号',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记（0-未删除 1-已删除）',
  PRIMARY KEY (`category_id`),
  UNIQUE KEY `uk_category_code` (`category_code`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜品分类表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_category`
--

LOCK TABLES `t_meal_category` WRITE;
/*!40000 ALTER TABLE `t_meal_category` DISABLE KEYS */;
INSERT INTO `t_meal_category` VALUES (1,'主食类','STAPLE',1,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0),(2,'肉类','MEAT',2,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0),(3,'蔬菜类','VEGETABLE',3,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0),(4,'汤品类','SOUP',4,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0),(5,'饮品类','DRINK',5,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0),(6,'水果类','FRUIT',6,1,NULL,'2025-12-27 13:24:53','2025-12-27 13:24:53',NULL,NULL,0);
/*!40000 ALTER TABLE `t_meal_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_inventory`
--

DROP TABLE IF EXISTS `t_meal_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_inventory` (
  `inventory_id` bigint NOT NULL AUTO_INCREMENT COMMENT '库存ID',
  `menu_id` bigint NOT NULL COMMENT '菜品ID',
  `inventory_date` date NOT NULL COMMENT '库存日期',
  `meal_type` tinyint NOT NULL COMMENT '餐别（1-早餐 2-午餐 3-晚餐）',
  `initial_quantity` int DEFAULT '0' COMMENT '初始数量',
  `sold_quantity` int DEFAULT '0' COMMENT '已售数量',
  `remaining_quantity` int DEFAULT '0' COMMENT '剩余数量',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-有效 0-无效）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`inventory_id`),
  UNIQUE KEY `uk_menu_date_type` (`menu_id`,`inventory_date`,`meal_type`),
  KEY `idx_inventory_date` (`inventory_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜品库存表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_inventory`
--

LOCK TABLES `t_meal_inventory` WRITE;
/*!40000 ALTER TABLE `t_meal_inventory` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_meal_inventory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_menu`
--

DROP TABLE IF EXISTS `t_meal_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_menu` (
  `menu_id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜品ID',
  `category_id` bigint NOT NULL COMMENT '分类ID',
  `menu_name` varchar(100) NOT NULL COMMENT '菜品名称',
  `menu_code` varchar(50) NOT NULL COMMENT '菜品编码',
  `menu_image` varchar(500) DEFAULT NULL COMMENT '菜品图片URL',
  `price` decimal(10,2) NOT NULL COMMENT '价格（元）',
  `original_price` decimal(10,2) DEFAULT NULL COMMENT '原价（元）',
  `unit` varchar(20) DEFAULT NULL COMMENT '单位（份/个/碗）',
  `description` varchar(500) DEFAULT NULL COMMENT '菜品描述',
  `ingredients` text COMMENT '食材清单（JSON格式）',
  `nutrition_info` text COMMENT '营养信息（JSON格式）',
  `spicy_level` tinyint DEFAULT '0' COMMENT '辣度（0-不辣 1-微辣 2-中辣 3-重辣）',
  `available_days` varchar(50) DEFAULT NULL COMMENT '供应日期（1,2,3,4,5代表周一到周五）',
  `available_start_time` time DEFAULT NULL COMMENT '供应开始时间',
  `available_end_time` time DEFAULT NULL COMMENT '供应结束时间',
  `max_daily_quantity` int DEFAULT '999' COMMENT '每日最大供应数量',
  `current_quantity` int DEFAULT '0' COMMENT '当前剩余数量',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-上架 0-下架）',
  `sort_order` int DEFAULT '0' COMMENT '排序号',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记（0-未删除 1-已删除）',
  PRIMARY KEY (`menu_id`),
  UNIQUE KEY `uk_menu_code` (`menu_code`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_available_days` (`available_days`),
  KEY `idx_menu_date_type` (`available_days`,`available_start_time`,`available_end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜品表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_menu`
--

LOCK TABLES `t_meal_menu` WRITE;
/*!40000 ALTER TABLE `t_meal_menu` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_meal_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_order`
--

DROP TABLE IF EXISTS `t_meal_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_order` (
  `order_id` bigint NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(50) NOT NULL COMMENT '订单号',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `user_name` varchar(50) DEFAULT NULL COMMENT '用户姓名',
  `user_phone` varchar(20) DEFAULT NULL COMMENT '用户手机号',
  `order_date` date NOT NULL COMMENT '订餐日期',
  `meal_type` tinyint NOT NULL COMMENT '餐别（1-早餐 2-午餐 3-晚餐）',
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '订单总额（元）',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '优惠金额（元）',
  `actual_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '实付金额（元）',
  `subsidy_amount` decimal(10,2) DEFAULT '0.00' COMMENT '补贴金额（元）',
  `order_status` tinyint DEFAULT '1' COMMENT '订单状态（1-待支付 2-已支付 3-已完成 4-已取消 5-已退款）',
  `payment_status` tinyint DEFAULT '0' COMMENT '支付状态（0-未支付 1-已支付 2-支付失败）',
  `payment_time` datetime DEFAULT NULL COMMENT '支付时间',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式（balance-余额 wechat-微信 alipay-支付宝）',
  `pickup_time` time DEFAULT NULL COMMENT '取餐时间',
  `pickup_location` varchar(100) DEFAULT NULL COMMENT '取餐地点',
  `special_requirements` varchar(500) DEFAULT NULL COMMENT '特殊要求',
  `cancel_reason` varchar(500) DEFAULT NULL COMMENT '取消原因',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `refund_amount` decimal(10,2) DEFAULT NULL COMMENT '退款金额',
  `refund_time` datetime DEFAULT NULL COMMENT '退款时间',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记（0-未删除 1-已删除）',
  PRIMARY KEY (`order_id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_date` (`order_date`),
  KEY `idx_order_status` (`order_status`),
  KEY `idx_meal_type` (`meal_type`),
  KEY `idx_order_user_date` (`user_id`,`order_date`),
  KEY `idx_order_status_date` (`order_status`,`order_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订餐订单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_order`
--

LOCK TABLES `t_meal_order` WRITE;
/*!40000 ALTER TABLE `t_meal_order` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_meal_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_order_config`
--

DROP TABLE IF EXISTS `t_meal_order_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_order_config` (
  `config_id` bigint NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `config_value` varchar(500) DEFAULT NULL COMMENT '配置值',
  `config_desc` varchar(200) DEFAULT NULL COMMENT '配置描述',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`config_id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订餐配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_order_config`
--

LOCK TABLES `t_meal_order_config` WRITE;
/*!40000 ALTER TABLE `t_meal_order_config` DISABLE KEYS */;
INSERT INTO `t_meal_order_config` VALUES (1,'order.advance.minutes','60','提前订餐分钟数',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(2,'order.cancel.minutes','30','订单取消截止时间（分钟）',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(3,'order.max.daily','10','每日最大订餐数量',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(4,'payment.enable.balance','1','是否启用余额支付',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(5,'payment.enable.wechat','0','是否启用微信支付',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(6,'payment.enable.alipay','0','是否启用支付宝支付',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(7,'subsidy.enable','1','是否启用补贴',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(8,'subsidy.breakfast.amount','5.00','早餐补贴金额',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(9,'subsidy.lunch.amount','15.00','午餐补贴金额',1,'2025-12-27 13:24:53','2025-12-27 13:24:53'),(10,'subsidy.dinner.amount','10.00','晚餐补贴金额',1,'2025-12-27 13:24:53','2025-12-27 13:24:53');
/*!40000 ALTER TABLE `t_meal_order_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_meal_order_item`
--

DROP TABLE IF EXISTS `t_meal_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_meal_order_item` (
  `item_id` bigint NOT NULL AUTO_INCREMENT COMMENT '明细ID',
  `order_id` bigint NOT NULL COMMENT '订单ID',
  `menu_id` bigint NOT NULL COMMENT '菜品ID',
  `menu_name` varchar(100) NOT NULL COMMENT '菜品名称',
  `menu_code` varchar(50) DEFAULT NULL COMMENT '菜品编码',
  `menu_image` varchar(500) DEFAULT NULL COMMENT '菜品图片URL',
  `unit_price` decimal(10,2) NOT NULL COMMENT '单价（元）',
  `quantity` int NOT NULL DEFAULT '1' COMMENT '数量',
  `subtotal` decimal(10,2) NOT NULL COMMENT '小计（元）',
  `special_requirements` varchar(500) DEFAULT NULL COMMENT '特殊要求',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`item_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订餐订单明细表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_meal_order_item`
--

LOCK TABLES `t_meal_order_item` WRITE;
/*!40000 ALTER TABLE `t_meal_order_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_meal_order_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_migration_history`
--

DROP TABLE IF EXISTS `t_migration_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_migration_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ä¸»é”®ID',
  `version` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ç‰ˆæœ¬å·',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'ç‰ˆæœ¬æè¿°',
  `script_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è„šæœ¬æ–‡ä»¶å',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'SUCCESS' COMMENT 'æ‰§è¡ŒçŠ¶æ€',
  `start_time` datetime DEFAULT NULL COMMENT 'å¼€å§‹æ—¶é—´',
  `end_time` datetime DEFAULT NULL COMMENT 'ç»“æŸæ—¶é—´',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_version` (`version`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='æ•°æ®åº“è¿ç§»åŽ†å²è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_migration_history`
--

LOCK TABLES `t_migration_history` WRITE;
/*!40000 ALTER TABLE `t_migration_history` DISABLE KEYS */;
INSERT INTO `t_migration_history` VALUES (1,'V1.0.0','æ•°æ®åº“åˆå§‹æž¶æž„ - åˆ›å»ºæ‰€æœ‰åŸºç¡€è¡¨ç»“æž„','01-ioedream-schema.sql','SUCCESS','2025-12-15 20:20:54','2025-12-15 20:20:54','2025-12-15 20:20:54'),(2,'V1.1.0','æ•°æ®åº“åˆå§‹æ•°æ® - åˆå§‹åŒ–ç”¨æˆ·ã€è§’è‰²ã€æƒé™ã€å­—å…¸ç­‰åŸºç¡€æ•°æ®','02-ioedream-data.sql','SUCCESS','2025-12-15 20:20:55','2025-12-15 20:44:16','2025-12-15 20:20:55'),(3,'V1.1.0-DEV','å¼€å‘çŽ¯å¢ƒåˆå§‹æ•°æ® - åŒ…å«æµ‹è¯•ç”¨æˆ·å’Œæµ‹è¯•æ•°æ®','02-ioedream-data-dev.sql','SUCCESS','2025-12-15 20:20:55','2025-12-15 20:44:16','2025-12-15 20:20:55');
/*!40000 ALTER TABLE `t_migration_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_quality_alarm`
--

DROP TABLE IF EXISTS `t_quality_alarm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_quality_alarm` (
  `alarm_id` bigint NOT NULL AUTO_INCREMENT COMMENT '告警ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `device_name` varchar(200) DEFAULT NULL COMMENT '设备名称',
  `rule_id` bigint DEFAULT NULL COMMENT '触发的规则ID',
  `alarm_level` int DEFAULT NULL COMMENT '告警级别(1-低 2-中 3-高 4-紧急)',
  `alarm_title` varchar(200) NOT NULL COMMENT '告警标题',
  `alarm_content` text COMMENT '告警内容',
  `alarm_status` tinyint DEFAULT '1' COMMENT '告警状态(1-待处理 2-处理中 3-已处理)',
  `handle_result` text COMMENT '处理结果',
  `handle_user_id` bigint DEFAULT NULL COMMENT '处理人ID',
  `handle_time` datetime DEFAULT NULL COMMENT '处理时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`alarm_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_alarm_level` (`alarm_level`),
  KEY `idx_alarm_status` (`alarm_status`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_handle_user` (`handle_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备质量告警表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_quality_alarm`
--

LOCK TABLES `t_quality_alarm` WRITE;
/*!40000 ALTER TABLE `t_quality_alarm` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_quality_alarm` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_quality_diagnosis_rule`
--

DROP TABLE IF EXISTS `t_quality_diagnosis_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_quality_diagnosis_rule` (
  `rule_id` bigint NOT NULL AUTO_INCREMENT COMMENT '规则ID',
  `rule_name` varchar(200) NOT NULL COMMENT '规则名称',
  `rule_code` varchar(100) NOT NULL COMMENT '规则编码',
  `device_type` int DEFAULT NULL COMMENT '设备类型(1-门禁 2-考勤 3-消费 4-视频 5-访客)',
  `metric_type` varchar(50) DEFAULT NULL COMMENT '指标类型',
  `rule_expression` varchar(500) DEFAULT NULL COMMENT '规则表达式',
  `threshold_value` decimal(10,2) DEFAULT NULL COMMENT '阈值',
  `alarm_level` int DEFAULT NULL COMMENT '告警级别(1-低 2-中 3-高 4-紧急)',
  `rule_status` tinyint DEFAULT '1' COMMENT '规则状态(1-启用 0-禁用)',
  `rule_description` varchar(500) DEFAULT NULL COMMENT '规则描述',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`rule_id`),
  UNIQUE KEY `uk_rule_code` (`rule_code`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_rule_status` (`rule_status`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='质量诊断规则表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_quality_diagnosis_rule`
--

LOCK TABLES `t_quality_diagnosis_rule` WRITE;
/*!40000 ALTER TABLE `t_quality_diagnosis_rule` DISABLE KEYS */;
INSERT INTO `t_quality_diagnosis_rule` VALUES (1,'门禁设备离线告警','ACCESS_OFFLINE',1,'online_status','eq',0.00,3,1,'门禁设备离线超过30分钟','2025-12-27 15:23:27'),(2,'门禁响应延迟告警','ACCESS_DELAY',1,'delay','gt',3000.00,2,1,'门禁响应延迟超过3秒','2025-12-27 15:23:27'),(3,'门禁温度告警','ACCESS_TEMP',1,'temperature','gt',70.00,3,1,'门禁设备温度超过70℃','2025-12-27 15:23:27'),(4,'考勤设备离线告警','ATTENDANCE_OFFLINE',2,'online_status','eq',0.00,3,1,'考勤设备离线超过30分钟','2025-12-27 15:23:27'),(5,'考勤识别成功率低','ATTENDANCE_SUCCESS_RATE',2,'success_rate','lt',90.00,2,1,'考勤识别成功率低于90%','2025-12-27 15:23:27'),(6,'考勤设备温度告警','ATTENDANCE_TEMP',2,'temperature','gt',70.00,3,1,'考勤设备温度超过70℃','2025-12-27 15:23:27'),(7,'消费设备离线告警','CONSUME_OFFLINE',3,'online_status','eq',0.00,3,1,'消费设备离线超过30分钟','2025-12-27 15:23:27'),(8,'消费交易响应延迟','CONSUME_DELAY',3,'delay','gt',5000.00,2,1,'消费交易响应延迟超过5秒','2025-12-27 15:23:27'),(9,'消费设备温度告警','CONSUME_TEMP',3,'temperature','gt',70.00,3,1,'消费设备温度超过70℃','2025-12-27 15:23:27'),(10,'视频设备离线告警','VIDEO_OFFLINE',4,'online_status','eq',0.00,3,1,'视频设备离线超过30分钟','2025-12-27 15:23:27'),(11,'视频码流异常','VIDEO_STREAM',4,'stream_status','eq',0.00,2,1,'视频码流异常','2025-12-27 15:23:27'),(12,'视频存储空间告警','VIDEO_STORAGE',4,'storage_usage','gt',90.00,3,1,'视频存储空间使用率超过90%','2025-12-27 15:23:27'),(13,'访客设备离线告警','VISITOR_OFFLINE',5,'online_status','eq',0.00,3,1,'访客设备离线超过30分钟','2025-12-27 15:23:27'),(14,'访客识别成功率低','VISITOR_SUCCESS_RATE',5,'success_rate','lt',90.00,2,1,'访客识别成功率低于90%','2025-12-27 15:23:27'),(15,'访客设备温度告警','VISITOR_TEMP',5,'temperature','gt',70.00,3,1,'访客设备温度超过70℃','2025-12-27 15:23:27');
/*!40000 ALTER TABLE `t_quality_diagnosis_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_category`
--

DROP TABLE IF EXISTS `t_report_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_category` (
  `category_id` bigint NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `category_name` varchar(100) NOT NULL COMMENT '分类名称',
  `category_code` varchar(50) NOT NULL COMMENT '分类编码',
  `parent_id` bigint DEFAULT '0' COMMENT '父分类ID',
  `sort_order` int DEFAULT '0' COMMENT '排序号',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记',
  PRIMARY KEY (`category_id`),
  UNIQUE KEY `uk_category_code` (`category_code`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表分类表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_category`
--

LOCK TABLES `t_report_category` WRITE;
/*!40000 ALTER TABLE `t_report_category` DISABLE KEYS */;
INSERT INTO `t_report_category` VALUES (1,'门禁报表','ACCESS',0,1,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0),(2,'考勤报表','ATTENDANCE',0,2,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0),(3,'消费报表','CONSUME',0,3,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0),(4,'访客报表','VISITOR',0,4,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0),(5,'视频报表','VIDEO',0,5,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0),(6,'综合报表','COMPREHENSIVE',0,6,1,'2025-12-27 15:23:26','2025-12-27 15:23:26',0);
/*!40000 ALTER TABLE `t_report_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_definition`
--

DROP TABLE IF EXISTS `t_report_definition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_definition` (
  `report_id` bigint NOT NULL AUTO_INCREMENT COMMENT '报表ID',
  `report_name` varchar(200) NOT NULL COMMENT '报表名称',
  `report_code` varchar(100) NOT NULL COMMENT '报表编码',
  `report_type` tinyint NOT NULL COMMENT '报表类型（1-列表 2-汇总 3-图表 4-交叉表）',
  `business_module` varchar(50) DEFAULT NULL COMMENT '业务模块（access/attendance/consume等）',
  `category_id` bigint DEFAULT NULL COMMENT '分类ID',
  `data_source_type` tinyint NOT NULL COMMENT '数据源类型（1-SQL 2-API 3-静态）',
  `data_source_config` text COMMENT '数据源配置（JSON）',
  `template_type` tinyint DEFAULT NULL COMMENT '模板类型（1-Excel 2-PDF 3-Word）',
  `template_config` text COMMENT '模板配置（JSON）',
  `export_formats` varchar(100) DEFAULT NULL COMMENT '导出格式（excel,pdf,word,csv）',
  `description` text COMMENT '报表描述',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `sort_order` int DEFAULT '0' COMMENT '排序号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记',
  PRIMARY KEY (`report_id`),
  UNIQUE KEY `uk_report_code` (`report_code`),
  KEY `idx_business_module` (`business_module`),
  KEY `idx_status` (`status`),
  KEY `idx_report_module_type` (`business_module`,`report_type`),
  KEY `idx_report_status_sort` (`status`,`sort_order`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表定义表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_definition`
--

LOCK TABLES `t_report_definition` WRITE;
/*!40000 ALTER TABLE `t_report_definition` DISABLE KEYS */;
INSERT INTO `t_report_definition` VALUES (1,'每日考勤汇总表','DAILY_ATTENDANCE_SUMMARY',2,'attendance',2,1,'{\"sql\":\"SELECT * FROM t_attendance_record WHERE record_date = #{date}\"}',1,NULL,'excel,pdf','统计每日考勤打卡情况',1,0,'2025-12-27 15:23:26','2025-12-27 15:23:26',NULL,NULL,0),(2,'月度消费统计表','MONTHLY_CONSUME_STATS',2,'consume',3,1,'{\"sql\":\"SELECT * FROM t_consume_record WHERE MONTH(consume_time) = #{month}\"}',1,NULL,'excel,pdf','统计月度消费数据',1,0,'2025-12-27 15:23:26','2025-12-27 15:23:26',NULL,NULL,0),(3,'门禁通行记录表','ACCESS_RECORD_LIST',1,'access',1,1,'{\"sql\":\"SELECT * FROM t_access_record WHERE access_time BETWEEN #{startTime} AND #{endTime}\"}',1,NULL,'excel,csv','查询门禁通行记录',1,0,'2025-12-27 15:23:26','2025-12-27 15:23:26',NULL,NULL,0);
/*!40000 ALTER TABLE `t_report_definition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_generation`
--

DROP TABLE IF EXISTS `t_report_generation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_generation` (
  `generation_id` bigint NOT NULL AUTO_INCREMENT COMMENT '生成记录ID',
  `report_id` bigint NOT NULL COMMENT '报表ID',
  `report_name` varchar(200) DEFAULT NULL COMMENT '报表名称',
  `parameters` text COMMENT '请求参数（JSON）',
  `generate_type` tinyint DEFAULT NULL COMMENT '生成方式（1-手动 2-定时 3-API）',
  `file_type` varchar(20) DEFAULT NULL COMMENT '文件类型（excel/pdf/word/csv）',
  `file_path` varchar(500) DEFAULT NULL COMMENT '文件路径',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小（字节）',
  `status` tinyint DEFAULT NULL COMMENT '状态（1-生成中 2-成功 3-失败）',
  `error_message` text COMMENT '错误信息',
  `generate_time` datetime DEFAULT NULL COMMENT '生成时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`generation_id`),
  KEY `idx_report_id` (`report_id`),
  KEY `idx_generate_time` (`generate_time`),
  KEY `idx_status` (`status`),
  KEY `idx_generation_report_time` (`report_id`,`generate_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表生成记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_generation`
--

LOCK TABLES `t_report_generation` WRITE;
/*!40000 ALTER TABLE `t_report_generation` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_report_generation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_parameter`
--

DROP TABLE IF EXISTS `t_report_parameter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_parameter` (
  `parameter_id` bigint NOT NULL AUTO_INCREMENT COMMENT '参数ID',
  `report_id` bigint NOT NULL COMMENT '报表ID',
  `parameter_name` varchar(100) NOT NULL COMMENT '参数名称',
  `parameter_code` varchar(50) NOT NULL COMMENT '参数编码',
  `parameter_type` varchar(50) NOT NULL COMMENT '参数类型（String/Integer/Date等）',
  `default_value` varchar(500) DEFAULT NULL COMMENT '默认值',
  `required` tinyint DEFAULT '0' COMMENT '是否必填（1-是 0-否）',
  `validation_rule` varchar(500) DEFAULT NULL COMMENT '验证规则（正则表达式）',
  `sort_order` int DEFAULT '0' COMMENT '排序号',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`parameter_id`),
  KEY `idx_report_id` (`report_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表参数表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_parameter`
--

LOCK TABLES `t_report_parameter` WRITE;
/*!40000 ALTER TABLE `t_report_parameter` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_report_parameter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_schedule`
--

DROP TABLE IF EXISTS `t_report_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_schedule` (
  `schedule_id` bigint NOT NULL AUTO_INCREMENT COMMENT '调度ID',
  `report_id` bigint NOT NULL COMMENT '报表ID',
  `schedule_name` varchar(200) NOT NULL COMMENT '调度名称',
  `cron_expression` varchar(100) NOT NULL COMMENT 'Cron表达式',
  `parameters` text COMMENT '调度参数（JSON）',
  `notification_config` text COMMENT '通知配置（邮件、消息等）',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `last_execute_time` datetime DEFAULT NULL COMMENT '最后执行时间',
  `next_execute_time` datetime DEFAULT NULL COMMENT '下次执行时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记',
  PRIMARY KEY (`schedule_id`),
  KEY `idx_report_id` (`report_id`),
  KEY `idx_status` (`status`),
  KEY `idx_next_execute_time` (`next_execute_time`),
  KEY `idx_schedule_next_time` (`status`,`next_execute_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表调度任务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_schedule`
--

LOCK TABLES `t_report_schedule` WRITE;
/*!40000 ALTER TABLE `t_report_schedule` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_report_schedule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_report_template`
--

DROP TABLE IF EXISTS `t_report_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_report_template` (
  `template_id` bigint NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `report_id` bigint NOT NULL COMMENT '报表ID',
  `template_name` varchar(200) NOT NULL COMMENT '模板名称',
  `template_type` tinyint NOT NULL COMMENT '模板类型（1-Excel 2-PDF 3-Word）',
  `file_path` varchar(500) NOT NULL COMMENT '模板文件路径',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小（字节）',
  `version` varchar(50) DEFAULT NULL COMMENT '版本号',
  `description` text COMMENT '模板描述',
  `status` tinyint DEFAULT '1' COMMENT '状态（1-启用 0-禁用）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标记',
  PRIMARY KEY (`template_id`),
  KEY `idx_report_id` (`report_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='报表模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_report_template`
--

LOCK TABLES `t_report_template` WRITE;
/*!40000 ALTER TABLE `t_report_template` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_report_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_response_format_validation`
--

DROP TABLE IF EXISTS `t_response_format_validation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_response_format_validation` (
  `format_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'æ ¼å¼ID',
  `client_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å®¢æˆ·ç«¯ç±»åž‹ï¼šSMART_ADMIN-ç®¡ç†ç«¯ MOBILE-ç§»åŠ¨ç«¯',
  `response_format` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'å“åº”æ ¼å¼ï¼šIOE_DREAM-IOEæ ¼å¼ SMART_ADMIN-æ™ºèƒ½æ ¼å¼',
  `field_mapping` text COLLATE utf8mb4_unicode_ci COMMENT 'å­—æ®µæ˜ å°„å…³ç³»',
  `format_compatible` tinyint DEFAULT '1' COMMENT 'æ ¼å¼å…¼å®¹ï¼š1-å…¼å®¹ 0-ä¸å…¼å®¹',
  `auto_conversion_support` tinyint DEFAULT '1' COMMENT 'è‡ªåŠ¨è½¬æ¢æ”¯æŒï¼š1-æ”¯æŒ 0-ä¸æ”¯æŒ',
  `validation_date` date NOT NULL COMMENT 'éªŒè¯æ—¥æœŸ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`format_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='å“åº”æ ¼å¼éªŒè¯è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_response_format_validation`
--

LOCK TABLES `t_response_format_validation` WRITE;
/*!40000 ALTER TABLE `t_response_format_validation` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_response_format_validation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_system_config`
--

DROP TABLE IF EXISTS `t_system_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_system_config` (
  `config_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'é…ç½®ID',
  `config_key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'é…ç½®é”®',
  `config_value` text COLLATE utf8mb4_unicode_ci COMMENT 'é…ç½®å€¼',
  `config_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'é…ç½®åç§°',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT 'æè¿°',
  `config_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'SYSTEM' COMMENT 'é…ç½®ç±»åž‹ï¼šSYSTEM-ç³»ç»Ÿ BUSINESS-ä¸šåŠ¡',
  `is_encrypted` tinyint DEFAULT '0' COMMENT 'æ˜¯å¦åŠ å¯†ï¼š0-å¦ 1-æ˜¯',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`config_id`),
  UNIQUE KEY `uk_config_key` (`config_key`,`deleted_flag`),
  KEY `idx_config_name` (`config_name`),
  KEY `idx_config_type` (`config_type`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ç³»ç»Ÿé…ç½®è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_system_config`
--

LOCK TABLES `t_system_config` WRITE;
/*!40000 ALTER TABLE `t_system_config` DISABLE KEYS */;
INSERT INTO `t_system_config` VALUES (1,'system.name','IOE-DREAMæ™ºæ…§å›­åŒºä¸€å¡é€šç®¡ç†å¹³å°','ç³»ç»Ÿåç§°','ç³»ç»Ÿæ˜¾ç¤ºåç§°','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'system.version','1.0.0','ç³»ç»Ÿç‰ˆæœ¬','å½“å‰ç³»ç»Ÿç‰ˆæœ¬å·','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'system.company','IOE-DREAMç§‘æŠ€æœ‰é™å…¬å¸','å…¬å¸åç§°','æ‰€å±žå…¬å¸','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(4,'system.logo','/static/logo.png','ç³»ç»ŸLogo','ç³»ç»ŸLogoå›¾ç‰‡è·¯å¾„','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(5,'system.favicon','/static/favicon.ico','ç½‘ç«™å›¾æ ‡','ç½‘ç«™å›¾æ ‡è·¯å¾„','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(6,'system.copyright','Â©2025 IOE-DREAM. All rights reserved.','ç‰ˆæƒä¿¡æ¯','ç³»ç»Ÿç‰ˆæƒä¿¡æ¯','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(7,'user.password.min.length','6','å¯†ç æœ€å°é•¿åº¦','ç”¨æˆ·å¯†ç æœ€å°é•¿åº¦è¦æ±‚','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(8,'user.password.max.length','20','å¯†ç æœ€å¤§é•¿åº¦','ç”¨æˆ·å¯†ç æœ€å¤§é•¿åº¦è¦æ±‚','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(9,'session.timeout','30','ä¼šè¯è¶…æ—¶æ—¶é—´','ç”¨æˆ·ä¼šè¯è¶…æ—¶æ—¶é—´ï¼ˆåˆ†é’Ÿï¼‰','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(10,'file.upload.max.size','10485760','æ–‡ä»¶ä¸Šä¼ å¤§å°é™åˆ¶','æ–‡ä»¶ä¸Šä¼ æœ€å¤§å¤§å°ï¼ˆå­—èŠ‚ï¼‰','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(11,'file.upload.allowed.types','jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx','å…è®¸ä¸Šä¼ çš„æ–‡ä»¶ç±»åž‹','å…è®¸ä¸Šä¼ çš„æ–‡ä»¶ç±»åž‹åˆ—è¡¨','SYSTEM',0,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_system_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_system_department`
--

DROP TABLE IF EXISTS `t_system_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_system_department` (
  `department_id` bigint NOT NULL AUTO_INCREMENT,
  `department_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` bigint DEFAULT NULL,
  `sort_order` int DEFAULT '0',
  `status` int DEFAULT '1',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `update_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deleted_flag` int DEFAULT '0',
  `create_user_id` bigint DEFAULT NULL,
  `update_user_id` bigint DEFAULT NULL,
  `deleted` int DEFAULT '0',
  `version` int DEFAULT '0',
  PRIMARY KEY (`department_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_system_department`
--

LOCK TABLES `t_system_department` WRITE;
/*!40000 ALTER TABLE `t_system_department` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_system_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_system_position`
--

DROP TABLE IF EXISTS `t_system_position`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_system_position` (
  `position_id` bigint NOT NULL AUTO_INCREMENT,
  `position_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `position_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort_order` int DEFAULT '0',
  `status` int DEFAULT '1',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `update_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deleted_flag` int DEFAULT '0',
  `create_user_id` bigint DEFAULT NULL,
  `update_user_id` bigint DEFAULT NULL,
  `deleted` int DEFAULT '0',
  `version` int DEFAULT '0',
  PRIMARY KEY (`position_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_system_position`
--

LOCK TABLES `t_system_position` WRITE;
/*!40000 ALTER TABLE `t_system_position` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_system_position` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_system_role`
--

DROP TABLE IF EXISTS `t_system_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_system_role` (
  `role_id` bigint NOT NULL AUTO_INCREMENT,
  `role_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort_order` int DEFAULT '0',
  `status` int DEFAULT '1',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `update_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `deleted_flag` int DEFAULT '0',
  `create_user_id` bigint DEFAULT NULL,
  `update_user_id` bigint DEFAULT NULL,
  `deleted` int DEFAULT '0',
  `version` int DEFAULT '0',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_system_role`
--

LOCK TABLES `t_system_role` WRITE;
/*!40000 ALTER TABLE `t_system_role` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_system_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_ai_model`
--

DROP TABLE IF EXISTS `t_video_ai_model`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_ai_model` (
  `model_id` bigint NOT NULL AUTO_INCREMENT COMMENT '模型ID',
  `model_name` varchar(200) NOT NULL COMMENT '模型名称',
  `model_version` varchar(50) NOT NULL COMMENT '模型版本号（语义化版本）',
  `model_type` varchar(50) NOT NULL COMMENT '模型类型（FACE_DETECTION/FALL_DETECTION/BEHAVIOR_DETECTION）',
  `file_path` varchar(500) DEFAULT NULL COMMENT 'MinIO文件路径',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小（字节）',
  `file_md5` varchar(32) DEFAULT NULL COMMENT '文件MD5值',
  `model_status` tinyint NOT NULL DEFAULT '0' COMMENT '模型状态（0-草稿 1-已发布 2-已废弃）',
  `supported_events` varchar(500) DEFAULT NULL COMMENT '支持的事件类型（JSON数组）',
  `model_metadata` text COMMENT '模型元数据（JSON格式）',
  `accuracy_rate` decimal(5,4) DEFAULT NULL COMMENT '模型准确率',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `published_by` bigint DEFAULT NULL COMMENT '发布人ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记（0-未删除 1-已删除）',
  PRIMARY KEY (`model_id`),
  UNIQUE KEY `uk_model_version` (`model_name`,`model_version`,`deleted_flag`),
  KEY `idx_model_type` (`model_type`),
  KEY `idx_model_status` (`model_status`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_model_type_status` (`model_type`,`model_status`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='AI模型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_ai_model`
--

LOCK TABLES `t_video_ai_model` WRITE;
/*!40000 ALTER TABLE `t_video_ai_model` DISABLE KEYS */;
INSERT INTO `t_video_ai_model` VALUES (1,'跌倒检测模型','1.0.0','FALL_DETECTION',NULL,NULL,NULL,1,'[\"FALL_DETECTION\"]',NULL,0.9200,NULL,NULL,'2025-12-29 20:08:12','2025-12-29 20:08:12',0),(2,'人脸检测模型','1.5.0','FACE_DETECTION',NULL,NULL,NULL,1,'[\"FACE_DETECTION\"]',NULL,0.9500,NULL,NULL,'2025-12-29 20:08:12','2025-12-29 20:08:12',0),(3,'徘徊检测模型','1.0.0','LOITERING_DETECTION',NULL,NULL,NULL,0,'[\"LOITERING_DETECTION\"]',NULL,0.8800,NULL,NULL,'2025-12-29 20:08:12','2025-12-29 20:08:12',0);
/*!40000 ALTER TABLE `t_video_ai_model` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_alarm_record`
--

DROP TABLE IF EXISTS `t_video_alarm_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_alarm_record` (
  `alarm_id` varchar(64) NOT NULL COMMENT '告警ID',
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `rule_name` varchar(128) NOT NULL COMMENT '规则名称',
  `event_id` varchar(64) NOT NULL COMMENT '事件ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `device_code` varchar(128) NOT NULL COMMENT '设备编码',
  `event_type` varchar(64) NOT NULL COMMENT '事件类型',
  `alarm_level` tinyint NOT NULL COMMENT '告警级别: 1-低, 2-中, 3-高, 4-紧急',
  `alarm_status` tinyint NOT NULL DEFAULT '0' COMMENT '告警状态: 0-待处理, 1-处理中, 2-已处理, 3-已忽略',
  `confidence` decimal(5,4) NOT NULL COMMENT '置信度',
  `bbox` varchar(256) DEFAULT NULL COMMENT '边界框(JSON格式)',
  `snapshot_url` varchar(512) DEFAULT NULL COMMENT '抓拍图片URL',
  `alarm_message` varchar(512) NOT NULL COMMENT '告警消息',
  `alarm_time` datetime NOT NULL COMMENT '告警时间',
  `handler_id` bigint DEFAULT NULL COMMENT '处理人ID',
  `handler_name` varchar(64) DEFAULT NULL COMMENT '处理人姓名',
  `handle_time` datetime DEFAULT NULL COMMENT '处理时间',
  `handle_remark` varchar(512) DEFAULT NULL COMMENT '处理备注',
  `notification_sent` tinyint NOT NULL DEFAULT '0' COMMENT '是否已推送通知: 1-是, 0-否',
  `notification_time` datetime DEFAULT NULL COMMENT '通知推送时间',
  `extended_attributes` text COMMENT '扩展属性(JSON格式)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
  PRIMARY KEY (`alarm_id`),
  KEY `idx_device_time` (`device_id`,`alarm_time`),
  KEY `idx_event_type_time` (`event_type`,`alarm_time`),
  KEY `idx_alarm_status` (`alarm_status`),
  KEY `idx_alarm_level` (`alarm_level`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_event_id` (`event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警记录表（边缘计算架构）';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_alarm_record`
--

LOCK TABLES `t_video_alarm_record` WRITE;
/*!40000 ALTER TABLE `t_video_alarm_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_video_alarm_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_alarm_rule`
--

DROP TABLE IF EXISTS `t_video_alarm_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_alarm_rule` (
  `rule_id` bigint NOT NULL AUTO_INCREMENT COMMENT '规则ID',
  `rule_name` varchar(128) NOT NULL COMMENT '规则名称',
  `rule_type` varchar(64) NOT NULL COMMENT '规则类型: FALL_DETECTION-跌倒检测, LOITERING_DETECTION-徘徊检测, GATHERING_DETECTION-聚集检测, FIGHTING_DETECTION-打架检测, INTRUSION_DETECTION-入侵检测',
  `event_type` varchar(64) NOT NULL COMMENT '事件类型',
  `confidence_threshold` decimal(5,4) NOT NULL DEFAULT '0.8000' COMMENT '置信度阈值: 0.0000-1.0000',
  `area_id` bigint DEFAULT NULL COMMENT '区域ID（可选）',
  `device_id` varchar(64) DEFAULT NULL COMMENT '设备ID（可选，空表示所有设备）',
  `rule_status` tinyint NOT NULL DEFAULT '1' COMMENT '规则状态: 1-启用, 0-禁用',
  `effective_start_time` time DEFAULT NULL COMMENT '生效时间开始',
  `effective_end_time` time DEFAULT NULL COMMENT '生效时间结束',
  `alarm_level` tinyint NOT NULL DEFAULT '2' COMMENT '告警级别: 1-低, 2-中, 3-高, 4-紧急',
  `push_notification` tinyint NOT NULL DEFAULT '1' COMMENT '是否推送通知: 1-是, 0-否',
  `notification_methods` varchar(256) DEFAULT NULL COMMENT '通知方式(JSON): {"email":true,"sms":false,"websocket":true}',
  `alarm_message_template` varchar(512) DEFAULT NULL COMMENT '告警消息模板',
  `priority` int NOT NULL DEFAULT '0' COMMENT '规则优先级（数字越大优先级越高）',
  `extended_config` text COMMENT '扩展配置(JSON格式)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
  PRIMARY KEY (`rule_id`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_rule_status` (`rule_status`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_priority` (`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警规则表（边缘计算架构）';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_alarm_rule`
--

LOCK TABLES `t_video_alarm_rule` WRITE;
/*!40000 ALTER TABLE `t_video_alarm_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_video_alarm_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_device`
--

DROP TABLE IF EXISTS `t_video_device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_device` (
  `device_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®¾å¤‡ID',
  `device_no` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¾å¤‡ç¼–å·',
  `device_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è®¾å¤‡åç§°',
  `device_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¾å¤‡ç±»åž‹ï¼šCAMERA-æ‘„åƒå¤´ NVR-å½•åƒæœº DVR-ç¡¬ç›˜å½•åƒæœº',
  `area_id` bigint DEFAULT NULL COMMENT 'æ‰€åœ¨åŒºåŸŸID',
  `ip_address` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IPåœ°å€',
  `port` int DEFAULT NULL COMMENT 'ç«¯å£',
  `username` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'ç”¨æˆ·å',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'å¯†ç ï¼ˆåŠ å¯†ï¼‰',
  `rtsp_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'RTSPåœ°å€',
  `hls_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'HLSåœ°å€',
  `status` tinyint DEFAULT '1' COMMENT 'çŠ¶æ€ï¼š1-åœ¨çº¿ 2-ç¦»çº¿ 3-æ•…éšœ',
  `resolution` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'åˆ†è¾¨çŽ‡',
  `frame_rate` int DEFAULT NULL COMMENT 'å¸§çŽ‡',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'æ›´æ–°æ—¶é—´',
  `create_user_id` bigint DEFAULT NULL COMMENT 'åˆ›å»ºäººID',
  `update_user_id` bigint DEFAULT NULL COMMENT 'æ›´æ–°äººID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT 'åˆ é™¤æ ‡è®°ï¼š0-æœªåˆ é™¤ 1-å·²åˆ é™¤',
  `version` int DEFAULT '0' COMMENT 'ä¹è§‚é”ç‰ˆæœ¬å·',
  PRIMARY KEY (`device_id`),
  UNIQUE KEY `uk_device_no` (`device_no`,`deleted_flag`),
  KEY `idx_device_name` (`device_name`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è§†é¢‘è®¾å¤‡è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_device`
--

LOCK TABLES `t_video_device` WRITE;
/*!40000 ALTER TABLE `t_video_device` DISABLE KEYS */;
INSERT INTO `t_video_device` VALUES (1,'VIDEO_CAM001','å¤§é—¨ç›‘æŽ§æ‘„åƒå¤´1','CAMERA',5,'192.168.1.101',554,'admin','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','rtsp://admin:admin123@192.168.1.101:554/cam/realtime?channel=1&subtype=0',NULL,1,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(2,'VIDEO_CAM002','å¤§é—¨ç›‘æŽ§æ‘„åƒå¤´2','CAMERA',5,'192.168.1.102',554,'admin','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','rtsp://admin:admin123@192.168.1.102:554/cam/realtime?channel=1&subtype=0',NULL,1,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0),(3,'VIDEO_CAM003','é£Ÿå ‚ç›‘æŽ§æ‘„åƒå¤´','CAMERA',8,'192.168.1.103',554,'admin','$2a$10$7JB720yubVSd.TbpngVKpONx1YH8hEIXy2B/S8RDnm6FSL.NGxqa','rtsp://admin:admin123@192.168.1.103:554/cam/realtime?channel=1&subtype=0',NULL,1,NULL,NULL,'2025-12-15 20:20:55','2025-12-15 20:20:55',NULL,NULL,0,0);
/*!40000 ALTER TABLE `t_video_device` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_device_ai_event`
--

DROP TABLE IF EXISTS `t_video_device_ai_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_device_ai_event` (
  `event_id` varchar(64) NOT NULL COMMENT '事件ID',
  `device_id` varchar(64) NOT NULL COMMENT '设备ID',
  `device_code` varchar(128) NOT NULL COMMENT '设备编码',
  `event_type` varchar(64) NOT NULL COMMENT '事件类型: FALL_DETECTION-跌倒检测, LOITERING_DETECTION-徘徊检测, GATHERING_DETECTION-聚集检测, FIGHTING_DETECTION-打架检测, RUNNING_DETECTION-奔跑检测, CLIMBING_DETECTION-攀爬检测, FACE_DETECTION-人脸检测, INTRUSION_DETECTION-入侵检测',
  `confidence` decimal(5,4) NOT NULL COMMENT '置信度: 0.0000-1.0000',
  `bbox` varchar(256) DEFAULT NULL COMMENT '边界框(JSON格式): {"x":100,"y":150,"width":200,"height":300}',
  `snapshot` longblob COMMENT '抓拍图片',
  `event_time` datetime NOT NULL COMMENT '事件时间',
  `extended_attributes` text COMMENT '扩展属性(JSON格式)',
  `event_status` tinyint NOT NULL DEFAULT '0' COMMENT '事件状态: 0-待处理, 1-已处理, 2-已忽略',
  `process_time` datetime DEFAULT NULL COMMENT '处理时间',
  `alarm_id` varchar(64) DEFAULT NULL COMMENT '告警ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除, 1-已删除',
  PRIMARY KEY (`event_id`),
  KEY `idx_device_time` (`device_id`,`event_time`),
  KEY `idx_event_type_time` (`event_type`,`event_time`),
  KEY `idx_event_status` (`event_status`),
  KEY `idx_device_code` (`device_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备AI事件表（边缘计算架构）';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_device_ai_event`
--

LOCK TABLES `t_video_device_ai_event` WRITE;
/*!40000 ALTER TABLE `t_video_device_ai_event` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_video_device_ai_event` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_device_model_sync`
--

DROP TABLE IF EXISTS `t_video_device_model_sync`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_device_model_sync` (
  `sync_id` bigint NOT NULL AUTO_INCREMENT COMMENT '同步ID',
  `model_id` bigint NOT NULL COMMENT '模型ID',
  `device_id` varchar(100) NOT NULL COMMENT '设备ID',
  `sync_status` tinyint NOT NULL DEFAULT '0' COMMENT '同步状态（0-待同步 1-同步中 2-成功 3-失败）',
  `sync_progress` int DEFAULT '0' COMMENT '同步进度（0-100）',
  `sync_start_time` datetime DEFAULT NULL COMMENT '同步开始时间',
  `sync_end_time` datetime DEFAULT NULL COMMENT '同步结束时间',
  `error_message` varchar(1000) DEFAULT NULL COMMENT '错误信息',
  `retry_count` int NOT NULL DEFAULT '0' COMMENT '重试次数',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`sync_id`),
  UNIQUE KEY `uk_model_device` (`model_id`,`device_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_sync_status` (`sync_status`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_device_status_time` (`device_id`,`sync_status`,`create_time`),
  CONSTRAINT `fk_sync_model` FOREIGN KEY (`model_id`) REFERENCES `t_video_ai_model` (`model_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='设备模型同步表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_device_model_sync`
--

LOCK TABLES `t_video_device_model_sync` WRITE;
/*!40000 ALTER TABLE `t_video_device_model_sync` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_video_device_model_sync` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_video_record`
--

DROP TABLE IF EXISTS `t_video_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_video_record` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'å½•åƒID',
  `device_id` bigint NOT NULL COMMENT 'è®¾å¤‡ID',
  `start_time` datetime NOT NULL COMMENT 'å¼€å§‹æ—¶é—´',
  `end_time` datetime NOT NULL COMMENT 'ç»“æŸæ—¶é—´',
  `record_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'CONTINUOUS' COMMENT 'å½•åƒç±»åž‹ï¼šCONTINUOUS-è¿žç»­ MOTION-åŠ¨æ£€ MANUAL-æ‰‹åŠ¨',
  `file_path` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'æ–‡ä»¶è·¯å¾„',
  `file_size` bigint DEFAULT NULL COMMENT 'æ–‡ä»¶å¤§å°ï¼ˆå­—èŠ‚ï¼‰',
  `duration` int DEFAULT NULL COMMENT 'å½•åƒæ—¶é•¿ï¼ˆç§’ï¼‰',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`record_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_start_end_time` (`start_time`,`end_time`),
  KEY `idx_record_type` (`record_type`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è§†é¢‘å½•åƒè¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_video_record`
--

LOCK TABLES `t_video_record` WRITE;
/*!40000 ALTER TABLE `t_video_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_video_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_appointment`
--

DROP TABLE IF EXISTS `t_visitor_appointment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_appointment` (
  `appointment_id` bigint NOT NULL AUTO_INCREMENT,
  `visitor_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `visitor_phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `visitor_id_card` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `company` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `visit_reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `visit_date` date NOT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `interviewee_id` bigint NOT NULL,
  `interviewee_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `area_id` bigint DEFAULT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'PENDING',
  `approve_time` datetime DEFAULT NULL,
  `approver_id` bigint DEFAULT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_user_id` bigint DEFAULT NULL,
  `update_user_id` bigint DEFAULT NULL,
  `deleted_flag` tinyint DEFAULT '0',
  `version` int DEFAULT '0',
  `approval_required` tinyint DEFAULT '0' COMMENT '是否需要审批：0-否 1-是',
  `approval_workflow_id` bigint DEFAULT NULL COMMENT '审批工作流ID',
  `current_approval_level` tinyint DEFAULT '0' COMMENT '当前审批级别',
  `approval_status` tinyint DEFAULT '0' COMMENT '审批状态：0-待审批 1-已通过 2-已拒绝 3-进行中',
  `priority_level` tinyint DEFAULT '1' COMMENT '优先级：1-普通 2-重要 3-紧急',
  `special_requirements` text COLLATE utf8mb4_unicode_ci COMMENT '特殊要求（JSON格式）',
  `background_check_required` tinyint DEFAULT '0' COMMENT '是否需要背景调查：0-否 1-是',
  `security_level` tinyint DEFAULT '1' COMMENT '安全等级：1-普通 2-高级 3-机密',
  `access_permissions` text COLLATE utf8mb4_unicode_ci COMMENT '访问权限（JSON格式）',
  `escorts_required` tinyint DEFAULT '0' COMMENT '是否需要陪同：0-否 1-是',
  `escort_info` text COLLATE utf8mb4_unicode_ci COMMENT '陪同人员信息（JSON格式）',
  `auto_check_in` tinyint DEFAULT '0' COMMENT '自动签到：0-否 1-是',
  `auto_check_out` tinyint DEFAULT '0' COMMENT '自动签出：0-否 1-是',
  `access_card_issued` tinyint DEFAULT '0' COMMENT '是否已发访客卡：0-否 1-是',
  `card_activation_time` datetime DEFAULT NULL COMMENT '访客卡激活时间',
  `temporary_credentials` text COLLATE utf8mb4_unicode_ci COMMENT '临时凭证（JSON格式）',
  `biometric_enrolled` tinyint DEFAULT '0' COMMENT '生物识别已注册：0-否 1-是',
  `face_features` text COLLATE utf8mb4_unicode_ci COMMENT '人脸特征数据（Base64编码）',
  `fingerprint_data` text COLLATE utf8mb4_unicode_ci COMMENT '指纹数据（Base64编码）',
  `verification_methods` text COLLATE utf8mb4_unicode_ci COMMENT '验证方式（JSON格式）',
  `visit_history` text COLLATE utf8mb4_unicode_ci COMMENT '历史访问记录（JSON格式）',
  `compliance_check` tinyint DEFAULT '1' COMMENT '合规检查通过：0-否 1-是',
  `risk_assessment` text COLLATE utf8mb4_unicode_ci COMMENT '风险评估（JSON格式）',
  `emergency_contact` text COLLATE utf8mb4_unicode_ci COMMENT '紧急联系人（JSON格式）',
  `health_status` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'NORMAL' COMMENT '健康状态',
  `health_declaration` text COLLATE utf8mb4_unicode_ci COMMENT '健康声明（JSON格式）',
  `temperature_check` tinyint DEFAULT '0' COMMENT '体温检测：0-否 1-是',
  `temperature_value` decimal(3,1) DEFAULT NULL COMMENT '体温值',
  `mask_provided` tinyint DEFAULT '0' COMMENT '是否提供口罩：0-否 1-是',
  `sanitization_required` tinyint DEFAULT '1' COMMENT '是否需要消毒：0-否 1-是',
  `custom_fields` text COLLATE utf8mb4_unicode_ci COMMENT '自定义字段（JSON格式）',
  PRIMARY KEY (`appointment_id`),
  KEY `idx_visitor_phone` (`visitor_phone`),
  KEY `idx_visitor_id_card` (`visitor_id_card`),
  KEY `idx_visit_date` (`visit_date`),
  KEY `idx_interviewee_id` (`interviewee_id`),
  KEY `idx_area_id` (`area_id`),
  KEY `idx_status` (`status`),
  KEY `idx_approver_id` (`approver_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è®¿å®¢é¢„çº¦ç”³è¯·è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_appointment`
--

LOCK TABLES `t_visitor_appointment` WRITE;
/*!40000 ALTER TABLE `t_visitor_appointment` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_appointment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_approval_record`
--

DROP TABLE IF EXISTS `t_visitor_approval_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_approval_record` (
  `approval_id` bigint NOT NULL AUTO_INCREMENT COMMENT '审批记录ID',
  `appointment_id` bigint NOT NULL COMMENT '预约ID（外键关联）',
  `visitor_id` bigint NOT NULL COMMENT '访客ID（外键关联）',
  `approval_level` tinyint NOT NULL DEFAULT '1' COMMENT '审批级别：1-一级审批 2-二级审批 3-三级审批',
  `approval_step` tinyint NOT NULL DEFAULT '1' COMMENT '审批步骤：第几步审批',
  `total_steps` tinyint NOT NULL DEFAULT '1' COMMENT '总审批步骤数',
  `approval_role` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批角色：RECEPTION-前台 SECURITY-安保 ADMIN-管理员 MANAGER-经理',
  `approver_id` bigint NOT NULL COMMENT '审批人ID',
  `approver_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '审批人姓名',
  `department_id` bigint DEFAULT NULL COMMENT '审批部门ID',
  `approval_decision` tinyint NOT NULL COMMENT '审批决定：1-通过 2-拒绝 3-转交 4-待补充',
  `approval_result` tinyint DEFAULT '0' COMMENT '审批结果：0-待定 1-成功 2-失败',
  `approval_time` datetime DEFAULT NULL COMMENT '审批时间',
  `processing_duration` int DEFAULT NULL COMMENT '处理耗时（秒）',
  `approval_comment` text COLLATE utf8mb4_unicode_ci COMMENT '审批意见',
  `rejection_reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '拒绝原因',
  `required_conditions` text COLLATE utf8mb4_unicode_ci COMMENT '必要条件（JSON格式）',
  `checklist_items` text COLLATE utf8mb4_unicode_ci COMMENT '检查清单项目（JSON格式）',
  `attachment_urls` text COLLATE utf8mb4_unicode_ci COMMENT '附件URL（JSON格式）',
  `risk_assessment` text COLLATE utf8mb4_unicode_ci COMMENT '风险评估（JSON格式）',
  `priority_level` tinyint DEFAULT '1' COMMENT '优先级：1-普通 2-重要 3-紧急',
  `auto_approval` tinyint DEFAULT '0' COMMENT '是否自动审批：0-否 1-是',
  `auto_approval_rules` text COLLATE utf8mb4_unicode_ci COMMENT '自动审批规则（JSON格式）',
  `timeout_minutes` int DEFAULT '720' COMMENT '超时时间（分钟）',
  `reminder_sent` tinyint DEFAULT '0' COMMENT '是否已发送提醒：0-否 1-是',
  `next_approver_id` bigint DEFAULT NULL COMMENT '下一级审批人ID',
  `parallel_approval` tinyint DEFAULT '0' COMMENT '是否并行审批：0-否 1-是',
  `approval_flow` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT 'STANDARD' COMMENT '审批流程：STANDARD-标准 EMERGENCY-紧急 SIMPLIFIED-简化',
  `device_verification` tinyint DEFAULT '0' COMMENT '是否设备验证：0-否 1-是',
  `biometric_verification` tinyint DEFAULT '0' COMMENT '是否生物识别：0-否 1-是',
  `background_check` tinyint DEFAULT '0' COMMENT '是否背景调查：0-否 1-是',
  `special_permissions` text COLLATE utf8mb4_unicode_ci COMMENT '特殊权限（JSON格式）',
  `access_granted_areas` text COLLATE utf8mb4_unicode_ci COMMENT '授权访问区域（JSON格式）',
  `access_duration` int DEFAULT NULL COMMENT '访问时长（分钟）',
  `temporary_credentials` text COLLATE utf8mb4_unicode_ci COMMENT '临时凭证信息（JSON格式）',
  `monitoring_required` tinyint DEFAULT '0' COMMENT '是否需要监控：0-否 1-是',
  `escalation_rules` text COLLATE utf8mb4_unicode_ci COMMENT '升级规则（JSON格式）',
  `custom_workflow_data` text COLLATE utf8mb4_unicode_ci COMMENT '自定义工作流数据（JSON格式）',
  `integration_data` text COLLATE utf8mb4_unicode_ci COMMENT '集成数据（JSON格式）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标识：0-未删除 1-已删除',
  PRIMARY KEY (`approval_id`),
  UNIQUE KEY `uk_appointment_level_step` (`appointment_id`,`approval_level`,`approval_step`,`deleted_flag`),
  KEY `idx_visitor_approval` (`visitor_id`,`approval_decision`,`deleted_flag`),
  KEY `idx_approver_status` (`approver_id`,`approval_result`,`deleted_flag`),
  KEY `idx_approval_time` (`approval_time`,`deleted_flag`),
  KEY `idx_priority_time` (`priority_level`,`approval_time`,`deleted_flag`),
  KEY `idx_appointment` (`appointment_id`,`deleted_flag`),
  KEY `idx_create_time` (`create_time`,`deleted_flag`),
  KEY `idx_result_decision` (`approval_result`,`approval_decision`,`deleted_flag`),
  KEY `idx_visitor_approval_appointment` (`appointment_id`,`deleted_flag`),
  KEY `idx_visitor_approval_visitor` (`visitor_id`,`deleted_flag`),
  KEY `idx_visitor_approval_approver` (`approver_id`,`approval_result`,`deleted_flag`),
  KEY `idx_visitor_approval_time` (`approval_time`,`deleted_flag`),
  KEY `idx_visitor_approval_decision` (`approval_decision`,`approval_result`,`deleted_flag`),
  KEY `idx_visitor_approval_priority` (`priority_level`,`approval_time`,`deleted_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='访客审批记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_approval_record`
--

LOCK TABLES `t_visitor_approval_record` WRITE;
/*!40000 ALTER TABLE `t_visitor_approval_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_approval_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_blacklist`
--

DROP TABLE IF EXISTS `t_visitor_blacklist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_blacklist` (
  `blacklist_id` bigint NOT NULL AUTO_INCREMENT COMMENT '黑名单ID',
  `visitor_id` bigint DEFAULT NULL COMMENT '访客ID（可关联到t_visitor表）',
  `blacklist_type` tinyint NOT NULL DEFAULT '1' COMMENT '黑名单类型：1-永久黑名单 2-临时黑名单 3-监控黑名单',
  `blacklist_level` tinyint NOT NULL DEFAULT '1' COMMENT '黑名单级别：1-普通 2-重要 3-紧急',
  `blacklist_reason` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '黑名单原因',
  `id_card` varchar(18) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '身份证号（索引字段）',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号（索引字段）',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱地址',
  `face_features` text COLLATE utf8mb4_unicode_ci COMMENT '人脸特征数据（Base64编码）',
  `device_fingerprint` text COLLATE utf8mb4_unicode_ci COMMENT '设备指纹信息',
  `risk_score` int DEFAULT '50' COMMENT '风险评分（0-100）',
  `effective_time` datetime DEFAULT NULL COMMENT '生效时间',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间（NULL表示永久）',
  `blacklist_status` tinyint NOT NULL DEFAULT '1' COMMENT '黑名单状态：1-有效 2-暂停 3-已解除 4-待审核',
  `verification_count` int DEFAULT '0' COMMENT '违规次数',
  `last_violation_time` datetime DEFAULT NULL COMMENT '最后违规时间',
  `auto_clear_time` int DEFAULT NULL COMMENT '自动清除时间（小时）',
  `notify_enabled` tinyint DEFAULT '1' COMMENT '是否启用通知：0-否 1-是',
  `source` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'MANUAL' COMMENT '来源：MANUAL-手动 SYSTEM-系统 AUTO-自动 AI-人工智能',
  `operator_id` bigint DEFAULT NULL COMMENT '操作员ID',
  `operator_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作员姓名',
  `approval_required` tinyint DEFAULT '0' COMMENT '是否需要审批：0-否 1-是',
  `approval_status` tinyint DEFAULT '0' COMMENT '审批状态：0-待审批 1-已通过 2-已拒绝',
  `approver_id` bigint DEFAULT NULL COMMENT '审批人ID',
  `approval_time` datetime DEFAULT NULL COMMENT '审批时间',
  `approval_remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '审批备注',
  `related_incidents` text COLLATE utf8mb4_unicode_ci COMMENT '相关违规事件（JSON格式）',
  `block_devices` text COLLATE utf8mb4_unicode_ci COMMENT '拦截设备列表（JSON格式）',
  `whitelist_devices` text COLLATE utf8mb4_unicode_ci COMMENT '白名单设备（JSON格式）',
  `geographic_restrictions` text COLLATE utf8mb4_unicode_ci COMMENT '地域限制（JSON格式）',
  `time_restrictions` text COLLATE utf8mb4_unicode_ci COMMENT '时间限制（JSON格式）',
  `custom_attributes` text COLLATE utf8mb4_unicode_ci COMMENT '自定义属性（JSON格式）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标识：0-未删除 1-已删除',
  PRIMARY KEY (`blacklist_id`),
  UNIQUE KEY `uk_visitor_card` (`visitor_id`,`id_card`,`deleted_flag`),
  UNIQUE KEY `uk_visitor_phone` (`visitor_id`,`phone`,`deleted_flag`),
  KEY `idx_blacklist_type` (`blacklist_type`,`blacklist_status`,`deleted_flag`),
  KEY `idx_risk_score` (`risk_score`,`deleted_flag`),
  KEY `idx_effective_expire` (`effective_time`,`expire_time`,`deleted_flag`),
  KEY `idx_source_status` (`source`,`approval_status`,`deleted_flag`),
  KEY `idx_operator` (`operator_id`,`deleted_flag`),
  KEY `idx_create_time` (`create_time`,`deleted_flag`),
  KEY `idx_violation_time` (`last_violation_time`,`deleted_flag`),
  KEY `idx_visitor_blacklist_type_status` (`blacklist_type`,`blacklist_status`,`deleted_flag`),
  KEY `idx_visitor_blacklist_risk` (`risk_score`,`deleted_flag`),
  KEY `idx_visitor_blacklist_time` (`effective_time`,`expire_time`,`deleted_flag`),
  KEY `idx_visitor_blacklist_source` (`source`,`approval_status`,`deleted_flag`),
  KEY `idx_visitor_blacklist_operator` (`operator_id`,`deleted_flag`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='访客黑名单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_blacklist`
--

LOCK TABLES `t_visitor_blacklist` WRITE;
/*!40000 ALTER TABLE `t_visitor_blacklist` DISABLE KEYS */;
INSERT INTO `t_visitor_blacklist` VALUES (1,NULL,2,3,'临时访客违规示例',NULL,NULL,NULL,NULL,NULL,50,NULL,NULL,2,0,NULL,NULL,1,'SYSTEM',1,'系统管理员',0,0,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-12-27 11:01:53','2025-12-27 11:01:53',NULL,NULL,0),(2,NULL,1,2,'安全威胁访客',NULL,NULL,NULL,NULL,NULL,50,NULL,NULL,1,0,NULL,NULL,1,'MANUAL',1,'安全主管',1,0,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-12-27 11:01:53','2025-12-27 11:01:53',NULL,NULL,0),(3,NULL,1,3,'严重违规访客',NULL,NULL,NULL,NULL,NULL,50,NULL,NULL,3,0,NULL,NULL,1,'MANUAL',1,'安全管理员',1,0,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2025-12-27 11:01:53','2025-12-27 11:01:53',NULL,NULL,0);
/*!40000 ALTER TABLE `t_visitor_blacklist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_duration_statistics`
--

DROP TABLE IF EXISTS `t_visitor_duration_statistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_duration_statistics` (
  `statistics_id` bigint NOT NULL AUTO_INCREMENT COMMENT '统计ID',
  `stat_date` date NOT NULL COMMENT '统计日期',
  `total_visitors` int NOT NULL DEFAULT '0' COMMENT '总访客数',
  `average_duration` int NOT NULL DEFAULT '0' COMMENT '平均访问时长(分钟)',
  `min_duration` int NOT NULL DEFAULT '0' COMMENT '最短访问时长(分钟)',
  `max_duration` int NOT NULL DEFAULT '0' COMMENT '最长访问时长(分钟)',
  `overtime_count` int NOT NULL DEFAULT '0' COMMENT '超时访客数',
  `overtime_rate` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '超时率(%)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`statistics_id`),
  UNIQUE KEY `uk_duration_stat_date` (`stat_date`),
  KEY `idx_duration_stat_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='访问时长统计表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_duration_statistics`
--

LOCK TABLES `t_visitor_duration_statistics` WRITE;
/*!40000 ALTER TABLE `t_visitor_duration_statistics` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_duration_statistics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_record`
--

DROP TABLE IF EXISTS `t_visitor_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_record` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'è®°å½•ID',
  `appointment_id` bigint DEFAULT NULL COMMENT 'é¢„çº¦ID',
  `visitor_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¿å®¢å§“å',
  `visitor_phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è®¿å®¢æ‰‹æœºå·',
  `interviewee_id` bigint NOT NULL COMMENT 'è¢«è®¿äººID',
  `interviewee_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'è¢«è®¿äººå§“å',
  `area_id` bigint DEFAULT NULL COMMENT 'è®¿é—®åŒºåŸŸID',
  `actual_arrive_time` datetime DEFAULT NULL COMMENT 'å®žé™…åˆ°è¾¾æ—¶é—´',
  `actual_leave_time` datetime DEFAULT NULL COMMENT 'å®žé™…ç¦»å¼€æ—¶é—´',
  `access_card_no` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'è®¿é—®å¡å·',
  `access_result` tinyint DEFAULT '1' COMMENT 'è®¿é—®ç»“æžœï¼š1-æˆåŠŸ 2-å¤±è´¥',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'åˆ›å»ºæ—¶é—´',
  PRIMARY KEY (`record_id`),
  KEY `idx_appointment_id` (`appointment_id`),
  KEY `idx_visitor_phone` (`visitor_phone`),
  KEY `idx_interviewee_id` (`interviewee_id`),
  KEY `idx_actual_arrive_time` (`actual_arrive_time`),
  KEY `idx_access_result` (`access_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='è®¿å®¢è®°å½•è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_record`
--

LOCK TABLES `t_visitor_record` WRITE;
/*!40000 ALTER TABLE `t_visitor_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_satisfaction_statistics`
--

DROP TABLE IF EXISTS `t_visitor_satisfaction_statistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_satisfaction_statistics` (
  `statistics_id` bigint NOT NULL AUTO_INCREMENT COMMENT '统计ID',
  `stat_date` date NOT NULL COMMENT '统计日期',
  `total_count` int NOT NULL DEFAULT '0' COMMENT '总评价数',
  `average_score` decimal(3,2) NOT NULL DEFAULT '0.00' COMMENT '平均评分',
  `score_5_count` int NOT NULL DEFAULT '0' COMMENT '5分评价数',
  `score_4_count` int NOT NULL DEFAULT '0' COMMENT '4分评价数',
  `score_3_count` int NOT NULL DEFAULT '0' COMMENT '3分评价数',
  `score_2_count` int NOT NULL DEFAULT '0' COMMENT '2分评价数',
  `score_1_count` int NOT NULL DEFAULT '0' COMMENT '1分评价数',
  `satisfaction_rate` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '满意度(%)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`statistics_id`),
  UNIQUE KEY `uk_stat_date` (`stat_date`),
  KEY `idx_stat_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='访客满意度统计表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_satisfaction_statistics`
--

LOCK TABLES `t_visitor_satisfaction_statistics` WRITE;
/*!40000 ALTER TABLE `t_visitor_satisfaction_statistics` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_satisfaction_statistics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_self_check_out`
--

DROP TABLE IF EXISTS `t_visitor_self_check_out`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_self_check_out` (
  `check_out_id` bigint NOT NULL AUTO_INCREMENT COMMENT '签离记录ID',
  `registration_id` bigint NOT NULL COMMENT '登记ID',
  `visitor_code` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访客码',
  `visitor_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访客姓名',
  `check_out_time` datetime NOT NULL COMMENT '签离时间',
  `visit_duration` int NOT NULL COMMENT '访问时长(分钟)',
  `is_overtime` tinyint NOT NULL DEFAULT '0' COMMENT '是否超时: 0-否 1-是',
  `overtime_duration` int NOT NULL DEFAULT '0' COMMENT '超时时长(分钟)',
  `check_out_method` tinyint NOT NULL COMMENT '签离方式: 1-自助签离 2-人工签离',
  `check_out_status` tinyint NOT NULL DEFAULT '0' COMMENT '签离状态: 0-待签离 1-已完成 2-已取消',
  `terminal_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '签离终端ID',
  `terminal_location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '签离终端位置',
  `card_return_status` tinyint NOT NULL COMMENT '卡归还状态: 0-未归还 1-已归还 2-卡遗失',
  `visitor_card` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访客卡号',
  `operator_id` bigint DEFAULT NULL COMMENT '操作人ID(人工签离时)',
  `operator_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作人姓名(人工签离时)',
  `manual_reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '人工签离原因',
  `satisfaction_score` tinyint DEFAULT NULL COMMENT '满意度评分(1-5分)',
  `visitor_feedback` text COLLATE utf8mb4_unicode_ci COMMENT '访客反馈',
  `feedback_time` datetime DEFAULT NULL COMMENT '反馈时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除 1-已删除',
  PRIMARY KEY (`check_out_id`),
  UNIQUE KEY `uk_checkout_visitor_code` (`visitor_code`),
  KEY `idx_checkout_registration` (`registration_id`),
  KEY `idx_checkout_time` (`check_out_time`),
  KEY `idx_checkout_status` (`check_out_status`),
  KEY `idx_checkout_overtime` (`is_overtime`),
  CONSTRAINT `t_visitor_self_check_out_ibfk_1` FOREIGN KEY (`registration_id`) REFERENCES `t_visitor_self_service_registration` (`registration_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='自助签离记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_self_check_out`
--

LOCK TABLES `t_visitor_self_check_out` WRITE;
/*!40000 ALTER TABLE `t_visitor_self_check_out` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_self_check_out` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_self_service_registration`
--

DROP TABLE IF EXISTS `t_visitor_self_service_registration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_self_service_registration` (
  `registration_id` bigint NOT NULL AUTO_INCREMENT COMMENT '登记ID',
  `registration_code` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登记码',
  `visitor_code` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访客码',
  `visitor_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '访客姓名',
  `id_card_type` tinyint NOT NULL COMMENT '证件类型: 1-身份证 2-护照 3-其他',
  `id_card` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '证件号码',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号码',
  `visitor_type` tinyint NOT NULL COMMENT '访客类型: 1-临时访客 2-常客 3-VIP访客',
  `visit_purpose` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '来访目的',
  `interviewee_id` bigint NOT NULL COMMENT '被访人ID',
  `interviewee_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '被访人姓名',
  `interviewee_department` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '被访人部门',
  `visit_date` date NOT NULL COMMENT '访问日期',
  `expected_enter_time` datetime NOT NULL COMMENT '预计进入时间',
  `expected_leave_time` datetime NOT NULL COMMENT '预计离开时间',
  `face_photo_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '人脸照片URL',
  `terminal_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '终端ID',
  `terminal_location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '终端位置',
  `registration_status` tinyint NOT NULL DEFAULT '0' COMMENT '登记状态: 0-待审批 1-审批通过 2-审批拒绝 3-已签到 4-已完成',
  `approver_id` bigint DEFAULT NULL COMMENT '审批人ID',
  `approver_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '审批人姓名',
  `approval_time` datetime DEFAULT NULL COMMENT '审批时间',
  `approval_comment` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '审批意见',
  `check_in_time` datetime DEFAULT NULL COMMENT '签到时间',
  `check_in_terminal` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '签到终端',
  `check_out_time` datetime DEFAULT NULL COMMENT '签离时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_flag` tinyint NOT NULL DEFAULT '0' COMMENT '删除标记: 0-未删除 1-已删除',
  PRIMARY KEY (`registration_id`),
  UNIQUE KEY `uk_registration_code` (`registration_code`),
  UNIQUE KEY `uk_visitor_code` (`visitor_code`),
  KEY `idx_registration_visitor` (`visitor_name`,`phone`),
  KEY `idx_registration_interviewee` (`interviewee_id`),
  KEY `idx_registration_date` (`visit_date`),
  KEY `idx_registration_status` (`registration_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='自助访客登记表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_self_service_registration`
--

LOCK TABLES `t_visitor_self_service_registration` WRITE;
/*!40000 ALTER TABLE `t_visitor_self_service_registration` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_self_service_registration` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_visitor_vehicle`
--

DROP TABLE IF EXISTS `t_visitor_vehicle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_visitor_vehicle` (
  `vehicle_id` bigint NOT NULL AUTO_INCREMENT COMMENT '车辆ID',
  `visitor_id` bigint DEFAULT NULL COMMENT '访客ID',
  `plate_no` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车牌号',
  `vehicle_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车辆类型',
  `vehicle_color` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车辆颜色',
  `vehicle_model` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车辆型号',
  `vin` varchar(17) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车辆识别码',
  `insurance_policy` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '保险单号',
  `insurance_expire_date` date DEFAULT NULL COMMENT '保险到期日期',
  `registration_date` date DEFAULT NULL COMMENT '注册日期',
  `vehicle_owner` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车主姓名',
  `owner_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '车主电话',
  `owner_relation` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '与访客关系',
  `transportation_purpose` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '运输目的',
  `cargo_description` text COLLATE utf8mb4_unicode_ci COMMENT '货物描述（JSON格式）',
  `weight_limit` decimal(10,2) DEFAULT NULL COMMENT '重量限制（公斤）',
  `hazardous_materials` text COLLATE utf8mb4_unicode_ci COMMENT '危险品信息（JSON格式）',
  `transport_permit` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '运输许可证号',
  `permit_valid_until` date DEFAULT NULL COMMENT '许可证有效期',
  `parking_allowed` tinyint DEFAULT '1' COMMENT '是否允许停车：0-否 1-是',
  `parking_area` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '停车区域',
  `parking_spot` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '停车位置',
  `entrance_permission` tinyint DEFAULT '1' COMMENT '入口权限：0-无 1-行人 2-车辆 3-混合',
  `exit_permission` tinyint DEFAULT '1' COMMENT '出口权限：0-无 1-行人 2-车辆 3-混合',
  `vehicle_images` text COLLATE utf8mb4_unicode_ci COMMENT '车辆照片（JSON格式）',
  `document_urls` text COLLATE utf8mb4_unicode_ci COMMENT '证件URL（JSON格式）',
  `security_check` tinyint DEFAULT '1' COMMENT '安全检查：0-未检查 1-已通过 2-未通过',
  `check_results` text COLLATE utf8mb4_unicode_ci COMMENT '检查结果（JSON格式）',
  `access_log` text COLLATE utf8mb4_unicode_ci COMMENT '访问日志（JSON格式）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_user_id` bigint DEFAULT NULL COMMENT '创建人ID',
  `update_user_id` bigint DEFAULT NULL COMMENT '更新人ID',
  `deleted_flag` tinyint DEFAULT '0' COMMENT '删除标识：0-未删除 1-已删除',
  PRIMARY KEY (`vehicle_id`),
  KEY `idx_visitor_id` (`visitor_id`,`deleted_flag`),
  KEY `idx_plate_no` (`plate_no`,`deleted_flag`),
  KEY `idx_vin` (`vin`,`deleted_flag`),
  KEY `idx_create_time` (`create_time`,`deleted_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='访客车辆信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_visitor_vehicle`
--

LOCK TABLES `t_visitor_vehicle` WRITE;
/*!40000 ALTER TABLE `t_visitor_vehicle` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_visitor_vehicle` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-30  4:11:49
