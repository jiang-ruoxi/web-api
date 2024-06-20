package search_service

import "github.com/jiang-ruoxi/gopkg/es"

type Boolean struct {
	Must               []es.QueryMap //逻辑与（AND）
	Should             []es.QueryMap //逻辑或（OR）
	MustNot            []es.QueryMap //逻辑非（NOT）
	Filter             []es.QueryMap //过滤条件
	MinimumShouldMatch int           // 最少匹配规则数量
}

func (cls Boolean) EncodeToQueryMap() es.QueryMap {
	var res = es.QueryMap{}
	//must 条件是必须满足的，相当于逻辑与（AND）
	if 0 != len(cls.Must) {
		res["must"] = cls.Must
	}
	//should 条件是可以满足也可以不满足的，相当于逻辑或（OR）
	if 0 != len(cls.Should) {
		res["should"] = cls.Should
	}
	//must_not 条件是必须不满足的，相当于逻辑非（NOT）
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
