package worker

import "golang-awesome/crawler/engine"

type CrawlService struct {
}

func (CrawlService) Process(req Request,result *ParseResult) error{
	request, err := DeserializeRequest(req)
	if err != nil {
		return err
	}

	parserResult, err := engine.Worker(request)
	*result = SerializeResult(parserResult)
	return nil
}
