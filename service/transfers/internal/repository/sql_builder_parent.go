package repository

// tables: transfers, documents, stores
const (
	// QueryParentCount is a query for parent count
	QueryParentCount = `
	SELECT COUNT(transfers.id) as Count 
	FROM transfers 
	INNER JOIN documents ON transfers.document_id = documents.id  
	INNER JOIN stores ON transfers.store_id = stores.id 
	WHERE 1=1 `

	// QueryParent is a query for parent
	QueryParent = `
	SELECT 
		transfers.*
	FROM transfers 
	INNER JOIN documents ON transfers.document_id = documents.id  
	INNER JOIN stores ON transfers.store_id = stores.id 
	WHERE 1=1 `

	// QueryParentJoins is a query for parent joins
	QueryParentJoins = `
INNER JOIN documents ON transfers.document_id = documents.id
INNER JOIN stores ON transfers.store_id = stores.id
`
)
