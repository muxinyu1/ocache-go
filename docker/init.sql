CREATE DATABASE IF NOT EXISTS ocache;
USE ocache;

CREATE TABLE IF NOT EXISTS kv_store (
    `group_name` VARCHAR(64) NOT NULL,
    `k` VARCHAR(64) NOT NULL,
    `v` TEXT,
    PRIMARY KEY (`group_name`, `k`)
);

DELIMITER $$
DROP PROCEDURE IF EXISTS generate_data;
CREATE PROCEDURE generate_data()
BEGIN
  DECLARE i INT DEFAULT 0;
  
  -- Scores
  SET i = 0;
  WHILE i < 2500 DO
    INSERT INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', CONCAT('key_', i), CONCAT('val_', i, '_', UUID()));
    SET i = i + 1;
  END WHILE;

  -- Tsinghua
  SET i = 0;
  WHILE i < 2500 DO
    INSERT INTO kv_store (`group_name`, `k`, `v`) VALUES ('Tsinghua', CONCAT('key_', i), CONCAT('val_', i, '_', UUID()));
    SET i = i + 1;
  END WHILE;

    -- ByteDance
  SET i = 0;
  WHILE i < 2500 DO
    INSERT INTO kv_store (`group_name`, `k`, `v`) VALUES ('ByteDance', CONCAT('key_', i), CONCAT('val_', i, '_', UUID()));
    SET i = i + 1;
  END WHILE;

    -- 学生
  SET i = 0;
  WHILE i < 2500 DO
    INSERT INTO kv_store (`group_name`, `k`, `v`) VALUES ('学生', CONCAT('key_', i), CONCAT('val_', i, '_', UUID()));
    SET i = i + 1;
  END WHILE;
  
  REPLACE INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', 'mxy', 'value_mxy');
  REPLACE INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', 'mqs', 'value_mqs');
  REPLACE INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', 'sheeep', 'value_sheeep');
  REPLACE INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', 'random', 'value_random');
  REPLACE INTO kv_store (`group_name`, `k`, `v`) VALUES ('Scores', '主席', 'value_zhu_xi');
  
END$$
DELIMITER ;

CALL generate_data();
