package variables

type Jwt struct {
	SecretKey         string
	Issuer            string
	ExpirationMinutes int64
}
