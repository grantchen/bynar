package repository

const (
	QueryParentCount = `
	SELECT COUNT(transfers.id) as rowCount 
	FROM transfers 
	INNER JOIN documents ON transfers.document_type_uuid = documents.id  
	INNER JOIN stores ON transfers.store_origin_uuid = stores.id 
	INNER JOIN stores ss ON transfers.store_destination_uuid = ss.id 
	INNER JOIN warehouses wh_origin ON transfers.warehouse_origin_uuid = wh_origin.id 
	INNER JOIN warehouses wh_destination ON transfers.warehouse_destination_uuid = wh_destination.id  
	INNER JOIN responsibility_center ON transfers.responsibility_center_uuid = responsibility_center.id 
	WHERE 1=1 `

	QueryParent = `
	SELECT 
		transfers.id,  
		transfers.document_no, 
		transfers.document_date,  
		transfers.posting_date, 
		transfers.entry_date,  
		transfers.delivery_date, 
		documents.document_type AS document_type,  
		documents.document_abbrevation AS document_abbrevation, 
		stores.code AS store_origin_code, 
		wh_origin.code AS warehouse_origin_code, 
		wh_destination.code AS warehouse_destination_code, 
		ss.code AS store_destination_code,  
		responsibility_center.code AS responsibility_center, 
		transfers.document_type_uuid, 
		transfers.store_origin_uuid, 
		transfers.warehouse_origin_uuid, 
		transfers.warehouse_destination_uuid,  
		transfers.responsibility_center_uuid, 
		transfers.warehouseman_destination_approve, 
		transfers.has_child FROM transfers 
	INNER JOIN documents ON transfers.document_type_uuid = documents.id  
	INNER JOIN stores ON transfers.store_origin_uuid = stores.id 
	INNER JOIN stores ss ON transfers.store_destination_uuid = ss.id  
	INNER JOIN warehouses wh_origin ON transfers.warehouse_origin_uuid = wh_origin.id 
	INNER JOIN warehouses wh_destination ON transfers.warehouse_destination_uuid = wh_destination.id  
	INNER JOIN responsibility_center ON transfers.responsibility_center_uuid = responsibility_center.id 
	WHERE 1=1 `

	QueryParentJoins = `
	INNER JOIN documents ON transfers.document_type_uuid = documents.id  
	INNER JOIN stores ON transfers.store_origin_uuid = stores.id 
	INNER JOIN stores ss ON transfers.store_destination_uuid = ss.id  
	INNER JOIN warehouses wh_origin ON transfers.warehouse_origin_uuid = wh_origin.id 
	INNER JOIN warehouses wh_destination ON transfers.warehouse_destination_uuid = wh_destination.id  
	INNER JOIN responsibility_center ON transfers.responsibility_center_uuid = responsibility_center.id 
	`
	QueryParentBuild = `
	SELECT * 
	FROM transfers `
)
