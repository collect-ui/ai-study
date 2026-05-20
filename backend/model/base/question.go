package base

const TableNameQuestionItem = "question_item"

type QuestionItem struct {
	QuestionID        string `gorm:"column:question_id;primaryKey" json:"question_id"`
	QuestionCode      string `gorm:"column:question_code" json:"question_code"`
	Title             string `gorm:"column:title" json:"title"`
	Subject           string `gorm:"column:subject" json:"subject"`
	Stage             string `gorm:"column:stage" json:"stage"`
	Grade             string `gorm:"column:grade" json:"grade"`
	TextbookVersion   string `gorm:"column:textbook_version" json:"textbook_version"`
	UnitID            string `gorm:"column:unit_id" json:"unit_id"`
	UnitCode          string `gorm:"column:unit_code" json:"unit_code"`
	UnitName          string `gorm:"column:unit_name" json:"unit_name"`
	QuestionType      string `gorm:"column:question_type" json:"question_type"`
	QuestionCategory  string `gorm:"column:question_category" json:"question_category"`
	Difficulty        string `gorm:"column:difficulty" json:"difficulty"`
	Score             int    `gorm:"column:score" json:"score"`
	DurationSeconds   int    `gorm:"column:duration_seconds" json:"duration_seconds"`
	SequenceNo        int    `gorm:"column:sequence_no" json:"sequence_no"`
	StemHTML          string `gorm:"column:stem_html" json:"stem_html"`
	StemText          string `gorm:"column:stem_text" json:"stem_text"`
	AnalysisHTML      string `gorm:"column:analysis_html" json:"analysis_html"`
	AnalysisText      string `gorm:"column:analysis_text" json:"analysis_text"`
	AnalysisMediaURL  string `gorm:"column:analysis_media_url" json:"analysis_media_url"`
	AnalysisMediaName string `gorm:"column:analysis_media_name" json:"analysis_media_name"`
	AnalysisMediaType string `gorm:"column:analysis_media_type" json:"analysis_media_type"`
	OptionCount       int    `gorm:"column:option_count" json:"option_count"`
	BlankCount        int    `gorm:"column:blank_count" json:"blank_count"`
	AssetCount        int    `gorm:"column:asset_count" json:"asset_count"`
	ContentHash       string `gorm:"column:content_hash" json:"content_hash"`
	Source            string `gorm:"column:source" json:"source"`
	Status            string `gorm:"column:status" json:"status"`
	Version           int    `gorm:"column:version" json:"version"`
	PublishTime       string `gorm:"column:publish_time" json:"publish_time"`
	PublishUser       string `gorm:"column:publish_user" json:"publish_user"`
	LastReviewID      string `gorm:"column:last_review_id" json:"last_review_id"`
	Remark            string `gorm:"column:remark" json:"remark"`
	IsDelete          string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime        string `gorm:"column:create_time" json:"create_time"`
	CreateUser        string `gorm:"column:create_user" json:"create_user"`
	ModifyTime        string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser        string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionItem) TableName() string {
	return TableNameQuestionItem
}

func (*QuestionItem) PrimaryKey() []string {
	return []string{"question_id"}
}

const TableNameQuestionOption = "question_option"

type QuestionOption struct {
	OptionID    string `gorm:"column:option_id;primaryKey" json:"option_id"`
	QuestionID  string `gorm:"column:question_id" json:"question_id"`
	OptionKey   string `gorm:"column:option_key" json:"option_key"`
	OptionOrder int    `gorm:"column:option_order" json:"option_order"`
	ContentMode string `gorm:"column:content_mode" json:"content_mode"`
	OptionHTML  string `gorm:"column:option_html" json:"option_html"`
	OptionText  string `gorm:"column:option_text" json:"option_text"`
	IsCorrect   string `gorm:"column:is_correct" json:"is_correct"`
	AssetCount  int    `gorm:"column:asset_count" json:"asset_count"`
	IsDelete    string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime  string `gorm:"column:create_time" json:"create_time"`
	CreateUser  string `gorm:"column:create_user" json:"create_user"`
	ModifyTime  string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser  string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionOption) TableName() string {
	return TableNameQuestionOption
}

func (*QuestionOption) PrimaryKey() []string {
	return []string{"option_id"}
}

const TableNameQuestionAnswer = "question_answer"

type QuestionAnswer struct {
	AnswerID         string `gorm:"column:answer_id;primaryKey" json:"answer_id"`
	QuestionID       string `gorm:"column:question_id" json:"question_id"`
	AnswerType       string `gorm:"column:answer_type" json:"answer_type"`
	AnswerValue      string `gorm:"column:answer_value" json:"answer_value"`
	AnswerText       string `gorm:"column:answer_text" json:"answer_text"`
	ReferenceText    string `gorm:"column:reference_text" json:"reference_text"`
	CaseSensitive    string `gorm:"column:case_sensitive" json:"case_sensitive"`
	AllowOrderChange string `gorm:"column:allow_order_change" json:"allow_order_change"`
	AutoGradingRule  string `gorm:"column:auto_grading_rule" json:"auto_grading_rule"`
	IsDelete         string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime       string `gorm:"column:create_time" json:"create_time"`
	CreateUser       string `gorm:"column:create_user" json:"create_user"`
	ModifyTime       string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser       string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionAnswer) TableName() string {
	return TableNameQuestionAnswer
}

func (*QuestionAnswer) PrimaryKey() []string {
	return []string{"answer_id"}
}

const TableNameQuestionBlankAnswer = "question_blank_answer"

type QuestionBlankAnswer struct {
	BlankAnswerID     string `gorm:"column:blank_answer_id;primaryKey" json:"blank_answer_id"`
	QuestionID        string `gorm:"column:question_id" json:"question_id"`
	BlankIndex        int    `gorm:"column:blank_index" json:"blank_index"`
	StandardAnswer    string `gorm:"column:standard_answer" json:"standard_answer"`
	AlternativeAnswer string `gorm:"column:alternative_answers" json:"alternative_answers"`
	Score             int    `gorm:"column:score" json:"score"`
	MatchMode         string `gorm:"column:match_mode" json:"match_mode"`
	CaseSensitive     string `gorm:"column:case_sensitive" json:"case_sensitive"`
	IsDelete          string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime        string `gorm:"column:create_time" json:"create_time"`
	CreateUser        string `gorm:"column:create_user" json:"create_user"`
	ModifyTime        string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser        string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionBlankAnswer) TableName() string {
	return TableNameQuestionBlankAnswer
}

func (*QuestionBlankAnswer) PrimaryKey() []string {
	return []string{"blank_answer_id"}
}

const TableNameQuestionScoringPoint = "question_scoring_point"

type QuestionScoringPoint struct {
	ScoringPointID string `gorm:"column:scoring_point_id;primaryKey" json:"scoring_point_id"`
	QuestionID     string `gorm:"column:question_id" json:"question_id"`
	PointIndex     int    `gorm:"column:point_index" json:"point_index"`
	PointText      string `gorm:"column:point_text" json:"point_text"`
	Score          int    `gorm:"column:score" json:"score"`
	Keywords       string `gorm:"column:keywords" json:"keywords"`
	IsRequired     string `gorm:"column:is_required" json:"is_required"`
	IsDelete       string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime     string `gorm:"column:create_time" json:"create_time"`
	CreateUser     string `gorm:"column:create_user" json:"create_user"`
	ModifyTime     string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser     string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionScoringPoint) TableName() string {
	return TableNameQuestionScoringPoint
}

func (*QuestionScoringPoint) PrimaryKey() []string {
	return []string{"scoring_point_id"}
}

const TableNameQuestionGrade = "question_grade"

type QuestionGrade struct {
	GradeID    string `gorm:"column:grade_id;primaryKey" json:"grade_id"`
	Stage      string `gorm:"column:stage" json:"stage"`
	GradeCode  string `gorm:"column:grade_code" json:"grade_code"`
	GradeName  string `gorm:"column:grade_name" json:"grade_name"`
	OrderIndex int    `gorm:"column:order_index" json:"order_index"`
	Status     string `gorm:"column:status" json:"status"`
	IsDelete   string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime string `gorm:"column:create_time" json:"create_time"`
	CreateUser string `gorm:"column:create_user" json:"create_user"`
	ModifyTime string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionGrade) TableName() string {
	return TableNameQuestionGrade
}

func (*QuestionGrade) PrimaryKey() []string {
	return []string{"grade_id"}
}

const TableNameQuestionSubject = "question_subject"

type QuestionSubject struct {
	SubjectID   string `gorm:"column:subject_id;primaryKey" json:"subject_id"`
	SubjectCode string `gorm:"column:subject_code" json:"subject_code"`
	SubjectName string `gorm:"column:subject_name" json:"subject_name"`
	OrderIndex  int    `gorm:"column:order_index" json:"order_index"`
	Status      string `gorm:"column:status" json:"status"`
	IsDelete    string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime  string `gorm:"column:create_time" json:"create_time"`
	CreateUser  string `gorm:"column:create_user" json:"create_user"`
	ModifyTime  string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser  string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionSubject) TableName() string {
	return TableNameQuestionSubject
}

func (*QuestionSubject) PrimaryKey() []string {
	return []string{"subject_id"}
}

const TableNameQuestionUnit = "question_unit"

type QuestionUnit struct {
	UnitID          string `gorm:"column:unit_id;primaryKey" json:"unit_id"`
	Subject         string `gorm:"column:subject" json:"subject"`
	Stage           string `gorm:"column:stage" json:"stage"`
	Grade           string `gorm:"column:grade" json:"grade"`
	TextbookVersion string `gorm:"column:textbook_version" json:"textbook_version"`
	ParentID        string `gorm:"column:parent_id" json:"parent_id"`
	UnitCode        string `gorm:"column:unit_code" json:"unit_code"`
	UnitName        string `gorm:"column:unit_name" json:"unit_name"`
	OrderIndex      int    `gorm:"column:order_index" json:"order_index"`
	Status          string `gorm:"column:status" json:"status"`
	IsDelete        string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime      string `gorm:"column:create_time" json:"create_time"`
	CreateUser      string `gorm:"column:create_user" json:"create_user"`
	ModifyTime      string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser      string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionUnit) TableName() string {
	return TableNameQuestionUnit
}

func (*QuestionUnit) PrimaryKey() []string {
	return []string{"unit_id"}
}

const TableNameQuestionKnowledge = "question_knowledge"

type QuestionKnowledge struct {
	KnowledgeID   string `gorm:"column:knowledge_id;primaryKey" json:"knowledge_id"`
	Subject       string `gorm:"column:subject" json:"subject"`
	Stage         string `gorm:"column:stage" json:"stage"`
	Grade         string `gorm:"column:grade" json:"grade"`
	ParentID      string `gorm:"column:parent_id" json:"parent_id"`
	KnowledgeCode string `gorm:"column:knowledge_code" json:"knowledge_code"`
	KnowledgeName string `gorm:"column:knowledge_name" json:"knowledge_name"`
	OrderIndex    int    `gorm:"column:order_index" json:"order_index"`
	Status        string `gorm:"column:status" json:"status"`
	IsDelete      string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime    string `gorm:"column:create_time" json:"create_time"`
	CreateUser    string `gorm:"column:create_user" json:"create_user"`
	ModifyTime    string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser    string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionKnowledge) TableName() string {
	return TableNameQuestionKnowledge
}

func (*QuestionKnowledge) PrimaryKey() []string {
	return []string{"knowledge_id"}
}

const TableNameQuestionKnowledgeRel = "question_knowledge_rel"

type QuestionKnowledgeRel struct {
	RelID         string `gorm:"column:rel_id;primaryKey" json:"rel_id"`
	QuestionID    string `gorm:"column:question_id" json:"question_id"`
	KnowledgeID   string `gorm:"column:knowledge_id" json:"knowledge_id"`
	KnowledgeName string `gorm:"column:knowledge_name" json:"knowledge_name"`
	OrderIndex    int    `gorm:"column:order_index" json:"order_index"`
	CreateTime    string `gorm:"column:create_time" json:"create_time"`
	CreateUser    string `gorm:"column:create_user" json:"create_user"`
}

func (*QuestionKnowledgeRel) TableName() string {
	return TableNameQuestionKnowledgeRel
}

func (*QuestionKnowledgeRel) PrimaryKey() []string {
	return []string{"rel_id"}
}

const TableNameQuestionAsset = "question_asset"

type QuestionAsset struct {
	AssetID    string `gorm:"column:asset_id;primaryKey" json:"asset_id"`
	QuestionID string `gorm:"column:question_id" json:"question_id"`
	UsageType  string `gorm:"column:usage_type" json:"usage_type"`
	UsageRef   string `gorm:"column:usage_ref" json:"usage_ref"`
	AssetURL   string `gorm:"column:asset_url" json:"asset_url"`
	AssetName  string `gorm:"column:asset_name" json:"asset_name"`
	MimeType   string `gorm:"column:mime_type" json:"mime_type"`
	FileSize   int    `gorm:"column:file_size" json:"file_size"`
	Sha256     string `gorm:"column:sha256" json:"sha256"`
	Status     string `gorm:"column:status" json:"status"`
	IsDelete   string `gorm:"column:is_delete" json:"is_delete"`
	CreateTime string `gorm:"column:create_time" json:"create_time"`
	CreateUser string `gorm:"column:create_user" json:"create_user"`
}

func (*QuestionAsset) TableName() string {
	return TableNameQuestionAsset
}

func (*QuestionAsset) PrimaryKey() []string {
	return []string{"asset_id"}
}

const TableNameQuestionReviewRecord = "question_review_record"

type QuestionReviewRecord struct {
	ReviewID      string `gorm:"column:review_id;primaryKey" json:"review_id"`
	QuestionID    string `gorm:"column:question_id" json:"question_id"`
	FromStatus    string `gorm:"column:from_status" json:"from_status"`
	ToStatus      string `gorm:"column:to_status" json:"to_status"`
	ReviewResult  string `gorm:"column:review_result" json:"review_result"`
	ReviewComment string `gorm:"column:review_comment" json:"review_comment"`
	ReviewUser    string `gorm:"column:review_user" json:"review_user"`
	ReviewTime    string `gorm:"column:review_time" json:"review_time"`
	IsDelete      string `gorm:"column:is_delete" json:"is_delete"`
}

func (*QuestionReviewRecord) TableName() string {
	return TableNameQuestionReviewRecord
}

func (*QuestionReviewRecord) PrimaryKey() []string {
	return []string{"review_id"}
}

const TableNameQuestionChangeLog = "question_change_log"

type QuestionChangeLog struct {
	LogID      string `gorm:"column:log_id;primaryKey" json:"log_id"`
	QuestionID string `gorm:"column:question_id" json:"question_id"`
	OpType     string `gorm:"column:op_type" json:"op_type"`
	BeforeJSON string `gorm:"column:before_json" json:"before_json"`
	AfterJSON  string `gorm:"column:after_json" json:"after_json"`
	OpUser     string `gorm:"column:op_user" json:"op_user"`
	OpTime     string `gorm:"column:op_time" json:"op_time"`
	Remark     string `gorm:"column:remark" json:"remark"`
}

func (*QuestionChangeLog) TableName() string {
	return TableNameQuestionChangeLog
}

func (*QuestionChangeLog) PrimaryKey() []string {
	return []string{"log_id"}
}

const TableNameQuestionImportBatch = "question_import_batch"

type QuestionImportBatch struct {
	BatchID         string `gorm:"column:batch_id;primaryKey" json:"batch_id"`
	FileName        string `gorm:"column:file_name" json:"file_name"`
	FileURL         string `gorm:"column:file_url" json:"file_url"`
	Subject         string `gorm:"column:subject" json:"subject"`
	Stage           string `gorm:"column:stage" json:"stage"`
	Grade           string `gorm:"column:grade" json:"grade"`
	TextbookVersion string `gorm:"column:textbook_version" json:"textbook_version"`
	Status          string `gorm:"column:status" json:"status"`
	TotalCount      int    `gorm:"column:total_count" json:"total_count"`
	SuccessCount    int    `gorm:"column:success_count" json:"success_count"`
	FailCount       int    `gorm:"column:fail_count" json:"fail_count"`
	ErrorSummary    string `gorm:"column:error_summary" json:"error_summary"`
	CreateTime      string `gorm:"column:create_time" json:"create_time"`
	CreateUser      string `gorm:"column:create_user" json:"create_user"`
	ModifyTime      string `gorm:"column:modify_time" json:"modify_time"`
	ModifyUser      string `gorm:"column:modify_user" json:"modify_user"`
}

func (*QuestionImportBatch) TableName() string {
	return TableNameQuestionImportBatch
}

func (*QuestionImportBatch) PrimaryKey() []string {
	return []string{"batch_id"}
}

const TableNameQuestionImportRow = "question_import_row"

type QuestionImportRow struct {
	RowID          string `gorm:"column:row_id;primaryKey" json:"row_id"`
	BatchID        string `gorm:"column:batch_id" json:"batch_id"`
	RowIndex       int    `gorm:"column:row_index" json:"row_index"`
	RawJSON        string `gorm:"column:raw_json" json:"raw_json"`
	ParsedJSON     string `gorm:"column:parsed_json" json:"parsed_json"`
	ValidateStatus string `gorm:"column:validate_status" json:"validate_status"`
	ErrorMsg       string `gorm:"column:error_msg" json:"error_msg"`
	QuestionID     string `gorm:"column:question_id" json:"question_id"`
	CreateTime     string `gorm:"column:create_time" json:"create_time"`
	CreateUser     string `gorm:"column:create_user" json:"create_user"`
}

func (*QuestionImportRow) TableName() string {
	return TableNameQuestionImportRow
}

func (*QuestionImportRow) PrimaryKey() []string {
	return []string{"row_id"}
}
