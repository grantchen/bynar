package repository

// tables: transfers, documents, stores
const (
	QueryParentCount = `
	SELECT COUNT(transfers.id) as Count 
	FROM transfers 
	INNER JOIN documents ON transfers.document_id = documents.id  
	INNER JOIN stores ON transfers.store_id = stores.id 
	WHERE 1=1 `

	QueryParent = `
	SELECT 
		transfers.*
	FROM transfers 
	INNER JOIN documents ON transfers.document_id = documents.id  
	INNER JOIN stores ON transfers.store_id = stores.id 
	WHERE 1=1 `

	QueryParentJoins = `
INNER JOIN documents ON transfers.document_id = documents.id
INNER JOIN stores ON transfers.store_id = stores.id
`
	QueryParentBuild = `
	SELECT * 
	FROM transfers `
)
