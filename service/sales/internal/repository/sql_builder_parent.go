package repository

// tables: sales, documents, stores
const (
	// QueryParentCount is a query for getting parent count
	QueryParentCount = `
	SELECT COUNT(sales.id) as Count 
	FROM sales 
	INNER JOIN documents ON sales.document_id = documents.id  
	INNER JOIN stores ON sales.store_id = stores.id 
	WHERE 1=1 `

	// QueryParent is a query for getting parent
	QueryParent = `
	SELECT 
		sales.*
	FROM sales 
	INNER JOIN documents ON sales.document_id = documents.id  
	INNER JOIN stores ON sales.store_id = stores.id 
	WHERE 1=1 `

	// QueryParentJoins is a query for getting parent joins
	QueryParentJoins = `
INNER JOIN documents ON sales.document_id = documents.id
INNER JOIN stores ON sales.store_id = stores.id
`
)
