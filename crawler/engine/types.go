package engine

//函数类型，只要参类型、个数、顺序相同就可以
type ParserFunc func(body []byte, url string) ParserResult

//解析函数接口
type Parser interface {
	Parse(contents []byte, url string) ParserResult
	Serialize() (name string, args interface{})
}

//请求，包括url和指定的解析函数
type Request struct {
	Url   string
	Parse Parser
}

//解析结果，包括爬到的内容和url
type ParserResult struct {
	Requests []Request
	Items    []Item
}

//空的解析方法
func NilParseFunc(body []byte, url string) ParserResult {
	return ParserResult{}
}

//一个页面的对象
type Item struct {
	Url     string
	Id      string
	Type    string
	Payload interface{}
}

type NilParse struct {
}

func (NilParse) Parse(contents []byte, url string) ParserResult {
	return ParserResult{}
}

func (NilParse) Serialize() (name string, args interface{}) {
	return "NilParse", nil
}

type FuncParser struct {
	parser ParserFunc
	name   string
}

func NewFuncParser(parser ParserFunc, name string) *FuncParser {
	return &FuncParser{parser: parser, name: name}
}

func (p *FuncParser) Parse(contents []byte, url string) ParserResult {
	return p.parser(contents, url)
}

//TODO what can we return after serialized
func (p *FuncParser) Serialize() (name string, args interface{}) {
	return p.name, nil
}
