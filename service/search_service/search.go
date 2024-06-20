package search_service

import (
	"api/common"
	"api/model"
	"api/utils/errs"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jiang-ruoxi/gopkg/es"
	"github.com/jiang-ruoxi/gopkg/log"
	"github.com/jiang-ruoxi/gopkg/server/api"
	"github.com/spf13/viper"
	"io"
)

func NewSearchService() *SearchService {
	engine, err := es.Get("engine")
	if nil != err {
		panic(err)
	}
	return &SearchService{
		engine: engine,
	}
}

type SearchService struct {
	engine *elasticsearch.Client
}

func (srv *SearchService) SearchList(ctx context.Context, req *common.SearchRequest) (response common.SearchListResponse, apiErr api.Error) {
	var sectionList []model.Section
	if err := model.DefaultWeb().Model(&model.Section{}).WithContext(ctx).Find(&sectionList).Error; err != nil {
		log.SugarContext(ctx).Errorw("SearchService Section List Find", "error", err)
	}
	// 获取左侧分组栏目
	groupDataList, err := srv.fetchGroupData(ctx, req, sectionList)
	if err != nil {
		log.SugarContext(ctx).Errorw("SearchService fetchGroupData 失败", "err", err)
		return response, errs.Failed
	}
	response.LibTypeList = groupDataList.LibTypeList
	response.SectionList = groupDataList.SectionList
	response.YearList = groupDataList.YearList
	response.CountryList = groupDataList.CountryList
	response.KeywordList = groupDataList.KeywordList

	// 获取右侧数据库列表
	dataList, count, err := srv.fetchSearchData(ctx, req, sectionList)
	if err != nil {
		log.SugarContext(ctx).Errorw("SearchService fetchSearchData 失败", "err", err)
		return response, errs.Failed
	}
	response.List = dataList
	response.P = req.P
	response.N = req.N
	response.Count = count

	return response, nil
}

func (srv *SearchService) fetchGroupData(ctx context.Context, req *common.SearchRequest, sectionList []model.Section) (statisticsItem common.StatisticsItem, err error) {
	var queryGroup Boolean
	if 0 != len(srv.searchEngineGroup(req)) {
		queryGroup.Must = append(queryGroup.Must, srv.langConditions(srv.searchEngineGroup(req)))
	}
	queryGroup.Must = append(queryGroup.Must, srv.requiredSearchCondition()...)
	queryGroup.Must = append(queryGroup.Must, srv.langQueryMustConditions(req))

	var bodyGroupConditions bytes.Buffer
	if err := json.NewEncoder(&bodyGroupConditions).Encode(srv.searchGroupQueryBody(queryGroup)); nil != err {
		log.SugarContext(ctx).Errorw("SearchService.fetchGroupData json.NewEncoder(&bodyGroupConditions) 失败", "err", err)
		return statisticsItem, err
	}

	var contentGroupConditions common.SearchResponseBody
	responseGroup, err := srv.engine.Search(
		srv.engine.Search.WithContext(ctx),
		srv.engine.Search.WithIndex(viper.GetStringSlice(viper.GetString("elasticsearch.engine.lang_indices"))...),
		srv.engine.Search.WithIgnoreUnavailable(true),
		srv.engine.Search.WithBody(&bodyGroupConditions),
		srv.engine.Search.WithFrom(searchFrom),
		srv.engine.Search.WithSize(1),
	)

	if err != nil {
		log.SugarContext(ctx).Errorw("SearchService.fetchGroupData bodyGroupConditions 失败", "err", err)
		return statisticsItem, err
	}

	defer func() {
		_ = responseGroup.Body.Close()
	}()

	if responseGroup.IsError() {
		bs, err := io.ReadAll(responseGroup.Body)
		if nil != err {
			log.SugarContext(ctx).Errorw("SearchService.fetchGroupData 读取ES响应错误信息失败", "err", err)
			return statisticsItem, err
		}

		log.SugarContext(ctx).Errorw("SearchService.fetchGroupData 搜索失败", "err", string(bs))
		return statisticsItem, err
	}

	if err := json.NewDecoder(responseGroup.Body).Decode(&contentGroupConditions); nil != err {
		log.SugarContext(ctx).Errorw("SearchService.fetchGroupData 解析ES响应失败", "err", err)
		return statisticsItem, err
	}

	libTypeDataList := srv.makeLibTypeData()
	sectionListData, yearDataListData, countryDataListData, keywordListData, libTypeDataListNew := srv.dealGroupStatisticsCount(contentGroupConditions, libTypeDataList, sectionList)
	srv.sortSliceList(libTypeDataListNew, sectionListData, yearDataListData, countryDataListData, keywordListData)
	statisticsItem.LibTypeList = libTypeDataListNew  //资源库列表
	statisticsItem.SectionList = sectionListData     // 栏目/类型列表
	statisticsItem.YearList = yearDataListData       // 发表年份列表
	statisticsItem.CountryList = countryDataListData // 国家列表
	statisticsItem.KeywordList = keywordListData     // 关键词列表
	return
}

func (srv *SearchService) fetchSearchData(ctx context.Context, req *common.SearchRequest, sectionList []model.Section) (searchData []common.SearchResponseItem, count int64, err error) {
	var query Boolean
	if 0 != len(srv.searchEngineGroup(req)) {
		query.Must = append(query.Must, srv.langConditions(srv.searchEngineGroup(req)))
	}
	query.Must = append(query.Must, srv.requiredSearchCondition()...)
	query.Must = append(query.Must, srv.langQueryMustConditions(req))
	query.Filter = append(query.Filter, srv.langQueryFilterConditions(req))
	sortOrder := srv.sortOrderMap(req.Sort)
	if len(req.Filters) <= 0 {
		if req.Sort == -1 || req.Sort == 1 {
			sortOrder = []string{"publish_time_order:desc"}
		}
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(srv.searchQueryBody(query)); nil != err {
		log.SugarContext(ctx).Errorw("SearchService.fetchSearchData json.NewEncoder(&body) 失败", "err", err)
		return searchData, count, err
	}
	var content common.SearchResponseBody
	response, err := srv.engine.Search(
		srv.engine.Search.WithContext(ctx),
		srv.engine.Search.WithIndex(viper.GetStringSlice(viper.GetString("elasticsearch.engine.lang_indices"))...),
		srv.engine.Search.WithIgnoreUnavailable(true),
		srv.engine.Search.WithBody(&body),
		srv.engine.Search.WithSort(sortOrder...),
		srv.engine.Search.WithFrom((req.P-1)*req.N),
		srv.engine.Search.WithSize(req.N),
	)
	if err != nil {
		log.SugarContext(ctx).Errorw("SearchService.fetchSearchData Search 失败", "err", err)
		return searchData, count, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.IsError() {
		bs, err := io.ReadAll(response.Body)
		if nil != err {
			log.SugarContext(ctx).Errorw("SearchService.fetchSearchData 读取ES响应错误信息失败", "err", err)
			return searchData, count, err
		}

		log.SugarContext(ctx).Errorw("SearchService.fetchSearchData 搜索失败", "err", string(bs))
		return searchData, count, err
	}

	if err := json.NewDecoder(response.Body).Decode(&content); nil != err {
		log.SugarContext(ctx).Errorw("SearchService.fetchSearchData 解析ES响应失败", "err", err)
		return searchData, count, err
	}
	searchData = srv.formatSearchData(content, srv.initLibTypeList(), sectionList)

	count = content.Hits.Total.Value
	return
}

func (srv *SearchService) formatSearchData(content common.SearchResponseBody, libTypeList []common.SearchLibTypeItem, sectionList []model.Section) (data []common.SearchResponseItem) {
	for _, hit := range content.Hits.Hits {
		var searchResponseItem = new(common.SearchResponseItem)
		var item = hit.Source
		searchResponseItem.LibType = item.LibType //库
		for idx, _ := range libTypeList {
			if libTypeList[idx].EName == item.LibType {
				searchResponseItem.LibTypeId = libTypeList[idx].Id     //库ID
				searchResponseItem.LibTypeName = libTypeList[idx].Name //库名称
			}
		}
		searchResponseItem.ID = item.ID       // id
		searchResponseItem.Title = item.Title // 标题
		if hit.Highlight != nil && len(hit.Highlight["title"]) > 0 {
			searchResponseItem.TitleHighlight = hit.Highlight["title"][0]
		}
		searchResponseItem.Author = item.Author // 作者
		if hit.Highlight != nil && len(hit.Highlight["author"]) > 0 {
			searchResponseItem.AuthorHighlight = hit.Highlight["author"][0]
		}
		searchResponseItem.SectionID = item.SectionID // 栏目ID
		for idx, _ := range sectionList {
			if sectionList[idx].ID == item.SectionID {
				searchResponseItem.SectionName = sectionList[idx].Name //库名称
			}
		}
		searchResponseItem.Introduction = item.Introduction // 简介
		if hit.Highlight != nil && len(hit.Highlight["introduction"]) > 0 {
			searchResponseItem.IntroductionHighlight = hit.Highlight["introduction"][0]
		}
		searchResponseItem.CoverImage = item.CoverImage               // 封面图
		searchResponseItem.PublishTimeOrder = item.PublishTimeOrder   // 出版时间
		searchResponseItem.PublishTime = item.PublishTime             // 出版时间
		searchResponseItem.PublishPeriodical = item.PublishPeriodical // 总刊期
		searchResponseItem.PublicationTitle = item.PublicationTitle   // 出版物题目
		searchResponseItem.Publisher = item.Publisher                 // 出版社
		searchResponseItem.Year = item.Year                           // 年
		searchResponseItem.Country = item.Country                     // 国家
		searchResponseItem.Keywords = item.Keywords                   // 关键词
		data = append(data, *searchResponseItem)
	}
	return
}
