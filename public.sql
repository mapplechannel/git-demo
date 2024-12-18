/*
 Navicat Premium Data Transfer

 Source Server         : postgres
 Source Server Type    : PostgreSQL
 Source Server Version : 100023 (100023)
 Source Host           : localhost:5432
 Source Catalog        : postgres
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 100023 (100023)
 File Encoding         : 65001

 Date: 18/12/2024 17:46:09
*/


-- ----------------------------
-- Sequence structure for job_serial
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."job_serial";
CREATE SEQUENCE "public"."job_serial" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."job_serial" OWNER TO "postgres";

-- ----------------------------
-- Table structure for abnormal
-- ----------------------------
DROP TABLE IF EXISTS "public"."abnormal";
CREATE TABLE "public"."abnormal" (
  "id" int4 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "abnormal_num" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."abnormal" OWNER TO "postgres";

-- ----------------------------
-- Records of abnormal
-- ----------------------------
BEGIN;
INSERT INTO "public"."abnormal" ("id", "name", "abnormal_num") VALUES (1, '工艺', '79');
INSERT INTO "public"."abnormal" ("id", "name", "abnormal_num") VALUES (2, '物料', '58');
INSERT INTO "public"."abnormal" ("id", "name", "abnormal_num") VALUES (3, '设备', '26');
INSERT INTO "public"."abnormal" ("id", "name", "abnormal_num") VALUES (4, '安全', '28');
COMMIT;

-- ----------------------------
-- Table structure for acount
-- ----------------------------
DROP TABLE IF EXISTS "public"."acount";
CREATE TABLE "public"."acount" (
  "aot" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."acount" OWNER TO "postgres";

-- ----------------------------
-- Records of acount
-- ----------------------------
BEGIN;
INSERT INTO "public"."acount" ("aot") VALUES ('sss');
INSERT INTO "public"."acount" ("aot") VALUES ('sss');
INSERT INTO "public"."acount" ("aot") VALUES ('sss');
COMMIT;

-- ----------------------------
-- Table structure for devicename
-- ----------------------------
DROP TABLE IF EXISTS "public"."devicename";
CREATE TABLE "public"."devicename" (
  "id" int4 NOT NULL,
  "device" varchar(255) COLLATE "pg_catalog"."default",
  "point" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."devicename" OWNER TO "postgres";

-- ----------------------------
-- Records of devicename
-- ----------------------------
BEGIN;
INSERT INTO "public"."devicename" ("id", "device", "point") VALUES (3, '厂房11_采集箱4', '15');
INSERT INTO "public"."devicename" ("id", "device", "point") VALUES (4, '厂房11_采集箱5', '16');
INSERT INTO "public"."devicename" ("id", "device", "point") VALUES (5, '厂房11_采集箱6', '17');
INSERT INTO "public"."devicename" ("id", "device", "point") VALUES (1, '厂房11_采集箱2', '11');
INSERT INTO "public"."devicename" ("id", "device", "point") VALUES (2, '厂房11_采集箱3', '22');
COMMIT;

-- ----------------------------
-- Table structure for executors
-- ----------------------------
DROP TABLE IF EXISTS "public"."executors";
CREATE TABLE "public"."executors" (
  "id" varchar(19) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "createuser" varchar(255) COLLATE "pg_catalog"."default",
  "editetime" timestamp(6),
  "desc" varchar(500) COLLATE "pg_catalog"."default",
  "max_tasks" int2,
  "exist_tasks" int2,
  "is_default" bool
)
;
ALTER TABLE "public"."executors" OWNER TO "postgres";

-- ----------------------------
-- Records of executors
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for materialalter
-- ----------------------------
DROP TABLE IF EXISTS "public"."materialalter";
CREATE TABLE "public"."materialalter" (
  "id" int4 NOT NULL,
  "materialid" int4,
  "changetype" varchar(255) COLLATE "pg_catalog"."default",
  "olddata" jsonb,
  "newdata" jsonb,
  "changereason" varchar(255) COLLATE "pg_catalog"."default",
  "changeedby" varchar(255) COLLATE "pg_catalog"."default",
  "changetimestamp" timestamp(6),
  "version" int4,
  "isdetele" int4
)
;
ALTER TABLE "public"."materialalter" OWNER TO "postgres";

-- ----------------------------
-- Records of materialalter
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for plan_comp
-- ----------------------------
DROP TABLE IF EXISTS "public"."plan_comp";
CREATE TABLE "public"."plan_comp" (
  "id" int4 NOT NULL,
  "时间" varchar COLLATE "pg_catalog"."default",
  "完成率" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."plan_comp" OWNER TO "postgres";

-- ----------------------------
-- Records of plan_comp
-- ----------------------------
BEGIN;
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (2, '10:00:10', '10');
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (3, '10:00:15', '15');
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (5, '10:00:25', '25');
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (6, '10:00:30', '30');
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (1, '10:00:05', '5');
INSERT INTO "public"."plan_comp" ("id", "时间", "完成率") VALUES (4, '10:00:20', '20');
COMMIT;

-- ----------------------------
-- Table structure for production
-- ----------------------------
DROP TABLE IF EXISTS "public"."production";
CREATE TABLE "public"."production" (
  "序号" int4 NOT NULL,
  "工单编号" varchar(255) COLLATE "pg_catalog"."default",
  "批次号" varchar(255) COLLATE "pg_catalog"."default",
  "计划开始时间" varchar(255) COLLATE "pg_catalog"."default",
  "计划结束时间" varchar(255) COLLATE "pg_catalog"."default",
  "实际开始时间" varchar(255) COLLATE "pg_catalog"."default",
  "实际结束时间" varchar(255) COLLATE "pg_catalog"."default",
  "工单名称" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."production" OWNER TO "postgres";

-- ----------------------------
-- Records of production
-- ----------------------------
BEGIN;
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (1, 'A11101', '1121', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_6板制作');
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (3, 'A11103', '1120', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_1板制作');
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (4, 'A11104', '1120', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_2板制作');
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (5, 'A11105', '1120', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_3板制作');
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (6, 'A11106', '1120', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_4板制作');
INSERT INTO "public"."production" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "工单名称") VALUES (2, 'A11102', '1121', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', 'PCB_5板制作');
COMMIT;

-- ----------------------------
-- Table structure for pronum
-- ----------------------------
DROP TABLE IF EXISTS "public"."pronum";
CREATE TABLE "public"."pronum" (
  "proname" varchar(255) COLLATE "pg_catalog"."default",
  "pronum" int4
)
;
ALTER TABLE "public"."pronum" OWNER TO "postgres";

-- ----------------------------
-- Records of pronum
-- ----------------------------
BEGIN;
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('半成品人工1号线', 79);
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('半成品人工2号线', 58);
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('半成品人工3号线', 26);
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('柔性线1号线', 28);
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('柔性线2号线', 49);
INSERT INTO "public"."pronum" ("proname", "pronum") VALUES ('柔性线3号线', 62);
COMMIT;

-- ----------------------------
-- Table structure for prooutput
-- ----------------------------
DROP TABLE IF EXISTS "public"."prooutput";
CREATE TABLE "public"."prooutput" (
  "id" int4 NOT NULL,
  "plan" int4,
  "actual" int4,
  "name" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."prooutput" OWNER TO "postgres";

-- ----------------------------
-- Records of prooutput
-- ----------------------------
BEGIN;
INSERT INTO "public"."prooutput" ("id", "plan", "actual", "name") VALUES (1, 79, 58, '水涂料');
INSERT INTO "public"."prooutput" ("id", "plan", "actual", "name") VALUES (2, 26, 28, '瓷涂料');
COMMIT;

-- ----------------------------
-- Table structure for quyukongtiao
-- ----------------------------
DROP TABLE IF EXISTS "public"."quyukongtiao";
CREATE TABLE "public"."quyukongtiao" (
  "id" int4 NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "count" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."quyukongtiao" OWNER TO "postgres";

-- ----------------------------
-- Records of quyukongtiao
-- ----------------------------
BEGIN;
INSERT INTO "public"."quyukongtiao" ("id", "name", "count") VALUES (2, '模块区', '16');
INSERT INTO "public"."quyukongtiao" ("id", "name", "count") VALUES (1, '模板区', '25');
INSERT INTO "public"."quyukongtiao" ("id", "name", "count") VALUES (3, '检测区', '12');
INSERT INTO "public"."quyukongtiao" ("id", "name", "count") VALUES (4, '办公区', '8');
COMMIT;

-- ----------------------------
-- Table structure for qwercount
-- ----------------------------
DROP TABLE IF EXISTS "public"."qwercount";
CREATE TABLE "public"."qwercount" (
  "aot" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."qwercount" OWNER TO "postgres";

-- ----------------------------
-- Records of qwercount
-- ----------------------------
BEGIN;
INSERT INTO "public"."qwercount" ("aot") VALUES ('55');
INSERT INTO "public"."qwercount" ("aot") VALUES ('55');
INSERT INTO "public"."qwercount" ("aot") VALUES ('55');
COMMIT;

-- ----------------------------
-- Table structure for tb_namespace
-- ----------------------------
DROP TABLE IF EXISTS "public"."tb_namespace";
CREATE TABLE "public"."tb_namespace" (
  "id" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "version" varchar(1024) COLLATE "pg_catalog"."default",
  "namespace" varchar(1024) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."tb_namespace" OWNER TO "postgres";

-- ----------------------------
-- Records of tb_namespace
-- ----------------------------
BEGIN;
INSERT INTO "public"."tb_namespace" ("id", "version", "namespace") VALUES ('d521795911e1e934ad99c5f9d9a4fb4e', 'V20241112190157', 'devicedevice');
INSERT INTO "public"."tb_namespace" ("id", "version", "namespace") VALUES ('qwer', 'fajogfag', 'ffffffff');
COMMIT;

-- ----------------------------
-- Table structure for testbugjiaoyan
-- ----------------------------
DROP TABLE IF EXISTS "public"."testbugjiaoyan";
CREATE TABLE "public"."testbugjiaoyan" (
  "testt" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."testbugjiaoyan" OWNER TO "postgres";

-- ----------------------------
-- Records of testbugjiaoyan
-- ----------------------------
BEGIN;
INSERT INTO "public"."testbugjiaoyan" ("testt") VALUES ('89');
INSERT INTO "public"."testbugjiaoyan" ("testt") VALUES ('89');
INSERT INTO "public"."testbugjiaoyan" ("testt") VALUES ('89');
COMMIT;

-- ----------------------------
-- Table structure for testcn
-- ----------------------------
DROP TABLE IF EXISTS "public"."testcn";
CREATE TABLE "public"."testcn" (
  "id" int4 NOT NULL,
  "real1" int4,
  "day1" int4,
  "month1" int4,
  "real2" int4,
  "day2" int4,
  "month2" int4,
  "real3" int4,
  "day3" int4,
  "month3" int4,
  "real4" int4,
  "day4" int4,
  "month4" int4,
  "lwater" int4,
  "pwater" int4
)
;
ALTER TABLE "public"."testcn" OWNER TO "postgres";

-- ----------------------------
-- Records of testcn
-- ----------------------------
BEGIN;
INSERT INTO "public"."testcn" ("id", "real1", "day1", "month1", "real2", "day2", "month2", "real3", "day3", "month3", "real4", "day4", "month4", "lwater", "pwater") VALUES (2, 223, 850, 8900, 10, 800, 7543, 52, 861, 9024, 201, 897, 8875, 63, 115);
INSERT INTO "public"."testcn" ("id", "real1", "day1", "month1", "real2", "day2", "month2", "real3", "day3", "month3", "real4", "day4", "month4", "lwater", "pwater") VALUES (3, 223, 850, 8900, 15, 800, 7543, 52, 861, 9024, 201, 897, 8875, 63, 115);
INSERT INTO "public"."testcn" ("id", "real1", "day1", "month1", "real2", "day2", "month2", "real3", "day3", "month3", "real4", "day4", "month4", "lwater", "pwater") VALUES (1, 117, 850, 8900, 248, 800, 7543, 182, 887, 9038, 231, 825, 9534, 25, 64);
COMMIT;

-- ----------------------------
-- Table structure for testlist
-- ----------------------------
DROP TABLE IF EXISTS "public"."testlist";
CREATE TABLE "public"."testlist" (
  "id" int4,
  "name" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."testlist" OWNER TO "postgres";

-- ----------------------------
-- Records of testlist
-- ----------------------------
BEGIN;
INSERT INTO "public"."testlist" ("id", "name") VALUES (1, 's');
INSERT INTO "public"."testlist" ("id", "name") VALUES (2, 'b');
INSERT INTO "public"."testlist" ("id", "name") VALUES (3, 'a');
INSERT INTO "public"."testlist" ("id", "name") VALUES (4, 'ss');
INSERT INTO "public"."testlist" ("id", "name") VALUES (5, 'bb');
INSERT INTO "public"."testlist" ("id", "name") VALUES (6, 'aa');
COMMIT;

-- ----------------------------
-- Table structure for testlqbz
-- ----------------------------
DROP TABLE IF EXISTS "public"."testlqbz";
CREATE TABLE "public"."testlqbz" (
  "id" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "vesion" varchar COLLATE "pg_catalog"."default",
  "namesapce" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."testlqbz" OWNER TO "postgres";

-- ----------------------------
-- Records of testlqbz
-- ----------------------------
BEGIN;
INSERT INTO "public"."testlqbz" ("id", "vesion", "namesapce") VALUES ('d521795911e1e934ad99c5f9d9a4fb4e', 'V20241112190157', 'devicedevice');
INSERT INTO "public"."testlqbz" ("id", "vesion", "namesapce") VALUES ('qwer', 'fajogfag', 'ffffffff');
COMMIT;

-- ----------------------------
-- Table structure for testnamespace
-- ----------------------------
DROP TABLE IF EXISTS "public"."testnamespace";
CREATE TABLE "public"."testnamespace" (
  "id" varchar COLLATE "pg_catalog"."default",
  "version" varchar COLLATE "pg_catalog"."default",
  "namespace" varchar COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."testnamespace" OWNER TO "postgres";

-- ----------------------------
-- Records of testnamespace
-- ----------------------------
BEGIN;
INSERT INTO "public"."testnamespace" ("id", "version", "namespace") VALUES ('d521795911e1e934ad99c5f9d9a4fb4e', 'V20241112190157', 'devicedevice');
INSERT INTO "public"."testnamespace" ("id", "version", "namespace") VALUES ('qwer', 'fajogfag', 'ffffffff');
COMMIT;

-- ----------------------------
-- Table structure for work_order_detail
-- ----------------------------
DROP TABLE IF EXISTS "public"."work_order_detail";
CREATE TABLE "public"."work_order_detail" (
  "序号" int4 NOT NULL,
  "工单编号" varchar(255) COLLATE "pg_catalog"."default",
  "批次号" varchar(255) COLLATE "pg_catalog"."default",
  "计划开始时间" varchar(255) COLLATE "pg_catalog"."default",
  "计划结束时间" varchar(255) COLLATE "pg_catalog"."default",
  "实际开始时间" varchar(255) COLLATE "pg_catalog"."default",
  "实际结束时间" varchar(255) COLLATE "pg_catalog"."default",
  "延期时长" varchar(255) COLLATE "pg_catalog"."default",
  "工单名称" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."work_order_detail" OWNER TO "postgres";

-- ----------------------------
-- Records of work_order_detail
-- ----------------------------
BEGIN;
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (2, 'A11102', '1121', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '78425', 'PCB_6板制作');
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (1, 'A11101', '1121', '2024-12-18 10:00:00', '10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '0', 'PCB_1板制作');
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (3, 'A11103', '1120', '2024-12-18 10:00:00', '10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '125', 'PCB_2板制作');
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (4, 'A11104', '1120', '2024-12-18 10:00:00', '10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '20', 'PCB_3板制作');
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (5, 'A11105', '1120', '2024-12-18 10:00:00', '10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '50000', 'PCB_4板制作');
INSERT INTO "public"."work_order_detail" ("序号", "工单编号", "批次号", "计划开始时间", "计划结束时间", "实际开始时间", "实际结束时间", "延期时长", "工单名称") VALUES (6, 'A11106', '1120', '2024-12-18 10:00:00', '10:00:00', '2024-12-18 10:00:00', '2024-12-19 10:00:00', '554', 'PCB_5板制作');
COMMIT;

-- ----------------------------
-- Table structure for zhaopan
-- ----------------------------
DROP TABLE IF EXISTS "public"."zhaopan";
CREATE TABLE "public"."zhaopan" (
  "id" int4 NOT NULL,
  "job_name" varchar(255) COLLATE "pg_catalog"."default",
  "11111" varchar(255) COLLATE "pg_catalog"."default",
  "345_GF" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."zhaopan" OWNER TO "postgres";
COMMENT ON COLUMN "public"."zhaopan"."id" IS '主键';
COMMENT ON COLUMN "public"."zhaopan"."job_name" IS '作业名称';

-- ----------------------------
-- Records of zhaopan
-- ----------------------------
BEGIN;
INSERT INTO "public"."zhaopan" ("id", "job_name", "11111", "345_GF") VALUES (1, '89', '1111', 'aaa');
COMMIT;

-- ----------------------------
-- Function structure for show_create_table
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."show_create_table"("in_schema_name" varchar, "in_table_name" varchar);
CREATE OR REPLACE FUNCTION "public"."show_create_table"("in_schema_name" varchar, "in_table_name" varchar)
  RETURNS "pg_catalog"."text" AS $BODY$ DECLARE-- the ddl we're building
	v_table_ddl TEXT;
-- data about the target table
v_table_oid INT;
v_table_type CHAR;
v_partition_key VARCHAR;
v_table_comment VARCHAR;
-- records for looping
v_column_record record;
v_constraint_record record;
v_index_record record;
v_column_comment_record record;
v_index_comment_record record;
v_constraint_comment_record record;
BEGIN-- grab the oid of the table; https://www.postgresql.org/docs/8.3/catalog-pg-class.html
	SELECT C
		.OID,
		C.relkind INTO v_table_oid,
		v_table_type 
	FROM
		pg_catalog.pg_class
		C LEFT JOIN pg_catalog.pg_namespace n ON n.OID = C.relnamespace 
	WHERE
		C.relkind IN ( 'r', 'p' ) 
		AND C.relname = in_table_name -- the table name
		
		AND n.nspname = in_schema_name;-- the schema
-- throw an error if table was not found
	IF
		( v_table_oid IS NULL ) THEN
			RAISE EXCEPTION 'table does not exist';
		
	END IF;
-- start the create definition
	v_table_ddl := 'CREATE TABLE "' || in_schema_name || '"."' || in_table_name || '" (' || E'\n';
-- define all of the columns in the table; https://stackoverflow.com/a/8153081/3068233
	FOR v_column_record IN SELECT C
	.COLUMN_NAME,
	C.data_type,
	C.character_maximum_length,
	C.is_nullable,
	C.column_default 
	FROM
		information_schema.COLUMNS C 
	WHERE
		( table_schema, TABLE_NAME ) = ( in_schema_name, in_table_name ) 
	ORDER BY
		ordinal_position
		LOOP
		v_table_ddl := v_table_ddl || '  ' -- note: two char spacer to start, to indent the column
		|| '"' || v_column_record.COLUMN_NAME || '" ' || v_column_record.data_type ||
	CASE
			
			WHEN v_column_record.character_maximum_length IS NOT NULL THEN
			( '(' || v_column_record.character_maximum_length || ')' ) ELSE'' 
		END || ' ' ||
CASE
	
	WHEN v_column_record.is_nullable = 'NO' THEN
	'NOT NULL' ELSE'NULL' 
	END ||
CASE
	
	WHEN v_column_record.column_default IS NOT NULL THEN
	( ' DEFAULT ' || v_column_record.column_default ) ELSE'' 
	END || ',' || E'\n';

END LOOP;
-- define all the constraints in the; https://www.postgresql.org/docs/9.1/catalog-pg-constraint.html && https://dba.stackexchange.com/a/214877/75296
FOR v_constraint_record IN SELECT
con.conname AS CONSTRAINT_NAME,
con.contype AS constraint_type,
CASE
		
		WHEN con.contype = 'p' THEN
		1 -- primary key constraint
		
		WHEN con.contype = 'u' THEN
		2 -- unique constraint
		
		WHEN con.contype = 'f' THEN
		3 -- foreign key constraint
		
		WHEN con.contype = 'c' THEN
		4 ELSE 5 
	END AS type_rank,
	pg_get_constraintdef ( con.OID ) AS constraint_definition 
FROM
	pg_catalog.pg_constraint con
	JOIN pg_catalog.pg_class rel ON rel.OID = con.conrelid
	JOIN pg_catalog.pg_namespace nsp ON nsp.OID = connamespace 
WHERE
	nsp.nspname = in_schema_name 
	AND rel.relname = in_table_name 
ORDER BY
	type_rank
	LOOP
IF
	v_constraint_record.constraint_type = 'p' THEN
		v_table_ddl := v_table_ddl || '  ' || v_constraint_record.constraint_definition || ',' || E'\n';
	ELSE v_table_ddl := v_table_ddl || '  ' -- note: two char spacer to start, to indent the column
	|| 'CONSTRAINT' || ' ' || '"' || v_constraint_record.CONSTRAINT_NAME || '" ' || v_constraint_record.constraint_definition || ',' || E'\n';
	
END IF;

END LOOP;
-- drop the last comma before ending the create statement
v_table_ddl = substr( v_table_ddl, 0, LENGTH ( v_table_ddl ) - 1 ) || E'\n';
-- end the create definition
v_table_ddl := v_table_ddl || ')';
IF
	v_table_type = 'p' THEN
	SELECT
		pg_get_partkeydef ( v_table_oid ) INTO v_partition_key;
	IF
		v_partition_key IS NOT NULL THEN
			v_table_ddl := v_table_ddl || ' PARTITION BY ' || v_partition_key;
		
	END IF;
	
END IF;
v_table_ddl := v_table_ddl || ';' || E'\n';
-- suffix create statement with all of the indexes on the table
FOR v_index_record IN SELECT
regexp_replace( indexdef, ' "?' || schemaname || '"?\.', ' ' ) AS indexdef 
FROM
	pg_catalog.pg_indexes 
WHERE
	( schemaname, tablename ) = ( in_schema_name, in_table_name ) 
	AND indexname NOT IN (
	SELECT
		con.conname 
	FROM
		pg_catalog.pg_constraint con
		JOIN pg_catalog.pg_class rel ON rel.OID = con.conrelid
		JOIN pg_catalog.pg_namespace nsp ON nsp.OID = connamespace 
	WHERE
		nsp.nspname = in_schema_name 
		AND rel.relname = in_table_name 
	)
	LOOP
	v_table_ddl := v_table_ddl || v_index_record.indexdef || ';' || E'\n';

END LOOP;
-- comment on table
SELECT
	description INTO v_table_comment 
FROM
	pg_catalog.pg_description 
WHERE
	objoid = v_table_oid 
	AND objsubid = 0;
IF
	v_table_comment IS NOT NULL THEN
		v_table_ddl := v_table_ddl || 'COMMENT ON TABLE "' || in_table_name || '" IS ''' || REPLACE ( v_table_comment, '''', '''''' ) || ''';' || E'\n';
	
END IF;
-- comment on column
FOR v_column_comment_record IN SELECT
col.COLUMN_NAME,
d.description 
FROM
	information_schema.COLUMNS col
	JOIN pg_catalog.pg_class C ON C.relname = col.
	TABLE_NAME JOIN pg_catalog.pg_namespace nsp ON nsp.OID = C.relnamespace 
	AND col.table_schema = nsp.nspname
	JOIN pg_catalog.pg_description d ON d.objoid = C.OID 
	AND d.objsubid = col.ordinal_position 
WHERE
	C.OID = v_table_oid 
ORDER BY
	col.ordinal_position
	LOOP
	v_table_ddl := v_table_ddl || 'COMMENT ON COLUMN "' || in_table_name || '"."' || v_column_comment_record.COLUMN_NAME || '" IS ''' || REPLACE ( v_column_comment_record.description, '''', '''''' ) || ''';' || E'\n';

END LOOP;
-- comment on index
FOR v_index_comment_record IN SELECT C
.relname,
d.description 
FROM
	pg_catalog.pg_index idx
	JOIN pg_catalog.pg_class C ON idx.indexrelid = C.
	OID JOIN pg_catalog.pg_description d ON idx.indexrelid = d.objoid 
WHERE
	idx.indrelid = v_table_oid
	LOOP
	v_table_ddl := v_table_ddl || 'COMMENT ON INDEX "' || v_index_comment_record.relname || '" IS ''' || REPLACE ( v_index_comment_record.description, '''', '''''' ) || ''';' || E'\n';

END LOOP;
-- comment on constraint
FOR v_constraint_comment_record IN SELECT
con.conname,
pg_description.description 
FROM
	pg_catalog.pg_constraint con
	JOIN pg_catalog.pg_class rel ON rel.OID = con.conrelid
	JOIN pg_catalog.pg_namespace nsp ON nsp.OID = connamespace
	JOIN pg_catalog.pg_description ON pg_description.objoid = con.OID 
WHERE
	rel.OID = v_table_oid
	LOOP
	v_table_ddl := v_table_ddl || 'COMMENT ON CONSTRAINT "' || v_constraint_comment_record.conname || '" ON "' || in_table_name || '" IS ''' || REPLACE ( v_constraint_comment_record.description, '''', '''''' ) || ''';' || E'\n';

END LOOP;
-- return the ddl
RETURN v_table_ddl;

END $BODY$
  LANGUAGE plpgsql VOLATILE
  COST 100;
ALTER FUNCTION "public"."show_create_table"("in_schema_name" varchar, "in_table_name" varchar) OWNER TO "postgres";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."job_serial"', 1, false);

-- ----------------------------
-- Primary Key structure for table abnormal
-- ----------------------------
ALTER TABLE "public"."abnormal" ADD CONSTRAINT "abnormal_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table devicename
-- ----------------------------
ALTER TABLE "public"."devicename" ADD CONSTRAINT "devicename_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table executors
-- ----------------------------
ALTER TABLE "public"."executors" ADD CONSTRAINT "executors_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table materialalter
-- ----------------------------
ALTER TABLE "public"."materialalter" ADD CONSTRAINT "MaterialAlter_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table plan_comp
-- ----------------------------
ALTER TABLE "public"."plan_comp" ADD CONSTRAINT "plan_comp_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table production
-- ----------------------------
ALTER TABLE "public"."production" ADD CONSTRAINT "develop_pkey" PRIMARY KEY ("序号");

-- ----------------------------
-- Primary Key structure for table prooutput
-- ----------------------------
ALTER TABLE "public"."prooutput" ADD CONSTRAINT "pruoutput_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table quyukongtiao
-- ----------------------------
ALTER TABLE "public"."quyukongtiao" ADD CONSTRAINT "quyukongtiao_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table tb_namespace
-- ----------------------------
ALTER TABLE "public"."tb_namespace" ADD CONSTRAINT "tb_namespace_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table testcn
-- ----------------------------
ALTER TABLE "public"."testcn" ADD CONSTRAINT "testcn_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table testlqbz
-- ----------------------------
ALTER TABLE "public"."testlqbz" ADD CONSTRAINT "testlqbz_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table work_order_detail
-- ----------------------------
ALTER TABLE "public"."work_order_detail" ADD CONSTRAINT "production_copy1_pkey" PRIMARY KEY ("序号");

-- ----------------------------
-- Primary Key structure for table zhaopan
-- ----------------------------
ALTER TABLE "public"."zhaopan" ADD CONSTRAINT "zhaopan_pkey" PRIMARY KEY ("id");
