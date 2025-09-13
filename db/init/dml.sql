-- -----------------------------------------------------
-- Master Data for `items` table
-- -----------------------------------------------------
INSERT INTO `items` (`id`, `name`, `rarity`) VALUES
(1, 'ひのきのぼう', 1),
(2, 'こんぼう', 1),
(3, 'どうのつるぎ', 2),
(4, 'てつのやり', 2),
(5, 'はがねのけん', 3),
(6, 'まじんのオノ', 3);


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
