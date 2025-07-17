delimiter ;;
CREATE DATABASE /*!32312 IF NOT EXISTS*/ `new_db_name` /*!40100 DEFAULT CHARACTER SET utf8*/;

USE `new_db_name`;

SET @queryCount = (
SELECT COUNT(*) FROM information_schema.TABLES
WHERE table_schema IN ('old_db_name'));

DROP PROCEDURE  IF EXISTS test_move;
CREATE PROCEDURE test_move ()
BEGIN

DECLARE i INT DEFAULT 0;
WHILE i < @queryCount 
DO
	SET @queryString := (
	SELECT CONCAT('RENAME TABLE ',table_schema,'.',table_name, ' TO ','new_db_name.',table_name,';')
		FROM information_schema.TABLES WHERE table_schema IN ('old_db_name') ORDER BY table_name LIMIT 1 );
	PREPARE stmt FROM @queryString;
	SELECT @queryString as ``;
	SELECT i as ``;
	EXECUTE stmt;
	SET i=i+1;
END WHILE ;
commit;

END;;
delimiter ;
CALL test_move();

