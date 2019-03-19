package config


type UOSCheck struct{
	Disable bool `json:"disable"`
	UrlForGetTableRows string `json:"url_for_get_table_rows"`
	TableName string `json:"table_name"`
	Code string `json:"code"`
}
