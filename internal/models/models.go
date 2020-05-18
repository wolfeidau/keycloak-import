package models

type User struct {
	Department        string
	DistinguishedName string
	Enabled           string
	GivenName         string
	Mail              string `mapstructure:"mail"`
	Manager           string
	Name              string
	ObjectClass       string
	ObjectGUID        string
	OfficePhone       string
	SamAccountName    string
	SID               string
	sn                string
	Surname           string
	Title             string
	UserPrincipalName string
}