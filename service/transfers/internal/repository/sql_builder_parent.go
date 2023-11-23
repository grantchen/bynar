package repository

// tables: transfers, documents, stores
const (
	// QueryParentCount is a query for parent count
	QueryParentCount = `
		SELECT COUNT(transfers.id) as Count
		FROM transfers
				 INNER JOIN documents ON transfers.document_id = documents.id
				 INNER JOIN stores ON transfers.store_id = stores.id
		WHERE 1=1
		`

	// QueryParent is a query for parent
	QueryParent = `
		SELECT transfers.*,
			   COALESCE(lines_t.Count, 0) AS Count
		FROM transfers
				 INNER JOIN documents ON transfers.document_id = documents.id
				 INNER JOIN stores ON transfers.store_id = stores.id
				 LEFT JOIN (SELECT COUNT(transfer_lines.id) AS Count,
								   transfer_lines.parent_id AS parent_id
							FROM transfer_lines
							WHERE 2=2
							GROUP BY parent_id) lines_t
						   ON lines_t.parent_id = transfers.id
		WHERE 1=1
		`

	// QueryParentJoins is a query for parent joins
	QueryParentJoins = `
		INNER JOIN documents ON transfers.document_id = documents.id
		INNER JOIN stores ON transfers.store_id = stores.id
		`
)
