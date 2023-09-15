-- bynar.absences definition

CREATE TABLE `absences` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `employee_code` int(11) DEFAULT NULL,
  `from_date` datetime DEFAULT NULL,
  `to_date` datetime DEFAULT NULL,
  `absence_cause_uuid` int(11) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `quantity` decimal(25,4) DEFAULT NULL,
  `unit_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.account_category definition

CREATE TABLE `account_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.account_groups definition

CREATE TABLE `account_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.account_subcategory definition

CREATE TABLE `account_subcategory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `base_category_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.account_users definition

CREATE TABLE `account_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `organisation_uuid` bigint(20) DEFAULT NULL,
  `account_uuid` bigint(20) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `user_group` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.accounting_periods definition

CREATE TABLE `accounting_periods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `starting_date` datetime DEFAULT NULL,
  `no_periods` int(11) DEFAULT NULL,
  `period_length` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.accounting_periods_items definition

CREATE TABLE `accounting_periods_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `starting_date` datetime DEFAULT NULL,
  `new_fiscal_year` tinyint(4) DEFAULT NULL,
  `closed` tinyint(4) DEFAULT NULL,
  `data_locked` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.accountings definition

CREATE TABLE `accountings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `Parent` varchar(33) COLLATE utf8_unicode_ci NOT NULL,
  `Def` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse_code` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company_vat_no` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_groups` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `currency` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `currency_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `currency_rate` decimal(25,4) DEFAULT NULL,
  `debit` decimal(25,4) DEFAULT NULL,
  `credit` decimal(25,4) DEFAULT NULL,
  `base_currency` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `base_currency_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accountant` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accounting_department` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accountant_approve` tinyint(1) DEFAULT '1',
  `note` varchar(250) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `warehouse_uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `base_currency_value` decimal(25,4) DEFAULT NULL,
  `base_currency_debit` decimal(25,4) DEFAULT NULL,
  `base_currency_credit` decimal(25,4) DEFAULT NULL,
  `Date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `items_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `items_id_uindex` (`id`),
  KEY `recipes_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.accounts definition

CREATE TABLE `accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `address_2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `postal_code` int(11) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `cognito_user_groups` enum('PrimaryOwner','Owner','Administrator','Users') DEFAULT 'PrimaryOwner',
  `organization_id` int(11) DEFAULT NULL,
  `organization_account` tinyint(4) DEFAULT NULL,
  `account_confirm_status` tinyint(1) DEFAULT '0',
  `language_preference` varchar(20) DEFAULT 'english',
  `management` tinyint(4) DEFAULT '0',
  `administrator` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.`accounts-old` definition

CREATE TABLE `accounts-old` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.accounts_cards definition

CREATE TABLE `accounts_cards` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_payment_gateway_id` text NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `is_default` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.accounts_management definition

CREATE TABLE `accounts_management` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account_id` int(11) DEFAULT NULL,
  `organisation_id` int(11) DEFAULT NULL,
  `user_group` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.activities definition

CREATE TABLE `activities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `activities_code_uindex` (`code`),
  UNIQUE KEY `activities_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `activities_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.activities_items definition

CREATE TABLE `activities_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `priority` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date_formula` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `activities_items_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `activities_items_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.addresses definition

CREATE TABLE `addresses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `no` varchar(255) DEFAULT NULL,
  `abbrevation` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `trading` varchar(255) DEFAULT NULL,
  `billing` varchar(255) DEFAULT NULL,
  `shipping` varchar(255) DEFAULT NULL,
  `sail` varchar(255) DEFAULT NULL,
  `registered_office` varchar(255) DEFAULT NULL,
  `postal_code` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `e_contact` varchar(255) DEFAULT NULL,
  `e_billing` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  `responsibility_center` varchar(255) DEFAULT NULL,
  `homepage` varchar(255) DEFAULT NULL,
  `company_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.approval_groups definition

CREATE TABLE `approval_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  `module_id` int(11) DEFAULT NULL,
  `active` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.approval_items definition

CREATE TABLE `approval_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `document_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `approval_order` tinyint(4) DEFAULT NULL,
  `approval_status` varchar(255) DEFAULT NULL,
  `request_change` tinyint(4) DEFAULT NULL,
  `request_delete` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.approvals definition

CREATE TABLE `approvals` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.assets definition

CREATE TABLE `assets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `class_uuid` int(11) DEFAULT NULL,
  `subclass_uuid` int(11) DEFAULT NULL,
  `depreciation_method` tinyint(4) DEFAULT NULL,
  `vendor_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `maintence_vendor_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `under_maintenance` tinyint(1) DEFAULT NULL,
  `next_service_date` datetime DEFAULT NULL,
  `warranty_date` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `insured` datetime DEFAULT NULL,
  `depreciations_documents` int(11) DEFAULT NULL,
  `depreciations_documents_total` decimal(25,4) DEFAULT NULL,
  `posting_group_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `transaction_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `budgeted_asset` tinyint(4) DEFAULT NULL,
  `serial_no` int(11) DEFAULT NULL,
  `asset_component` tinyint(4) DEFAULT NULL,
  `component_uuid` int(11) DEFAULT NULL,
  `acquired` tinyint(4) DEFAULT NULL,
  `last_date_modified` datetime DEFAULT NULL,
  `depreciation_book_uuid` int(11) DEFAULT NULL,
  `depreciation_starting` datetime DEFAULT NULL,
  `depreciation_ending` datetime DEFAULT NULL,
  `depreciation_no` decimal(25,4) DEFAULT NULL,
  `book_value` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `assets_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.assets_entry definition

CREATE TABLE `assets_entry` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `Parent` varchar(33) COLLATE utf8_unicode_ci NOT NULL,
  `Def` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(250) COLLATE utf8_unicode_ci DEFAULT NULL,
  `uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouseman` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouseman_approve` tinyint(1) DEFAULT NULL,
  `warehouseman_department` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse_uuid` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse_code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `subtransaction` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `store_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `assets_entry_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `assets_entry_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.assets_entry_items definition

CREATE TABLE `assets_entry_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` varchar(30) COLLATE utf8_unicode_ci NOT NULL,
  `Parent` varchar(33) COLLATE utf8_unicode_ci NOT NULL,
  `Def` varchar(11) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_type` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_unit_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_unit_discount_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_unit_tax_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_unit_tax` decimal(25,4) DEFAULT NULL,
  `item_value` decimal(25,4) DEFAULT NULL,
  `input_quantity` decimal(10,4) DEFAULT NULL,
  `item_quantity_unit` decimal(10,4) DEFAULT NULL,
  `item_quantity` decimal(10,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_tempory` int(11) NOT NULL,
  `item_uuid` varchar(40) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `assets_entry_items_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `assets_entry_items_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.auxiliary_inboud definition

CREATE TABLE `auxiliary_inboud` (
  `id` int(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `document_no` varchar(20) DEFAULT NULL,
  `document_date` varchar(20) DEFAULT NULL,
  `purchase_document_no` varchar(20) DEFAULT NULL,
  `responsability_center` varchar(20) DEFAULT NULL,
  `document_uuid` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.auxiliary_inboud_items definition

CREATE TABLE `auxiliary_inboud_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) DEFAULT NULL,
  `item_document_no` varchar(20) DEFAULT NULL,
  `item_code` varchar(50) DEFAULT NULL,
  `item_code_2` varchar(50) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.auxiliary_outbound definition

CREATE TABLE `auxiliary_outbound` (
  `id` int(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `document_no` varchar(20) DEFAULT NULL,
  `document_date` varchar(20) DEFAULT NULL,
  `sales_document_no` varchar(20) DEFAULT NULL,
  `responsability_center` varchar(20) DEFAULT NULL,
  `document_uuid` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.auxiliary_outbound_items definition

CREATE TABLE `auxiliary_outbound_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) DEFAULT NULL,
  `item_document_no` varchar(20) DEFAULT NULL,
  `item_code` varchar(50) DEFAULT NULL,
  `item_code_2` varchar(50) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.bank_account_posting_groups definition

CREATE TABLE `bank_account_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `gl_bank_account_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.bank_clearing_standards definition

CREATE TABLE `bank_clearing_standards` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.banks definition

CREATE TABLE `banks` (
  `id` int(11) NOT NULL,
  `type` varchar(20) NOT NULL,
  `description` varchar(30) DEFAULT NULL,
  `currency` varchar(20) DEFAULT NULL,
  `number` varchar(20) DEFAULT NULL,
  `swift` varchar(20) DEFAULT NULL,
  `iban` varchar(20) DEFAULT NULL,
  `address` varchar(150) DEFAULT NULL,
  `contact` varchar(20) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `code_2` varchar(50) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `balance` decimal(25,4) DEFAULT NULL,
  `currency_code` varchar(20) DEFAULT NULL,
  `responsibility_center` varchar(20) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `no` varchar(20) DEFAULT NULL,
  `subsidiaries` varchar(20) DEFAULT NULL,
  `uuid` varchar(40) DEFAULT NULL,
  `company` varchar(20) DEFAULT NULL,
  `abbrevation` varchar(20) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `banks_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.base_calendars definition

CREATE TABLE `base_calendars` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.billing_information definition

CREATE TABLE `billing_information` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account_id` varchar(255) DEFAULT NULL,
  `company_id` varchar(255) DEFAULT NULL,
  `account_name` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `address_2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `postal_code` int(11) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `organizaton_name` varchar(255) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `business_account` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.bins definition

CREATE TABLE `bins` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(20) DEFAULT NULL,
  `code_2` varchar(50) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `section_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `bins_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `bins_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.brands definition

CREATE TABLE `brands` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.budgets definition

CREATE TABLE `budgets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_type` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.business_relation_type definition

CREATE TABLE `business_relation_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `type` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `table_name_code_uindex` (`code`),
  UNIQUE KEY `table_name_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `table_name_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.bynar_addresses definition

CREATE TABLE `bynar_addresses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `company_id` int(11) NOT NULL,
  `line1` varchar(50) NOT NULL,
  `line2` varchar(50) DEFAULT NULL,
  `city` varchar(25) NOT NULL,
  `postal_code` varchar(20) DEFAULT NULL,
  `state` varchar(25) NOT NULL,
  `country` varchar(50) NOT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `company_id` (`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_adjustment_items definition

CREATE TABLE `bynar_adjustment_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `adjustment_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `option_id` int(11) DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `serial_no` varchar(255) DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `adjustment_id` (`adjustment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_adjustments definition

CREATE TABLE `bynar_adjustments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `reference_no` varchar(55) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `note` text DEFAULT NULL,
  `attachment` varchar(55) DEFAULT NULL,
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `count_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_brands definition

CREATE TABLE `bynar_brands` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `image` varchar(50) DEFAULT NULL,
  `slug` varchar(55) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_calendar definition

CREATE TABLE `bynar_calendar` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(55) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `start` datetime NOT NULL,
  `end` datetime DEFAULT NULL,
  `color` varchar(7) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_captcha definition

CREATE TABLE `bynar_captcha` (
  `captcha_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `captcha_time` int(10) unsigned NOT NULL,
  `ip_address` varchar(16) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL DEFAULT '0',
  `word` varchar(20) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL,
  PRIMARY KEY (`captcha_id`) /*T![clustered_index] CLUSTERED */,
  KEY `word` (`word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_categories definition

CREATE TABLE `bynar_categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(55) NOT NULL,
  `name` varchar(55) NOT NULL,
  `image` varchar(55) DEFAULT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `slug` varchar(55) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_combo_items definition

CREATE TABLE `bynar_combo_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `item_code` varchar(20) NOT NULL,
  `quantity` decimal(12,4) NOT NULL,
  `unit_price` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_companies definition

CREATE TABLE `bynar_companies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(10) unsigned DEFAULT NULL,
  `group_name` varchar(20) NOT NULL,
  `customer_group_id` int(11) DEFAULT NULL,
  `customer_group_name` varchar(100) DEFAULT NULL,
  `name` varchar(55) NOT NULL,
  `company` varchar(255) NOT NULL,
  `vat_no` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `city` varchar(55) DEFAULT NULL,
  `state` varchar(55) DEFAULT NULL,
  `postal_code` varchar(8) DEFAULT NULL,
  `country` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `cf1` varchar(100) DEFAULT NULL,
  `cf2` varchar(100) DEFAULT NULL,
  `cf3` varchar(100) DEFAULT NULL,
  `cf4` varchar(100) DEFAULT NULL,
  `cf5` varchar(100) DEFAULT NULL,
  `cf6` varchar(100) DEFAULT NULL,
  `invoice_footer` text DEFAULT NULL,
  `payment_term` int(11) DEFAULT '0',
  `logo` varchar(255) DEFAULT 'logo.png',
  `award_points` int(11) DEFAULT '0',
  `deposit_amount` decimal(25,4) DEFAULT NULL,
  `price_group_id` int(11) DEFAULT NULL,
  `price_group_name` varchar(50) DEFAULT NULL,
  `gst_no` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `group_id` (`group_id`),
  KEY `group_id_2` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_costing definition

CREATE TABLE `bynar_costing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `product_id` int(11) DEFAULT NULL,
  `sale_item_id` int(11) NOT NULL,
  `sale_id` int(11) DEFAULT NULL,
  `purchase_item_id` int(11) DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `purchase_net_unit_cost` decimal(25,4) DEFAULT NULL,
  `purchase_unit_cost` decimal(25,4) DEFAULT NULL,
  `sale_net_unit_price` decimal(25,4) NOT NULL,
  `sale_unit_price` decimal(25,4) NOT NULL,
  `quantity_balance` decimal(15,4) DEFAULT NULL,
  `inventory` tinyint(1) DEFAULT '0',
  `overselling` tinyint(1) DEFAULT '0',
  `option_id` int(11) DEFAULT NULL,
  `purchase_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_currencies definition

CREATE TABLE `bynar_currencies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(5) NOT NULL,
  `name` varchar(55) NOT NULL,
  `rate` decimal(12,4) NOT NULL,
  `auto_update` tinyint(1) NOT NULL DEFAULT '0',
  `symbol` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_customer_groups definition

CREATE TABLE `bynar_customer_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `percent` int(11) NOT NULL,
  `discount` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_date_format definition

CREATE TABLE `bynar_date_format` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `js` varchar(20) NOT NULL,
  `php` varchar(20) NOT NULL,
  `sql` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_deliveries definition

CREATE TABLE `bynar_deliveries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sale_id` int(11) NOT NULL,
  `do_reference_no` varchar(50) NOT NULL,
  `sale_reference_no` varchar(50) NOT NULL,
  `customer` varchar(55) NOT NULL,
  `address` varchar(1000) NOT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `status` varchar(15) DEFAULT NULL,
  `attachment` varchar(50) DEFAULT NULL,
  `delivered_by` varchar(50) DEFAULT NULL,
  `received_by` varchar(50) DEFAULT NULL,
  `created_by` int(11) DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_deposits definition

CREATE TABLE `bynar_deposits` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `company_id` int(11) NOT NULL,
  `amount` decimal(25,4) NOT NULL,
  `paid_by` varchar(50) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_expense_categories definition

CREATE TABLE `bynar_expense_categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(55) NOT NULL,
  `name` varchar(55) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_expenses definition

CREATE TABLE `bynar_expenses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reference` varchar(50) NOT NULL,
  `amount` decimal(25,4) NOT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `created_by` varchar(55) NOT NULL,
  `attachment` varchar(55) DEFAULT NULL,
  `category_id` int(11) DEFAULT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_gift_card_topups definition

CREATE TABLE `bynar_gift_card_topups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `card_id` int(11) NOT NULL,
  `amount` decimal(15,4) NOT NULL,
  `created_by` int(11) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `card_id` (`card_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_gift_cards definition

CREATE TABLE `bynar_gift_cards` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `card_no` varchar(20) NOT NULL,
  `value` decimal(25,4) NOT NULL,
  `customer_id` int(11) DEFAULT NULL,
  `customer` varchar(255) DEFAULT NULL,
  `balance` decimal(25,4) NOT NULL,
  `expiry` date DEFAULT NULL,
  `created_by` varchar(55) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `card_no` (`card_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_groups definition

CREATE TABLE `bynar_groups` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `description` varchar(100) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_login_attempts definition

CREATE TABLE `bynar_login_attempts` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `ip_address` varbinary(16) NOT NULL,
  `login` varchar(100) NOT NULL,
  `time` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_logs definition

CREATE TABLE `bynar_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `detail` varchar(190) NOT NULL,
  `model` longtext DEFAULT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_migrations definition

CREATE TABLE `bynar_migrations` (
  `version` bigint(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_notifications definition

CREATE TABLE `bynar_notifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `comment` text NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `from_date` datetime DEFAULT NULL,
  `till_date` datetime DEFAULT NULL,
  `scope` tinyint(1) NOT NULL DEFAULT '3',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_order_ref definition

CREATE TABLE `bynar_order_ref` (
  `ref_id` int(11) NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `so` int(11) NOT NULL DEFAULT '1',
  `qu` int(11) NOT NULL DEFAULT '1',
  `po` int(11) NOT NULL DEFAULT '1',
  `to` int(11) NOT NULL DEFAULT '1',
  `pos` int(11) NOT NULL DEFAULT '1',
  `do` int(11) NOT NULL DEFAULT '1',
  `pay` int(11) NOT NULL DEFAULT '1',
  `re` int(11) NOT NULL DEFAULT '1',
  `rep` int(11) NOT NULL DEFAULT '1',
  `ex` int(11) NOT NULL DEFAULT '1',
  `ppay` int(11) NOT NULL DEFAULT '1',
  `qa` int(11) DEFAULT '1',
  PRIMARY KEY (`ref_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_payments definition

CREATE TABLE `bynar_payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp DEFAULT CURRENT_TIMESTAMP,
  `sale_id` int(11) DEFAULT NULL,
  `return_id` int(11) DEFAULT NULL,
  `purchase_id` int(11) DEFAULT NULL,
  `reference_no` varchar(50) NOT NULL,
  `transaction_id` varchar(50) DEFAULT NULL,
  `paid_by` varchar(20) NOT NULL,
  `cheque_no` varchar(20) DEFAULT NULL,
  `cc_no` varchar(20) DEFAULT NULL,
  `cc_holder` varchar(25) DEFAULT NULL,
  `cc_month` varchar(2) DEFAULT NULL,
  `cc_year` varchar(4) DEFAULT NULL,
  `cc_type` varchar(20) DEFAULT NULL,
  `amount` decimal(25,4) NOT NULL,
  `currency` varchar(3) DEFAULT NULL,
  `created_by` int(11) NOT NULL,
  `attachment` varchar(55) DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `pos_paid` decimal(25,4) DEFAULT '0.0000',
  `pos_balance` decimal(25,4) DEFAULT '0.0000',
  `approval_code` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_paypal_old definition

CREATE TABLE `bynar_paypal_old` (
  `id` int(11) NOT NULL,
  `active` tinyint(4) NOT NULL,
  `account_email` varchar(255) NOT NULL,
  `paypal_currency` varchar(3) NOT NULL DEFAULT 'USD',
  `fixed_charges` decimal(25,4) NOT NULL DEFAULT '2.0000',
  `extra_charges_my` decimal(25,4) NOT NULL DEFAULT '3.9000',
  `extra_charges_other` decimal(25,4) NOT NULL DEFAULT '4.4000',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_permissions definition

CREATE TABLE `bynar_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) NOT NULL,
  `products-index` tinyint(1) DEFAULT '0',
  `products-add` tinyint(1) DEFAULT '0',
  `products-edit` tinyint(1) DEFAULT '0',
  `products-delete` tinyint(1) DEFAULT '0',
  `products-cost` tinyint(1) DEFAULT '0',
  `products-price` tinyint(1) DEFAULT '0',
  `quotes-index` tinyint(1) DEFAULT '0',
  `quotes-add` tinyint(1) DEFAULT '0',
  `quotes-edit` tinyint(1) DEFAULT '0',
  `quotes-pdf` tinyint(1) DEFAULT '0',
  `quotes-email` tinyint(1) DEFAULT '0',
  `quotes-delete` tinyint(1) DEFAULT '0',
  `sales-index` tinyint(1) DEFAULT '0',
  `sales-add` tinyint(1) DEFAULT '0',
  `sales-edit` tinyint(1) DEFAULT '0',
  `sales-pdf` tinyint(1) DEFAULT '0',
  `sales-email` tinyint(1) DEFAULT '0',
  `sales-delete` tinyint(1) DEFAULT '0',
  `purchases-index` tinyint(1) DEFAULT '0',
  `purchases-add` tinyint(1) DEFAULT '0',
  `purchases-edit` tinyint(1) DEFAULT '0',
  `purchases-pdf` tinyint(1) DEFAULT '0',
  `purchases-email` tinyint(1) DEFAULT '0',
  `purchases-delete` tinyint(1) DEFAULT '0',
  `transfers-index` tinyint(1) DEFAULT '0',
  `transfers-add` tinyint(1) DEFAULT '0',
  `transfers-edit` tinyint(1) DEFAULT '0',
  `transfers-pdf` tinyint(1) DEFAULT '0',
  `transfers-email` tinyint(1) DEFAULT '0',
  `transfers-delete` tinyint(1) DEFAULT '0',
  `customers-index` tinyint(1) DEFAULT '0',
  `customers-add` tinyint(1) DEFAULT '0',
  `customers-edit` tinyint(1) DEFAULT '0',
  `customers-delete` tinyint(1) DEFAULT '0',
  `suppliers-index` tinyint(1) DEFAULT '0',
  `suppliers-add` tinyint(1) DEFAULT '0',
  `suppliers-edit` tinyint(1) DEFAULT '0',
  `suppliers-delete` tinyint(1) DEFAULT '0',
  `sales-deliveries` tinyint(1) DEFAULT '0',
  `sales-add_delivery` tinyint(1) DEFAULT '0',
  `sales-edit_delivery` tinyint(1) DEFAULT '0',
  `sales-delete_delivery` tinyint(1) DEFAULT '0',
  `sales-email_delivery` tinyint(1) DEFAULT '0',
  `sales-pdf_delivery` tinyint(1) DEFAULT '0',
  `sales-gift_cards` tinyint(1) DEFAULT '0',
  `sales-add_gift_card` tinyint(1) DEFAULT '0',
  `sales-edit_gift_card` tinyint(1) DEFAULT '0',
  `sales-delete_gift_card` tinyint(1) DEFAULT '0',
  `pos-index` tinyint(1) DEFAULT '0',
  `sales-return_sales` tinyint(1) DEFAULT '0',
  `reports-index` tinyint(1) DEFAULT '0',
  `reports-warehouse_stock` tinyint(1) DEFAULT '0',
  `reports-quantity_alerts` tinyint(1) DEFAULT '0',
  `reports-expiry_alerts` tinyint(1) DEFAULT '0',
  `reports-products` tinyint(1) DEFAULT '0',
  `reports-daily_sales` tinyint(1) DEFAULT '0',
  `reports-monthly_sales` tinyint(1) DEFAULT '0',
  `reports-sales` tinyint(1) DEFAULT '0',
  `reports-payments` tinyint(1) DEFAULT '0',
  `reports-purchases` tinyint(1) DEFAULT '0',
  `reports-profit_loss` tinyint(1) DEFAULT '0',
  `reports-customers` tinyint(1) DEFAULT '0',
  `reports-suppliers` tinyint(1) DEFAULT '0',
  `reports-staff` tinyint(1) DEFAULT '0',
  `reports-register` tinyint(1) DEFAULT '0',
  `sales-payments` tinyint(1) DEFAULT '0',
  `purchases-payments` tinyint(1) DEFAULT '0',
  `purchases-expenses` tinyint(1) DEFAULT '0',
  `products-adjustments` tinyint(1) NOT NULL DEFAULT '0',
  `bulk_actions` tinyint(1) NOT NULL DEFAULT '0',
  `customers-deposits` tinyint(1) NOT NULL DEFAULT '0',
  `customers-delete_deposit` tinyint(1) NOT NULL DEFAULT '0',
  `products-barcode` tinyint(1) NOT NULL DEFAULT '0',
  `purchases-return_purchases` tinyint(1) NOT NULL DEFAULT '0',
  `reports-expenses` tinyint(1) NOT NULL DEFAULT '0',
  `reports-daily_purchases` tinyint(1) DEFAULT '0',
  `reports-monthly_purchases` tinyint(1) DEFAULT '0',
  `products-stock_count` tinyint(1) DEFAULT '0',
  `edit_price` tinyint(1) DEFAULT '0',
  `returns-index` tinyint(1) DEFAULT '0',
  `returns-add` tinyint(1) DEFAULT '0',
  `returns-edit` tinyint(1) DEFAULT '0',
  `returns-delete` tinyint(1) DEFAULT '0',
  `returns-email` tinyint(1) DEFAULT '0',
  `returns-pdf` tinyint(1) DEFAULT '0',
  `reports-tax` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_pos_register definition

CREATE TABLE `bynar_pos_register` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` int(11) NOT NULL,
  `cash_in_hand` decimal(25,4) NOT NULL,
  `status` varchar(10) NOT NULL,
  `total_cash` decimal(25,4) DEFAULT NULL,
  `total_cheques` int(11) DEFAULT NULL,
  `total_cc_slips` int(11) DEFAULT NULL,
  `total_cash_submitted` decimal(25,4) DEFAULT NULL,
  `total_cheques_submitted` int(11) DEFAULT NULL,
  `total_cc_slips_submitted` int(11) DEFAULT NULL,
  `note` text DEFAULT NULL,
  `closed_at` timestamp NULL DEFAULT NULL,
  `transfer_opened_bills` varchar(50) DEFAULT NULL,
  `closed_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_pos_settings definition

CREATE TABLE `bynar_pos_settings` (
  `pos_id` int(11) NOT NULL,
  `cat_limit` int(11) NOT NULL,
  `pro_limit` int(11) NOT NULL,
  `default_category` int(11) NOT NULL,
  `default_customer` int(11) NOT NULL,
  `default_biller` int(11) NOT NULL,
  `display_time` varchar(3) NOT NULL DEFAULT 'yes',
  `cf_title1` varchar(255) DEFAULT NULL,
  `cf_title2` varchar(255) DEFAULT NULL,
  `cf_value1` varchar(255) DEFAULT NULL,
  `cf_value2` varchar(255) DEFAULT NULL,
  `receipt_printer` varchar(55) DEFAULT NULL,
  `cash_drawer_codes` varchar(55) DEFAULT NULL,
  `focus_add_item` varchar(55) DEFAULT NULL,
  `add_manual_product` varchar(55) DEFAULT NULL,
  `customer_selection` varchar(55) DEFAULT NULL,
  `add_customer` varchar(55) DEFAULT NULL,
  `toggle_category_slider` varchar(55) DEFAULT NULL,
  `toggle_subcategory_slider` varchar(55) DEFAULT NULL,
  `cancel_sale` varchar(55) DEFAULT NULL,
  `suspend_sale` varchar(55) DEFAULT NULL,
  `print_items_list` varchar(55) DEFAULT NULL,
  `finalize_sale` varchar(55) DEFAULT NULL,
  `today_sale` varchar(55) DEFAULT NULL,
  `open_hold_bills` varchar(55) DEFAULT NULL,
  `close_register` varchar(55) DEFAULT NULL,
  `keyboard` tinyint(1) NOT NULL,
  `pos_printers` varchar(255) DEFAULT NULL,
  `java_applet` tinyint(1) NOT NULL,
  `product_button_color` varchar(20) NOT NULL DEFAULT 'default',
  `tooltips` tinyint(1) DEFAULT '1',
  `paypal_pro` tinyint(1) DEFAULT '0',
  `stripe` tinyint(1) DEFAULT '0',
  `rounding` tinyint(1) DEFAULT '0',
  `char_per_line` tinyint(4) DEFAULT '42',
  `pin_code` varchar(20) DEFAULT NULL,
  `purchase_code` varchar(100) DEFAULT 'purchase_code',
  `envato_username` varchar(50) DEFAULT 'envato_username',
  `version` varchar(10) DEFAULT '3.4.38',
  `after_sale_page` tinyint(1) DEFAULT '0',
  `item_order` tinyint(1) DEFAULT '0',
  `authorize` tinyint(1) DEFAULT '0',
  `toggle_brands_slider` varchar(55) DEFAULT NULL,
  `remote_printing` tinyint(1) DEFAULT '1',
  `printer` int(11) DEFAULT NULL,
  `order_printers` varchar(55) DEFAULT NULL,
  `auto_print` tinyint(1) DEFAULT '0',
  `customer_details` tinyint(1) DEFAULT NULL,
  `local_printers` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`pos_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_price_groups definition

CREATE TABLE `bynar_price_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_printers definition

CREATE TABLE `bynar_printers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(55) NOT NULL,
  `type` varchar(25) NOT NULL,
  `profile` varchar(25) NOT NULL,
  `char_per_line` tinyint(3) unsigned DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `ip_address` varbinary(45) DEFAULT NULL,
  `port` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_product_photos definition

CREATE TABLE `bynar_product_photos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `photo` varchar(100) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_product_prices definition

CREATE TABLE `bynar_product_prices` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `price_group_id` int(11) NOT NULL,
  `price` decimal(25,4) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `price_group_id` (`price_group_id`),
  KEY `product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_product_variants definition

CREATE TABLE `bynar_product_variants` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `name` varchar(55) NOT NULL,
  `cost` decimal(25,4) DEFAULT NULL,
  `price` decimal(25,4) DEFAULT NULL,
  `quantity` decimal(15,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `unique_product_id_name` (`product_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_products definition

CREATE TABLE `bynar_products` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(255) NOT NULL,
  `unit` int(11) DEFAULT NULL,
  `cost` decimal(25,4) DEFAULT NULL,
  `price` decimal(25,4) NOT NULL,
  `alert_quantity` decimal(15,4) DEFAULT '20.0000',
  `image` varchar(255) DEFAULT 'no_image.png',
  `category_id` int(11) NOT NULL,
  `subcategory_id` int(11) DEFAULT NULL,
  `cf1` varchar(255) DEFAULT NULL,
  `cf2` varchar(255) DEFAULT NULL,
  `cf3` varchar(255) DEFAULT NULL,
  `cf4` varchar(255) DEFAULT NULL,
  `cf5` varchar(255) DEFAULT NULL,
  `cf6` varchar(255) DEFAULT NULL,
  `quantity` decimal(15,4) DEFAULT '0.0000',
  `tax_rate` int(11) DEFAULT NULL,
  `track_quantity` tinyint(1) DEFAULT '1',
  `details` varchar(1000) DEFAULT NULL,
  `warehouse` int(11) DEFAULT NULL,
  `barcode_symbology` varchar(55) NOT NULL DEFAULT 'code128',
  `file` varchar(100) DEFAULT NULL,
  `product_details` text DEFAULT NULL,
  `tax_method` tinyint(1) DEFAULT '0',
  `type` varchar(55) NOT NULL DEFAULT 'standard',
  `supplier1` int(11) DEFAULT NULL,
  `supplier1price` decimal(25,4) DEFAULT NULL,
  `supplier2` int(11) DEFAULT NULL,
  `supplier2price` decimal(25,4) DEFAULT NULL,
  `supplier3` int(11) DEFAULT NULL,
  `supplier3price` decimal(25,4) DEFAULT NULL,
  `supplier4` int(11) DEFAULT NULL,
  `supplier4price` decimal(25,4) DEFAULT NULL,
  `supplier5` int(11) DEFAULT NULL,
  `supplier5price` decimal(25,4) DEFAULT NULL,
  `promotion` tinyint(1) DEFAULT '0',
  `promo_price` decimal(25,4) DEFAULT NULL,
  `start_date` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `supplier1_part_no` varchar(50) DEFAULT NULL,
  `supplier2_part_no` varchar(50) DEFAULT NULL,
  `supplier3_part_no` varchar(50) DEFAULT NULL,
  `supplier4_part_no` varchar(50) DEFAULT NULL,
  `supplier5_part_no` varchar(50) DEFAULT NULL,
  `sale_unit` int(11) DEFAULT NULL,
  `purchase_unit` int(11) DEFAULT NULL,
  `brand` int(11) DEFAULT NULL,
  `slug` varchar(55) DEFAULT NULL,
  `featured` tinyint(1) DEFAULT NULL,
  `weight` decimal(10,4) DEFAULT NULL,
  `hsn_code` int(11) DEFAULT NULL,
  `views` int(11) NOT NULL DEFAULT '0',
  `hide` tinyint(1) NOT NULL DEFAULT '0',
  `second_name` varchar(255) DEFAULT NULL,
  `hide_pos` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `code` (`code`),
  KEY `brand` (`brand`),
  KEY `category_id` (`category_id`),
  KEY `category_id_2` (`category_id`),
  KEY `id` (`id`),
  KEY `id_2` (`id`),
  KEY `unit` (`unit`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_promos definition

CREATE TABLE `bynar_promos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `product2buy` int(11) NOT NULL,
  `product2get` int(11) NOT NULL,
  `start_date` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `description` text DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_purchase_items definition

CREATE TABLE `bynar_purchase_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `purchase_id` int(11) DEFAULT NULL,
  `transfer_id` int(11) DEFAULT NULL,
  `product_id` int(11) NOT NULL,
  `product_code` varchar(50) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `option_id` int(11) DEFAULT NULL,
  `net_unit_cost` decimal(25,4) NOT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(20) DEFAULT NULL,
  `discount` varchar(20) DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `expiry` date DEFAULT NULL,
  `subtotal` decimal(25,4) NOT NULL,
  `quantity_balance` decimal(15,4) DEFAULT '0.0000',
  `date` date NOT NULL,
  `status` varchar(50) NOT NULL,
  `unit_cost` decimal(25,4) DEFAULT NULL,
  `real_unit_cost` decimal(25,4) DEFAULT NULL,
  `quantity_received` decimal(15,4) DEFAULT NULL,
  `supplier_part_no` varchar(50) DEFAULT NULL,
  `purchase_item_id` int(11) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  `base_unit_cost` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `purchase_id` (`purchase_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_purchases definition

CREATE TABLE `bynar_purchases` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `reference_no` varchar(55) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `supplier_id` int(11) NOT NULL,
  `supplier` varchar(55) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `note` varchar(1000) NOT NULL,
  `total` decimal(25,4) DEFAULT NULL,
  `product_discount` decimal(25,4) DEFAULT NULL,
  `order_discount_id` varchar(20) DEFAULT NULL,
  `order_discount` decimal(25,4) DEFAULT NULL,
  `total_discount` decimal(25,4) DEFAULT NULL,
  `product_tax` decimal(25,4) DEFAULT NULL,
  `order_tax_id` int(11) DEFAULT NULL,
  `order_tax` decimal(25,4) DEFAULT NULL,
  `total_tax` decimal(25,4) DEFAULT '0.0000',
  `shipping` decimal(25,4) DEFAULT '0.0000',
  `grand_total` decimal(25,4) NOT NULL,
  `paid` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `status` varchar(55) DEFAULT '',
  `payment_status` varchar(20) DEFAULT 'pending',
  `created_by` int(11) DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `attachment` varchar(55) DEFAULT NULL,
  `payment_term` tinyint(4) DEFAULT NULL,
  `due_date` date DEFAULT NULL,
  `return_id` int(11) DEFAULT NULL,
  `surcharge` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `return_purchase_ref` varchar(55) DEFAULT NULL,
  `purchase_id` int(11) DEFAULT NULL,
  `return_purchase_total` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_quote_items definition

CREATE TABLE `bynar_quote_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `quote_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `product_code` varchar(55) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `product_type` varchar(20) DEFAULT NULL,
  `option_id` int(11) DEFAULT NULL,
  `net_unit_price` decimal(25,4) NOT NULL,
  `unit_price` decimal(25,4) DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(55) DEFAULT NULL,
  `discount` varchar(55) DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `subtotal` decimal(25,4) NOT NULL,
  `serial_no` varchar(255) DEFAULT NULL,
  `real_unit_price` decimal(25,4) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `quote_id` (`quote_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_quotes definition

CREATE TABLE `bynar_quotes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reference_no` varchar(55) NOT NULL,
  `customer_id` int(11) NOT NULL,
  `customer` varchar(55) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `biller_id` int(11) NOT NULL,
  `biller` varchar(55) NOT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `internal_note` varchar(1000) DEFAULT NULL,
  `total` decimal(25,4) NOT NULL,
  `product_discount` decimal(25,4) DEFAULT '0.0000',
  `order_discount` decimal(25,4) DEFAULT NULL,
  `order_discount_id` varchar(20) DEFAULT NULL,
  `total_discount` decimal(25,4) DEFAULT '0.0000',
  `product_tax` decimal(25,4) DEFAULT '0.0000',
  `order_tax_id` int(11) DEFAULT NULL,
  `order_tax` decimal(25,4) DEFAULT NULL,
  `total_tax` decimal(25,4) DEFAULT NULL,
  `shipping` decimal(25,4) DEFAULT '0.0000',
  `grand_total` decimal(25,4) NOT NULL,
  `status` varchar(20) DEFAULT NULL,
  `created_by` int(11) DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `attachment` varchar(55) DEFAULT NULL,
  `supplier_id` int(11) DEFAULT NULL,
  `supplier` varchar(55) DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_return_items definition

CREATE TABLE `bynar_return_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `return_id` int(10) unsigned NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `product_code` varchar(55) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `product_type` varchar(20) DEFAULT NULL,
  `option_id` int(11) DEFAULT NULL,
  `net_unit_price` decimal(25,4) NOT NULL,
  `unit_price` decimal(25,4) DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(55) DEFAULT NULL,
  `discount` varchar(55) DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `subtotal` decimal(25,4) NOT NULL,
  `serial_no` varchar(255) DEFAULT NULL,
  `real_unit_price` decimal(25,4) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `product_id_2` (`product_id`,`return_id`),
  KEY `return_id` (`return_id`),
  KEY `return_id_2` (`return_id`,`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_returns definition

CREATE TABLE `bynar_returns` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reference_no` varchar(55) NOT NULL,
  `customer_id` int(11) NOT NULL,
  `customer` varchar(55) NOT NULL,
  `biller_id` int(11) NOT NULL,
  `biller` varchar(55) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `staff_note` varchar(1000) DEFAULT NULL,
  `total` decimal(25,4) NOT NULL,
  `product_discount` decimal(25,4) DEFAULT '0.0000',
  `order_discount_id` varchar(20) DEFAULT NULL,
  `total_discount` decimal(25,4) DEFAULT '0.0000',
  `order_discount` decimal(25,4) DEFAULT '0.0000',
  `product_tax` decimal(25,4) DEFAULT '0.0000',
  `order_tax_id` int(11) DEFAULT NULL,
  `order_tax` decimal(25,4) DEFAULT '0.0000',
  `total_tax` decimal(25,4) DEFAULT '0.0000',
  `grand_total` decimal(25,4) NOT NULL,
  `created_by` int(11) DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `total_items` smallint(6) DEFAULT NULL,
  `paid` decimal(25,4) DEFAULT '0.0000',
  `surcharge` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `attachment` varchar(55) DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  `shipping` decimal(25,4) DEFAULT '0.0000',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_sale_items definition

CREATE TABLE `bynar_sale_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sale_id` int(10) unsigned NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `product_code` varchar(55) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `product_type` varchar(20) DEFAULT NULL,
  `option_id` int(11) DEFAULT NULL,
  `net_unit_price` decimal(25,4) NOT NULL,
  `unit_price` decimal(25,4) DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(55) DEFAULT NULL,
  `discount` varchar(55) DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `subtotal` decimal(25,4) NOT NULL,
  `serial_no` varchar(255) DEFAULT NULL,
  `real_unit_price` decimal(25,4) DEFAULT NULL,
  `sale_item_id` int(11) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `product_id_2` (`product_id`,`sale_id`),
  KEY `sale_id` (`sale_id`),
  KEY `sale_id_2` (`sale_id`,`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_sales definition

CREATE TABLE `bynar_sales` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reference_no` varchar(55) NOT NULL,
  `customer_id` int(11) NOT NULL,
  `customer` varchar(55) NOT NULL,
  `biller_id` int(11) NOT NULL,
  `biller` varchar(55) NOT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `staff_note` varchar(1000) DEFAULT NULL,
  `total` decimal(25,4) NOT NULL,
  `product_discount` decimal(25,4) DEFAULT '0.0000',
  `order_discount_id` varchar(20) DEFAULT NULL,
  `total_discount` decimal(25,4) DEFAULT '0.0000',
  `order_discount` decimal(25,4) DEFAULT '0.0000',
  `product_tax` decimal(25,4) DEFAULT '0.0000',
  `order_tax_id` int(11) DEFAULT NULL,
  `order_tax` decimal(25,4) DEFAULT '0.0000',
  `total_tax` decimal(25,4) DEFAULT '0.0000',
  `shipping` decimal(25,4) DEFAULT '0.0000',
  `grand_total` decimal(25,4) NOT NULL,
  `sale_status` varchar(20) DEFAULT NULL,
  `payment_status` varchar(20) DEFAULT NULL,
  `payment_term` tinyint(4) DEFAULT NULL,
  `due_date` date DEFAULT NULL,
  `created_by` int(11) DEFAULT NULL,
  `updated_by` int(11) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `total_items` smallint(6) DEFAULT NULL,
  `pos` tinyint(1) NOT NULL DEFAULT '0',
  `paid` decimal(25,4) DEFAULT '0.0000',
  `return_id` int(11) DEFAULT NULL,
  `surcharge` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `attachment` varchar(55) DEFAULT NULL,
  `return_sale_ref` varchar(55) DEFAULT NULL,
  `sale_id` int(11) DEFAULT NULL,
  `return_sale_total` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `rounding` decimal(10,4) DEFAULT NULL,
  `suspend_note` varchar(255) DEFAULT NULL,
  `api` tinyint(1) DEFAULT '0',
  `shop` tinyint(1) DEFAULT '0',
  `address_id` int(11) DEFAULT NULL,
  `reserve_id` int(11) DEFAULT NULL,
  `hash` varchar(255) DEFAULT NULL,
  `manual_payment` varchar(55) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  `payment_method` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_sessions definition

CREATE TABLE `bynar_sessions` (
  `id` varchar(40) NOT NULL,
  `ip_address` varchar(45) NOT NULL,
  `timestamp` int(10) unsigned NOT NULL DEFAULT '0',
  `data` blob NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `ci_sessions_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_settings definition

CREATE TABLE `bynar_settings` (
  `setting_id` int(11) NOT NULL,
  `logo` varchar(255) NOT NULL,
  `logo2` varchar(255) NOT NULL,
  `site_name` varchar(55) NOT NULL,
  `language` varchar(20) NOT NULL,
  `default_warehouse` int(11) NOT NULL,
  `accounting_method` tinyint(4) NOT NULL DEFAULT '0',
  `default_currency` varchar(3) NOT NULL,
  `default_tax_rate` int(11) NOT NULL,
  `rows_per_page` int(11) NOT NULL,
  `version` varchar(10) NOT NULL DEFAULT '1.0',
  `default_tax_rate2` int(11) NOT NULL DEFAULT '0',
  `dateformat` int(11) NOT NULL,
  `sales_prefix` varchar(20) DEFAULT NULL,
  `quote_prefix` varchar(20) DEFAULT NULL,
  `purchase_prefix` varchar(20) DEFAULT NULL,
  `transfer_prefix` varchar(20) DEFAULT NULL,
  `delivery_prefix` varchar(20) DEFAULT NULL,
  `payment_prefix` varchar(20) DEFAULT NULL,
  `return_prefix` varchar(20) DEFAULT NULL,
  `returnp_prefix` varchar(20) DEFAULT NULL,
  `expense_prefix` varchar(20) DEFAULT NULL,
  `item_addition` tinyint(1) NOT NULL DEFAULT '0',
  `theme` varchar(20) NOT NULL,
  `product_serial` tinyint(4) NOT NULL,
  `default_discount` int(11) NOT NULL,
  `product_discount` tinyint(1) NOT NULL DEFAULT '0',
  `discount_method` tinyint(4) NOT NULL,
  `tax1` tinyint(4) NOT NULL,
  `tax2` tinyint(4) NOT NULL,
  `overselling` tinyint(1) NOT NULL DEFAULT '0',
  `restrict_user` tinyint(4) NOT NULL DEFAULT '0',
  `restrict_calendar` tinyint(4) NOT NULL DEFAULT '0',
  `timezone` varchar(100) DEFAULT NULL,
  `iwidth` int(11) NOT NULL DEFAULT '0',
  `iheight` int(11) NOT NULL,
  `twidth` int(11) NOT NULL,
  `theight` int(11) NOT NULL,
  `watermark` tinyint(1) DEFAULT NULL,
  `reg_ver` tinyint(1) DEFAULT NULL,
  `allow_reg` tinyint(1) DEFAULT NULL,
  `reg_notification` tinyint(1) DEFAULT NULL,
  `auto_reg` tinyint(1) DEFAULT NULL,
  `protocol` varchar(20) NOT NULL DEFAULT 'mail',
  `mailpath` varchar(55) DEFAULT '/usr/sbin/sendmail',
  `smtp_host` varchar(100) DEFAULT NULL,
  `smtp_user` varchar(100) DEFAULT NULL,
  `smtp_pass` varchar(255) DEFAULT NULL,
  `smtp_port` varchar(10) DEFAULT '25',
  `smtp_crypto` varchar(10) DEFAULT NULL,
  `corn` datetime DEFAULT NULL,
  `customer_group` int(11) NOT NULL,
  `default_email` varchar(100) NOT NULL,
  `mmode` tinyint(1) NOT NULL,
  `bc_fix` tinyint(4) NOT NULL DEFAULT '0',
  `auto_detect_barcode` tinyint(1) NOT NULL DEFAULT '0',
  `captcha` tinyint(1) NOT NULL DEFAULT '1',
  `reference_format` tinyint(1) NOT NULL DEFAULT '1',
  `racks` tinyint(1) DEFAULT '0',
  `attributes` tinyint(1) NOT NULL DEFAULT '0',
  `product_expiry` tinyint(1) NOT NULL DEFAULT '0',
  `decimals` tinyint(4) NOT NULL DEFAULT '2',
  `qty_decimals` tinyint(4) NOT NULL DEFAULT '2',
  `decimals_sep` varchar(2) NOT NULL DEFAULT '.',
  `thousands_sep` varchar(2) NOT NULL DEFAULT ',',
  `invoice_view` tinyint(1) DEFAULT '0',
  `default_biller` int(11) DEFAULT NULL,
  `envato_username` varchar(50) DEFAULT NULL,
  `purchase_code` varchar(100) DEFAULT NULL,
  `rtl` tinyint(1) DEFAULT '0',
  `each_spent` decimal(15,4) DEFAULT NULL,
  `ca_point` tinyint(4) DEFAULT NULL,
  `each_sale` decimal(15,4) DEFAULT NULL,
  `sa_point` tinyint(4) DEFAULT NULL,
  `update` tinyint(1) DEFAULT '0',
  `sac` tinyint(1) DEFAULT '0',
  `display_all_products` tinyint(1) DEFAULT '0',
  `display_symbol` tinyint(1) DEFAULT NULL,
  `symbol` varchar(50) DEFAULT NULL,
  `remove_expired` tinyint(1) DEFAULT '0',
  `barcode_separator` varchar(2) NOT NULL DEFAULT '-',
  `set_focus` tinyint(1) NOT NULL DEFAULT '0',
  `price_group` int(11) DEFAULT NULL,
  `barcode_img` tinyint(1) NOT NULL DEFAULT '1',
  `ppayment_prefix` varchar(20) DEFAULT 'POP',
  `disable_editing` smallint(6) DEFAULT '90',
  `qa_prefix` varchar(55) DEFAULT NULL,
  `update_cost` tinyint(1) DEFAULT NULL,
  `apis` tinyint(1) NOT NULL DEFAULT '0',
  `state` varchar(100) DEFAULT NULL,
  `pdf_lib` varchar(20) DEFAULT 'dompdf',
  `use_code_for_slug` tinyint(1) DEFAULT NULL,
  `ws_barcode_type` varchar(10) DEFAULT 'weight',
  `ws_barcode_chars` tinyint(4) DEFAULT NULL,
  `flag_chars` tinyint(4) DEFAULT NULL,
  `item_code_start` tinyint(4) DEFAULT NULL,
  `item_code_chars` tinyint(4) DEFAULT NULL,
  `price_start` tinyint(4) DEFAULT NULL,
  `price_chars` tinyint(4) DEFAULT NULL,
  `price_divide_by` int(11) DEFAULT NULL,
  `weight_start` tinyint(4) DEFAULT NULL,
  `weight_chars` tinyint(4) DEFAULT NULL,
  `weight_divide_by` int(11) DEFAULT NULL,
  `site_id` varchar(50) DEFAULT NULL,
  `address` varchar(100) DEFAULT NULL,
  `address2` varchar(100) DEFAULT NULL,
  `city` varchar(30) DEFAULT NULL,
  `sat_zip_code` varchar(30) DEFAULT NULL,
  `country_code` varchar(20) DEFAULT NULL,
  `contact_name` varchar(50) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `federal_id_no` varchar(20) DEFAULT NULL,
  `eroi_no` varchar(20) DEFAULT NULL,
  `email` varchar(30) DEFAULT NULL,
  `homepage` varchar(30) DEFAULT NULL,
  `bank_name` varchar(30) DEFAULT NULL,
  `bank_branch_no` varchar(20) DEFAULT NULL,
  `ship_to_name` varchar(30) DEFAULT NULL,
  `ship_to_address` varchar(100) DEFAULT NULL,
  `ship_to_address2` varchar(100) DEFAULT NULL,
  `base_calendar_code` varchar(20) DEFAULT NULL,
  `ship_to_state` varchar(30) DEFAULT NULL,
  `ship_to_zip_code` varchar(20) DEFAULT NULL,
  `ship_to_contact` varchar(30) DEFAULT NULL,
  `location_code` varchar(10) DEFAULT NULL,
  `tax_area_code` varchar(20) DEFAULT NULL,
  `tax_exemption_no` varchar(20) DEFAULT NULL,
  `tax_scheme` varchar(20) DEFAULT NULL,
  `gln` varchar(30) DEFAULT NULL,
  `bank_account_no` varchar(30) DEFAULT NULL,
  `payment_routing_no` varchar(30) DEFAULT NULL,
  `giro_no` varchar(30) DEFAULT NULL,
  `sales_prefix_2` varchar(5) DEFAULT NULL,
  `sales_template_prefix` varchar(10) DEFAULT NULL,
  `sales_draft_prefix` varchar(10) DEFAULT NULL,
  `sales_order_prefix` varchar(10) DEFAULT NULL,
  `sales_invoice_prefix` varchar(10) DEFAULT NULL,
  `sales_contract_prefix` varchar(10) DEFAULT NULL,
  `sales_recurring_invoice` varchar(10) DEFAULT NULL,
  `return_sales_invoice` varchar(10) DEFAULT NULL,
  `return_sales_order_prefix` varchar(10) DEFAULT NULL,
  `purchases_template_prefix` varchar(10) DEFAULT NULL,
  `purchases_quote_prefix` varchar(10) DEFAULT NULL,
  `purchases_order_prefix` varchar(10) DEFAULT NULL,
  `purchases_invoice_prefix` varchar(10) DEFAULT NULL,
  `purchases_contract_prefix` varchar(10) DEFAULT NULL,
  `purchases_recurring_invoice_prefix` varchar(10) DEFAULT NULL,
  `return_purchases_invoice_prefix` varchar(10) DEFAULT NULL,
  `receipt_template_prefix` varchar(10) DEFAULT NULL,
  `receipt_draft_prefix` varchar(10) DEFAULT NULL,
  `receipt_order_prefix` varchar(10) DEFAULT NULL,
  `receipt_documents_prefix` varchar(10) DEFAULT NULL,
  `production_template_prefix` varchar(10) DEFAULT NULL,
  `production_draft_prefix` varchar(10) DEFAULT NULL,
  `production_order_prefix` varchar(10) DEFAULT NULL,
  `production_document_prefix` varchar(10) DEFAULT NULL,
  `receipts_order_prefix` varchar(10) DEFAULT NULL,
  `receipts_document_prefix` varchar(10) DEFAULT NULL,
  `payments_order_prefix` varchar(10) DEFAULT NULL,
  `payments_document_prefix` varchar(10) DEFAULT NULL,
  `return_purchases_draft_prefix` varchar(10) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`setting_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_skrill definition

CREATE TABLE `bynar_skrill` (
  `id` int(11) NOT NULL,
  `active` tinyint(4) NOT NULL,
  `account_email` varchar(255) NOT NULL DEFAULT 'testaccount2@moneybookers.com',
  `secret_word` varchar(20) NOT NULL DEFAULT 'mbtest',
  `skrill_currency` varchar(3) NOT NULL DEFAULT 'USD',
  `fixed_charges` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `extra_charges_my` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `extra_charges_other` decimal(25,4) NOT NULL DEFAULT '0.0000',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_stock_count_items definition

CREATE TABLE `bynar_stock_count_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `stock_count_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `product_code` varchar(50) DEFAULT NULL,
  `product_name` varchar(255) DEFAULT NULL,
  `product_variant` varchar(55) DEFAULT NULL,
  `product_variant_id` int(11) DEFAULT NULL,
  `expected` decimal(15,4) NOT NULL,
  `counted` decimal(15,4) NOT NULL,
  `cost` decimal(25,4) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `stock_count_id` (`stock_count_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_suspended_bills definition

CREATE TABLE `bynar_suspended_bills` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `customer_id` int(11) NOT NULL,
  `customer` varchar(55) DEFAULT NULL,
  `count` int(11) NOT NULL,
  `order_discount_id` varchar(20) DEFAULT NULL,
  `order_tax_id` int(11) DEFAULT NULL,
  `total` decimal(25,4) NOT NULL,
  `biller_id` int(11) DEFAULT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `created_by` int(11) NOT NULL,
  `suspend_note` varchar(255) DEFAULT NULL,
  `shipping` decimal(15,4) DEFAULT '0.0000',
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_suspended_items definition

CREATE TABLE `bynar_suspended_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `suspend_id` int(10) unsigned NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `product_code` varchar(55) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `net_unit_price` decimal(25,4) NOT NULL,
  `unit_price` decimal(25,4) NOT NULL,
  `quantity` decimal(15,4) DEFAULT '0.0000',
  `warehouse_id` int(11) DEFAULT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(55) DEFAULT NULL,
  `discount` varchar(55) DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `subtotal` decimal(25,4) NOT NULL,
  `serial_no` varchar(255) DEFAULT NULL,
  `option_id` int(11) DEFAULT NULL,
  `product_type` varchar(20) DEFAULT NULL,
  `real_unit_price` decimal(25,4) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_tax_rates definition

CREATE TABLE `bynar_tax_rates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(55) NOT NULL,
  `code` varchar(10) DEFAULT NULL,
  `rate` decimal(12,4) NOT NULL,
  `type` varchar(50) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_transfer_items definition

CREATE TABLE `bynar_transfer_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `transfer_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `product_code` varchar(55) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `option_id` int(11) DEFAULT NULL,
  `expiry` date DEFAULT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `tax_rate_id` int(11) DEFAULT NULL,
  `tax` varchar(55) DEFAULT NULL,
  `item_tax` decimal(25,4) DEFAULT NULL,
  `net_unit_cost` decimal(25,4) DEFAULT NULL,
  `subtotal` decimal(25,4) DEFAULT NULL,
  `quantity_balance` decimal(15,4) NOT NULL,
  `unit_cost` decimal(25,4) DEFAULT NULL,
  `real_unit_cost` decimal(25,4) DEFAULT NULL,
  `date` date DEFAULT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `product_unit_id` int(11) DEFAULT NULL,
  `product_unit_code` varchar(10) DEFAULT NULL,
  `unit_quantity` decimal(15,4) NOT NULL,
  `gst` varchar(20) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `transfer_id` (`transfer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_transfers definition

CREATE TABLE `bynar_transfers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `transfer_no` varchar(55) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `from_warehouse_id` int(11) NOT NULL,
  `from_warehouse_code` varchar(55) NOT NULL,
  `from_warehouse_name` varchar(55) NOT NULL,
  `to_warehouse_id` int(11) NOT NULL,
  `to_warehouse_code` varchar(55) NOT NULL,
  `to_warehouse_name` varchar(55) NOT NULL,
  `note` varchar(1000) DEFAULT NULL,
  `total` decimal(25,4) DEFAULT NULL,
  `total_tax` decimal(25,4) DEFAULT NULL,
  `grand_total` decimal(25,4) DEFAULT NULL,
  `created_by` varchar(255) DEFAULT NULL,
  `status` varchar(55) NOT NULL DEFAULT 'pending',
  `shipping` decimal(25,4) NOT NULL DEFAULT '0.0000',
  `attachment` varchar(55) DEFAULT NULL,
  `cgst` decimal(25,4) DEFAULT NULL,
  `sgst` decimal(25,4) DEFAULT NULL,
  `igst` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_units definition

CREATE TABLE `bynar_units` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(55) NOT NULL,
  `base_unit` int(11) DEFAULT NULL,
  `operator` varchar(1) DEFAULT NULL,
  `unit_value` varchar(55) DEFAULT NULL,
  `operation_value` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `base_unit` (`base_unit`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_user_logins definition

CREATE TABLE `bynar_user_logins` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `company_id` int(11) DEFAULT NULL,
  `ip_address` varbinary(16) NOT NULL,
  `login` varchar(100) NOT NULL,
  `time` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_users definition

CREATE TABLE `bynar_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `last_ip_address` varbinary(45) DEFAULT NULL,
  `ip_address` varbinary(45) NOT NULL,
  `username` varchar(100) NOT NULL,
  `password` varchar(40) NOT NULL,
  `salt` varchar(40) DEFAULT NULL,
  `email` varchar(100) NOT NULL,
  `activation_code` varchar(40) DEFAULT NULL,
  `forgotten_password_code` varchar(40) DEFAULT NULL,
  `forgotten_password_time` int(10) unsigned DEFAULT NULL,
  `remember_code` varchar(40) DEFAULT NULL,
  `created_on` int(10) unsigned NOT NULL,
  `last_login` int(10) unsigned DEFAULT NULL,
  `active` tinyint(3) unsigned DEFAULT NULL,
  `first_name` varchar(50) DEFAULT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  `company` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `avatar` varchar(55) DEFAULT NULL,
  `gender` varchar(20) DEFAULT NULL,
  `group_id` int(10) unsigned NOT NULL,
  `warehouse_id` int(10) unsigned DEFAULT NULL,
  `biller_id` int(10) unsigned DEFAULT NULL,
  `company_id` int(11) DEFAULT NULL,
  `show_cost` tinyint(1) DEFAULT '0',
  `show_price` tinyint(1) DEFAULT '0',
  `award_points` int(11) DEFAULT '0',
  `view_right` tinyint(1) NOT NULL DEFAULT '0',
  `edit_right` tinyint(1) NOT NULL DEFAULT '0',
  `allow_discount` tinyint(1) DEFAULT '0',
  `middle_name` varchar(50) DEFAULT NULL,
  `department` varchar(50) DEFAULT NULL,
  `cost_center_code` varchar(150) DEFAULT NULL,
  `warehouse_code` varchar(150) DEFAULT NULL,
  `address` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `grids_permissions_code` varchar(15) DEFAULT NULL,
  `grids_permissions_grids` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `group_id` (`group_id`,`warehouse_id`,`biller_id`),
  KEY `group_id_2` (`group_id`,`company_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_variants definition

CREATE TABLE `bynar_variants` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(55) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_warehouses definition

CREATE TABLE `bynar_warehouses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `map` varchar(255) DEFAULT NULL,
  `phone` varchar(55) DEFAULT NULL,
  `email` varchar(55) DEFAULT NULL,
  `price_group_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_warehouses_products definition

CREATE TABLE `bynar_warehouses_products` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `rack` varchar(55) DEFAULT NULL,
  `avg_cost` decimal(25,4) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `product_id` (`product_id`),
  KEY `warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.bynar_warehouses_products_variants definition

CREATE TABLE `bynar_warehouses_products_variants` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `option_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `warehouse_id` int(11) NOT NULL,
  `quantity` decimal(15,4) NOT NULL,
  `rack` varchar(55) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `option_id` (`option_id`),
  KEY `product_id` (`product_id`),
  KEY `warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.campaign_statuses definition

CREATE TABLE `campaign_statuses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.campaigns definition

CREATE TABLE `campaigns` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status_code` varchar(255) DEFAULT NULL,
  `starting_date` datetime DEFAULT NULL,
  `ending_date` datetime DEFAULT NULL,
  `salesperson_uuid` int(11) DEFAULT NULL,
  `department_uuid` int(11) DEFAULT NULL,
  `project_uuid` int(11) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.cash definition

CREATE TABLE `cash` (
  `id` int(11) NOT NULL,
  `type` varchar(55) NOT NULL,
  `currency` varchar(55) DEFAULT NULL,
  `number` varchar(55) DEFAULT NULL,
  `swift` varchar(55) DEFAULT NULL,
  `iban` varchar(55) DEFAULT NULL,
  `address` varchar(150) DEFAULT NULL,
  `phone` varchar(55) DEFAULT NULL,
  `email` varchar(55) DEFAULT NULL,
  `city` varchar(55) DEFAULT NULL,
  `state` varchar(55) DEFAULT NULL,
  `contact` varchar(55) DEFAULT NULL,
  `country` varchar(55) DEFAULT NULL,
  `webpage` varchar(50) DEFAULT NULL,
  `postal_code` varchar(50) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `note` varchar(255) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `address2` varchar(150) DEFAULT NULL,
  `grid_id` int(11) NOT NULL,
  `Parent` varchar(30) DEFAULT NULL,
  `Def` varchar(20) DEFAULT NULL,
  `balance` decimal(25,4) DEFAULT NULL,
  `currency_code` varchar(20) DEFAULT NULL,
  `payments_documents` int(11) DEFAULT NULL,
  `payments_total` decimal(25,4) DEFAULT NULL,
  `receipts_documents` int(11) DEFAULT NULL,
  `receipts_total` decimal(25,4) DEFAULT NULL,
  `cashier_in_charge` varchar(50) DEFAULT NULL,
  `cashier_in_charge_department` varchar(50) DEFAULT NULL,
  `fax` varchar(20) DEFAULT NULL,
  `payments_documents_draft` int(11) DEFAULT NULL,
  `payments_documents_draft_total` decimal(25,4) DEFAULT NULL,
  `receipts_documents_draft` int(11) DEFAULT NULL,
  `receitps_documents_draft_total` decimal(25,4) DEFAULT NULL,
  `deposit_document_draft` int(11) DEFAULT NULL,
  `deposit_document_draft_total` decimal(25,4) DEFAULT NULL,
  `deposit_document` int(11) DEFAULT NULL,
  `deposit_document_total` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `bankscash_uuid_uindex` (`uuid`),
  UNIQUE KEY `cash_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `cash_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.cash_flow definition

CREATE TABLE `cash_flow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `bank_cash_id` int(11) DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `value` decimal(25,10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.cash_management definition

CREATE TABLE `cash_management` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `bank_cash_id` int(11) DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `value` decimal(25,10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.catalogue_items definition

CREATE TABLE `catalogue_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `Parent` int(11) DEFAULT NULL,
  `item_name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_code` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_barcode` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_discount` decimal(25,4) DEFAULT NULL,
  `item_discount_value` decimal(25,4) DEFAULT NULL,
  `item_price` decimal(25,4) DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_type` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_brand` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_category` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_subcategory` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `price_list_items_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.catalogues definition

CREATE TABLE `catalogues` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(11) DEFAULT NULL,
  `Parent` varchar(10) DEFAULT NULL,
  `document_type` varchar(30) DEFAULT NULL,
  `document_abbrevation` varchar(30) DEFAULT NULL,
  `document_no` varchar(11) DEFAULT NULL,
  `document_date` varchar(30) DEFAULT NULL,
  `posting_date` varchar(11) DEFAULT NULL,
  `available_from_date` varchar(11) DEFAULT NULL,
  `available_till_date` varchar(11) DEFAULT NULL,
  `approval_manager` varchar(30) DEFAULT NULL,
  `approval_manager_department` varchar(30) DEFAULT NULL,
  `approval_manager_status` tinyint(1) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `price_lists_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.categories definition

CREATE TABLE `categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `units_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.causes_of_absences definition

CREATE TABLE `causes_of_absences` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `unit` decimal(25,4) DEFAULT NULL,
  `total_absence` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.causes_of_inactivity definition

CREATE TABLE `causes_of_inactivity` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.chart_of_accounts definition

CREATE TABLE `chart_of_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `income_balance` tinyint(4) DEFAULT NULL,
  `debit_credit` tinyint(4) DEFAULT NULL,
  `account_type` tinyint(4) DEFAULT NULL,
  `totaling` decimal(25,4) DEFAULT NULL,
  `net_change` decimal(25,4) DEFAULT NULL,
  `balance` decimal(25,4) DEFAULT NULL,
  `consol_credit_acc` varchar(255) DEFAULT NULL,
  `consol_debit_acc` varchar(255) DEFAULT NULL,
  `consol_translation_method` tinyint(4) DEFAULT NULL,
  `exchange_rate_adjustment` tinyint(4) DEFAULT NULL,
  `cost_type_no` int(11) DEFAULT NULL,
  `direct_posting` tinyint(4) DEFAULT NULL,
  `gen_posting_type_uuid` int(11) DEFAULT NULL,
  `gen_bus_posting_type_uuid` int(11) DEFAULT NULL,
  `gen_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `default_ic_partner_gl_acc_no_uuid` int(11) DEFAULT NULL,
  `default_deferral_template_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.chart_of_accounts_items definition

CREATE TABLE `chart_of_accounts_items` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) DEFAULT NULL,
  `main_account` varchar(15) DEFAULT NULL,
  `description` varchar(30) DEFAULT NULL,
  `main_account_type` varchar(30) DEFAULT NULL,
  `main_acccount_category` varchar(30) DEFAULT NULL,
  `Parent` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `chart_of_accounts_items_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.charts_of_accounts definition

CREATE TABLE `charts_of_accounts` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(4) DEFAULT NULL,
  `country_code` varchar(10) DEFAULT NULL,
  `name` varchar(30) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `charts_of_accounts_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.classes definition

CREATE TABLE `classes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.companies definition

CREATE TABLE `companies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `category_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.company_categories definition

CREATE TABLE `company_categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.confidentials definition

CREATE TABLE `confidentials` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.contact_information definition

CREATE TABLE `contact_information` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `organisation_id` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `address_2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `postal_code` int(11) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.contacts definition

CREATE TABLE `contacts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `surname` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `middle_name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `fax` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `email` varchar(30) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code_2` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `language` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `position` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `department` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `contacts_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.cost_centers definition

CREATE TABLE `cost_centers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.costings definition

CREATE TABLE `costings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_id` int(11) DEFAULT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `transaction_id` int(11) DEFAULT NULL,
  `costing_method` tinyint(4) DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `quantity` decimal(25,10) DEFAULT NULL,
  `value` decimal(25,10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.countries definition

CREATE TABLE `countries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `country` varchar(50) DEFAULT NULL,
  `alpha_2_code` varchar(2) DEFAULT NULL,
  `alpha_3_code` varchar(3) DEFAULT NULL,
  `numeric` int(11) DEFAULT NULL,
  `fips` varchar(5) DEFAULT NULL,
  `licence_plate` varchar(5) DEFAULT NULL,
  `domain` varchar(5) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `countries_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.countriess definition

CREATE TABLE `countriess` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `iso_code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `iso_code_no` int(11) DEFAULT NULL,
  `address_format` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contact_address_format` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `table_name_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `table_name_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.currencies definition

CREATE TABLE `currencies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `state` varchar(40) DEFAULT NULL,
  `description` varchar(30) DEFAULT NULL,
  `code` varchar(5) DEFAULT NULL,
  `no_code` int(11) DEFAULT NULL,
  `rate` decimal(12,4) NOT NULL,
  `exchange_rate_date` varchar(20) DEFAULT NULL,
  `symbol` varchar(3) DEFAULT NULL,
  `responsibility_center` varchar(20) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `currencies_id_uindex` (`id`),
  UNIQUE KEY `currencies_iso_code_uindex` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.customer_bank_accounts definition

CREATE TABLE `customer_bank_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `desription` varchar(255) DEFAULT NULL,
  `bank_account_no` varchar(255) DEFAULT NULL,
  `transit_no` varchar(255) DEFAULT NULL,
  `bank_clearing_standart_code` varchar(255) DEFAULT NULL,
  `swift` varchar(255) DEFAULT NULL,
  `iban` varchar(255) DEFAULT NULL,
  `bank_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `bank_clearing_standart_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_groups definition

CREATE TABLE `customer_groups` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) NOT NULL,
  `Parent` int(11) DEFAULT NULL,
  `code` varchar(15) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `payment_terms` varchar(20) DEFAULT NULL,
  `invoice_payment_date` varchar(30) DEFAULT NULL,
  `tax_group` varchar(15) DEFAULT NULL,
  `customer_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`grid_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_invoicing definition

CREATE TABLE `customer_invoicing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `vat_no_uuid` int(11) DEFAULT NULL,
  `quote_from` int(11) DEFAULT NULL,
  `gen_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `customer_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `customer_price_group_uuid` int(11) DEFAULT NULL,
  `customer_discount_group_uuid` int(11) DEFAULT NULL,
  `invoice_discount_code_uuid` int(11) DEFAULT NULL,
  `line_discount` tinyint(4) DEFAULT NULL,
  `prices_vat` tinyint(4) DEFAULT NULL,
  `customer_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_list_items definition

CREATE TABLE `customer_list_items` (
  `id` int(11) DEFAULT NULL,
  `Parent` int(11) DEFAULT NULL,
  `item_code` varchar(20) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_lists definition

CREATE TABLE `customer_lists` (
  `id` int(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `document_no` varchar(20) DEFAULT NULL,
  `document_date` varchar(20) DEFAULT NULL,
  `customer_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `customer_uuid` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_payments definition

CREATE TABLE `customer_payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `payments` decimal(25,4) DEFAULT NULL,
  `application_method` tinyint(4) DEFAULT NULL,
  `partner_type` tinyint(4) DEFAULT NULL,
  `payment_terms_uuid` int(11) DEFAULT NULL,
  `payment_method_uuid` int(11) DEFAULT NULL,
  `reminder_term_uuid` int(11) DEFAULT NULL,
  `finance_charge_term_uuid` int(11) DEFAULT NULL,
  `cash_flow_payment_term_uuid` int(11) DEFAULT NULL,
  `preferred_bank_account_uuid` int(11) DEFAULT NULL,
  `company_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_posting_groups definition

CREATE TABLE `customer_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `receivables_account` int(11) DEFAULT NULL,
  `service_charge_account` int(11) DEFAULT NULL,
  `payment_disc_debit_account` int(11) DEFAULT NULL,
  `payment_disc_credit_account` int(11) DEFAULT NULL,
  `invoice_rounding_account` int(11) DEFAULT NULL,
  `debit_curr_appln_rndg_account` int(11) DEFAULT NULL,
  `credit_curr_appln_rndg_account` int(11) DEFAULT NULL,
  `debit_rounding_account` int(11) DEFAULT NULL,
  `credit_rounding_account` int(11) DEFAULT NULL,
  `used_customers` int(11) DEFAULT NULL,
  `used_ledger` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_price_groups definition

CREATE TABLE `customer_price_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `item_discount` tinyint(4) DEFAULT NULL,
  `invoice_discount` tinyint(4) DEFAULT NULL,
  `include_vat` tinyint(4) DEFAULT NULL,
  `vat_business_posting_group_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customer_shipping definition

CREATE TABLE `customer_shipping` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `combine_shipments` tinyint(4) DEFAULT NULL,
  `reserve` int(11) DEFAULT NULL,
  `shipping_advise` tinyint(4) DEFAULT NULL,
  `shipment_method_uuid` int(11) DEFAULT NULL,
  `base_calendar_uuid` int(11) DEFAULT NULL,
  `agent_uuid` int(11) DEFAULT NULL,
  `agent_service_uuid` int(11) DEFAULT NULL,
  `shipping_time_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.customers definition

CREATE TABLE `customers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `company_id` int(11) DEFAULT NULL,
  `ic_partner_id` int(11) DEFAULT NULL,
  `service_zone_id` int(11) DEFAULT NULL,
  `document_sending_profile_id` int(11) DEFAULT NULL,
  `sales_person_id` int(11) DEFAULT NULL,
  `responsibility_center_id` int(11) DEFAULT NULL,
  `address_id` int(11) DEFAULT NULL,
  `contact_id` int(11) DEFAULT NULL,
  `vat` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `gln` int(11) DEFAULT NULL,
  `invoicing_id` int(11) DEFAULT NULL,
  `payment_id` int(11) DEFAULT NULL,
  `shipment_id` int(11) DEFAULT NULL,
  `bill_to_customer_id` int(11) DEFAULT NULL,
  `copy_sell_to` tinyint(4) DEFAULT NULL,
  `invoice_copies` int(11) DEFAULT NULL,
  `currency_id` int(11) DEFAULT NULL,
  `price_group_id` int(11) DEFAULT NULL,
  `discount_group_id` int(11) DEFAULT NULL,
  `invoice_discount_code` int(11) DEFAULT NULL,
  `line_discount` tinyint(4) DEFAULT NULL,
  `prices_vat` tinyint(4) DEFAULT NULL,
  `payment_terms_id` int(11) DEFAULT NULL,
  `warehouse_id` int(11) DEFAULT NULL,
  `combine_shipments` tinyint(4) DEFAULT NULL,
  `reserve` tinyint(4) DEFAULT NULL,
  `shipping_advice` tinyint(4) DEFAULT NULL,
  `shipment_method_id` int(11) DEFAULT NULL,
  `base_calendar_id` int(11) DEFAULT NULL,
  `customized_calendar` tinyint(4) DEFAULT NULL,
  `shipment_agent_id` int(11) DEFAULT NULL,
  `shipment_agent_service_id` int(11) DEFAULT NULL,
  `shipping_time` int(11) DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `gen_bus_posting_group_id` int(11) DEFAULT NULL,
  `vat_bus_posting_group_id` int(11) DEFAULT NULL,
  `customer_posting_group_id` int(11) DEFAULT NULL,
  `prepayment` decimal(25,4) DEFAULT NULL,
  `partner_type` int(11) DEFAULT NULL,
  `reminder_terms_uuid` int(11) DEFAULT NULL,
  `charge_terms_uuid` int(11) DEFAULT NULL,
  `cash_flow_payment_terms_uuid` int(11) DEFAULT NULL,
  `print_statement` tinyint(4) DEFAULT NULL,
  `last_statement_no` int(11) DEFAULT NULL,
  `block_payment_tolerance` tinyint(4) DEFAULT NULL,
  `bank_id` int(11) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `document_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `customers_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.customers_discount_groups definition

CREATE TABLE `customers_discount_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.defferals definition

CREATE TABLE `defferals` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `defferral_account` int(11) DEFAULT NULL,
  `defferal_percentage` int(11) DEFAULT NULL,
  `calculation_method` varchar(255) DEFAULT NULL,
  `start_date` varchar(255) DEFAULT NULL,
  `no_of_periods` int(11) DEFAULT NULL,
  `period_description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.departments definition

CREATE TABLE `departments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code_2` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `bynar_all_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.depraciations definition

CREATE TABLE `depraciations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) NOT NULL,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_depreciation_total` decimal(25,4) DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `warehouseman_approve` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `depraciations_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.depraciations_items definition

CREATE TABLE `depraciations_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `depreciation_rate` decimal(25,4) DEFAULT NULL,
  `book_value` decimal(25,4) DEFAULT NULL,
  `item_depreciation_total` decimal(25,4) NOT NULL,
  `depreciation_method_uuid` int(11) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_tempory` int(11) NOT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `depraciations_items_id_uindex` (`id`),
  KEY `depraciations_items_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.dimensions definition

CREATE TABLE `dimensions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `dimensions_Id_uindex` (`id`),
  UNIQUE KEY `dimensions_code_uindex` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.discount_groups definition

CREATE TABLE `discount_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `discount_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.discounts definition

CREATE TABLE `discounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `value` decimal(25,4) DEFAULT NULL,
  `operation` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `discounts_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.document_flow_lines definition

CREATE TABLE `document_flow_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `document_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `order_flow` int(11) DEFAULT NULL,
  `update` tinyint(4) DEFAULT NULL,
  `delete` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.document_sending_profiles definition

CREATE TABLE `document_sending_profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `printer` tinyint(4) DEFAULT NULL,
  `email` tinyint(4) DEFAULT NULL,
  `disc` tinyint(4) DEFAULT NULL,
  `electronic_document` tinyint(4) DEFAULT NULL,
  `other` tinyint(4) DEFAULT NULL,
  `default` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.documents definition

CREATE TABLE `documents` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_type` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `workspace` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `series_id` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.documents_workflow definition

CREATE TABLE `documents_workflow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.employee_posting_groups definition

CREATE TABLE `employee_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `payables_account` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.employee_statistics_groups definition

CREATE TABLE `employee_statistics_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.employees definition

CREATE TABLE `employees` (
  `id` int(11) NOT NULL,
  `surname` varchar(55) DEFAULT NULL,
  `name` varchar(55) DEFAULT NULL,
  `middle_name` varchar(55) DEFAULT NULL,
  `position` varchar(55) DEFAULT NULL,
  `proffesion` varchar(55) DEFAULT NULL,
  `education` varchar(55) DEFAULT NULL,
  `title` varchar(55) DEFAULT NULL,
  `birthday` varchar(55) DEFAULT NULL,
  `birthplace` varchar(55) DEFAULT NULL,
  `citizenship` varchar(55) DEFAULT NULL,
  `identification_number` varchar(55) DEFAULT NULL,
  `payroll_bruto` varchar(55) DEFAULT NULL,
  `payroll_contributions` varchar(55) DEFAULT NULL,
  `departments` varchar(55) DEFAULT NULL,
  `bank` varchar(55) DEFAULT NULL,
  `bank_account` varchar(55) DEFAULT NULL,
  `warehouse` varchar(55) DEFAULT NULL,
  `company` varchar(55) DEFAULT NULL,
  `address` varchar(155) DEFAULT NULL,
  `phone` varchar(55) DEFAULT NULL,
  `email` varchar(55) DEFAULT NULL,
  `description` varchar(155) DEFAULT NULL,
  `contract` varchar(55) DEFAULT NULL,
  `grid_id` varchar(30) NOT NULL,
  `Parent` varchar(30) DEFAULT NULL,
  `Def` varchar(10) DEFAULT NULL,
  `note` varchar(250) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `vat_no` varchar(50) DEFAULT NULL,
  `postal_code` varchar(50) DEFAULT NULL,
  `country` varchar(50) DEFAULT NULL,
  `state` varchar(50) DEFAULT NULL,
  `city` varchar(50) DEFAULT NULL,
  `posting_groups` varchar(50) DEFAULT NULL,
  `payroll_account` varchar(50) DEFAULT NULL,
  `department` varchar(50) DEFAULT NULL,
  `contract_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`,`grid_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.employees_contracts definition

CREATE TABLE `employees_contracts` (
  `id` int(11) NOT NULL,
  `surname` varchar(55) DEFAULT NULL,
  `name` varchar(55) DEFAULT NULL,
  `middle_name` varchar(55) DEFAULT NULL,
  `position` varchar(55) DEFAULT NULL,
  `proffesion` varchar(55) DEFAULT NULL,
  `education` varchar(55) DEFAULT NULL,
  `title` varchar(55) DEFAULT NULL,
  `birthday` varchar(55) DEFAULT NULL,
  `birthplace` varchar(55) DEFAULT NULL,
  `citizenship` varchar(55) DEFAULT NULL,
  `identification_number` varchar(55) DEFAULT NULL,
  `payroll_bruto` varchar(55) DEFAULT NULL,
  `payroll_contributions` varchar(55) DEFAULT NULL,
  `departments` varchar(55) DEFAULT NULL,
  `bank` varchar(55) DEFAULT NULL,
  `bank_account` varchar(55) DEFAULT NULL,
  `warehouse` varchar(55) DEFAULT NULL,
  `address` varchar(155) DEFAULT NULL,
  `phone` varchar(55) DEFAULT NULL,
  `email` varchar(55) DEFAULT NULL,
  `description` varchar(155) DEFAULT NULL,
  `contract` varchar(55) DEFAULT NULL,
  `grid_id` varchar(30) NOT NULL,
  `Parent` varchar(30) DEFAULT NULL,
  `Def` varchar(10) DEFAULT NULL,
  `note` varchar(250) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `vat_no` varchar(50) DEFAULT NULL,
  `postal_code` varchar(50) DEFAULT NULL,
  `country` varchar(50) DEFAULT NULL,
  `state` varchar(50) DEFAULT NULL,
  `city` varchar(50) DEFAULT NULL,
  `posting_groups` varchar(50) DEFAULT NULL,
  `payroll_account` varchar(50) DEFAULT NULL,
  `department` varchar(50) DEFAULT NULL,
  `contract_uuid` int(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(10) DEFAULT NULL,
  `document_no` varchar(15) DEFAULT NULL,
  `document_date` varchar(10) DEFAULT NULL,
  `posting_date` varchar(10) DEFAULT NULL,
  `valid_from_date` varchar(10) DEFAULT NULL,
  `valid_till_date` varchar(10) DEFAULT NULL,
  `job_position` varchar(50) DEFAULT NULL,
  `salary_structure_type` varchar(20) DEFAULT NULL,
  `wage_type` varchar(20) DEFAULT NULL,
  `schedule_pay` varchar(10) DEFAULT NULL,
  `default_working_hours` varchar(20) DEFAULT NULL,
  `work_entry` varchar(10) DEFAULT NULL,
  `trial_period_end` varchar(10) DEFAULT NULL,
  `hr_responsible_code` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`,`grid_id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `employees_contracts_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.employment_contracts definition

CREATE TABLE `employment_contracts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `contract_no` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.fa_posting_groups definition

CREATE TABLE `fa_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `acquisition_cost_account` int(11) DEFAULT NULL,
  `accum_depreciation_account` int(11) DEFAULT NULL,
  `acq_cost_account_disposal` int(11) DEFAULT NULL,
  `accum_depreciation_account_disposal` int(11) DEFAULT NULL,
  `gains_account_disposal` int(11) DEFAULT NULL,
  `loses_account_disposal` int(11) DEFAULT NULL,
  `maintenance_expense_account` int(11) DEFAULT NULL,
  `acquisition_cost_balance_account` int(11) DEFAULT NULL,
  `depreciation_cost_account` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.financial_institutions definition

CREATE TABLE `financial_institutions` (
  `id` int(11) DEFAULT NULL,
  `grid_id` varchar(30) DEFAULT NULL,
  `Parent` varchar(33) DEFAULT NULL,
  `Def` varchar(11) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `vat_no` varchar(50) DEFAULT NULL,
  `address` varchar(250) DEFAULT NULL,
  `address2` varchar(250) DEFAULT NULL,
  `postal_code` varchar(50) DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `country` varchar(100) DEFAULT NULL,
  `state` varchar(100) DEFAULT NULL,
  `phone` varchar(30) DEFAULT NULL,
  `email` varchar(30) DEFAULT NULL,
  `webpage` varchar(30) DEFAULT NULL,
  `contact` varchar(100) DEFAULT NULL,
  `note` varchar(500) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `fax` varchar(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.financial_institutions_contracts definition

CREATE TABLE `financial_institutions_contracts` (
  `id` int(11) DEFAULT NULL,
  `grid_id` varchar(30) DEFAULT NULL,
  `Parent` varchar(33) DEFAULT NULL,
  `Def` varchar(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `posting_date` varchar(15) DEFAULT NULL,
  `contract_date` varchar(15) DEFAULT NULL,
  `contract_no` varchar(20) DEFAULT NULL,
  `contract_intermediate` varchar(50) DEFAULT NULL,
  `financial_institution_description` varchar(30) DEFAULT NULL,
  `financial_institution_vat_no` varchar(20) DEFAULT NULL,
  `loan_amount` decimal(25,4) DEFAULT NULL,
  `loan_period` int(11) DEFAULT NULL,
  `loan_start_date` varchar(15) DEFAULT NULL,
  `scheduled_payment` decimal(25,4) DEFAULT NULL,
  `extra_payment` decimal(25,4) DEFAULT NULL,
  `note` varchar(500) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `interest_rate` varchar(25) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.general_business_posting_groups definition

CREATE TABLE `general_business_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `default_vat_business_posting_group_uuid` int(11) DEFAULT NULL,
  `auto_insert_default` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.general_journal definition

CREATE TABLE `general_journal` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `debit` decimal(25,4) DEFAULT NULL,
  `credit` decimal(25,4) DEFAULT NULL,
  `net_change` decimal(25,4) DEFAULT NULL,
  `document_uuid` bigint(20) DEFAULT NULL,
  `account_uuid` int(11) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.general_ledger definition

CREATE TABLE `general_ledger` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_id` int(11) DEFAULT NULL,
  `transaction_id` int(11) DEFAULT NULL,
  `account_id` int(11) DEFAULT NULL,
  `value` decimal(25,10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.general_posting_setup definition

CREATE TABLE `general_posting_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sales_account` int(11) DEFAULT NULL,
  `sales_credit_memo_account` int(11) DEFAULT NULL,
  `sales_line_discount_account` int(11) DEFAULT NULL,
  `sales_invoice_discount_account` int(11) DEFAULT NULL,
  `sales_payment_discount_debit_account` int(11) DEFAULT NULL,
  `sales_payment_discount_credit_account` int(11) DEFAULT NULL,
  `sales_prepayments_account` int(11) DEFAULT NULL,
  `purchase_account` int(11) DEFAULT NULL,
  `purchase_credit_memo_account` int(11) DEFAULT NULL,
  `purchase_line_discount_account` int(11) DEFAULT NULL,
  `purchase_invoice_discount_account` int(11) DEFAULT NULL,
  `purchase_payment_discount_debit_account` int(11) DEFAULT NULL,
  `purchase_payment_discount_credit_account` int(11) DEFAULT NULL,
  `purchase_prepayments_account` int(11) DEFAULT NULL,
  `cogs_account` int(11) DEFAULT NULL,
  `inventory_adjustment_account` int(11) DEFAULT NULL,
  `direct_cost_applied_account` int(11) DEFAULT NULL,
  `overhead_applied_account` int(11) DEFAULT NULL,
  `purchase_variance_account` int(11) DEFAULT NULL,
  `used_ledger` int(11) DEFAULT NULL,
  `general_business_posting_group_id` int(11) DEFAULT NULL,
  `general_product_posting_group_id` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `archived` tinyint(4) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.general_product_posting_groups definition

CREATE TABLE `general_product_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `default_vat_product_posting_group_uuid` int(11) DEFAULT NULL,
  `auto_insert_default` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.gl_ledger definition

CREATE TABLE `gl_ledger` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `document_no` varchar(255) DEFAULT NULL,
  `document_uuid` varchar(255) DEFAULT NULL,
  `document_date` varchar(255) DEFAULT NULL,
  `posting_date` varchar(255) DEFAULT NULL,
  `entry_date` varchar(255) DEFAULT NULL,
  `delivery_date` varchar(255) DEFAULT NULL,
  `total_lcy` decimal(25,4) DEFAULT NULL,
  `responsability_center_uuid` int(11) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.gl_ledger_items definition

CREATE TABLE `gl_ledger_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Parent` int(11) DEFAULT NULL,
  `debit` decimal(25,4) DEFAULT NULL,
  `credit` decimal(25,4) DEFAULT NULL,
  `net` decimal(25,4) DEFAULT NULL,
  `exchange_rate` decimal(25,4) DEFAULT NULL,
  `total_lcy` decimal(25,4) DEFAULT NULL,
  `account_uuid` int(11) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.grids_permission definition

CREATE TABLE `grids_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `description` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `grids` varchar(500) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `grids_permission_code_uindex` (`code`),
  UNIQUE KEY `grids_permission_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `grids_permission_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.grounds_for_termination definition

CREATE TABLE `grounds_for_termination` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.hr_units definition

CREATE TABLE `hr_units` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `quantity` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inbound_flow definition

CREATE TABLE `inbound_flow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `module_id` int(11) DEFAULT NULL,
  `item_id` int(11) DEFAULT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `location_id` int(11) DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `quantity` decimal(15,5) DEFAULT NULL,
  `value` decimal(20,5) DEFAULT NULL,
  `outbound_quantity` decimal(15,5) DEFAULT NULL,
  `outbound_value` decimal(20,5) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.income_outcome definition

CREATE TABLE `income_outcome` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `warehouse_id` int(11) DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `income_quantity` decimal(25,10) DEFAULT NULL,
  `income_value` decimal(25,10) DEFAULT NULL,
  `outcome_quantity` decimal(25,10) DEFAULT NULL,
  `outcome_value` decimal(25,10) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.industry_groups definition

CREATE TABLE `industry_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `industry_groups_code_uindex` (`code`),
  UNIQUE KEY `industry_groups_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `industry_groups_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.insurances definition

CREATE TABLE `insurances` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `effective_date` datetime DEFAULT NULL,
  `expiration_date` datetime DEFAULT NULL,
  `policy_no` varchar(255) DEFAULT NULL,
  `annual_premium` decimal(25,4) DEFAULT NULL,
  `policy_coverage` decimal(25,4) DEFAULT NULL,
  `last_date_modified` datetime DEFAULT NULL,
  `vendor_uuid` int(11) DEFAULT NULL,
  `insurance_type_uuid` int(11) DEFAULT NULL,
  `class_uuid` int(11) DEFAULT NULL,
  `subclass_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.interaction_groups definition

CREATE TABLE `interaction_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `interaction_groups_code_uindex` (`code`),
  UNIQUE KEY `interaction_groups_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `interaction_groups_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.intercompany_partners definition

CREATE TABLE `intercompany_partners` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `company_uuid` int(11) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  `customer_uuid` int(11) DEFAULT NULL,
  `receivables_account` int(11) DEFAULT NULL,
  `vendor_uuid` int(11) DEFAULT NULL,
  `payables_account` int(11) DEFAULT NULL,
  `cost_distribution_lcy` tinyint(4) DEFAULT NULL,
  `outbound_sales_item_type` tinyint(4) DEFAULT NULL,
  `outbound_purchase_item_type` tinyint(4) DEFAULT NULL,
  `autoaccept_transaction` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inventiry_entry_items definition

CREATE TABLE `inventiry_entry_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_tempory` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `inventiry_entry_items_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.inventiry_exit_items definition

CREATE TABLE `inventiry_exit_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_cost` decimal(25,4) DEFAULT NULL,
  `item_cost_quantity` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_tempory` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.inventory definition

CREATE TABLE `inventory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `location_id` int(11) DEFAULT NULL,
  `item_id` int(11) DEFAULT NULL,
  `quantity` decimal(20,5) DEFAULT NULL,
  `value` decimal(20,5) DEFAULT NULL,
  `value_fifo` decimal(20,5) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inventory_counting definition

CREATE TABLE `inventory_counting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `count_frequency` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inventory_entry definition

CREATE TABLE `inventory_entry` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` bigint(20) NOT NULL,
  `document_date` datetime DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `entry_date` datetime DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `warehouseman_approve` tinyint(1) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.inventory_exit definition

CREATE TABLE `inventory_exit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` datetime DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `entry_date` datetime DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `warehouseman_approve` tinyint(1) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.inventory_location_posting_setup definition

CREATE TABLE `inventory_location_posting_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `location_code` varchar(20) DEFAULT NULL,
  `invt_posting_group_code` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inventory_posting_groups definition

CREATE TABLE `inventory_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `used_items` int(11) DEFAULT NULL,
  `used_values` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.inventory_posting_setup definition

CREATE TABLE `inventory_posting_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `inventory_account` int(11) DEFAULT NULL,
  `inventory_account_interim` int(11) DEFAULT NULL,
  `wip_account` int(11) DEFAULT NULL,
  `material_variance_account` int(11) DEFAULT NULL,
  `capacity_variance_account` int(11) DEFAULT NULL,
  `subcontracted_variance_account` int(11) DEFAULT NULL,
  `cap_overhead_variance_account` int(11) DEFAULT NULL,
  `mfg_overhead_variance_account` int(11) DEFAULT NULL,
  `used_ledger` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `inventory_posting_group_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.item_cost_tracking definition

CREATE TABLE `item_cost_tracking` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_id` int(11) DEFAULT NULL,
  `item_transaction_id` int(11) DEFAULT NULL,
  `inbound_id` int(11) DEFAULT NULL,
  `outbound_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.item_trackings definition

CREATE TABLE `item_trackings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.item_types definition

CREATE TABLE `item_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.items definition

CREATE TABLE `items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `abbrevation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `costing_method` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code_2` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `low_quantity` decimal(25,4) DEFAULT NULL,
  `type_uuid` int(11) DEFAULT NULL,
  `brand_uuid` int(11) DEFAULT NULL,
  `category_uuid` int(11) DEFAULT NULL,
  `subcategory_uuid` int(11) DEFAULT NULL,
  `unit_uuid` int(11) DEFAULT NULL,
  `vat_uuid` int(11) DEFAULT NULL,
  `storage_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `posting_group_uuid` int(11) DEFAULT NULL,
  `unit_volume` decimal(25,4) DEFAULT NULL,
  `sales_unit_uuid` int(11) DEFAULT NULL,
  `shelf_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `net_weight` decimal(25,4) DEFAULT NULL,
  `gross_weight` decimal(25,4) DEFAULT NULL,
  `gen_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `inventory_posting_group` int(11) DEFAULT NULL,
  `tariff_uuid` int(11) DEFAULT NULL,
  `origin_code` int(11) DEFAULT NULL,
  `price_profit_calculation` tinyint(4) DEFAULT NULL,
  `allow_invoice_discount` tinyint(4) DEFAULT NULL,
  `inventory_effect` tinyint(1) DEFAULT NULL,
  `item_discount_group_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `items_type_uuid` (`type_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.items_price definition

CREATE TABLE `items_price` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) NOT NULL,
  `Parent` varchar(10) DEFAULT NULL,
  `Def` varchar(5) DEFAULT NULL,
  `document_type` varchar(50) DEFAULT NULL,
  `document_abbrevation` varchar(30) DEFAULT NULL,
  `document_date` varchar(30) DEFAULT NULL,
  `document_no` varchar(15) DEFAULT NULL,
  `item_code` varchar(15) DEFAULT NULL,
  `item_price` varchar(20) DEFAULT NULL,
  `responsible_code` varchar(30) DEFAULT NULL,
  `responsible_status` tinyint(1) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `valid_from_date` varchar(10) DEFAULT NULL,
  `valid_till_date` varchar(10) DEFAULT NULL,
  UNIQUE KEY `warehouses_items_price_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `warehouses_items_price_pk` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.items_price_items definition

CREATE TABLE `items_price_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) NOT NULL,
  `Parent` varchar(10) DEFAULT NULL,
  `Def` varchar(5) DEFAULT NULL,
  `item_description` varchar(50) DEFAULT NULL,
  `item_code` varchar(15) DEFAULT NULL,
  `item_barcode` varchar(30) DEFAULT NULL,
  `item_price` decimal(25,4) DEFAULT NULL,
  `responsible_code` varchar(15) DEFAULT NULL,
  `responsible_status` tinyint(1) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  UNIQUE KEY `warehouses_items_price_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `warehouses_items_price_pk` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.job_posting_groups definition

CREATE TABLE `job_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `wip_accounts` int(11) DEFAULT NULL,
  `wip_accured_costs_account` int(11) DEFAULT NULL,
  `job_costs_applied_account` int(11) DEFAULT NULL,
  `item_costs_applied_account` int(11) DEFAULT NULL,
  `resource_costs_applied_account` int(11) DEFAULT NULL,
  `gl_costs_applied_account` int(11) DEFAULT NULL,
  `job_costs_adjustmet_account` int(11) DEFAULT NULL,
  `gl_expense_account` int(11) DEFAULT NULL,
  `wip_accured_sales_account` int(11) DEFAULT NULL,
  `wip_invoiced_sales_account` int(11) DEFAULT NULL,
  `job_sales_adjustment_account` int(11) DEFAULT NULL,
  `recognized_costs_account` int(11) DEFAULT NULL,
  `recognized_sales_account` int(11) DEFAULT NULL,
  `job_sales_applied_account` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.job_responsibilities definition

CREATE TABLE `job_responsibilities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `job_responsibilities_code_uindex` (`code`),
  UNIQUE KEY `job_responsibilities_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `job_responsibilities_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.language_items definition

CREATE TABLE `language_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.languages definition

CREATE TABLE `languages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `country` varchar(30) DEFAULT NULL,
  `language` varchar(40) DEFAULT NULL,
  `two_letters` varchar(10) DEFAULT NULL,
  `three_letters` varchar(10) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.ledger_accounts definition

CREATE TABLE `ledger_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` varchar(11) DEFAULT NULL,
  `Parent` varchar(20) DEFAULT NULL,
  `Def` varchar(10) DEFAULT NULL,
  `code` varchar(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `statement` varchar(50) DEFAULT NULL,
  `head` varchar(50) DEFAULT NULL,
  `category` varchar(50) DEFAULT NULL,
  `subcategory` varchar(50) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `currencies_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;


-- bynar.loan_management definition

CREATE TABLE `loan_management` (
  `id` int(11) DEFAULT NULL,
  `grid_id` varchar(30) DEFAULT NULL,
  `Parent` varchar(33) DEFAULT NULL,
  `Def` varchar(11) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `posting_date` varchar(15) DEFAULT NULL,
  `contract_date` varchar(15) DEFAULT NULL,
  `contract_no` varchar(20) DEFAULT NULL,
  `contract_intermediate` varchar(50) DEFAULT NULL,
  `financial_institution_description` varchar(30) DEFAULT NULL,
  `financial_institution_vat_no` varchar(20) DEFAULT NULL,
  `loan_amount` decimal(25,4) DEFAULT NULL,
  `loan_period` int(11) DEFAULT NULL,
  `loan_start_date` varchar(15) DEFAULT NULL,
  `scheduled_payment` decimal(25,4) DEFAULT NULL,
  `extra_payment` decimal(25,4) DEFAULT NULL,
  `note` varchar(500) DEFAULT NULL,
  `uuid` varchar(50) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.location_group_items definition

CREATE TABLE `location_group_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `location_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.locations definition

CREATE TABLE `locations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `postal_code` varchar(20) DEFAULT NULL,
  `city` varchar(30) DEFAULT NULL,
  `country` varchar(30) DEFAULT NULL,
  `state` varchar(30) DEFAULT NULL,
  `continent` varchar(30) DEFAULT NULL,
  `time_zone` varchar(20) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `locations_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.locations_group definition

CREATE TABLE `locations_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.mailing_groups definition

CREATE TABLE `mailing_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.main_account_categories definition

CREATE TABLE `main_account_categories` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) DEFAULT NULL,
  `main_account_category` varchar(30) DEFAULT NULL,
  `description` varchar(30) DEFAULT NULL,
  `main_account_type` varchar(30) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `closed` tinyint(1) DEFAULT NULL,
  UNIQUE KEY `main_account_categories_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.main_accounts definition

CREATE TABLE `main_accounts` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(4) DEFAULT NULL,
  `main_account` varchar(20) DEFAULT NULL,
  `description` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `main_accounts_grid_id_uindex` (`grid_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.misc_articles definition

CREATE TABLE `misc_articles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.number_series definition

CREATE TABLE `number_series` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `starting_date` date DEFAULT NULL,
  `starting_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `ending_no` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `number_series_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `number_series_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.opportunities definition

CREATE TABLE `opportunities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `salesperson_uuid` int(11) DEFAULT NULL,
  `sales_cycle_uuid` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.organisation_accounts definition

CREATE TABLE `organisation_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `seller_id` int(11) DEFAULT NULL,
  `account_uuid` bigint(20) DEFAULT NULL,
  `user_group` varchar(255) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `address_2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `postal_code` int(11) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `organisation_name` int(11) DEFAULT NULL,
  `organisation_uuid` bigint(20) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `closed` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.organisations definition

CREATE TABLE `organisations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `seller_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.organizational_levels definition

CREATE TABLE `organizational_levels` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `organizational_levels_code_uindex` (`code`),
  UNIQUE KEY `organizational_levels_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `organizational_levels_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.organizations definition

CREATE TABLE `organizations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  `user_group_int` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.outbound_flow definition

CREATE TABLE `outbound_flow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `module_id` int(11) DEFAULT NULL,
  `location_id` int(11) DEFAULT NULL,
  `item_id` int(11) DEFAULT NULL,
  `parent_id` int(11) DEFAULT NULL,
  `transaction_id` int(11) DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `quantity` decimal(15,5) DEFAULT NULL,
  `value_avco` decimal(20,5) DEFAULT NULL,
  `value_fifo` decimal(20,5) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.payment definition

CREATE TABLE `payment` (
  `id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.payment_method definition

CREATE TABLE `payment_method` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `balance_account_type_uuid` int(11) DEFAULT NULL,
  `balance_account_uuid` int(11) DEFAULT NULL,
  `debit_term_uuid` int(11) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.payment_terms definition

CREATE TABLE `payment_terms` (
  `id` int(11) NOT NULL,
  `code` varchar(15) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `days` int(11) DEFAULT NULL,
  `discounts` decimal(25,4) DEFAULT NULL,
  `discount_date_calculation` varchar(255) DEFAULT NULL,
  `cash_payment` tinyint(1) DEFAULT NULL,
  `default_term` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.payments definition

CREATE TABLE `payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_id` int(11) DEFAULT NULL,
  `document_no` varchar(255) DEFAULT NULL,
  `transaction_no` bigint(20) DEFAULT NULL,
  `store_id` int(11) DEFAULT NULL,
  `document_date` datetime DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `project_code` int(11) DEFAULT NULL,
  `department_code` int(11) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `batch_id` int(11) DEFAULT NULL,
  `account_type` int(11) DEFAULT NULL,
  `account_no` int(11) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `amount` decimal(20,5) DEFAULT NULL,
  `amount_lcy` decimal(20,5) DEFAULT NULL,
  `applies_to_document_type` int(11) DEFAULT NULL,
  `applies_to_document_no` int(11) DEFAULT NULL,
  `balance_account_type` int(11) DEFAULT NULL,
  `balance_account_number` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.payments_bills_old definition

CREATE TABLE `payments_bills_old` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `purchase_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `payments_bills_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.payments_old definition

CREATE TABLE `payments_old` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_currency_exchange_rate` decimal(25,4) DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `bank_uuid` int(11) DEFAULT NULL,
  `account_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `cashier_approve` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.payroll_items definition

CREATE TABLE `payroll_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Parent` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `id_number` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `employee_code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contract_no` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `working_hours` decimal(10,4) DEFAULT NULL,
  `working_days` int(11) DEFAULT NULL,
  `gross_salary` decimal(25,4) DEFAULT NULL,
  `taxable_salary` decimal(25,4) DEFAULT NULL,
  `social_insurance_company` decimal(25,4) DEFAULT NULL,
  `social_insurance_employee` decimal(25,4) DEFAULT NULL,
  `extra_contribution` decimal(25,4) DEFAULT NULL,
  `health_insurance_company` decimal(25,4) DEFAULT NULL,
  `health_insurance_employee` decimal(25,4) DEFAULT NULL,
  `personal_income_tax` decimal(25,4) DEFAULT NULL,
  `net_salary` decimal(25,4) DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `payroll_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `payroll_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.payrolls definition

CREATE TABLE `payrolls` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Parent` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbrevation` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_no` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `tax_period` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contract_no` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `working_hours` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `working_days` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `gross_salary` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `taxable_salary` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `social_insurance_company` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `social_insurance_employee` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `extra_contribution` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `health_insurance_company` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `health_insurance_employee` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `personal_income_tax` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `net_salary` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accountant_code` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `accountant_approve` tinyint(1) DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `payroll_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `payroll_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.personalized_transactions definition

CREATE TABLE `personalized_transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `subtransaction` varchar(20) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `document_type` varchar(20) DEFAULT NULL,
  `store_code` varchar(20) DEFAULT NULL,
  `warehouse_code` varchar(20) DEFAULT NULL,
  `account_type` varchar(20) DEFAULT NULL,
  `currency_code` varchar(20) DEFAULT NULL,
  `budget_code` varchar(20) DEFAULT NULL,
  `project_code` varchar(20) DEFAULT NULL,
  `workspace_code` varchar(20) DEFAULT NULL,
  `user_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `personalized_transactions_code_uindex` (`code`),
  UNIQUE KEY `personalized_transactions_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.policies definition

CREATE TABLE `policies` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `analytics` tinyint(4) DEFAULT NULL,
  `user_list` tinyint(4) DEFAULT NULL,
  `sales` tinyint(4) DEFAULT NULL,
  `sales_data` tinyint(4) DEFAULT NULL,
  `procurements` tinyint(4) DEFAULT NULL,
  `procurements_data` tinyint(4) DEFAULT NULL,
  `transfers` tinyint(4) DEFAULT NULL,
  `transfers_data` tinyint(4) DEFAULT NULL,
  `payments` tinyint(4) DEFAULT NULL,
  `payment_data` tinyint(4) DEFAULT NULL,
  `user_groups` tinyint(4) DEFAULT NULL,
  `organizations` tinyint(4) DEFAULT NULL,
  `organization_id` int(11) DEFAULT NULL,
  `org_id` binary(16) DEFAULT NULL,
  `organizations_data` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.posting_accounts definition

CREATE TABLE `posting_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `description` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.posting_groups definition

CREATE TABLE `posting_groups` (
  `id` int(11) NOT NULL,
  `grid_id` varchar(30) NOT NULL,
  `Parent` varchar(33) DEFAULT NULL,
  `Def` varchar(11) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `description` varchar(150) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `vendor` varchar(50) DEFAULT NULL,
  `customer` varchar(50) DEFAULT NULL,
  `bank` varchar(50) DEFAULT NULL,
  `inventory` varchar(50) DEFAULT NULL,
  `depreciation` varchar(50) DEFAULT NULL,
  `tax` varchar(50) DEFAULT NULL,
  `sales_account` varchar(50) DEFAULT NULL,
  `purchase_account` varchar(50) DEFAULT NULL,
  `production` varchar(50) DEFAULT NULL,
  `adjustment_plus` varchar(50) DEFAULT NULL,
  `adjustment_minus` varchar(50) DEFAULT NULL,
  `asset_depreciation` varchar(50) DEFAULT NULL,
  `sales_credit_memo` varchar(50) DEFAULT NULL,
  `purchase_credit_memo` varchar(50) DEFAULT NULL,
  `receipt` varchar(50) DEFAULT NULL,
  `payment` varchar(50) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `type` varchar(50) DEFAULT NULL,
  `no` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`,`grid_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.posting_profiles definition

CREATE TABLE `posting_profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `posting_group_vendor` varchar(20) DEFAULT NULL,
  `vat_posting_group_vendor` varchar(20) DEFAULT NULL,
  `gen_posting_group_vendor` varchar(20) DEFAULT NULL,
  `posting_group_customer` varchar(20) DEFAULT NULL,
  `vat_posting_group_customer` varchar(20) DEFAULT NULL,
  `gen_posting_group_customer` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.preferences definition

CREATE TABLE `preferences` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `description` varchar(20) DEFAULT NULL,
  `address_code` varchar(20) DEFAULT NULL,
  `contact_code` varchar(20) DEFAULT NULL,
  `delivery_code` varchar(20) DEFAULT NULL,
  `shipment_code` varchar(20) DEFAULT NULL,
  `payment_code` varchar(20) DEFAULT NULL,
  `bank_code` varchar(20) DEFAULT NULL,
  `currency_code` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.procurements definition

CREATE TABLE `procurements` (
  `id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.production_bom definition

CREATE TABLE `production_bom` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `unit_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.production_orders definition

CREATE TABLE `production_orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `recipe_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `warehouseman_approve` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `items_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.production_orders_items definition

CREATE TABLE `production_orders_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `raw_material_quantity` decimal(25,4) DEFAULT NULL,
  `temp_raw_material_quantity` decimal(25,4) NOT NULL,
  `product_quantity` decimal(25,4) DEFAULT NULL,
  `item_cost` decimal(25,4) DEFAULT NULL,
  `item_cost_quantity` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_inventory_effect` int(11) NOT NULL,
  `item_tempory` int(11) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `production_orders_items_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.purchases definition

CREATE TABLE `purchases` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_currency_exchange_rate` decimal(25,4) DEFAULT NULL,
  `item_unit_value` decimal(25,4) DEFAULT NULL,
  `item_unit_discount` decimal(25,4) DEFAULT NULL,
  `item_unit_total` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount` decimal(25,4) DEFAULT NULL,
  `item_quantity_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total` decimal(25,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `payments_doc` int(11) DEFAULT NULL,
  `payment_due_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_priority` int(11) DEFAULT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `remaining` decimal(25,4) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `account_uuid` int(11) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  `budget_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `document_manager_approval` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.purchases_items definition

CREATE TABLE `purchases_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `item_value` decimal(25,4) DEFAULT NULL,
  `item_unit_value` decimal(25,4) DEFAULT NULL,
  `item_unit_discount` decimal(25,4) DEFAULT NULL,
  `item_unit_tax` decimal(25,4) DEFAULT NULL,
  `item_unit_total` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total` decimal(25,4) DEFAULT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount` decimal(25,4) DEFAULT NULL,
  `item_quantity_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total` decimal(25,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_discount_uuid` int(11) DEFAULT NULL,
  `item_tax_uuid` int(11) DEFAULT NULL,
  `item_inventory_effect` tinyint(1) NOT NULL,
  `item_tempory` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `recipes_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.put_away definition

CREATE TABLE `put_away` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.qualifications definition

CREATE TABLE `qualifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.questionnaire_setup definition

CREATE TABLE `questionnaire_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `contact_type` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `business_relation_uuid` int(11) DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `questionnaire_setup_code_uindex` (`code`),
  UNIQUE KEY `questionnaire_setup_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `questionnaire_setup_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.reason_codes definition

CREATE TABLE `reason_codes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `reason_codes_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `reason_codes_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.receipts definition

CREATE TABLE `receipts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `document_date` datetime DEFAULT NULL,
  `entry_date` datetime DEFAULT NULL,
  `delivery_date` datetime DEFAULT NULL,
  `document_currency_exchange_rate` decimal(25,4) DEFAULT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `bank_uuid` int(11) DEFAULT NULL,
  `account_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `cashier_approve` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.receipts_bills definition

CREATE TABLE `receipts_bills` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `sales_uuid` int(11) DEFAULT NULL,
  `item_tempory` int(11) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `receipts_bills_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.recipes definition

CREATE TABLE `recipes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.recipes_items definition

CREATE TABLE `recipes_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `raw_material_quantity` decimal(25,4) DEFAULT NULL,
  `product_quantity` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_inventory_effect` int(11) NOT NULL,
  `item_tempory` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `recipes_items_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.relatives definition

CREATE TABLE `relatives` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.reminder_terms definition

CREATE TABLE `reminder_terms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.resource_groups definition

CREATE TABLE `resource_groups` (
  `id` int(11) NOT NULL,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `resource_groups_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.resources definition

CREATE TABLE `resources` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `type` tinyint(4) DEFAULT NULL,
  `unit_uuid` int(11) DEFAULT NULL,
  `group_uuid` int(11) DEFAULT NULL,
  `last_date_modified` datetime DEFAULT NULL,
  `use_time_sheet` tinyint(4) DEFAULT NULL,
  `unit_cost` decimal(25,4) DEFAULT NULL,
  `indirect_cost` decimal(25,4) DEFAULT NULL,
  `price_profit_calculation` tinyint(4) DEFAULT NULL,
  `profit` decimal(25,4) DEFAULT NULL,
  `unit_price` decimal(25,4) DEFAULT NULL,
  `gen_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_prod_posting_group` int(11) DEFAULT NULL,
  `deferral_templaye_uuid` int(11) DEFAULT NULL,
  `ic_partner_purchase_accunt` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contract_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.responsibility_center definition

CREATE TABLE `responsibility_center` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `workspace` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `responsibility_center_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.responsibility_center_items definition

CREATE TABLE `responsibility_center_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `user_code` varchar(20) DEFAULT NULL,
  `user_description` varchar(20) DEFAULT NULL,
  `user_department` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `responsibility_center_items_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.sales definition

CREATE TABLE `sales` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_currency_exchange_rate` decimal(25,4) DEFAULT NULL,
  `item_cost` decimal(25,4) DEFAULT NULL,
  `item_cost_quantity` decimal(25,4) DEFAULT NULL,
  `item_unit_value` decimal(25,4) DEFAULT NULL,
  `item_unit_discount` decimal(25,4) DEFAULT NULL,
  `item_unit_total` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount` decimal(25,4) DEFAULT NULL,
  `item_quantity_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total` decimal(25,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `paid` decimal(25,4) DEFAULT NULL,
  `remaining` decimal(25,4) DEFAULT NULL,
  `payments_doc` int(11) DEFAULT NULL,
  `payment_due_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_priority` int(11) DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_uuid` int(11) DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `account_uuid` int(11) DEFAULT NULL,
  `currency_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `document_manager_approval` tinyint(1) DEFAULT NULL,
  `sales_order_status` tinyint(1) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.sales_invoice_details definition

CREATE TABLE `sales_invoice_details` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `vat_uuid` int(11) DEFAULT NULL,
  `company_type_uuid` int(11) DEFAULT NULL,
  `price_group_uuid` int(11) DEFAULT NULL,
  `discount_group_uuid` int(11) DEFAULT NULL,
  `general_business_posting_group_uuid` int(11) DEFAULT NULL,
  `customer_posting_group_uuid` int(11) DEFAULT NULL,
  `transaction_type_uuid` int(11) DEFAULT NULL,
  `department_uuid` int(11) DEFAULT NULL,
  `project_uuid` int(11) DEFAULT NULL,
  `budget_uuid` int(11) DEFAULT NULL,
  `customer_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.sales_items definition

CREATE TABLE `sales_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `item_cost` decimal(25,4) DEFAULT NULL,
  `item_cost_quantity` decimal(25,4) DEFAULT NULL,
  `item_value` decimal(25,4) DEFAULT NULL,
  `item_unit_value` decimal(25,4) DEFAULT NULL,
  `item_unit_discount` decimal(25,4) DEFAULT NULL,
  `item_unit_tax` decimal(25,4) DEFAULT NULL,
  `item_unit_total` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total` decimal(25,4) DEFAULT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount` decimal(25,4) DEFAULT NULL,
  `item_quantity_total` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total` decimal(25,4) DEFAULT NULL,
  `item_unit_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_unit_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_discount_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_tax_value_base_currency` decimal(25,4) DEFAULT NULL,
  `item_quantity_grand_total_base_currency` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_discount_uuid` int(11) DEFAULT NULL,
  `item_tax_uuid` int(11) DEFAULT NULL,
  `item_inventory_effect` tinyint(1) NOT NULL,
  `item_tempory` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `recipes_Parent_index` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.salesperson_purchaser definition

CREATE TABLE `salesperson_purchaser` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `commission` decimal(25,4) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.salutations definition

CREATE TABLE `salutations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `salutations_code_uindex` (`code`),
  UNIQUE KEY `salutations_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `salutations_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.sections definition

CREATE TABLE `sections` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(20) DEFAULT NULL,
  `code_2` varchar(50) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `storage_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `sections_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `sections_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.segments definition

CREATE TABLE `segments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `salesperson_uuid` int(11) DEFAULT NULL,
  `date` datetime DEFAULT NULL,
  `interaction_uuid` int(11) DEFAULT NULL,
  `information_flow` int(11) DEFAULT NULL,
  `initiation` int(11) DEFAULT NULL,
  `unit_cost` decimal(25,4) DEFAULT NULL,
  `unit_duration` decimal(25,4) DEFAULT NULL,
  `campain_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.sellers definition

CREATE TABLE `sellers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `vat_number` varchar(255) DEFAULT NULL,
  `country` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.serial_numbers definition

CREATE TABLE `serial_numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `starting` varchar(255) DEFAULT NULL,
  `ending` varchar(255) DEFAULT NULL,
  `last_number_used` varchar(255) DEFAULT NULL,
  `last_date_used` datetime DEFAULT NULL,
  `default_nos` tinyint(4) DEFAULT NULL,
  `manual_nos` tinyint(4) DEFAULT NULL,
  `date_order` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.series definition

CREATE TABLE `series` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `default_nos` tinyint(4) DEFAULT NULL,
  `manual_nos` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.series_items definition

CREATE TABLE `series_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `starting_no` varchar(255) DEFAULT NULL,
  `increment_no` int(11) DEFAULT NULL,
  `last_date_used` date DEFAULT NULL,
  `last_no_used` varchar(255) DEFAULT NULL,
  `warning_no` varchar(255) DEFAULT NULL,
  `ending_no` varchar(255) DEFAULT NULL,
  `open` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.service_zones definition

CREATE TABLE `service_zones` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.settings_workflow definition

CREATE TABLE `settings_workflow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `dimension_id` int(11) DEFAULT NULL,
  `flow_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.shipment_methods definition

CREATE TABLE `shipment_methods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.shipping_agent_services definition

CREATE TABLE `shipping_agent_services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `shipping_time` int(11) DEFAULT NULL,
  `base_calendar_uuid` int(11) DEFAULT NULL,
  `customized_calendar` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.shipping_agents definition

CREATE TABLE `shipping_agents` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `tracking_address` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.shipping_companies definition

CREATE TABLE `shipping_companies` (
  `id` int(11) DEFAULT NULL,
  `description` varchar(55) DEFAULT NULL,
  `vat_no` varchar(55) DEFAULT NULL,
  `state` varchar(55) DEFAULT NULL,
  `city` varchar(55) DEFAULT NULL,
  `phone` varchar(55) DEFAULT NULL,
  `email` varchar(55) DEFAULT NULL,
  `method` varchar(55) DEFAULT NULL,
  `terms` varchar(55) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `address` varchar(250) DEFAULT NULL,
  `Def` varchar(10) DEFAULT NULL,
  `Parent` varchar(11) DEFAULT NULL,
  `grid_id` varchar(11) DEFAULT NULL,
  `country` varchar(55) DEFAULT NULL,
  `postal_code` varchar(30) DEFAULT NULL,
  `address2` varchar(250) DEFAULT NULL,
  `contacts` varchar(55) DEFAULT NULL,
  `website` varchar(50) DEFAULT NULL,
  `uuid` varchar(55) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `contract` varchar(50) DEFAULT NULL,
  `contract_uuid` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.shippings definition

CREATE TABLE `shippings` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `shipping_address` varchar(255) DEFAULT NULL,
  `billing_address` int(11) DEFAULT NULL,
  `base_calendar_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `shipping_agent_uuid` int(11) DEFAULT NULL,
  `company_uuid` int(11) DEFAULT NULL,
  `shipment_method_uuid` int(11) DEFAULT NULL,
  `billing` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.sites definition

CREATE TABLE `sites` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `subsidiaries_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `sites_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.source_codes definition

CREATE TABLE `source_codes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `Def` varchar(5) COLLATE utf8_unicode_ci DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `source_codes_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `source_codes_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.special_equipments definition

CREATE TABLE `special_equipments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.statuses definition

CREATE TABLE `statuses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.storages definition

CREATE TABLE `storages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `transaction_code` varchar(255) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.store_group_items definition

CREATE TABLE `store_group_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `location_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.store_groups definition

CREATE TABLE `store_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.stores definition

CREATE TABLE `stores` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `catalogue_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `location_group_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.sub_transactions definition

CREATE TABLE `sub_transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_type` varchar(20) DEFAULT NULL,
  `subtransactions` varchar(10) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `store_code` varchar(20) DEFAULT NULL,
  `warehouse_code` varchar(20) DEFAULT NULL,
  `currency_code` varchar(20) DEFAULT NULL,
  `budget_code` varchar(20) DEFAULT NULL,
  `workspace_code` varchar(20) DEFAULT NULL,
  `user_code` varchar(250) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `subtransactions_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.subcategories definition

CREATE TABLE `subcategories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `category_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.subclasses definition

CREATE TABLE `subclasses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `class_uuid` int(11) DEFAULT NULL,
  `posting_group_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.subsidiaries definition

CREATE TABLE `subsidiaries` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.subtransactions definition

CREATE TABLE `subtransactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `subtransactions` varchar(15) DEFAULT NULL,
  `transaction_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `subtransactions_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.tariff_numbers definition

CREATE TABLE `tariff_numbers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `supplementary_units` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.tenant_management definition

CREATE TABLE `tenant_management` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `organisation_id` int(11) DEFAULT NULL,
  `application_id` int(11) DEFAULT NULL,
  `tenant_uuid` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transaction_approvals definition

CREATE TABLE `transaction_approvals` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `module_id` tinyint(4) DEFAULT NULL,
  `transaction_id` int(11) DEFAULT NULL,
  `approval_id` int(11) DEFAULT NULL,
  `approval_status` varchar(255) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transaction_documents definition

CREATE TABLE `transaction_documents` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `subtransaction` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `workspace_code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  UNIQUE KEY `documents_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.transaction_specifications definition

CREATE TABLE `transaction_specifications` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transaction_types definition

CREATE TABLE `transaction_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transactions definition

CREATE TABLE `transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_code` varchar(5) DEFAULT NULL,
  `operation_code` varchar(5) DEFAULT NULL,
  `document_code` varchar(5) DEFAULT NULL,
  `transaction_code` varchar(15) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `transactions_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transfer definition

CREATE TABLE `transfer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_id` int(11) DEFAULT NULL,
  `document_no` varchar(255) DEFAULT NULL,
  `transaction_no` int(11) DEFAULT NULL,
  `store_id` int(11) DEFAULT NULL,
  `document_date` datetime DEFAULT NULL,
  `posting_date` datetime DEFAULT NULL,
  `entry_date` datetime DEFAULT NULL,
  `shipment_date` datetime DEFAULT NULL,
  `project_id` int(11) DEFAULT NULL,
  `department_id` int(11) DEFAULT NULL,
  `in_transit_id` int(11) DEFAULT NULL,
  `shipment_method_id` int(11) DEFAULT NULL,
  `shipping_agent_id` int(11) DEFAULT NULL,
  `shipping_agent_service_id` int(11) DEFAULT NULL,
  `transaction_type_id` int(11) DEFAULT NULL,
  `transaction_specification_id` int(11) DEFAULT NULL,
  `area_id` int(11) DEFAULT NULL,
  `entry_exit_point_id` int(11) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  `location_origin_id` int(11) DEFAULT NULL,
  `location_destination_id` int(11) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transfer_lines definition

CREATE TABLE `transfer_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `item_id` int(11) DEFAULT NULL,
  `input_quantity` decimal(15,5) DEFAULT NULL,
  `item_unit_value` decimal(15,5) DEFAULT NULL,
  `quantity` decimal(15,5) DEFAULT NULL,
  `item_unit_id` int(11) DEFAULT NULL,
  `shipment_date` datetime DEFAULT NULL,
  `receipt_date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.transfers definition

CREATE TABLE `transfers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `delivery_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `entry_date` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_type_uuid` int(11) DEFAULT NULL,
  `store_origin_uuid` int(11) DEFAULT NULL,
  `store_destination_uuid` int(11) DEFAULT NULL,
  `warehouse_origin_uuid` int(11) DEFAULT NULL,
  `warehouse_destination_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `warehouseman_destination_approve` int(11) DEFAULT NULL,
  `has_child` tinyint(1) DEFAULT NULL,
  `group_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `transfers_document_no` (`document_no`),
  KEY `transfers_document_type_uuid` (`document_type_uuid`),
  KEY `transfers_responsibility_center_uuidd` (`warehouse_destination_uuid`),
  KEY `transfers_store_destination_uuid` (`store_destination_uuid`),
  KEY `transfers_store_origin_uuid` (`store_origin_uuid`),
  KEY `transfers_warehouse_destination_uuid` (`warehouse_destination_uuid`),
  KEY `transfers_warehouse_origin_uuid` (`warehouse_origin_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.transfers_items definition

CREATE TABLE `transfers_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Parent` int(11) NOT NULL,
  `input_quantity` decimal(25,4) DEFAULT NULL,
  `item_quantity_unit` decimal(25,4) DEFAULT NULL,
  `item_quantity` decimal(25,4) DEFAULT NULL,
  `item_cost` decimal(25,4) DEFAULT NULL,
  `item_cost_quantity` decimal(25,4) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL,
  `item_unit_uuid` int(11) DEFAULT NULL,
  `item_tempory` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `transfers_items_item_unit_uuid` (`item_unit_uuid`),
  KEY `transfers_items_item_uuid` (`item_uuid`),
  KEY `transfers_items_parent` (`Parent`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.trial_balance definition

CREATE TABLE `trial_balance` (
  `id` int(11) NOT NULL,
  `gl_account` varchar(10) DEFAULT NULL,
  `gl_account_description` varchar(30) DEFAULT NULL,
  `opening_balance` decimal(25,4) DEFAULT NULL,
  `debit` decimal(25,4) DEFAULT NULL,
  `credit` decimal(25,4) DEFAULT NULL,
  `net_change` decimal(25,4) DEFAULT NULL,
  `closing_balance` decimal(25,4) DEFAULT NULL,
  `document_date` varchar(10) DEFAULT NULL,
  `posting_date` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.unions definition

CREATE TABLE `unions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `phone` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.units definition

CREATE TABLE `units` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `base_unit` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `operation_value` decimal(25,4) DEFAULT NULL,
  `unit_value` decimal(25,4) DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `value` decimal(15,5) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.user_group_lines definition

CREATE TABLE `user_group_lines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.user_groups definition

CREATE TABLE `user_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.users definition

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(255) DEFAULT NULL,
  `full_name` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT '0',
  `language_preference` varchar(255) DEFAULT 'en',
  `policy_id` int(11) DEFAULT NULL,
  `theme` varchar(255) DEFAULT NULL,
  `profile_photo` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.users_setup definition

CREATE TABLE `users_setup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `allow_reading_from` datetime DEFAULT NULL,
  `allow_reading_to` datetime DEFAULT NULL,
  `allow_posting_from` datetime DEFAULT NULL,
  `allow_posting_to` datetime DEFAULT NULL,
  `register_time` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vat_business_posting_groups definition

CREATE TABLE `vat_business_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vat_clauses definition

CREATE TABLE `vat_clauses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vat_posting_setup definition

CREATE TABLE `vat_posting_setup` (
  `id` int(11) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `vat_identifier` varchar(255) DEFAULT NULL,
  `vat` decimal(25,4) DEFAULT NULL,
  `sales_vat_account` int(11) DEFAULT NULL,
  `purchase_vat_account` int(11) DEFAULT NULL,
  `reverse_chrg_vat_acc` int(11) DEFAULT NULL,
  `tax_category` int(11) DEFAULT NULL,
  `eu_service` tinyint(1) DEFAULT NULL,
  `used_ledger` int(11) DEFAULT NULL,
  `vat_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_prod_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_calculation_type` varchar(255) DEFAULT NULL,
  `vat_clause_uuid` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vat_product_posting_groups definition

CREATE TABLE `vat_product_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vats definition

CREATE TABLE `vats` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `value` decimal(25,4) DEFAULT NULL,
  `calculation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `postal_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `city` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `country` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `state` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `posting_groups_uuid` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.vendor_bank_accounts definition

CREATE TABLE `vendor_bank_accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `desription` varchar(255) DEFAULT NULL,
  `bank_account_no` varchar(255) DEFAULT NULL,
  `transit_no` varchar(255) DEFAULT NULL,
  `bank_clearing_standart_code` varchar(255) DEFAULT NULL,
  `swift` varchar(255) DEFAULT NULL,
  `iban` varchar(255) DEFAULT NULL,
  `bank_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `bank_clearing_standart_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_groups definition

CREATE TABLE `vendor_groups` (
  `id` int(11) NOT NULL,
  `grid_id` int(11) NOT NULL,
  `Def` varchar(5) DEFAULT NULL,
  `document_no` varchar(15) DEFAULT NULL,
  `code` varchar(15) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `payment_terms` varchar(20) DEFAULT NULL,
  `invoice_payment_date` varchar(30) DEFAULT NULL,
  `tax_group` varchar(15) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`grid_id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_invoicing definition

CREATE TABLE `vendor_invoicing` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `vat_no_uuid` int(11) DEFAULT NULL,
  `gln` int(11) DEFAULT NULL,
  `pay_to_vendor_no_uuid` int(11) DEFAULT NULL,
  `invoice_discount_uuid` int(11) DEFAULT NULL,
  `prices_including_vat` tinyint(4) DEFAULT NULL,
  `gen_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `vat_bus_posting_group_uuid` int(11) DEFAULT NULL,
  `vendor_posting_group_uuid` int(11) DEFAULT NULL,
  `foreign_trade_currency_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_list_items definition

CREATE TABLE `vendor_list_items` (
  `id` int(11) DEFAULT NULL,
  `Parent` int(11) DEFAULT NULL,
  `item_code` varchar(20) DEFAULT NULL,
  `item_uuid` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_lists definition

CREATE TABLE `vendor_lists` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_type` varchar(20) DEFAULT NULL,
  `document_abbrevation` varchar(20) DEFAULT NULL,
  `document_no` varchar(20) DEFAULT NULL,
  `document_date` varchar(20) DEFAULT NULL,
  `vendor_code` varchar(20) DEFAULT NULL,
  `note` varchar(150) DEFAULT NULL,
  `vendor_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_payments definition

CREATE TABLE `vendor_payments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `prepayment` decimal(25,4) DEFAULT NULL,
  `application_method` tinyint(4) DEFAULT NULL,
  `payment_terms_uuid` int(11) DEFAULT NULL,
  `payment_method_uuid` int(11) DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `block_payment_tolerance` int(11) DEFAULT NULL,
  `preferred_bank_account_uuid` int(11) DEFAULT NULL,
  `partner_type` tinyint(4) DEFAULT NULL,
  `cash_flow_payment_terms_uuid` int(11) DEFAULT NULL,
  `creditor_no` int(11) DEFAULT NULL,
  `vendor_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_posting_group definition

CREATE TABLE `vendor_posting_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(20) DEFAULT NULL,
  `description` varchar(50) DEFAULT NULL,
  `payables_account` int(11) DEFAULT NULL,
  `service_charge_acc` int(11) DEFAULT NULL,
  `payment_disc_debit_acc` int(11) DEFAULT NULL,
  `payment_disc_credit_acc` int(11) DEFAULT NULL,
  `debit_curr_appln_rndg_acc` int(11) DEFAULT NULL,
  `credit_curr_appln_rndg_acc` int(11) DEFAULT NULL,
  `debit_rounding_account` int(11) DEFAULT NULL,
  `credit_rounding_account` int(11) DEFAULT NULL,
  `payment_tolerance_debit_acc` int(11) DEFAULT NULL,
  `payment_tolerance_credit_acc` int(11) DEFAULT NULL,
  `invoice_rounding_account` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_posting_groups definition

CREATE TABLE `vendor_posting_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `payable_account` int(11) DEFAULT NULL,
  `service_charge_account` int(11) DEFAULT NULL,
  `payment_disc_debit_account` int(11) DEFAULT NULL,
  `payment_disc_credit_account` int(11) DEFAULT NULL,
  `invoice_rounding_account` int(11) DEFAULT NULL,
  `debit_curr_appln_rndg_account` int(11) DEFAULT NULL,
  `credit_curr_appln_rndg_account` int(11) DEFAULT NULL,
  `debit_rounding_account` int(11) DEFAULT NULL,
  `credit_rounding_account` int(11) DEFAULT NULL,
  `used_vendors` int(11) DEFAULT NULL,
  `used_ledger` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendor_receiving definition

CREATE TABLE `vendor_receiving` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `warehouse_uuid` int(11) DEFAULT NULL,
  `shipment_method_uuid` int(11) DEFAULT NULL,
  `base_calendar_uuid` int(11) DEFAULT NULL,
  `lead_time_calculation` decimal(25,4) DEFAULT NULL,
  `customized_calendar` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.vendors definition

CREATE TABLE `vendors` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `document_no` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `document_abbreviation` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  `receiving_uuid` int(11) DEFAULT NULL,
  `payment_uuid` int(11) DEFAULT NULL,
  `invoicing_uuid` int(11) DEFAULT NULL,
  `document_sending_profile_uuid` int(11) DEFAULT NULL,
  `ic_partner_uuid` int(11) DEFAULT NULL,
  `purchaser_uuid` int(11) DEFAULT NULL,
  `company_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.warehouse_classes definition

CREATE TABLE `warehouse_classes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.warehouses definition

CREATE TABLE `warehouses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `transaction_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `site_uuid` int(11) DEFAULT NULL,
  `address_uuid` int(11) DEFAULT NULL,
  `contact_uuid` int(11) DEFAULT NULL,
  `responsibility_center_uuid` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.warehouses_items_quantity definition

CREATE TABLE `warehouses_items_quantity` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `warehouse_uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_uuid` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `item_quanitity` decimal(25,4) DEFAULT NULL,
  `item_reserved_quantity` decimal(25,4) DEFAULT '0.0000',
  `item_stock_value` decimal(25,4) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `sma_all_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.web_sources definition

CREATE TABLE `web_sources` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `grid_id` int(11) DEFAULT NULL,
  `code` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `url` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `note` varchar(150) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `web_sources_code_uindex` (`code`),
  UNIQUE KEY `web_sources_grid_id_uindex` (`grid_id`),
  UNIQUE KEY `web_sources_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


-- bynar.workflow_items definition

CREATE TABLE `workflow_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) DEFAULT NULL,
  `account_id` bigint(20) DEFAULT NULL,
  `document_id` int(11) DEFAULT NULL,
  `approval_order` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.workflows definition

CREATE TABLE `workflows` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `module_id` int(11) DEFAULT NULL,
  `store_group_id` int(11) DEFAULT NULL,
  `user_group_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- bynar.zip_codes definition

CREATE TABLE `zip_codes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code` varchar(15) DEFAULT NULL,
  `city` varchar(50) DEFAULT NULL,
  `country` varchar(20) DEFAULT NULL,
  `state` varchar(50) DEFAULT NULL,
  `continent` varchar(20) DEFAULT NULL,
  `timezone` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `zip_codes_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;