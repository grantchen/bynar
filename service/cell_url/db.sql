-- phpMyAdmin SQL Dump
-- version 5.1.0
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 16, 2021 at 01:42 PM
-- Server version: 10.4.18-MariaDB
-- PHP Version: 7.4.16

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
use `test`;
--

-- --------------------------------------------------------

--
-- Table structure for table `inventory_managment`
--

CREATE TABLE `inventory_managment` (
  `id` int(11) NOT NULL,
  `product_uuid` varchar(50) NOT NULL,
  `product_description` varchar(150) NOT NULL,
  `product_barcode` varchar(100) DEFAULT NULL,
  `warehouse_uuid` varchar(50) NOT NULL,
  `warehouse_description` varchar(150) NOT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `avg_cost` decimal(25,4) NOT NULL
) ENGINE=InnoDB AVG_ROW_LENGTH=264 DEFAULT CHARSET=utf8;

--
-- Dumping data for table `inventory_managment`
--

INSERT INTO `inventory_managment` VALUES (96,'85f46a74-8291-4a17-9a9d-ddc3e53e85f1','Orange','87654321','2a3bc115-9c3c-4310-b708-99c04cda71d3','Warehouse 1',10.0000,2.0000),(97,'72b7e0da-5d34-4985-8952-030f8dd9ae31','Tomatoe','11223344','2a3bc115-9c3c-4310-b708-99c04cda71d3','Warehouse 1',5.0000,0.5000),(98,'85f46a74-8291-4a17-9a9d-ddc3e53e85f1','Orange','87654321','91bcb8b2-05cb-4f15-83c9-c79883c920f4','Warehouse 2',17.0000,1.2000),(99,'8beef1a2-5fa0-4afe-baf7-10282b33fcd1','Apple','12345678','78d3ac51-95b9-476a-bca1-f757706158a2 ','Warehouse 3',9.0000,0.7000),(100,'8beef1a2-5fa0-4afe-baf7-10282b33fcd1','Apple2','12345678','78d3ac51-95b9-476a-bca1-f757706158a2 ','Warehouse 3',9.0000,0.7000);

-- --------------------------------------------------------

--
-- Table structure for table `products`
--

CREATE TABLE `products` (
  `id` int(11) NOT NULL,
  `uuid` varchar(50) NOT NULL,
  `barcode` varchar(50) NOT NULL,
  `name` varchar(155) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `products`
--

INSERT INTO `products` (`id`, `uuid`, `barcode`, `name`) VALUES
(1, '8beef1a2-5fa0-4afe-baf7-10282b33fcd1', '12345678', 'Apple'),
(2, '85f46a74-8291-4a17-9a9d-ddc3e53e85f1', '87654321', 'Orange'),
(3, '72b7e0da-5d34-4985-8952-030f8dd9ae31', '11223344', 'Tomatoe');

-- --------------------------------------------------------

--
-- Table structure for table `warehouses`
--

CREATE TABLE `warehouses` (
  `id` int(11) NOT NULL,
  `uuid` varchar(50) NOT NULL,
  `name` varchar(155) NOT NULL,
  `costing_methods` varchar(20) NOT NULL
) ENGINE=InnoDB AVG_ROW_LENGTH=5461 DEFAULT CHARSET=utf8;

--
-- Dumping data for table `warehouses`
--

INSERT INTO `warehouses` (`id`, `uuid`, `name`, `costing_methods`) VALUES
(5, '2a3bc115-9c3c-4310-b708-99c04cda71d3', 'Warehouse 111', 'Avco'),
(6, '91bcb8b2-05cb-4f15-83c9-c79883c920f4', 'Warehouse 2', 'AVCO'),
(7, '78d3ac51-95b9-476a-bca1-f757706158a2 ', 'Warehouse 3', 'AVCO'),
(8, '78d3ac51-95b9-476a-bca1-f757706158a1 ', 'Warehouse 4', 'AVCO'),
(9, '78d3ac51-95b9-476a-bca1-f757706158a4', 'Warehouse 5', 'AVCO');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `inventory_managment`
--
ALTER TABLE `inventory_managment`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `products`
--
ALTER TABLE `products`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `warehouses`
--
ALTER TABLE `warehouses`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `inventory_managment`
--
ALTER TABLE `inventory_managment`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=100;

--
-- AUTO_INCREMENT for table `warehouses`
--
ALTER TABLE `warehouses`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
