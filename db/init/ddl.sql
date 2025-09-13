-- -----------------------------------------------------
-- Database game_api_db
-- -----------------------------------------------------

CREATE DATABASE IF NOT EXISTS `game_api_db` DEFAULT CHARACTER SET utf8mb4 ;
USE `game_api_db` ;

SET CHARSET utf8mb4;

-- -----------------------------------------------------
-- Table `users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` VARCHAR(30) NOT NULL COMMENT '名前',
  `highscore` INT NOT NULL DEFAULT 0 COMMENT 'ハイスコア',
  `coin` INT NOT NULL DEFAULT 0 COMMENT '所持コイン',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`))
ENGINE = InnoDB
COMMENT = 'ユーザー';

-- -----------------------------------------------------
-- Table `items`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `items` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` VARCHAR(30) NOT NULL COMMENT '名前',
  `rarity` INT NOT NULL COMMENT 'レアリティ(1=N, 2=R, 3=SR)',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`))
ENGINE = InnoDB
COMMENT = 'アイテム';

-- -----------------------------------------------------
-- Table `user_collections`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `user_collections` (
  `user_id` INT NOT NULL COMMENT 'ユーザーID',
  `item_id` INT NOT NULL COMMENT 'アイテムID',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`user_id`, `item_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (`item_id`) REFERENCES `items`(`id`) ON DELETE CASCADE ON UPDATE CASCADE)
ENGINE = InnoDB
COMMENT = 'ユーザーのアイテムコレクション';

-- -----------------------------------------------------
-- Table `gachas`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `gachas` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `item_id` INT NOT NULL COMMENT 'アイテムID',
  `weight` INT NOT NULL COMMENT '重み',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
  PRIMARY KEY (`id`),
  FOREIGN KEY (`item_id`) REFERENCES `items`(`id`) ON DELETE CASCADE ON UPDATE CASCADE)
ENGINE = InnoDB
COMMENT = 'ガチャの中身';
