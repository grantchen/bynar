package repository

var (
	// TransferFieldNames is a map of transfer field names
	TransferFieldNames = map[string][]string{
		"id": {
			"id",
			"transfers.id",
		},
		"document_id": {
			"document_id",
			"transfers.document_id",
		},
		"document_no": {
			"document_no",
			"transfers.document_no",
		},
		"transaction_no": {
			"transaction_no",
			"transfers.transaction_no",
		},
		"document_date": {
			"document_date",
			"transfers.document_date",
		},
		"posting_date": {
			"posting_date",
			"transfers.posting_date",
		},
		"entry_date": {
			"entry_date",
			"transfers.entry_date",
		},
		"shipment_date": {
			"shipment_date",
			"transfers.shipment_date",
		},
		"project_id": {
			"project_id",
			"transfers.project_id",
		},
		"area_id": {
			"area_id",
			"transfers.area_id",
		},
		"entry_exit_point_id": {
			"entry_exit_point_id",
			"transfers.entry_exit_point_id",
		},
		"department_id": {
			"department_id",
			"transfers.department_id",
		},
		"in_transit_id": {
			"in_transit_id",
			"transfers.in_transit_id",
		},
		"shipment_method_id": {
			"shipment_method_id",
			"transfers.shipment_method_id",
		},
		"shipping_agent_id": {
			"shipping_agent_id",
			"transfers.shipping_agent_id",
		},
		"shipping_agent_service_id": {
			"shipping_agent_service_id",
			"transfers.shipping_agent_service_id",
		},
		"transaction_type_id": {
			"transaction_type_id",
			"transfers.transaction_type_id",
		},
		"transaction_specification_id": {
			"transaction_specification_id",
			"transfers.transaction_specification_id",
		},
		"user_group_id": {
			"user_group_id",
			"transfers.user_group_id",
		},
		"store_id": {
			"store_id",
			"transfers.store_id",
		},
		"location_origin_id": {
			"location_origin_id",
			"transfers.location_origin_id",
		},
		"location_destination_id": {
			"location_destination_id",
			"transfers.location_destination_id",
		},
		"status": {
			"status",
			"transfers.status",
		},
	}

	// TransferLineFieldNames is a map of transfer line field names
	TransferLineFieldNames = map[string][]string{
		"id-line": {
			"id",
			"transfer_lines.id",
		},
		"Parent": {
			"parent_id",
			"transfer_lines.parent_id",
		},
		"item_id": {
			"item_id",
			"transfer_lines.item_id",
		},
		"input_quantity": {
			"input_quantity",
			"transfer_lines.input_quantity",
		},
		"item_unit_id": {
			"item_unit_id",
			"transfer_lines.item_unit_id",
		},
		"shipment_date": {
			"shipment_date",
			"transfer_lines.shipment_date",
		},
		"receipt_date": {
			"transfer_lines.receipt_date",
			"transfer_lines.receipt_date",
		},

		// need to insert, not in treegrid cols
		"item_unit_value": {
			"item_unit_value",
			"transfer_lines.item_unit_value",
		},
		"quantity": {
			"quantity",
			"transfer_lines.quantity",
		},
	}
)
