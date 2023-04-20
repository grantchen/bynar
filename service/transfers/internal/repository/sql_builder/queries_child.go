package sqlbuilder

const (
	QueryChildCount = `
SELECT COUNT(*) as rowCount 
FROM transfers_items 
INNER JOIN items ON transfers_items.item_uuid = items.id  
INNER JOIN units ON transfers_items.item_unit_uuid = units.id 
INNER JOIN item_types ON items.type_uuid = item_types.id  
where 1=1 `

	QueryChild = `
SELECT 
	transfers_items.id, 
	transfers_items.Parent, 
	item_types.code AS item_type, 
	items.no AS item_no, 
	items.description AS item_name, 
	units.code AS item_unit, 
	transfers_items.input_quantity, 
	transfers_items.item_quantity_unit, 
	transfers_items.item_quantity, 
	transfers_items.item_tempory, 
	transfers_items.item_uuid, 
	transfers_items.item_unit_uuid  
FROM transfers_items 
INNER JOIN items ON transfers_items.item_uuid = items.id  
INNER JOIN units ON transfers_items.item_unit_uuid = units.id 
INNER JOIN item_types ON items.type_uuid = item_types.id
`

	QueryChildJoins = `
INNER JOIN items ON transfers_items.item_uuid = items.id  
INNER JOIN units ON transfers_items.item_unit_uuid = units.id 
INNER JOIN item_types ON items.type_uuid = item_types.id`
)
