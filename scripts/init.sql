SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `real_name` VARCHAR(255) NOT NULL,
    `user_type` TINYINT NOT NULL DEFAULT 0 COMMENT '0:Student, 1:Staff',
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `buses` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `number` VARCHAR(255) NOT NULL UNIQUE,
    `origin` VARCHAR(255) NOT NULL,
    `dest` VARCHAR(255) NOT NULL,
    `start_time` DATETIME NOT NULL,
    `total_seat` INT NOT NULL,
    `left_seat` INT NOT NULL,
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL,
    INDEX `idx_start_time` (`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `orders` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `bus_id` BIGINT UNSIGNED NOT NULL,
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0:Pending, 1:Paid, 2:Cancelled',
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_bus_id` (`bus_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Initial data for testing
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) 
VALUES ('BUS-001', '校区 A', '校区 B', '2026-06-15 08:00:00', 50, 50, NOW(), NOW());
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) 
VALUES ('BUS-002', '校区 B', '校区 A', '2026-06-15 17:30:00', 50, 50, NOW(), NOW());
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) 
VALUES ('BUS-003', '校区 A', '高铁站', '2026-06-15 09:00:00', 50, 50, NOW(), NOW());
