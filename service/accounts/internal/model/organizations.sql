CREATE TABLE IF NOT EXISTS `accounting_periods_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `starting_date` datetime NOT NULL,
  `closed` tinyint(4) NOT NULL DEFAULT '0',
  `data_locked` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `areas` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `cash_managements` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) DEFAULT NULL,
  `bank_id` int(11) NOT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `amount` decimal(20,5) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `cash_receipt_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL,
  `applies_document_type` tinyint(4) NOT NULL,
  `applies_document_id` int(11) NOT NULL,
  `amount` decimal(20,5) NOT NULL,
  `amount_lcy` decimal(20,5) NOT NULL,
  `applied` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `cash_receipts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `batch_id` int(11) NOT NULL,
  `document_id` int(11) NOT NULL,
  `document_no` varchar(255) NOT NULL,
  `transaction_no` int(11) NOT NULL,
  `store_id` int(11) NOT NULL,
  `document_date` datetime NOT NULL,
  `posting_date` datetime NOT NULL,
  `entry_date` datetime NOT NULL,
  `account_type` tinyint(4) NOT NULL,
  `account_id` int(11) NOT NULL,
  `balance_account_type` tinyint(4) NOT NULL,
  `balance_account_id` int(11) NOT NULL,
  `amount` decimal(20,5) NOT NULL,
  `amount_lcy` decimal(20,5) NOT NULL,
  `currency_value` decimal(15,5) NOT NULL,
  `user_group_id` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `bank_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `currencies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exchange_rate` decimal(20,5) NOT NULL DEFAULT '0.00000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `departments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(20) NOT NULL DEFAULT '',
  `description` varchar(50) NOT NULL DEFAULT '',
  `address` varchar(20) NOT NULL DEFAULT '',
  `code` varchar(20) NOT NULL DEFAULT '',
  `code_2` varchar(20) NOT NULL DEFAULT '',
  `transaction_code` varchar(20) NOT NULL DEFAULT '',
  `responsibility_center` varchar(20) NOT NULL DEFAULT '',
  `note` varchar(150) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `discounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `value` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `percentage` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `general_ledger_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `local_currency_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `languages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `country` varchar(30) DEFAULT NULL,
  `language` varchar(40) DEFAULT NULL,
  `two_letters` varchar(10) DEFAULT NULL,
  `three_letters` varchar(10) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `location_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `locations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `postal_code` varchar(20) NOT NULL DEFAULT '',
  `city` varchar(30) NOT NULL DEFAULT '',
  `country` varchar(30) NOT NULL DEFAULT '',
  `state` varchar(30) NOT NULL DEFAULT '',
  `continent` varchar(30) NOT NULL DEFAULT '',
  `time_zone` varchar(20) NOT NULL DEFAULT '',
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `organizations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `payment_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL,
  `applies_document_type` tinyint(4) NOT NULL,
  `applies_document_id` int(11) NOT NULL,
  `payment_type_id` int(11) DEFAULT NULL,
  `amount` decimal(20,5) NOT NULL,
  `amount_lcy` decimal(20,5) NOT NULL,
  `applied` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `batch_id` int(11) NOT NULL,
  `document_id` int(11) NOT NULL,
  `document_no` varchar(255) NOT NULL,
  `external_document_no` varchar(255) NOT NULL,
  `transaction_no` int(11) NOT NULL,
  `store_id` int(11) NOT NULL,
  `document_date` datetime NOT NULL,
  `posting_date` datetime NOT NULL,
  `entry_date` datetime NOT NULL,
  `account_type` tinyint(4) NOT NULL,
  `account_id` int(11) NOT NULL,
  `recipient_bank_account_id` int(11) NOT NULL,
  `balance_account_type` tinyint(4) NOT NULL,
  `balance_account_id` int(11) NOT NULL,
  `amount` decimal(20,5) NOT NULL,
  `amount_lcy` decimal(20,5) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `currency_value` decimal(15,5) NOT NULL,
  `user_group_id` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `payment_method_id` int(11) NOT NULL,
  `payment_reference` varchar(255) NOT NULL,
  `creditor_no` varchar(255) NOT NULL,
  `bank_payment_type_id` int(11) NOT NULL,
  `bank_id` int(11) NOT NULL,
  `paid` decimal(20,5) NOT NULL,
  `remaining` decimal(20,5) NOT NULL,
  `paid_status` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `points` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `procurement_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL,
  `item_type` tinyint(4) NOT NULL,
  `item_id` int(11) NOT NULL,
  `location_id` int(11) NOT NULL,
  `input_quantity` decimal(15,5) NOT NULL,
  `item_unit_value` decimal(15,5) NOT NULL,
  `quantity` decimal(15,5) NOT NULL,
  `item_unit_id` int(11) NOT NULL,
  `discount_id` int(11) NOT NULL,
  `tax_area_id` int(11) NOT NULL,
  `vat_id` int(11) NOT NULL,
  `quantity_assign` decimal(15,5) NOT NULL,
  `quantity_assigned` decimal(15,5) NOT NULL,
  `subtotal_exclusive_vat` decimal(20,5) NOT NULL,
  `total_discount` decimal(20,5) NOT NULL,
  `total_exclusive_vat` decimal(20,5) NOT NULL,
  `total_vat` decimal(20,5) NOT NULL,
  `total_inclusive_vat` decimal(20,5) NOT NULL,
  `subtotal_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_discount_lcy` decimal(20,5) NOT NULL,
  `total_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_vat_lcy` decimal(20,5) NOT NULL,
  `total_inclusive_vat_lcy` decimal(20,5) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `procurements` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_id` int(11) NOT NULL,
  `document_no` varchar(255) NOT NULL,
  `transaction_no` int(11) NOT NULL,
  `store_id` int(11) NOT NULL,
  `document_date` datetime NOT NULL,
  `posting_date` datetime NOT NULL,
  `entry_date` datetime NOT NULL,
  `shipment_date` datetime NOT NULL,
  `project_id` int(11) NOT NULL,
  `department_id` int(11) NOT NULL,
  `contract_id` int(11) NOT NULL,
  `user_group_id` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `budget_id` int(11) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `currency_value` decimal(15,5) NOT NULL,
  `vendor_id` int(11) NOT NULL,
  `vendor_invoice_no` varchar(255) NOT NULL,
  `purchaser_id` int(11) NOT NULL,
  `responsibility_center_id` int(11) NOT NULL,
  `payment_terms_id` int(11) NOT NULL,
  `payment_method_id` int(11) NOT NULL,
  `transaction_type_id` int(11) NOT NULL,
  `payment_discount` decimal(20,5) NOT NULL,
  `shipment_method_id` int(11) NOT NULL,
  `payment_reference` int(11) NOT NULL,
  `creditor_no` varchar(255) NOT NULL,
  `on_hold` varchar(255) NOT NULL,
  `transaction_specification_id` int(11) NOT NULL,
  `transport_method_id` int(11) NOT NULL,
  `entry_point_id` int(11) NOT NULL,
  `campaign_id` int(11) NOT NULL,
  `area_id` int(11) NOT NULL,
  `vendor_shipment_no` varchar(255) NOT NULL,
  `subtotal_exclusive_vat` decimal(20,5) NOT NULL,
  `total_discount` decimal(20,5) NOT NULL,
  `total_exclusive_vat` decimal(20,5) NOT NULL,
  `total_vat` decimal(20,5) NOT NULL,
  `total_inclusive_vat` decimal(20,5) NOT NULL,
  `subtotal_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_discount_lcy` decimal(20,5) NOT NULL,
  `total_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_vat_lcy` decimal(20,5) NOT NULL,
  `total_inclusive_vat_lcy` decimal(20,5) NOT NULL,
  `discount_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `projects` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `sale_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL,
  `item_type` tinyint(4) NOT NULL,
  `item_id` int(11) NOT NULL,
  `location_id` int(11) NOT NULL,
  `input_quantity` decimal(15,5) NOT NULL,
  `item_unit_value` decimal(15,5) NOT NULL,
  `quantity` decimal(15,5) NOT NULL,
  `item_unit_id` int(11) NOT NULL,
  `discount_id` int(11) NOT NULL,
  `tax_area_id` int(11) NOT NULL,
  `vat_id` int(11) NOT NULL,
  `quantity_assign` decimal(15,5) NOT NULL,
  `quantity_assigned` decimal(15,5) NOT NULL,
  `subtotal_exclusive_vat` decimal(20,5) NOT NULL,
  `total_discount` decimal(20,5) NOT NULL,
  `total_exclusive_vat` decimal(20,5) NOT NULL,
  `total_vat` decimal(20,5) NOT NULL,
  `total_inclusive_vat` decimal(20,5) NOT NULL,
  `subtotal_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_discount_lcy` decimal(20,5) NOT NULL,
  `total_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_vat_lcy` decimal(20,5) NOT NULL,
  `total_inclusive_vat_lcy` decimal(20,5) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `sales` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_id` int(11) NOT NULL,
  `document_no` varchar(255) NOT NULL,
  `transaction_no` int(11) NOT NULL,
  `store_id` int(11) NOT NULL,
  `document_date` datetime NOT NULL,
  `posting_date` datetime NOT NULL,
  `entry_date` datetime NOT NULL,
  `shipment_date` datetime NOT NULL,
  `project_id` int(11) NOT NULL,
  `department_id` int(11) NOT NULL,
  `contract_id` int(11) NOT NULL,
  `user_group_id` int(11) NOT NULL,
  `status` tinyint(4) NOT NULL,
  `currency_id` int(11) NOT NULL,
  `currency_value` decimal(15,5) NOT NULL,
  `customer_id` int(11) NOT NULL,
  `salesperson_id` int(11) NOT NULL,
  `responsibility_center_id` int(11) NOT NULL,
  `payment_terms_id` int(11) NOT NULL,
  `payment_method_id` int(11) NOT NULL,
  `transaction_type_id` int(11) NOT NULL,
  `payment_discount` decimal(20,5) NOT NULL,
  `shipment_method_id` int(11) NOT NULL,
  `payment_reference` int(11) NOT NULL,
  `transaction_specification_id` int(11) NOT NULL,
  `transport_method_id` int(11) NOT NULL,
  `exit_point_id` int(11) NOT NULL,
  `campaign_id` int(11) NOT NULL,
  `area_id` int(11) NOT NULL,
  `package_tracking_no` varchar(255) NOT NULL,
  `subtotal_exclusive_vat` decimal(20,5) NOT NULL,
  `total_discount` decimal(20,5) NOT NULL,
  `total_exclusive_vat` decimal(20,5) NOT NULL,
  `total_vat` decimal(20,5) NOT NULL,
  `total_inclusive_vat` decimal(20,5) NOT NULL,
  `subtotal_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_discount_lcy` decimal(20,5) NOT NULL,
  `total_exclusive_vat_lcy` decimal(20,5) NOT NULL,
  `total_vat_lcy` decimal(20,5) NOT NULL,
  `total_inclusive_vat_lcy` decimal(20,5) NOT NULL,
  `direct_debit_mandate_id` int(11) NOT NULL,
  `agent_id` int(11) NOT NULL,
  `agent_service_id` int(11) NOT NULL,
  `discount_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `series` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `default_nos` tinyint(4) NOT NULL DEFAULT '0',
  `manual_nos` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `service_zones` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `shipment_methods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `shipping_agent_services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `shipping_time` int(11) NOT NULL DEFAULT '0',
  `base_calendar_uuid` int(11) NOT NULL DEFAULT '0',
  `customized_calendar` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `shipping_agents` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `tracking_address` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `transaction_specifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `transaction_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_code` varchar(5) NOT NULL DEFAULT '',
  `operation_code` varchar(5) NOT NULL DEFAULT '',
  `document_code` varchar(5) NOT NULL DEFAULT '',
  `transaction_code` varchar(15) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `units` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `base_unit` varchar(255) NOT NULL DEFAULT '',
  `operation_value` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `unit_value` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `transaction_code` varchar(255) NOT NULL DEFAULT '',
  `responsibility_center_uuid` int(11) NOT NULL DEFAULT '0',
  `value` decimal(15,5) NOT NULL DEFAULT '0.00000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `user_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `vats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `value` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `percentage` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `documents` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_type` varchar(255) NOT NULL DEFAULT '',
  `document_abbrevation` varchar(255) NOT NULL DEFAULT '',
  `workspace` varchar(255) NOT NULL DEFAULT '',
  `series_id` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `documents_ibfk_1` (`series_id`),
  KEY `documents_series_id` (`series_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `inventories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `location_id` int(11) NOT NULL DEFAULT '0',
  `item_id` int(11) NOT NULL DEFAULT '0',
  `quantity` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `value` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `value_fifo` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `value_lifo` decimal(15,5) DEFAULT '0.00000',
  PRIMARY KEY (`id`),
  KEY `inventories_ibfk_1` (`location_id`),
  KEY `inventories_item_id` (`item_id`),
  KEY `inventories_location_id` (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `location_group_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `location_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `location_group_items_ibfk_1` (`parent_id`),
  KEY `location_group_items_ibfk_2` (`location_id`),
  KEY `location_group_items_location_id` (`location_id`),
  KEY `location_group_items_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `outbound_flows` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `module_id` int(11) NOT NULL DEFAULT '0',
  `location_id` int(11) NOT NULL DEFAULT '0',
  `item_id` int(11) NOT NULL DEFAULT '0',
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `transaction_id` bigint(20) NOT NULL DEFAULT '0',
  `posting_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `quantity` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `value_avco` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `value_fifo` decimal(20,5) NOT NULL DEFAULT '0.00000',
  `value_lifo` decimal(15,5) DEFAULT '0.00000',
  PRIMARY KEY (`id`),
  KEY `outbound_flows_ibfk_1` (`location_id`),
  KEY `outbound_flows_ibfk_2` (`item_id`),
  KEY `outbound_flows_item_id` (`item_id`),
  KEY `outbound_flows_location_id` (`location_id`),
  KEY `outbound_flows_module_id` (`module_id`),
  KEY `outbound_flows_parent_id` (`parent_id`),
  KEY `outbound_flows_transaction_id` (`transaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `series_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `starting_no` varchar(255) NOT NULL DEFAULT '0',
  `increment_no` int(11) NOT NULL DEFAULT '0',
  `last_date_used` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_no_used` varchar(255) NOT NULL DEFAULT '',
  `warning_no` varchar(255) NOT NULL DEFAULT '',
  `ending_no` varchar(255) NOT NULL DEFAULT '',
  `open` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `series_items_ibfk_1` (`parent_id`),
  KEY `series_items_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `stores` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `transaction_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `warehouse_uuid` int(11) NOT NULL DEFAULT '0',
  `catalogue_uuid` int(11) NOT NULL DEFAULT '0',
  `address_uuid` int(11) NOT NULL DEFAULT '0',
  `contact_uuid` int(11) NOT NULL DEFAULT '0',
  `responsibility_center_uuid` int(11) NOT NULL DEFAULT '0',
  `location_group_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `stores_ibfk_1` (`location_group_id`),
  KEY `stores_location_group_id` (`location_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `transfers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_id` int(11) NOT NULL DEFAULT '0',
  `document_no` varchar(255) NOT NULL DEFAULT '',
  `transaction_no` bigint(20) NOT NULL DEFAULT '0',
  `store_id` int(11) NOT NULL DEFAULT '0',
  `document_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `posting_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `entry_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `shipment_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `project_id` int(11) NOT NULL DEFAULT '0',
  `department_id` int(11) NOT NULL DEFAULT '0',
  `in_transit_id` int(11) NOT NULL DEFAULT '0',
  `shipment_method_id` int(11) NOT NULL DEFAULT '0',
  `shipping_agent_id` int(11) NOT NULL DEFAULT '0',
  `shipping_agent_service_id` int(11) NOT NULL DEFAULT '0',
  `transaction_type_id` int(11) NOT NULL DEFAULT '0',
  `transaction_specification_id` int(11) NOT NULL DEFAULT '0',
  `area_id` int(11) NOT NULL DEFAULT '0',
  `entry_exit_point_id` int(11) NOT NULL DEFAULT '0',
  `user_group_id` int(11) NOT NULL DEFAULT '0',
  `location_origin_id` int(11) NOT NULL DEFAULT '0',
  `location_destination_id` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `transfers_ibfk_1` (`document_id`),
  KEY `transfers_ibfk_10` (`transaction_specification_id`),
  KEY `transfers_ibfk_11` (`area_id`),
  KEY `transfers_ibfk_12` (`entry_exit_point_id`),
  KEY `transfers_ibfk_13` (`user_group_id`),
  KEY `transfers_ibfk_14` (`location_origin_id`),
  KEY `transfers_ibfk_15` (`location_destination_id`),
  KEY `transfers_ibfk_2` (`store_id`),
  KEY `transfers_ibfk_3` (`project_id`),
  KEY `transfers_ibfk_4` (`department_id`),
  KEY `transfers_ibfk_5` (`in_transit_id`),
  KEY `transfers_ibfk_6` (`shipment_method_id`),
  KEY `transfers_ibfk_7` (`shipping_agent_id`),
  KEY `transfers_ibfk_8` (`shipping_agent_service_id`),
  KEY `transfers_ibfk_9` (`transaction_type_id`),
  KEY `transfers_area_id` (`area_id`),
  KEY `transfers_department_id` (`department_id`),
  KEY `transfers_document_id` (`document_id`),
  KEY `transfers_entry_exit_point_id` (`entry_exit_point_id`),
  KEY `transfers_in_transit_id` (`in_transit_id`),
  KEY `transfers_location_destination_id` (`location_destination_id`),
  KEY `transfers_location_origin_id` (`location_origin_id`),
  KEY `transfers_project_id` (`project_id`),
  KEY `transfers_shipment_method_id` (`shipment_method_id`),
  KEY `transfers_shipping_agent_id` (`shipping_agent_id`),
  KEY `transfers_shipping_agent_service_id` (`shipping_agent_service_id`),
  KEY `transfers_store_id` (`store_id`),
  KEY `transfers_transaction_specification_id` (`transaction_specification_id`),
  KEY `transfers_transaction_type_id` (`transaction_type_id`),
  KEY `transfers_user_group_id` (`user_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `workflows` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  `module_id` int(11) NOT NULL DEFAULT '0',
  `store_id` int(11) NOT NULL DEFAULT '0',
  `user_group_id` int(11) NOT NULL DEFAULT '0',
  `project_group_id` int(11) DEFAULT '0',
  `department_group_id` int(11) DEFAULT '0',
  `contract_group_id` int(11) DEFAULT '0',
  `currency_group_id` int(11) DEFAULT '0',
  `account_group_id` int(11) DEFAULT '0',
  `campaign_group_id` int(11) DEFAULT '0',
  `payment_method_group_id` int(11) DEFAULT '0',
  `area_group_id` int(11) DEFAULT '0',
  `point_group_id` int(11) DEFAULT '0',
  `agent_group_id` int(11) DEFAULT '0',
  `agent_service_group_id` int(11) DEFAULT '0',
  `payment_terms_group_id` int(11) DEFAULT '0',
  `customer_group_id` int(11) DEFAULT '0',
  `sales_person_group_id` int(11) DEFAULT '0',
  `responsibility_center_group_id` int(11) DEFAULT '0',
  `transaction_group_id` int(11) DEFAULT '0',
  `discount_availability` tinyint(4) DEFAULT '0',
  `vat_availability` tinyint(4) DEFAULT '0',
  `less_equal_total_inclusive_vat_lcy` decimal(20,5) DEFAULT '0.00000',
  `greater_equal_total_inclusive_vat_lcy` decimal(20,5) DEFAULT '0.00000',
  `total_inclusive_vat_lcy` varchar(255) DEFAULT '0',
  `line_item_group_id` int(11) DEFAULT '0',
  `line_item_type_group_id` int(11) DEFAULT '0',
  `line_item_quantity` varchar(255) DEFAULT '0',
  `budget_group_id` int(11) DEFAULT '0',
  `bank_group_id` int(11) DEFAULT '0',
  `amount_lcy` varchar(255) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `workflows_ibfk_1` (`store_id`),
  KEY `workflows_ibfk_2` (`user_group_id`),
  KEY `workflows_store_id` (`store_id`),
  KEY `workflows_user_group_id` (`user_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `transfer_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `item_id` int(11) NOT NULL DEFAULT '0',
  `input_quantity` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `item_unit_value` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `quantity` decimal(15,5) NOT NULL DEFAULT '0.00000',
  `item_unit_id` int(11) NOT NULL DEFAULT '0',
  `shipment_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `receipt_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `transfer_lines_ibfk_1` (`parent_id`),
  KEY `transfer_lines_ibfk_2` (`item_id`),
  KEY `transfer_lines_ibfk_3` (`item_unit_id`),
  KEY `transfer_lines_item_id` (`item_id`),
  KEY `transfer_lines_item_unit_id` (`item_unit_id`),
  KEY `transfer_lines_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `workflow_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `account_id` bigint(20) NOT NULL DEFAULT '0',
  `document_id` int(11) NOT NULL DEFAULT '0',
  `approval_order` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `workflow_items_ibfk_1` (`parent_id`),
  KEY `workflow_items_ibfk_2` (`document_id`),
  KEY `workflow_items_document_id` (`document_id`),
  KEY `workflow_items_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

ALTER TABLE `documents` ADD CONSTRAINT `documents_ibfk_1` FOREIGN KEY (`series_id`) REFERENCES `series` (`id`) ON UPDATE CASCADE;
ALTER TABLE `inventories` ADD CONSTRAINT `inventories_ibfk_1` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON UPDATE CASCADE;
ALTER TABLE `location_group_items` ADD CONSTRAINT `location_group_items_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `location_groups` (`id`) ON UPDATE CASCADE;
ALTER TABLE `location_group_items` ADD CONSTRAINT `location_group_items_ibfk_2` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON UPDATE CASCADE;
ALTER TABLE `outbound_flows` ADD CONSTRAINT `outbound_flows_ibfk_1` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON UPDATE CASCADE;
ALTER TABLE `outbound_flows` ADD CONSTRAINT `outbound_flows_ibfk_2` FOREIGN KEY (`item_id`) REFERENCES `inventories` (`item_id`) ON UPDATE CASCADE;
ALTER TABLE `series_items` ADD CONSTRAINT `series_items_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `series` (`id`) ON UPDATE CASCADE;
ALTER TABLE `stores` ADD CONSTRAINT `stores_ibfk_1` FOREIGN KEY (`location_group_id`) REFERENCES `location_groups` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_1` FOREIGN KEY (`document_id`) REFERENCES `documents` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_10` FOREIGN KEY (`transaction_specification_id`) REFERENCES `transaction_specifications` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_11` FOREIGN KEY (`area_id`) REFERENCES `areas` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_12` FOREIGN KEY (`entry_exit_point_id`) REFERENCES `points` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_13` FOREIGN KEY (`user_group_id`) REFERENCES `user_groups` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_14` FOREIGN KEY (`location_origin_id`) REFERENCES `locations` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_15` FOREIGN KEY (`location_destination_id`) REFERENCES `locations` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_2` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_3` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_4` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_5` FOREIGN KEY (`in_transit_id`) REFERENCES `transactions` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_6` FOREIGN KEY (`shipment_method_id`) REFERENCES `shipment_methods` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_7` FOREIGN KEY (`shipping_agent_id`) REFERENCES `shipping_agents` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_8` FOREIGN KEY (`shipping_agent_service_id`) REFERENCES `shipping_agent_services` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfers` ADD CONSTRAINT `transfers_ibfk_9` FOREIGN KEY (`transaction_type_id`) REFERENCES `transaction_types` (`id`) ON UPDATE CASCADE;
ALTER TABLE `workflows` ADD CONSTRAINT `workflows_ibfk_1` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE CASCADE;
ALTER TABLE `workflows` ADD CONSTRAINT `workflows_ibfk_2` FOREIGN KEY (`user_group_id`) REFERENCES `user_groups` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfer_lines` ADD CONSTRAINT `transfer_lines_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `transfers` (`id`) ON UPDATE CASCADE;
ALTER TABLE `transfer_lines` ADD CONSTRAINT `transfer_lines_ibfk_2` FOREIGN KEY (`item_id`) REFERENCES `inventories` (`item_id`) ON UPDATE CASCADE;
ALTER TABLE `transfer_lines` ADD CONSTRAINT `transfer_lines_ibfk_3` FOREIGN KEY (`item_unit_id`) REFERENCES `units` (`id`) ON UPDATE CASCADE;
ALTER TABLE `workflow_items` ADD CONSTRAINT `workflow_items_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `workflows` (`id`) ON UPDATE CASCADE;
ALTER TABLE `workflow_items` ADD CONSTRAINT `workflow_items_ibfk_2` FOREIGN KEY (`document_id`) REFERENCES `documents` (`id`) ON UPDATE CASCADE;
