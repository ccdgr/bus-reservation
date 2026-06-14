-- Ensure the session uses utf8mb4
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Explicitly set the database charset if it exists
USE `bus_reservation`;
ALTER DATABASE `bus_reservation` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(255) NOT NULL UNIQUE,
    `password` VARCHAR(255) NOT NULL,
    `real_name` VARCHAR(255) NOT NULL,
    `user_type` TINYINT NOT NULL DEFAULT 0 COMMENT '0:Student, 1:Staff',
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `orders` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `bus_id` BIGINT UNSIGNED NOT NULL,
    `payment_id` VARCHAR(255) NULL,
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0:Pending, 1:Paid, 2:Cancelled, 3:Expired, 4:Verified, 5:Refunding',
    `created_at` DATETIME(3) NULL,
    `updated_at` DATETIME(3) NULL,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_bus_id` (`bus_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Initial data for testing (Date: 2026-06-15)
-- Route 1: 校区 A -> 校区 B
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-101', '校区 A', '校区 B', '2026-06-15 07:00:00', 50, 50, NOW(), NOW()),
('BUS-102', '校区 A', '校区 B', '2026-06-15 08:30:00', 50, 50, NOW(), NOW()),
('BUS-103', '校区 A', '校区 B', '2026-06-15 10:00:00', 50, 50, NOW(), NOW()),
('BUS-104', '校区 A', '校区 B', '2026-06-15 11:30:00', 50, 50, NOW(), NOW()),
('BUS-105', '校区 A', '校区 B', '2026-06-15 13:30:00', 50, 50, NOW(), NOW()),
('BUS-106', '校区 A', '校区 B', '2026-06-15 15:00:00', 50, 50, NOW(), NOW()),
('BUS-107', '校区 A', '校区 B', '2026-06-15 16:30:00', 50, 50, NOW(), NOW()),
('BUS-108', '校区 A', '校区 B', '2026-06-15 18:00:00', 50, 50, NOW(), NOW()),
('BUS-109', '校区 A', '校区 B', '2026-06-15 20:00:00', 50, 50, NOW(), NOW());

-- Route 2: 校区 B -> 校区 A
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-201', '校区 B', '校区 A', '2026-06-15 07:15:00', 50, 50, NOW(), NOW()),
('BUS-202', '校区 B', '校区 A', '2026-06-15 08:45:00', 50, 50, NOW(), NOW()),
('BUS-203', '校区 B', '校区 A', '2026-06-15 10:15:00', 50, 50, NOW(), NOW()),
('BUS-204', '校区 B', '校区 A', '2026-06-15 11:45:00', 50, 50, NOW(), NOW()),
('BUS-205', '校区 B', '校区 A', '2026-06-15 13:45:00', 50, 50, NOW(), NOW()),
('BUS-206', '校区 B', '校区 A', '2026-06-15 15:15:00', 50, 50, NOW(), NOW()),
('BUS-207', '校区 B', '校区 A', '2026-06-15 16:45:00', 50, 50, NOW(), NOW()),
('BUS-208', '校区 B', '校区 A', '2026-06-15 18:15:00', 50, 50, NOW(), NOW()),
('BUS-209', '校区 B', '校区 A', '2026-06-15 20:15:00', 50, 50, NOW(), NOW());

-- Route 3: 校区 A -> 高铁站
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-301', '校区 A', '高铁站', '2026-06-15 06:30:00', 50, 50, NOW(), NOW()),
('BUS-302', '校区 A', '高铁站', '2026-06-15 08:00:00', 50, 50, NOW(), NOW()),
('BUS-303', '校区 A', '高铁站', '2026-06-15 09:30:00', 50, 50, NOW(), NOW()),
('BUS-304', '校区 A', '高铁站', '2026-06-15 11:00:00', 50, 50, NOW(), NOW()),
('BUS-305', '校区 A', '高铁站', '2026-06-15 13:00:00', 50, 50, NOW(), NOW()),
('BUS-306', '校区 A', '高铁站', '2026-06-15 14:30:00', 50, 50, NOW(), NOW()),
('BUS-307', '校区 A', '高铁站', '2026-06-15 16:00:00', 50, 50, NOW(), NOW()),
('BUS-308', '校区 A', '高铁站', '2026-06-15 17:30:00', 50, 50, NOW(), NOW()),
('BUS-309', '校区 A', '高铁站', '2026-06-15 19:30:00', 50, 50, NOW(), NOW());

-- Route 4: 高铁站 -> 校区 A
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-401', '高铁站', '校区 A', '2026-06-15 07:30:00', 50, 50, NOW(), NOW()),
('BUS-402', '高铁站', '校区 A', '2026-06-15 09:00:00', 50, 50, NOW(), NOW()),
('BUS-403', '高铁站', '校区 A', '2026-06-15 10:30:00', 50, 50, NOW(), NOW()),
('BUS-404', '高铁站', '校区 A', '2026-06-15 12:00:00', 50, 50, NOW(), NOW()),
('BUS-405', '高铁站', '校区 A', '2026-06-15 14:00:00', 50, 50, NOW(), NOW()),
('BUS-406', '高铁站', '校区 A', '2026-06-15 15:30:00', 50, 50, NOW(), NOW()),
('BUS-407', '高铁站', '校区 A', '2026-06-15 17:00:00', 50, 50, NOW(), NOW()),
('BUS-408', '高铁站', '校区 A', '2026-06-15 18:30:00', 50, 50, NOW(), NOW()),
('BUS-409', '高铁站', '校区 A', '2026-06-15 20:30:00', 50, 50, NOW(), NOW());

-- Route 5: 校区 B -> 高铁站
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-501', '校区 B', '高铁站', '2026-06-15 06:45:00', 50, 50, NOW(), NOW()),
('BUS-502', '校区 B', '高铁站', '2026-06-15 08:15:00', 50, 50, NOW(), NOW()),
('BUS-503', '校区 B', '高铁站', '2026-06-15 09:45:00', 50, 50, NOW(), NOW()),
('BUS-504', '校区 B', '高铁站', '2026-06-15 11:15:00', 50, 50, NOW(), NOW()),
('BUS-505', '校区 B', '高铁站', '2026-06-15 13:15:00', 50, 50, NOW(), NOW()),
('BUS-506', '校区 B', '高铁站', '2026-06-15 14:45:00', 50, 50, NOW(), NOW()),
('BUS-507', '校区 B', '高铁站', '2026-06-15 16:15:00', 50, 50, NOW(), NOW()),
('BUS-508', '校区 B', '高铁站', '2026-06-15 17:45:00', 50, 50, NOW(), NOW()),
('BUS-509', '校区 B', '高铁站', '2026-06-15 19:45:00', 50, 50, NOW(), NOW());

-- Route 6: 高铁站 -> 校区 B
INSERT INTO `buses` (`number`, `origin`, `dest`, `start_time`, `total_seat`, `left_seat`, `created_at`, `updated_at`) VALUES 
('BUS-601', '高铁站', '校区 B', '2026-06-15 07:45:00', 50, 50, NOW(), NOW()),
('BUS-602', '高铁站', '校区 B', '2026-06-15 09:15:00', 50, 50, NOW(), NOW()),
('BUS-603', '高铁站', '校区 B', '2026-06-15 10:45:00', 50, 50, NOW(), NOW()),
('BUS-604', '高铁站', '校区 B', '2026-06-15 12:15:00', 50, 50, NOW(), NOW()),
('BUS-605', '高铁站', '校区 B', '2026-06-15 14:15:00', 50, 50, NOW(), NOW()),
('BUS-606', '高铁站', '校区 B', '2026-06-15 15:45:00', 50, 50, NOW(), NOW()),
('BUS-607', '高铁站', '校区 B', '2026-06-15 17:15:00', 50, 50, NOW(), NOW()),
('BUS-608', '高铁站', '校区 B', '2026-06-15 18:45:00', 50, 50, NOW(), NOW()),
('BUS-609', '高铁站', '校区 B', '2026-06-15 20:45:00', 50, 50, NOW(), NOW());
