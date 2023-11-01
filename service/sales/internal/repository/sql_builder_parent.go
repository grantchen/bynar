package repository

// tables: sales, documents, stores
const (
	QueryParentCount = `
	SELECT COUNT(sales.id) as Count 
	FROM sales 
	INNER JOIN documents ON sales.document_id = documents.id  
	INNER JOIN stores ON sales.store_id = stores.id 
	WHERE 1=1 `

	QueryParent = `
	SELECT 
		sales.*
	FROM sales 
	INNER JOIN documents ON sales.document_id = documents.id  
	INNER JOIN stores ON sales.store_id = stores.id 
	WHERE 1=1 `

	QueryParentJoins = `
INNER JOIN documents ON sales.document_id = documents.id
INNER JOIN stores ON sales.store_id = stores.id
`
	QueryParentBuild = `
	SELECT * 
	FROM sales `
)
