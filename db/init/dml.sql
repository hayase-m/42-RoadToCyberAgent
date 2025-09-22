USE `game_api_db` ;

-- -----------------------------------------------------
-- Master Data for `items` table
-- -----------------------------------------------------
INSERT INTO `items` (`id`, `name`, `rarity`) VALUES
(1, 'Stick', 1),
(2, 'Club', 1),
(3, 'Sword', 2),
(4, 'Spear', 2),
(5, 'Steel Sword', 3),
(6, 'Axe', 3);


-- -----------------------------------------------------
-- Master Data for `gachas` table
-- -----------------------------------------------------
-- レアリティN (rarity=1) のアイテムは重み3
INSERT INTO `gachas` (`item_id`, `weight`) VALUES
(1, 3),
(2, 3);

-- レアリティR (rarity=2) のアイテムは重み2
INSERT INTO `gachas` (`item_id`, `weight`) VALUES
(3, 2),
(4, 2);

-- レアリティSR (rarity=3) のアイテムは重み1
INSERT INTO `gachas` (`item_id`, `weight`) VALUES
(5, 1),
(6, 1);
