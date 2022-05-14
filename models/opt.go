package models

type Opts struct {
	Location     string
	Proxy        string
	KeyHeader    string
	ValueHeader  string
	Elastic      bool
	ElasticURI   string
	ElasticUser  string
	ElasticPass  string
	ElasticIndex string
}
type Elastic struct {
	Endpoint   string
	Method     string
	ReqBody    string
	ResBody    string
	Headers    map[string]string
	RespStatus int
}
