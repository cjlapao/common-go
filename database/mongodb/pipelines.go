package mongodb

//TODO: Refactor implementation
import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	Type      pipelineType
	primitive primitive.D
}

type projectField struct {
	field       string
	projectedAs string
}

type sortField struct {
	field string
	order mongoSort
}

type mongoSort int

const (
	Asc  mongoSort = 1
	Desc mongoSort = -1
)

type pipelineType int

const (
	User pipelineType = iota
	Count
	Match
	Project
	ProjectField
	Skip
	Limit
	Sort
	SortAfter
)

func (s pipelineType) String() string {
	return toPipelineTypeString[s]
}

func (s pipelineType) FromString(key string) pipelineType {
	return toPipelineTypeID[key]
}

var toPipelineTypeString = map[pipelineType]string{
	Count:        "COUNT",
	Limit:        "LIMIT",
	Match:        "MATCH",
	Project:      "PROJECT",
	ProjectField: "PROJECT_FIELD",
	Skip:         "SKIP",
	Sort:         "SORT",
	SortAfter:    "SORT_AFTER",
	User:         "USER",
}

var toPipelineTypeID = map[string]pipelineType{
	"COUNT":      Count,
	"LIMIT":      Limit,
	"MATCH":      Match,
	"PROJECT":    Project,
	"SKIP":       Skip,
	"SORT":       Sort,
	"SORT_AFTER": SortAfter,
	"USER":       User,
}

type PipelineBuilder struct {
	collection      *mongo.Collection
	pipelines       []Pipeline
	filters         []filter
	sortingFields   []sortField
	projectedFields []projectField
}

// NewPipelineBuilder Creates a new pipeline builder for a specific collection
func NewPipelineBuilder(collection *mongoCollection) *PipelineBuilder {
	builder := PipelineBuilder{}
	builder.pipelines = make([]Pipeline, 0)
	builder.projectedFields = make([]projectField, 0)
	builder.filters = make([]filter, 0)
	builder.sortingFields = make([]sortField, 0)
	builder.collection = collection.coll

	return &builder
}

// Add Adds a user custom pipeline to the builder
func (pipelineBuilder *PipelineBuilder) Add(pipeline bson.D) *PipelineBuilder {
	pipelineEntry := Pipeline{
		Type:      User,
		primitive: pipeline,
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, pipelineEntry)
	return pipelineBuilder
}

// Page Creates a paging pipeline based in the page number and page size, this will
// generate a $skip and a $limit pipeline if the page and skip are filled in
func (pipelineBuilder *PipelineBuilder) Page(page int, pageSize int) *PipelineBuilder {
	if page == -1 || pageSize <= 0 {
		return pipelineBuilder
	}

	skip := 0
	if page > 0 {
		skip = page * pageSize
	}

	pipelineBuilder.Skip(skip)
	pipelineBuilder.Limit(pageSize)

	return pipelineBuilder
}

// Gets the pipeline count based in the current set of pipelines, this will take into
// consideration any filtering done by the user but not any system pipelines
func (pipelineBuilder *PipelineBuilder) GetCount() int {
	ctx := context.Background()
	cursor, err := pipelineBuilder.Count().AggregateUserWithCount()
	if err != nil {
		return -1
	}
	var element []map[string]interface{}
	err = cursor.cursor.All(ctx, &element)
	if err != nil || len(element) == 0 {
		return 0
	}

	return int(element[0]["count"].(int32))
}

// CountCollection This will count the pipeline collection excluding anything from the pipelines
// we can use this for odata responses or to count how many objects the collection has
func (pipelineBuilder *PipelineBuilder) CountCollection() int {
	ctx := context.Background()
	countDocument := bson.D{
		{
			Key:   "$count",
			Value: "count",
		},
	}

	pipeline := bson.A{}
	pipeline = append(pipeline, countDocument)
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return -1
	}
	var element []map[string]interface{}
	err = cursor.All(ctx, &element)
	if err != nil || len(element) == 0 {
		return 0
	}

	return int(element[0]["count"].(int32))
}

// Count Adds a count pipeline, this will overseed any other pipeline and will always return a count
// pipeline response with a value
func (pipelineBuilder *PipelineBuilder) Count() *PipelineBuilder {
	countPipeline := Pipeline{
		Type: Count,
		primitive: bson.D{
			{
				Key:   "$count",
				Value: "count",
			},
		},
	}

	has, index := pipelineBuilder.has(Count)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, countPipeline)
	} else {
		pipelineBuilder.pipelines[index] = countPipeline
	}

	return pipelineBuilder
}

// Match Adds a Match pipeline using a complex interface, this allows finer tunning in the required
// query, there is no validation of the passed value and it needs to be a valid query.
// Use the filter pipeline if you want to pass an odata type of query
func (pipelineBuilder *PipelineBuilder) Match(value interface{}) *PipelineBuilder {
	matchPipeline := Pipeline{
		Type: Match,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: value,
			},
		},
	}

	has, index := pipelineBuilder.has(Match)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchPipeline)
	} else {
		pipelineBuilder.pipelines[index] = matchPipeline
	}

	return pipelineBuilder
}

// FilterBy Adds a Match pipeline for a specific field, this allow simple operations always with
// and joining fields if they are more than one, use the Filter pipeline to pass in a more complex
// odata type of query for a better control of the filtering
func (pipelineBuilder *PipelineBuilder) FilterBy(field string, operation filterOperation, value interface{}) *PipelineBuilder {
	has, index := pipelineBuilder.hasFilteredField(field)
	if !has {
		filter := filter{
			field:     field,
			operation: operation,
			value:     value,
		}

		pipelineBuilder.filters = append(pipelineBuilder.filters, filter)
	} else {
		pipelineBuilder.filters[index].value = value
	}

	return pipelineBuilder
}

// Filter Adds a Match pipeline using a odata type of query for a easier and powerful query language
// If the expressing fails to parse it will not be added to the pipeline
func (pipelineBuilder *PipelineBuilder) Filter(filter string) *PipelineBuilder {
	filterParser := NewFilterParser(filter)
	parsedFilter, err := filterParser.Parse()

	if err != nil {
		return pipelineBuilder
	}
	matchPipeline := Pipeline{
		Type: Match,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: parsedFilter,
			},
		},
	}

	has, index := pipelineBuilder.has(Match)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchPipeline)
	} else {
		pipelineBuilder.pipelines[index] = matchPipeline
	}

	return pipelineBuilder
}

// Project Adds a projection pipeline using a complex interface, this allows finer tunning in the required
// query, there is no validation of the passed value and it needs to be a valid query.
// Use ProjectField to build it using individual fields
func (pipelineBuilder *PipelineBuilder) Project(fields interface{}) *PipelineBuilder {
	projectPipeline := Pipeline{
		Type: Project,
		primitive: bson.D{
			{
				Key:   "$project",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(Project)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, projectPipeline)
	} else {
		pipelineBuilder.pipelines[index] = projectPipeline
	}

	return pipelineBuilder
}

// ProjectField adds a field to be projected by the pipeline, this allows to easily build a projection
// of complex fields without adding a pre built interface.
func (pipelineBuilder *PipelineBuilder) ProjectField(field string) *PipelineBuilder {
	has, index := pipelineBuilder.hasProjectedField(field)
	if !has {
		projected := projectField{
			field:       field,
			projectedAs: field,
		}
		pipelineBuilder.projectedFields = append(pipelineBuilder.projectedFields, projected)
	} else {
		pipelineBuilder.projectedFields[index].projectedAs = field
	}

	return pipelineBuilder
}

// ProjectFieldAs adds a field to be projected by the pipeline and assign a different property name,
// this allows to easily build a projection of complex fields without adding a pre built interface.
func (pipelineBuilder *PipelineBuilder) ProjectFieldAs(field string, projectedAs string) *PipelineBuilder {
	has, index := pipelineBuilder.hasProjectedField(field)
	if !has {
		projected := projectField{
			field:       field,
			projectedAs: projectedAs,
		}
		pipelineBuilder.projectedFields = append(pipelineBuilder.projectedFields, projected)
	} else {
		pipelineBuilder.projectedFields[index].projectedAs = projectedAs
	}

	return pipelineBuilder
}

// Skip adds a skip pipeline, if the skip is lower than 0 then no pipeline is added
func (pipelineBuilder *PipelineBuilder) Skip(skip int) *PipelineBuilder {
	if skip < 0 {
		return pipelineBuilder
	}

	skipPipeline := Pipeline{
		Type: Skip,
		primitive: bson.D{
			{
				Key:   "$skip",
				Value: skip,
			},
		},
	}

	has, index := pipelineBuilder.has(Skip)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, skipPipeline)
	} else {
		pipelineBuilder.pipelines[index] = skipPipeline
	}

	return pipelineBuilder
}

// Limit Adds a limit pipeline, if the limit is lower or equal than 0 then tno pipeline is added
func (pipelineBuilder *PipelineBuilder) Limit(limit int) *PipelineBuilder {
	if limit <= 0 {
		return pipelineBuilder
	}

	limitPipeline := Pipeline{
		Type: Limit,
		primitive: bson.D{
			{
				Key:   "$limit",
				Value: limit,
			},
		},
	}

	has, index := pipelineBuilder.has(Limit)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, limitPipeline)
	} else {
		pipelineBuilder.pipelines[index] = limitPipeline
	}

	return pipelineBuilder
}

// Adds a sort pipeline, this allows fields to be sorted, be wary of the order sort is placed
// during the pipeline construction, this might impact the results
func (pipelineBuilder *PipelineBuilder) Sort(fields interface{}) *PipelineBuilder {
	sortPipeline := Pipeline{
		Type: Sort,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(Sort)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortPipeline)
	} else {
		pipelineBuilder.pipelines[index] = sortPipeline
	}

	return pipelineBuilder
}

// SortBy Adds a Sort pipeline using fields to sort, the order of the fields will be kept when building
// the pipeline and any repeating field will only update the order
func (pipelineBuilder *PipelineBuilder) SortBy(field string, order mongoSort) *PipelineBuilder {
	has, index := pipelineBuilder.hasSortingField(field)
	if !has {
		sortedField := sortField{
			field: field,
			order: order,
		}
		pipelineBuilder.sortingFields = append(pipelineBuilder.sortingFields, sortedField)
	} else {
		pipelineBuilder.sortingFields[index].order = order
	}

	return pipelineBuilder
}

// SortAtEnd Adds a sort pipeline that runs at the very end of each pipeline, this is commonly used
// for example during odata process where we might have a inner sorting but we want to sort it based
// on that query
func (pipelineBuilder *PipelineBuilder) SortAtEnd(fields interface{}) *PipelineBuilder {
	sortPipeline := Pipeline{
		Type: SortAfter,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(SortAfter)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortPipeline)
	} else {
		pipelineBuilder.pipelines[index] = sortPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) SortByAtEnd(field string, order mongoSort) *PipelineBuilder {
	sortPipeline := Pipeline{
		Type: SortAfter,
		primitive: bson.D{
			{
				Key: "$sort",
				Value: bson.D{
					{
						Key:   field,
						Value: order,
					},
				},
			},
		},
	}

	has, index := pipelineBuilder.has(SortAfter)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortPipeline)
	} else {
		pipelineBuilder.pipelines[index] = sortPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) GetUserPipeline() *bson.A {
	pipelines := bson.A{}
	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == User {
			pipelines = append(pipelines, pipeline.primitive)
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) GetUserPipelineWithCount() *bson.A {
	pipelines := bson.A{}
	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == User {
			pipelines = append(pipelines, pipeline.primitive)
		}
	}

	// Appending Match  if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Match {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending Match  if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Project {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending ProjectFields pipeline if it exists
	pipelineBuilder.getProjectedFieldsPipeline()

	// Appending the count pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Count {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) get() *bson.A {
	pipelines := bson.A{}
	// Processing the custom filter pipeline and adding it to the pipe
	pipelineBuilder.getFilteredPipeline()
	// Processing the custom projection pipeline and adding it to the pipe
	pipelineBuilder.getProjectedFieldsPipeline()
	// Processing the custom sorting pipeline and adding it to the pipe
	pipelineBuilder.getSortingFieldsPipeline()

	// Appending first the sort if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Sort {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == User {
			pipelines = append(pipelines, pipeline.primitive)
		}
	}

	// Appending first the sort if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == SortAfter {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending Match if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Match {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending Projection if it exists
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Project {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending the skip pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Skip {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending the limit pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Limit {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	// Appending the count pipelines if they exist
	for _, pipeline := range pipelineBuilder.pipelines {
		if pipeline.Type == Count {
			pipelines = append(pipelines, pipeline.primitive)
			break
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) Aggregate() (*mongoCursor, error) {
	ctx := context.Background()
	pipeline := pipelineBuilder.get()
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, *pipeline)

	return &mongoCursor{cursor: cursor}, err
}

func (pipelineBuilder *PipelineBuilder) AggregateUser() (*mongoCursor, error) {
	ctx := context.Background()
	pipeline := pipelineBuilder.GetUserPipeline()
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, *pipeline)

	return &mongoCursor{cursor: cursor}, err
}

func (pipelineBuilder *PipelineBuilder) AggregateUserWithCount() (*mongoCursor, error) {
	ctx := context.Background()
	pipeline := pipelineBuilder.GetUserPipelineWithCount()
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, *pipeline)

	return &mongoCursor{cursor: cursor}, err
}

func (pipelineBuilder *PipelineBuilder) has(key pipelineType) (bool, int) {
	for index, pipeline := range pipelineBuilder.pipelines {
		if key == pipeline.Type {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) hasFilteredField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.filters {
		if strings.EqualFold(field.field, fieldName) {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) hasProjectedField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.projectedFields {
		if strings.EqualFold(field.field, fieldName) {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) hasSortingField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.sortingFields {
		if strings.EqualFold(field.field, fieldName) {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) getFilteredPipeline() bool {
	if len(pipelineBuilder.filters) == 0 {
		return false
	}

	fields := primitive.M{}

	var stringFilter string
	for _, filter := range pipelineBuilder.filters {
		switch filter.operation {
		case "gt":
			stringFilter = fmt.Sprintf("%v gt '%v'", filter.field, filter.value)
		case "ge":
			stringFilter = fmt.Sprintf("%v ge '%v'", filter.field, filter.value)
		case "lt":
			stringFilter = fmt.Sprintf("%v lt '%v'", filter.field, filter.value)
		case "le":
			stringFilter = fmt.Sprintf("%v le '%v'", filter.field, filter.value)
		case "eq":
			stringFilter = fmt.Sprintf("%v eq '%v'", filter.field, filter.value)
		case "ne":
			stringFilter = fmt.Sprintf("%v ne '%v'", filter.field, filter.value)
		case "regex":
			stringFilter = fmt.Sprintf("%v regex %v", filter.field, filter.value)
		case "contains":
			stringFilter = fmt.Sprintf("contains(%v, '%v')", filter.field, filter.value)
		case "endswith":
			stringFilter = fmt.Sprintf("endswith(%v, '%v')", filter.field, filter.value)
		case "startswith":
			stringFilter = fmt.Sprintf("startswith(%v, '%v')", filter.field, filter.value)
		}

		filterParser := NewFilterParser(stringFilter)
		parsedFilter, err := filterParser.Parse()

		if err == nil {
			fields[filter.field] = parsedFilter.(primitive.M)[filter.field]
		}
	}

	matchPipeline := Pipeline{
		Type: Match,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(Match)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchPipeline)
	} else {
		pipelineBuilder.pipelines[index] = matchPipeline
	}

	return true
}

func (pipelineBuilder *PipelineBuilder) getProjectedFieldsPipeline() bool {
	if len(pipelineBuilder.projectedFields) == 0 {
		return false
	}

	fields := primitive.D{}

	for _, projectedField := range pipelineBuilder.projectedFields {
		primitiveField := bson.D{
			{
				Key:   projectedField.projectedAs,
				Value: "$" + projectedField.field,
			},
		}
		fields = append(fields, primitiveField...)
	}
	projectPipeline := Pipeline{
		Type: Project,
		primitive: bson.D{
			{
				Key:   "$project",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(Project)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, projectPipeline)
	} else {
		pipelineBuilder.pipelines[index] = projectPipeline
	}

	return true
}

func (pipelineBuilder *PipelineBuilder) getSortingFieldsPipeline() bool {
	if len(pipelineBuilder.sortingFields) == 0 {
		return false
	}

	fields := primitive.D{}

	for _, sortingField := range pipelineBuilder.sortingFields {
		primitiveField := bson.D{
			{
				Key:   sortingField.field,
				Value: sortingField.order,
			},
		}
		fields = append(fields, primitiveField...)
	}
	sortingPipeline := Pipeline{
		Type: Sort,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(Sort)
	if !has {
		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortingPipeline)
	} else {
		pipelineBuilder.pipelines[index] = sortingPipeline
	}

	return true
}
