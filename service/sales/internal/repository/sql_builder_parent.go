package repository

// tables: sales, documents, stores
const (
	// QueryParentCount is a query for getting parent count
	QueryParentCount = `
		SELECT COUNT(sales.id) as Count
		FROM sales
				 INNER JOIN documents ON sales.document_id = documents.id
				 INNER JOIN stores ON sales.store_id = stores.id
		WHERE 1=1
		`

	// QueryParent is a query for getting parent
	QueryParent = `
		SELECT sales.*,
			   COALESCE(lines_t.Count, 0) AS Count
		FROM sales
				 INNER JOIN documents ON sales.document_id = documents.id
				 INNER JOIN stores ON sales.store_id = stores.id
				 LEFT JOIN (SELECT COUNT(sale_lines.id) AS Count,
				                   sale_lines.parent_id AS parent_id
							FROM sale_lines
							WHERE 2=2
							GROUP BY parent_id) lines_t
						   ON lines_t.parent_id = sales.id
		WHERE 1=1
		`

	// QueryParentJoins is a query for getting parent joins
	QueryParentJoins = `
		INNER JOIN documents ON sales.document_id = documents.id
		INNER JOIN stores ON sales.store_id = stores.id
		`
)
