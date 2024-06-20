package common

type SearchRequest struct {
	SectionIDs []uint          `json:"section_ids"`            // 栏目ID数组
	Years      []string        `json:"years"`                  // 年份数组
	Keywords   []string        `json:"keywords"`               // 关键词数组
	Countrys   []string        `json:"countrys"`               // 国家数组
	Libs       []string        `json:"libs"`                   //搜索条件: 检索条件组
	Sort       int8            `json:"sort" enums:"1,-1,2,-2"` // 1相关度升序 -1相关度降序 2发表时间升序 -2发表时间降序
	Filters    []FilterRequest `json:"filters" binding:"dive"` // 搜索条件: 检索条件组
	FilterLibs []string        `json:"filter_libs"`            // 搜索条件: 检索条件库
	P          int             `json:"p"`                      // 分页
	N          int             `json:"n"`                      // 每页显示数量
}

type FilterRequest struct {
	Field   string `json:"field" enums:"all,title,author,introduction,content,full_content" validate:"required" binding:"required,oneof=all title author introduction content full_content"`
	Value   string `json:"value" validate:"optional" binding:"omitempty"`                                      // 关键字
	Logical string `json:"logical" enums:"and,or,not" validate:"required" binding:"required,oneof=and or not"` // 字段过滤逻辑 and与 or或 not非
	Type    uint8  `json:"type" binding:"required,gte=1,lte=2"`                                                // 1 模糊搜索 2精确搜索
}

type (
	SearchListResponse struct {
		Count int64                `json:"count"` // 总数
		P     int                  `json:"p"`     // 分页
		N     int                  `json:"n"`     // 每页显示数量
		List  []SearchResponseItem `json:"list"`
		StatisticsItem
	}
)

type StatisticsItem struct {
	LibTypeList []SearchLibTypeItem `json:"lib_type_list"` // 资源库列表
	SectionList []SearchSectionItem `json:"section_list"`  // 栏目/类型列表
	YearList    []SearchYearItem    `json:"year_list"`     // 发表年份列表
	CountryList []SearchCountryItem `json:"country_list"`  // 国家列表
	KeywordList []SearchKeywordItem `json:"keyword_list"`  // 关键词列表
}

type SearchResponseItem struct {
	ID                    uint        `json:"id"`                     // 主键ID
	SectionID             uint        `json:"section_id" `            // 类型ID
	Title                 string      `json:"title"`                  // 标题
	TitleHighlight        string      `json:"title_highlight"`        // 标题高亮
	LibType               string      `json:"lib_type"`               // 资源库类型
	LibTypeId             int         `json:"lib_type_id"`            // 资源库id
	LibTypeName           string      `json:"lib_type_name"`          // 资源库类型名称
	Author                string      `json:"author"`                 // 作者
	AuthorHighlight       string      `json:"author_highlight"`       // 作者
	SectionName           string      `json:"section_name"`           // 类型
	Introduction          string      `json:"introduction"`           // 简介
	IntroductionHighlight string      `json:"introduction_highlight"` // 简介高亮
	PublishTime           interface{} `json:"publish_time"`           // 出版时间
	PublishTimeOrder      int         `json:"publish_time_order"`
	CoverImage            string      `json:"cover_image"`        // 封面图
	PublishPeriodical     string      `json:"publish_periodical"` // 总刊期
	PublicationTitle      string      `json:"publication_title"`  // 出版物题目
	Publisher             string      `json:"publisher"`          //出版社
	Year                  string      `json:"year"`               // 年
	Country               string      `json:"country"`            // 国家
	Status                int         `json:"status"`
	Content               string      `json:"content"` //动态资讯中的简介
	DeleteTime            int         `json:"delete_time"`
	Keywords              []string    `json:"keywords"` // 关键词
}

type SearchResponseBody struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source    SearchResponseItem  `json:"_source"`
			Highlight map[string][]string `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]struct {
		Buckets []struct {
			Key      interface{} `json:"key"`
			DocCount float64     `json:"doc_count"`
		} `json:"buckets"`
	} `json:"aggregations"`
}

type SearchSectionItem struct {
	ID    uint   `json:"id"`    // 栏目/类型ID
	Name  string `json:"name"`  // 栏目/类型名称
	Count int64  `json:"count"` // 总数
}

type SearchLibTypeItem struct {
	Id    int    `json:"id"`     //id
	Name  string `json:"name"`   // 名称
	EName string `json:"e_name"` // 英文名称
	Count int64  `json:"count"`  // 总数
}

type SearchYearItem struct {
	Name  string `json:"name"`  // 名称
	Count int64  `json:"count"` // 总数
}

type SearchCountryItem struct {
	Name  string `json:"name"`  // 名称
	Count int64  `json:"count"` // 总数
}

type SearchKeywordItem struct {
	Keyword string `json:"keyword"` // 关键词
	Count   int64  `json:"count"`   // 总数
}
