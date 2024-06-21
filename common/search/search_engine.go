package search

const SearchEngineBehaviorMatch uint8 = 1 // 匹配
const SearchEngineBehaviorTerms uint8 = 2 // 精确

type SearchEngineLangMetadataType string

const (
	SearchEngineLangMetadataTypeYear          = SearchEngineLangMetadataType("year")           // 年
	SearchEngineLangMetadataTypeSubject       = SearchEngineLangMetadataType("subject")        // 学科分类
	SearchEngineLangMetadataTypeAuthorId      = SearchEngineLangMetadataType("author_id")      // 作者ID
	SearchEngineLangMetadataTypeAuthorCn      = SearchEngineLangMetadataType("author_cn")      // 作者
	SearchEngineLangMetadataTypeAuthorEn      = SearchEngineLangMetadataType("author_en")      // 作者英文
	SearchEngineLangMetadataTypeAffiliationId = SearchEngineLangMetadataType("affiliation_id") // 机构ID
	SearchEngineLangMetadataTypeAffiliationCn = SearchEngineLangMetadataType("affiliation_cn") // 机构
	SearchEngineLangMetadataTypeAffiliationEn = SearchEngineLangMetadataType("affiliation_en") // 机构
	SearchEngineLangMetadataTypeKeywordId     = SearchEngineLangMetadataType("keyword_id")     // 关键词ID
	SearchEngineLangMetadataTypeKeywordCn     = SearchEngineLangMetadataType("keyword_cn")     // 关键词
	SearchEngineLangMetadataTypeKeywordEn     = SearchEngineLangMetadataType("keyword_en")     // 关键词英文
	SearchEngineLangMetadataTypeJournal       = SearchEngineLangMetadataType("journal")        // 期刊
	SearchEngineLangMetadataTypeCategoryCn    = SearchEngineLangMetadataType("category_cn")    // 栏目中文
	SearchEngineLangMetadataTypeCategoryEn    = SearchEngineLangMetadataType("category_en")    // 栏目英文
	SearchEngineLangMetadataTypeFund          = SearchEngineLangMetadataType("fund")           // 基金
)

type SearchEngineLangRequestMetadataOrder string

const (
	SearchEngineLangRequestMetadataOrderDefault = SearchEngineLangRequestMetadataOrder("default")
	SearchEngineLangRequestMetadataOrderKeyAsc  = SearchEngineLangRequestMetadataOrder("key_asc")
	SearchEngineLangRequestMetadataOrderKeyDesc = SearchEngineLangRequestMetadataOrder("key_desc")
)

type SearchEngineLangRequestMetadata struct {
	Type  SearchEngineLangMetadataType         // 聚合类型
	Order SearchEngineLangRequestMetadataOrder // 排序条件
	Tag   string                               // 自定义标签，相当于聚合的名字
	Size  int                                  // 聚合结果的数量
}

type SearchEngineLangResponseMetadataItem struct {
	Name  string // 聚合的项名称
	Total int64  // 项的统计数量
}

type SearchEngineLangResponseMetadata struct {
	Type SearchEngineLangMetadataType           // 聚合类型
	Tag  string                                 // 自定义标签，相当于聚合的名字
	Data []SearchEngineLangResponseMetadataItem // 数据
}

type SearchEngineGroupType string

const (
	SearchEngineGroupTypeAll          = SearchEngineGroupType("all")          // 全部可用字段
	SearchEngineGroupTypeFC           = SearchEngineGroupType("full_content") // 全文
	SearchEngineGroupTypeTitle        = SearchEngineGroupType("title")        // 标题
	SearchEngineGroupTypeAuthor       = SearchEngineGroupType("author")       // 作者
	SearchEngineGroupTypeIntroduction = SearchEngineGroupType("introduction") // 简介
	SearchEngineGroupTypeContent      = SearchEngineGroupType("content")      // 内容
)

type SearchEngineGroupOrder string

const (
	SearchEngineGroupOrderDefault     = SearchEngineGroupOrder("id_desc")      // 默认排序
	SearchEngineGroupOrderScoreAsc    = SearchEngineGroupOrder("score_asc")    // 评分降序
	SearchEngineGroupOrderScoreDesc   = SearchEngineGroupOrder("score_desc")   // 评分降序
	SearchEngineGroupOrderPublishAsc  = SearchEngineGroupOrder("publish_asc")  // 发布时间升序
	SearchEngineGroupOrderPublishDesc = SearchEngineGroupOrder("publish_desc") // 发布时间降序
)

type SearchEngineGroupMode string

const (
	SearchEngineGroupModeAnd = SearchEngineGroupMode("and") // 与
	SearchEngineGroupModeOr  = SearchEngineGroupMode("or")  // 或
	SearchEngineGroupModeNot = SearchEngineGroupMode("not") // 非
)

type SearchEngineGroupBehavior string

const (
	SearchEngineGroupBehaviorMatch       = SearchEngineGroupBehavior("match")        // 匹配
	SearchEngineGroupBehaviorMatchPhrase = SearchEngineGroupBehavior("match_phrase") // 短语匹配
	SearchEngineGroupBehaviorTerms       = SearchEngineGroupBehavior("terms")        // 精确
	SearchEngineGroupBehaviorPrefix      = SearchEngineGroupBehavior("prefix")       // 前缀
	SearchEngineGroupBehaviorRange       = SearchEngineGroupBehavior("range")        // 范围
)

type SearchEngineGroup struct {
	Type     SearchEngineGroupType     // 类型
	Mode     SearchEngineGroupMode     // 组合模式
	Behavior SearchEngineGroupBehavior // 搜索行为
	Values   []string                  // 非range行为下的value
	RangeMin int64                     // range行为时的下限，指针为nil时不设置下限
	RangeMax int64                     // range行为时的上限，指针为nil时不设置上限
	Boost    int                       // 评分
	Sub      []SearchEngineGroup       // 子条件组，子条件组不为空时，其他字段仅mode字段有效
}
