package search_service

import "github.com/jiang-ruoxi/gopkg/es"

type Boolean struct {
	Must               []es.QueryMap
	Should             []es.QueryMap
	MustNot            []es.QueryMap
	Filter             []es.QueryMap
	MinimumShouldMatch int // 最少匹配规则数量
}

func (cls Boolean) EncodeToQueryMap() es.QueryMap {
	var res = es.QueryMap{}
	if 0 != len(cls.Must) {
		res["must"] = cls.Must
	}

	if 0 != len(cls.Should) {
		res["should"] = cls.Should
	}

	if 0 != len(cls.MustNot) {
		res["must_not"] = cls.MustNot
	}

	if 0 != len(cls.Filter) {
		res["filter"] = cls.Filter
	}

	if 0 != cls.MinimumShouldMatch {
		res["minimum_should_match"] = cls.MinimumShouldMatch
	}

	return es.QueryMap{"bool": res}
}

func (cls Boolean) IsEmpty() bool {
	return 0 == len(cls.Must) &&
		0 == len(cls.Should) &&
		0 == len(cls.MustNot) &&
		0 == len(cls.Filter)
}
