package search_service

import (
	"api/common"
	"api/model"
	"fmt"
	"github.com/jiang-ruoxi/gopkg/es"
	"sort"
)

const searchFrom int = 0
const statusEnable int = 1                      // 状态:1已发布 2未发布
const libTypeDynamicId int = 1                  //资讯类型ID
const libTypeDynamic string = "dynamic"         //资讯类型
const libTypeDynamicZH string = "动态资讯"          //资讯类型
const libTypeLawId int = 2                      //法律类型ID
const libTypeLaw string = "law"                 //法律类型
const libTypeLawZH string = "语言法律"              //法律类型
const libTypeCaseId int = 3                     //案例类型ID
const libTypeCase string = "case"               //案例类型
const libTypeCaseZH string = "语言生活"             //案例类型
const libTypeAchievementId int = 4              //成果库类型ID
const libTypeAchievement string = "achievement" //成果库类型
const libTypeAchievementZH string = "研究成果"      //成果库类型
const libTypeReportId int = 5                   //报告类型ID
const libTypeReport string = "report"           //报告类型
const libTypeReportZH string = "机构报告"           //报告类型

var LangFieldMapping = map[common.SearchEngineGroupType][]string{
	common.SearchEngineGroupTypeAll:          {"title", "author", "full_content", "content", "introduction"}, // all
	common.SearchEngineGroupTypeFC:           {"full_content"},                                               // 全文
	common.SearchEngineGroupTypeTitle:        {"title"},                                                      // 标题
	common.SearchEngineGroupTypeIntroduction: {"introduction"},                                               // 简介
	common.SearchEngineGroupTypeAuthor:       {"author"},                                                     // 作者
	common.SearchEngineGroupTypeContent:      {"content"},                                                    // 内容
}

var LangOrderMapping = map[common.SearchEngineGroupOrder]string{
	common.SearchEngineGroupOrderDefault:     "id:desc",
	common.SearchEngineGroupOrderScoreAsc:    "_score:asc",
	common.SearchEngineGroupOrderScoreDesc:   "_score:desc",
	common.SearchEngineGroupOrderPublishAsc:  "publish_time_order:asc",
	common.SearchEngineGroupOrderPublishDesc: "publish_time_order:desc",
}

func (srv *SearchService) sortOrderMap(sortOrder int8) (sort []string) {
	var langOrder string
	switch sortOrder {
	case 1:
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderScoreAsc]
	case -1:
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderScoreDesc]
	case 2:
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderPublishAsc]
	case -2:
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderPublishDesc]
	default:
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderDefault]
	}

	sort = append(sort, langOrder)
	if langOrder == "_score:desc" {
		langOrder = LangOrderMapping[common.SearchEngineGroupOrderPublishDesc]
		sort = append(sort, langOrder)
	}

	return
}

func (srv *SearchService) field(behavior common.SearchEngineGroupBehavior, boost int, field string, rangeMin int64, rangeMax int64, values []string) []es.QueryMap {
	var conditions []es.QueryMap

	switch behavior {
	case common.SearchEngineGroupBehaviorMatch:
		for _, value := range values {
			if boost > 1 {
				// 短语
				conditions = append(conditions, es.QueryMap{"match_phrase": es.QueryMap{fmt.Sprintf("%s", field): es.QueryMap{"query": value, "boost": boost * 128}}})
			} else {
				// 短语
				conditions = append(conditions, es.QueryMap{"match_phrase": es.QueryMap{fmt.Sprintf("%s", field): es.QueryMap{"query": value, "boost": boost * 16}}})

				// 分词
				conditions = append(conditions, es.QueryMap{"match": es.QueryMap{fmt.Sprintf("%s", field): es.QueryMap{"query": value, "boost": boost}}})
			}
		}
	case common.SearchEngineGroupBehaviorMatchPhrase:
		conditions = append(conditions, es.QueryMap{"terms": es.QueryMap{fmt.Sprintf("%s.keyword", field): values, "boost": boost * 128}})
		for _, value := range values {
			// 短语
			conditions = append(conditions, es.QueryMap{"match_phrase": es.QueryMap{fmt.Sprintf("%s.base.base", field): es.QueryMap{"query": value, "boost": boost}}})
		}
	case common.SearchEngineGroupBehaviorTerms:
		conditions = append(conditions, es.QueryMap{"terms": es.QueryMap{fmt.Sprintf("%s.keyword", field): values, "boost": boost * 128}})
	case common.SearchEngineGroupBehaviorPrefix:
		for _, value := range values {
			conditions = append(conditions, es.QueryMap{"prefix": es.QueryMap{fmt.Sprintf("%s.keyword", field): es.QueryMap{"value": value, "boost": boost}}})
		}
	case common.SearchEngineGroupBehaviorRange:
		if rangeMin > 0 {
			conditions = append(conditions, es.QueryMap{"range": es.QueryMap{fmt.Sprintf("%s.long", field): es.QueryMap{"gte": rangeMin, "boost": boost}}})
		}
		if rangeMax > 0 {
			conditions = append(conditions, es.QueryMap{"range": es.QueryMap{fmt.Sprintf("%s.long", field): es.QueryMap{"lt": rangeMax, "boost": boost}}})
		}
	}

	return conditions
}

func (srv *SearchService) makeLibTypeData() (libTypeDataList []common.SearchLibTypeItem) {
	libTypeList := srv.initLibTypeList()
	var libTypeData common.SearchLibTypeItem
	for _, item := range libTypeList {
		libTypeData.Id = item.Id
		libTypeData.Name = item.Name
		libTypeData.EName = item.EName
		libTypeData.Count = item.Count
		libTypeDataList = append(libTypeDataList, libTypeData)
	}
	return
}

func (srv *SearchService) sortSliceList(libTypeDataList []common.SearchLibTypeItem, sectionDataList []common.SearchSectionItem, yearDataList []common.SearchYearItem, countryDataList []common.SearchCountryItem, keywordDataList []common.SearchKeywordItem) {
	sort.Slice(libTypeDataList, func(i, j int) bool {
		return libTypeDataList[i].Name < libTypeDataList[j].Name
	})
	sort.Slice(sectionDataList, func(i, j int) bool {
		if sectionDataList[i].Count == sectionDataList[j].Count {
			return sectionDataList[i].Name < sectionDataList[j].Name
		}
		return sectionDataList[i].Count > sectionDataList[j].Count
	})
	sort.Slice(yearDataList, func(i, j int) bool {
		return yearDataList[i].Name < yearDataList[j].Name
	})
	sort.Slice(countryDataList, func(i, j int) bool {
		if countryDataList[i].Count == countryDataList[j].Count {
			return countryDataList[i].Name < countryDataList[j].Name
		}
		return countryDataList[i].Count > countryDataList[j].Count
	})
	sort.Slice(keywordDataList, func(i, j int) bool {
		if keywordDataList[i].Count == keywordDataList[j].Count {
			return keywordDataList[i].Keyword < keywordDataList[j].Keyword
		}
		return keywordDataList[i].Count > keywordDataList[j].Count
	})
}

func (srv *SearchService) dealGroupStatisticsCount(content common.SearchResponseBody, libTypeDataList []common.SearchLibTypeItem, sectionList []model.Section) (sectionDataList []common.SearchSectionItem, yearDataList []common.SearchYearItem, countryDataList []common.SearchCountryItem, keywordDataList []common.SearchKeywordItem, libTypeList []common.SearchLibTypeItem) {
	var aggregationsList = content.Aggregations
	keywordItemList := aggregationsList["group_by_keywords"].Buckets
	for _, item := range keywordItemList {
		if item.Key.(string) == "" {
			continue
		}
		keywordDataList = append(keywordDataList, common.SearchKeywordItem{
			Keyword: item.Key.(string),
			Count:   int64(item.DocCount),
		})
	}
	libTypeItemList := aggregationsList["group_by_lib_type"].Buckets
	for idx, _ := range libTypeDataList {
		var libType common.SearchLibTypeItem
		libType.Id = libTypeDataList[idx].Id
		libType.Name = libTypeDataList[idx].Name
		libType.EName = libTypeDataList[idx].EName
		libType.Count = 0
		for _, item := range libTypeItemList {
			if libTypeDataList[idx].EName == item.Key {
				libType.Id = libTypeDataList[idx].Id
				libType.Name = libTypeDataList[idx].Name
				libType.EName = libTypeDataList[idx].EName
				libType.Count = int64(item.DocCount)
			}
		}
		libTypeList = append(libTypeList, libType)
	}
	yearItemList := aggregationsList["group_by_year"].Buckets
	for _, item := range yearItemList {
		if item.Key.(string) == "" {
			continue
		}
		yearDataList = append(yearDataList, common.SearchYearItem{
			Name:  item.Key.(string),
			Count: int64(item.DocCount),
		})
	}
	countryItemList := aggregationsList["group_by_country"].Buckets
	for _, item := range countryItemList {
		if item.Key.(string) == "" {
			continue
		}
		countryDataList = append(countryDataList, common.SearchCountryItem{
			Name:  item.Key.(string),
			Count: int64(item.DocCount),
		})
	}
	sectionItemList := aggregationsList["group_by_section_id"].Buckets
	var sectionItem common.SearchSectionItem
	for _, item := range sectionItemList {
		for sIdx, _ := range sectionList {
			if item.Key.(float64) == float64(int(sectionList[sIdx].ID)) {
				sectionItem.ID = uint(int(sectionList[sIdx].ID))
				sectionItem.Name = sectionList[sIdx].Name
				sectionItem.Count = int64(item.DocCount)
			}
		}
		sectionDataList = append(sectionDataList, sectionItem)
	}
	return
}

func (srv *SearchService) searchGroupQueryBody(query Boolean) (queryBody es.QueryMap) {
	queryBody = es.QueryMap{
		"query":            query.EncodeToQueryMap(),
		"track_total_hits": true,
		"aggs": es.QueryMap{
			"group_by_keywords": es.QueryMap{
				"terms": es.QueryMap{
					"field": "keywords.keyword",
					"size":  1000,
				},
			},
			"group_by_country": es.QueryMap{
				"terms": es.QueryMap{
					"field": "country.keyword",
					"size":  1000,
				},
			},
			"group_by_year": es.QueryMap{
				"terms": es.QueryMap{
					"field": "year.keyword",
					"size":  1000,
				},
			},
			"group_by_section_id": es.QueryMap{
				"terms": es.QueryMap{
					"field": "section_id",
					"size":  1000,
				},
			},
			"group_by_lib_type": es.QueryMap{
				"terms": es.QueryMap{
					"field": "lib_type.keyword",
					"size":  1000,
				},
			},
		},
	}
	return
}

func (srv *SearchService) searchQueryBody(query Boolean) (queryBody es.QueryMap) {
	queryBody = es.QueryMap{
		"query":            query.EncodeToQueryMap(),
		"track_total_hits": true,
		"highlight": es.QueryMap{
			"fields": es.QueryMap{
				"title": es.QueryMap{
					"pre_tags":            "<em class='jp-s-em'>",
					"post_tags":           "</em>",
					"fragment_size":       0,          // 设置为0返回整个字段内容
					"number_of_fragments": 0,          // 设置为0返回整个字段内容
					"boundary_scanner":    "sentence", // 设置边界扫描为句子
				},
				"author": es.QueryMap{
					"pre_tags":            "<em class='jp-s-em'>",
					"post_tags":           "</em>",
					"fragment_size":       0,          // 设置为0返回整个字段内容
					"number_of_fragments": 0,          // 设置为0返回整个字段内容
					"boundary_scanner":    "sentence", // 设置边界扫描为句子
				},
				"introduction": es.QueryMap{
					"pre_tags":            "<em class='jp-s-em'>",
					"post_tags":           "</em>",
					"fragment_size":       0,          // 设置为0返回整个字段内容
					"number_of_fragments": 0,          // 设置为0返回整个字段内容
					"boundary_scanner":    "sentence", // 设置边界扫描为句子
				},
				"content": es.QueryMap{
					"pre_tags":            "<em class='jp-s-em'>",
					"post_tags":           "</em>",
					"fragment_size":       0,          // 设置为0返回整个字段内容
					"number_of_fragments": 0,          // 设置为0返回整个字段内容
					"boundary_scanner":    "sentence", // 设置边界扫描为句子
				},
			},
		},
	}
	return
}

func (srv *SearchService) searchEngineGroup(req *common.SearchRequest) (groups []common.SearchEngineGroup) {
	for n, item := range req.Filters {
		// 默认匹配
		var g = common.SearchEngineGroup{
			Type:     common.SearchEngineGroupType(item.Field),
			Mode:     common.SearchEngineGroupMode(item.Logical),
			Behavior: common.SearchEngineGroupBehaviorMatch,
			Values:   []string{item.Value},
			Boost:    1,
		}
		// 精确
		if item.Type == common.SearchEngineBehaviorTerms {
			g.Boost = 3
		}
		// 第一个条件必须是and，前端有可能传错，强制处理
		if 0 == n {
			g.Mode = common.SearchEngineGroupModeAnd
		}
		groups = append(groups, g) // 评分条件组
	}
	return
}

func (srv *SearchService) langQueryMustConditions(req *common.SearchRequest) es.QueryMap {
	var boolean = Boolean{}
	var queryItemList []es.QueryMap
	var queryItem = es.QueryMap{}
	if len(req.FilterLibs) > 0 {
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"lib_type": req.FilterLibs,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}

	boolean.Must = append(boolean.Must, queryItemList...)
	return boolean.EncodeToQueryMap()
}

func (srv *SearchService) langQueryFilterConditions(req *common.SearchRequest) es.QueryMap {
	var boolean = Boolean{}
	var queryItemList []es.QueryMap
	var queryItem = es.QueryMap{}
	if len(req.Libs) > 0 {
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"lib_type": req.Libs,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}
	if len(req.Years) > 0 {
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"year": req.Years,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}
	if len(req.Countrys) > 0 {
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"country.keyword": req.Countrys,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}

	if len(req.SectionIDs) > 0 {
		var sectionIdList []string
		sectionIds := req.SectionIDs
		for _, sectionId := range sectionIds {
			sectionIdList = append(sectionIdList, fmt.Sprintf("%d", sectionId))
		}
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"section_id": sectionIdList,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}

	if len(req.Keywords) > 0 {
		queryItem = es.QueryMap{
			"terms": map[string]interface{}{
				"keywords.keyword": req.Keywords,
			},
		}
		queryItemList = append(queryItemList, queryItem)
	}

	boolean.Filter = append(boolean.Filter, queryItemList...)
	return boolean.EncodeToQueryMap()
}

func (srv *SearchService) langConditions(groups []common.SearchEngineGroup) es.QueryMap {
	// 处理过滤，每个组的mode对前一个组生效，and与前一个组并列，or与前一个组取一，对应到es的
	// should + mini 1，not时，前一个组应该在must，当前组应该在must not

	var boolean = Boolean{}

	// filter用于过滤文档
	// 优化：连续的AND与OR可以放在一个层级里
	var before = common.SearchEngineGroupModeAnd
	for _, group := range groups {
		var condition es.QueryMap
		// 处理子条件组
		if nil == group.Sub {
			condition = Boolean{MinimumShouldMatch: 1, Should: srv.langQueryMap(group, group.Boost)}.EncodeToQueryMap()
		} else {
			condition = srv.langConditions(group.Sub)
		}

		// 其他组需要与前一个组进行组合
		switch group.Mode {
		case common.SearchEngineGroupModeAnd:
			if common.SearchEngineGroupModeAnd == before {
				// 直接合并条件
				boolean.Must = append(boolean.Must, condition)
			} else {
				// 创建新的对象并替换
				var newFilter = Boolean{}
				if !boolean.IsEmpty() {
					newFilter.Must = append(newFilter.Must, boolean.EncodeToQueryMap())
				}

				newFilter.Must = append(newFilter.Must, condition)
				boolean = newFilter
			}
		case common.SearchEngineGroupModeOr:
			if common.SearchEngineGroupModeOr == before {
				// 直接合并条件
				boolean.Should = append(boolean.Should, condition)
			} else {
				// 创建新的对象并替换
				var newBoolean = Boolean{}
				if !boolean.IsEmpty() {
					newBoolean.Should = append(newBoolean.Should, boolean.EncodeToQueryMap())
				}

				newBoolean.Should = append(newBoolean.Should, condition)
				boolean = newBoolean
			}

			// 最少匹配一个
			boolean.MinimumShouldMatch = 1
		case common.SearchEngineGroupModeNot:
			if common.SearchEngineGroupModeNot == before {
				// 直接合并条件
				boolean.MustNot = append(boolean.MustNot, condition)
			} else {
				// 创建新的对象并替换
				var newBoolean = Boolean{}
				if !boolean.IsEmpty() {
					newBoolean.Must = append(newBoolean.Must, boolean.EncodeToQueryMap())
				}

				newBoolean.MustNot = append(newBoolean.MustNot, condition)
				boolean = newBoolean
			}
		}

		before = group.Mode
	}

	return boolean.EncodeToQueryMap()
}

func (srv *SearchService) langQueryMap(group common.SearchEngineGroup, boost int) []es.QueryMap {
	var should []es.QueryMap

	for _, field := range LangFieldMapping[group.Type] {
		should = append(should, srv.field(group.Behavior, boost, field, group.RangeMin, group.RangeMax, group.Values)...)
	}

	return should
}

func (srv *SearchService) initLibTypeList() (libTypeList []common.SearchLibTypeItem) {
	libTypeList = append(libTypeList, common.SearchLibTypeItem{
		Id:    libTypeAchievementId,
		Name:  libTypeAchievementZH,
		EName: libTypeAchievement,
		Count: 0,
	}, common.SearchLibTypeItem{
		Id:    libTypeCaseId,
		Name:  libTypeCaseZH,
		EName: libTypeCase,
		Count: 0,
	}, common.SearchLibTypeItem{
		Id:    libTypeDynamicId,
		Name:  libTypeDynamicZH,
		EName: libTypeDynamic,
		Count: 0,
	}, common.SearchLibTypeItem{
		Id:    libTypeLawId,
		Name:  libTypeLawZH,
		EName: libTypeLaw,
		Count: 0,
	}, common.SearchLibTypeItem{
		Id:    libTypeReportId,
		Name:  libTypeReportZH,
		EName: libTypeReport,
		Count: 0,
	})
	return
}

func (srv *SearchService) requiredSearchCondition() (must []es.QueryMap) {
	// 状态必须是启用状态
	must = append(must, es.QueryMap{
		"term": es.QueryMap{
			"status": es.QueryMap{
				"value": statusEnable,
			},
		},
	})
	// 年份不能为空
	must = append(must, es.QueryMap{
		"range": es.QueryMap{
			"year": es.QueryMap{
				"gte": "",
			},
		},
	})
	// 不能为删除的数据
	must = append(must, es.QueryMap{
		"term": es.QueryMap{
			"delete_time": es.QueryMap{
				"value": 0,
			},
		},
	})
	// publish_time不能为空
	must = append(must, es.QueryMap{
		"range": es.QueryMap{
			"publish_time": es.QueryMap{
				"gte": "",
			},
		},
	})
	// lib_type不能为空
	must = append(must, es.QueryMap{
		"range": es.QueryMap{
			"lib_type": es.QueryMap{
				"gte": "",
			},
		},
	})
	return
}
