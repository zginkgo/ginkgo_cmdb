package impl

const (
	insertSecretSQL = `INSERT INTO secret (
		id,create_at,description,vendor,address,allow_regions,crendential_type,api_key,api_secret,request_rate
	) VALUES (?,?,?,?,?,?,?,?,?,?);`

	deleteSecretSQL = `DELETE FROM secret WHERE id = ?;`
	querySecretSQL  = `SELECT * FROM secret`
)
