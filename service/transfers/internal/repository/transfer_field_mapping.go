package repository

var (
	TreeGridServiceNames = []string{
		"Added",
		"Deleted",
		"Changed",
	}

	TransferFieldNames = map[string][]string{
		"id": {
			"id",
			"transfer.id",
		},
		// "document_id": {
		// 	"document_id",
		// 	"transfer.document_id",
		// },
		"document_no": {
			"document_no",
			"transfer.document_no",
		},
		// "transaction_no": {
		// 	"transaction_no",
		// 	"transfer.transaction_no",
		// },
		"document_date": {
			"document_date",
			"transfer.document_date",
		},
		"posting_date": {
			"posting_date",
			"transfer.posting_date",
		},
		"entry_date": {
			"entry_date",
			"transfer.entry_date",
		},
		// "shipment_date": {
		// 	"shipment_date",
		// 	"transfer.shipment_date",
		// },
		// "project_id": {
		// 	"project_id",
		// 	"transfer.project_id",
		// },
		// "in_transit_id": {
		// 	"in_transit_id",
		// 	"transfer.in_transit_id",
		// },
		// "shipment_method_id": {
		// 	"shipment_method_id",
		// 	"transfer.shipment_method_id",
		// },
		// "shipping_agent_id": {
		// 	"shipping_agent_id",
		// 	"transfer.shipping_agent_id",
		// },
		// "shipping_agent_service_id": {
		// 	"shipping_agent_service_id",
		// 	"transfer.shipping_agent_service_id",
		// },
		// "transaction_type_id": {
		// 	"transaction_type_id",
		// 	"transfer.transaction_type_id",
		// },
		// "transaction_specification_id": {
		// 	"transaction_specification_id",
		// 	"transfer.transaction_specification_id",
		// },
		// "user_group_id": {
		// 	"user_group_id",
		// 	"transfer.user_group_id",
		// },
		// "location_origin_id": {
		// 	"location_origin_id",
		// 	"transfer.location_origin_id",
		// },
		// "location_destination_id": {
		// 	"location_destination_id",
		// 	"transfer.location_destination_id",
		// },
		// "status": {
		// 	"status",
		// 	"transfer.status",
		// },
	}

	TransferLineFieldNames = map[string][]string{
		"id": {
			"transfer_lines.id",
		},
		"Parent": {
			"transfer_lines.parent_id",
		},
		"item_id": {
			"transfer_lines.item_id",
		},
		"input_quantity": {
			"transfer_lines.input_quantity",
		},
		"item_unit_value": {
			"transfer_lines.item_unit_value",
		},
		"quantity": {
			"transfer_lines.quantity",
		},
		"item_unit_id": {
			"transfer_lines.item_unit_id",
		},
		"shipment_date": {
			"transfer_lines.shipment_date",
		},
		"receipt_date": {
			"transfer_lines.receipt_date",
		},
	}
)
