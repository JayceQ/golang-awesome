package parser

import (
	"golang-awesome/crawler/engine"
	"golang-awesome/crawler/model"
	"regexp"
)
//<div class="m-btn purple" data-v-bff6f798="">离异</div>
//var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>(\d+)岁</td>`)
//var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>(\d+)CM</td>`)
//var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
//var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([^<]+)KG</span></td>`) //<td><span class="label">体重：</span><span field="">--</span></td>
//var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
//var xinzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
//var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
//var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
//var occupationRe = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
//var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
//var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
//var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
//var IdRe = regexp.MustCompile(`<div class="id" data-v-5b109fc3>ID：([^<]+)</div>`)
var Re = regexp.MustCompile(`<div class="m-btn purple" data-v-bff6f798>([^<]+)</div>`)


func ParseProfile(contents []byte, url, name string) engine.ParserResult{
	profile := &model.Profile{}
	profile.Name = name
	//idsubmatch := IdRe.FindSubmatch(contents)
	//profile.Id = string(idsubmatch[1])
	submatch := Re.FindAllSubmatch(contents, -1)
	if len(submatch) > 0{
		profile.Marriage = string(submatch[0][1])
		profile.Age = string(submatch[1][1])
		profile.Xingzuo = string(submatch[2][1])
		profile.Height = string(submatch[3][1])
		profile.Weight = string(submatch[4][1])
		profile.Hukou = string(submatch[5][1])
		profile.Income = string(submatch[6][1])
		profile.Education = string(submatch[len(submatch)-1][1])
	}


	item := engine.Item{
		Url:url,
		Payload:profile,
		Type:"zhenai",
		Id:profile.Id,
	}
	rs := engine.ParserResult{}
	rs.Items = []engine.Item{item}
	// 取本页面内，猜你喜欢的的
	//var guessRe = regexp.MustCompile(`href="(http://album.zhenai.com/u/\w+)"[^>]*>([^<]+)</a>`)
	//ms := guessRe.FindAllSubmatch(contents, -1)
	//for _, m := range ms {
	//	rs.Requests = append(rs.Requests, engine.Request{
	//		Url:   string(m[1]),
	//		Parse:  NewProfileParser(string(m[2])),
	//	})
	//}
	return rs
}

//func extractString(contents []byte, re *regexp.Regexp) string{
//	match := re.FindSubmatch(contents)
//	if len(match) >= 2 {
//		return string(match[1])
//	}else {
//		return ""
//
//	}
//}

type ProfileParser struct{
	userName string
}

func NewProfileParser(userName string) *ProfileParser {
	return &ProfileParser{userName:
		userName}
}

func (p *ProfileParser) Parse(contents []byte, url string) engine.ParserResult {
	return ParseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return "ProfileParser", p.userName
}

