package mongodb

//TODO: Implement a more dynamic way of building the pipes
import (
	"context"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	executionOrder int
	pipelineType   pipelineType
	primitive      primitive.D
}

type projectField struct {
	field       string
	projectedAs string
}

type addField struct {
	field string
	value interface{}
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
	Custom pipelineType = iota
	AddFields
	Count
	Limit
	Lookup
	Match
	MatchField
	Page
	Project
	ProjectField
	Skip
	Sort
	SortField
	SortAfter
	SortAfterField
	Unwind
)

func (s pipelineType) String() string {
	return toPipelineTypeString[s]
}

func (s pipelineType) FromString(key string) pipelineType {
	return toPipelineTypeID[key]
}

var toPipelineTypeString = map[pipelineType]string{
	AddFields:      "ADD_FIELDS",
	Count:          "COUNT",
	Custom:         "CUSTOM",
	Lookup:         "LOOKUP",
	Limit:          "LIMIT",
	Match:          "MATCH",
	MatchField:     "MATCH_FIELD",
	Page:           "PAGE",
	Project:        "PROJECT",
	ProjectField:   "PROJECT_FIELD",
	Skip:           "SKIP",
	Sort:           "SORT",
	SortField:      "SORT_FIELD",
	SortAfter:      "SORT_AFTER",
	SortAfterField: "SORT_AFTER_FIELD",
	Unwind:         "UNWIND",
}

var toPipelineTypeID = map[string]pipelineType{
	"ADD_FIELDS":       AddFields,
	"COUNT":            Count,
	"CUSTOM":           Custom,
	"LIMIT":            Limit,
	"LOOKUP":           Lookup,
	"MATCH":            Match,
	"MATCH_FIELD":      MatchField,
	"PAGE":             Page,
	"PROJECT":          Project,
	"PROJECT_FIELD":    ProjectField,
	"SKIP":             Skip,
	"SORT":             Sort,
	"SORT_FIELD":       SortField,
	"SORT_AFTER":       SortAfter,
	"SORT_AFTER_FIELD": SortAfterField,
	"UNWIND":           Unwind,
}

type PipelineOptions struct {
	IncludeAddFields  bool
	IncludeCount      bool
	IncludeLimit      bool
	IncludeMatch      bool
	IncludeProjection bool
	IncludeSkip       bool
	IncludeSort       bool
	IncludeUser       bool
	IncludeUnwind     bool
}

type PipelineBuilder struct {
	options          PipelineOptions
	collection       *mongo.Collection
	pipelines        []Pipeline
	filters          []filter
	sortingFields    []sortField
	sortingEndFields []sortField
	projectedFields  []projectField
	addFields        []addField
}

// NewEmptyPipeline Creates a new pipeline builder for a specific collection
func NewEmptyPipeline(collection *mongoCollection) *PipelineBuilder {
	builder := PipelineBuilder{}
	builder.pipelines = make([]Pipeline, 0)
	builder.projectedFields = make([]projectField, 0)
	builder.filters = make([]filter, 0)
	builder.sortingFields = make([]sortField, 0)
	builder.sortingEndFields = make([]sortField, 0)
	builder.addFields = make([]addField, 0)
	builder.collection = collection.coll
	builder.options = PipelineOptions{
		IncludeAddFields:  true,
		IncludeCount:      true,
		IncludeLimit:      true,
		IncludeMatch:      true,
		IncludeProjection: true,
		IncludeSkip:       true,
		IncludeSort:       true,
		IncludeUser:       true,
		IncludeUnwind:     true,
	}

	return &builder
}

// Add Adds a user custom pipeline to the builder, this can be any valid mongo pipeline
func (pipelineBuilder *PipelineBuilder) Add(pipeline bson.D) *PipelineBuilder {
	pipelineEntry := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Custom,
		primitive:      pipeline,
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

	pipelineEntry := Pipeline{
		executionOrder: 9999,
		pipelineType:   Page,
		primitive: bson.D{
			{
				Key:   "$skip",
				Value: skip,
			},
			{
				Key:   "$limit",
				Value: pageSize,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, pipelineEntry)
	return pipelineBuilder
}

// CountPipeline Gets the pipeline count based in the current set of pipelines, this will take into
// consideration any filtering done by the user but not any system pipelines
func (pipelineBuilder *PipelineBuilder) CountPipeline() int {
	currentPipelines := pipelineBuilder.pipelines
	ctx := context.Background()
	options := PipelineOptions{
		IncludeCount:      true,
		IncludeMatch:      true,
		IncludeProjection: true,
		IncludeUser:       true,
	}
	pipeline := pipelineBuilder.Count().buildPipeline(options)
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, *pipeline)

	if err != nil {
		return -1
	}
	var element []map[string]interface{}
	err = cursor.All(ctx, &element)
	if err != nil || len(element) == 0 {
		return 0
	}

	pipelineBuilder.pipelines = currentPipelines
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
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Count,
		primitive: bson.D{
			{
				Key:   "$count",
				Value: "count",
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, countPipeline)

	return pipelineBuilder
}

// Match Adds a Match pipeline using a complex interface, this allows finer tunning in the required
// query, there is no validation of the passed value and it needs to be a valid query.
// Use the filter pipeline if you want to pass an odata type of query
func (pipelineBuilder *PipelineBuilder) Match(value interface{}) *PipelineBuilder {
	matchPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Match,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: value,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchPipeline)

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

	hasPipeline, _ := pipelineBuilder.has(MatchField)
	if !hasPipeline {
		matchFieldPipeline := Pipeline{
			executionOrder: pipelineBuilder.getNextIndex(),
			pipelineType:   MatchField,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchFieldPipeline)
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
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Match,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: parsedFilter,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, matchPipeline)

	return pipelineBuilder
}

// Project Adds a projection pipeline using a complex interface, this allows finer tunning in the required
// query, there is no validation of the passed value and it needs to be a valid query.
// Use ProjectField to build it using individual fields
func (pipelineBuilder *PipelineBuilder) Project(fields interface{}) *PipelineBuilder {
	projectPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Project,
		primitive: bson.D{
			{
				Key:   "$project",
				Value: fields,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, projectPipeline)

	return pipelineBuilder
}

// Lookup Adds a lookup pipeline to join other collections into this one as subobjects
func (pipelineBuilder *PipelineBuilder) Lookup(from string, localField string, foreignField string, fieldAs string) *PipelineBuilder {
	lookupPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Lookup,
		primitive: bson.D{
			{
				Key: "$lookup",
				Value: bson.M{
					"from":         from,
					"localField":   localField,
					"foreignField": foreignField,
					"as":           fieldAs,
				},
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, lookupPipeline)

	return pipelineBuilder
}

// Unwind Adds an Unwind pipeline to flatten an array
func (pipelineBuilder *PipelineBuilder) Unwind(path string) *PipelineBuilder {
	lookupPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Lookup,
		primitive: bson.D{
			{
				Key: "$unwind",
				Value: bson.M{
					"path": path,
				},
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, lookupPipeline)

	return pipelineBuilder
}

// UnwindWidthIndex Adds an Unwind pipeline to flatten an array and includes the index in the object
func (pipelineBuilder *PipelineBuilder) UnwindWidthIndex(path string, includeArrayIndex string, preserveNullAndEmptyArrays bool) *PipelineBuilder {
	lookupPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Lookup,
		primitive: bson.D{
			{
				Key: "$unwind",
				Value: bson.M{
					"path":                       path,
					"includeArrayIndex":          includeArrayIndex,
					"preserveNullAndEmptyArrays": preserveNullAndEmptyArrays,
				},
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, lookupPipeline)

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

	hasPipeline, _ := pipelineBuilder.has(ProjectField)
	if !hasPipeline {
		projectFieldPipeline := Pipeline{
			executionOrder: pipelineBuilder.getNextIndex(),
			pipelineType:   ProjectField,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, projectFieldPipeline)
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

	hasPipeline, _ := pipelineBuilder.has(ProjectField)
	if !hasPipeline {
		projectFieldPipeline := Pipeline{
			executionOrder: pipelineBuilder.getNextIndex(),
			pipelineType:   ProjectField,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, projectFieldPipeline)
	}

	return pipelineBuilder
}

// AddField Adds a field to the result
func (pipelineBuilder *PipelineBuilder) AddField(field string, value interface{}) *PipelineBuilder {
	has, index := pipelineBuilder.hasAddField(field)
	if !has {
		addField := addField{
			field: field,
			value: value,
		}
		pipelineBuilder.addFields = append(pipelineBuilder.addFields, addField)
	} else {
		pipelineBuilder.addFields[index].value = value
	}

	hasPipeline, _ := pipelineBuilder.has(AddFields)
	if !hasPipeline {
		addFieldsPipeline := Pipeline{
			executionOrder: pipelineBuilder.getNextIndex(),
			pipelineType:   AddFields,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, addFieldsPipeline)
	}

	return pipelineBuilder
}

// Skip adds a skip pipeline, if the skip is lower than 0 then no pipeline is added
func (pipelineBuilder *PipelineBuilder) Skip(skip int) *PipelineBuilder {
	if skip < 0 {
		return pipelineBuilder
	}

	skipPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Skip,
		primitive: bson.D{
			{
				Key:   "$skip",
				Value: skip,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, skipPipeline)

	return pipelineBuilder
}

// Limit Adds a limit pipeline, if the limit is lower or equal than 0 then tno pipeline is added
func (pipelineBuilder *PipelineBuilder) Limit(limit int) *PipelineBuilder {
	if limit <= 0 {
		return pipelineBuilder
	}

	limitPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Limit,
		primitive: bson.D{
			{
				Key:   "$limit",
				Value: limit,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, limitPipeline)

	return pipelineBuilder
}

// Adds a sort pipeline, this allows fields to be sorted, be wary of the order sort is placed
// during the pipeline construction, this might impact the results
func (pipelineBuilder *PipelineBuilder) Sort(fields interface{}) *PipelineBuilder {
	sortPipeline := Pipeline{
		executionOrder: pipelineBuilder.getNextIndex(),
		pipelineType:   Sort,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortPipeline)

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

	hasPipeline, _ := pipelineBuilder.has(SortField)
	if !hasPipeline {
		sortFieldPipeline := Pipeline{
			executionOrder: pipelineBuilder.getNextIndex(),
			pipelineType:   SortField,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortFieldPipeline)
	}

	return pipelineBuilder
}

// SortAtEnd Adds a sort pipeline that runs at the very end of each pipeline, this is commonly used
// for example during odata process where we might have a inner sorting but we want to sort it based
// on that query
func (pipelineBuilder *PipelineBuilder) SortAtEnd(fields interface{}) *PipelineBuilder {
	sortPipeline := Pipeline{
		executionOrder: 9998,
		pipelineType:   SortAfter,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortPipeline)

	return pipelineBuilder
}

// SortBy SortAtEnd Adds a sort pipeline that runs at the very end of each pipeline, this is commonly used
// for example during odata process where we might have a inner sorting but we want to sort it based
// on that query using fields to sort, the order of the fields will be kept when building
// the pipeline and any repeating field will only update the order
func (pipelineBuilder *PipelineBuilder) SortByAtEnd(field string, order mongoSort) *PipelineBuilder {
	has, index := pipelineBuilder.hasSortingEndField(field)
	if !has {
		sortedField := sortField{
			field: field,
			order: order,
		}
		pipelineBuilder.sortingEndFields = append(pipelineBuilder.sortingEndFields, sortedField)
	} else {
		pipelineBuilder.sortingEndFields[index].order = order
	}

	hasPipeline, _ := pipelineBuilder.has(SortAfterField)
	if !hasPipeline {
		sortByAfterFieldPipeline := Pipeline{
			executionOrder: 9998,
			pipelineType:   SortAfterField,
		}

		pipelineBuilder.pipelines = append(pipelineBuilder.pipelines, sortByAfterFieldPipeline)
	}

	return pipelineBuilder
}

// WithOptions Sets the pipeline aggregation options, this allows for a more granular execution
// of some of the internal pipelines, for example execute a count with only the users assigned
// pipelines and none of the other ones
func (pipelineBuilder *PipelineBuilder) WithOptions(options PipelineOptions) *PipelineBuilder {
	pipelineBuilder.options = options
	return pipelineBuilder
}

// FilterByUserPipelines Sets the pipeline filter to only execute the pipelines created by the
// user, this allows the execution to skip any internal ones like MATCH or Project
func (pipelineBuilder *PipelineBuilder) FilterByUserPipelines() *PipelineBuilder {
	options := PipelineOptions{
		IncludeAddFields:  false,
		IncludeCount:      false,
		IncludeLimit:      false,
		IncludeMatch:      false,
		IncludeProjection: false,
		IncludeSkip:       false,
		IncludeSort:       false,
		IncludeUser:       true,
		IncludeUnwind:     false,
	}

	pipelineBuilder.options = options
	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Aggregate() (*mongoCursor, error) {
	ctx := context.Background()

	pipeline := pipelineBuilder.buildPipeline()
	cursor, err := pipelineBuilder.collection.Aggregate(ctx, *pipeline)

	return &mongoCursor{cursor: cursor}, err
}

func (pipelineBuilder *PipelineBuilder) buildPipeline(options ...PipelineOptions) *bson.A {
	var builderOptions PipelineOptions
	if len(options) == 0 {
		builderOptions = pipelineBuilder.options
	} else {
		builderOptions = options[0]
	}

	pipelines := bson.A{}
	if builderOptions.IncludeMatch {
		// Processing the custom filter pipeline and adding it to the pipe
		pipelineBuilder.getFilteredPipeline()
	}
	if builderOptions.IncludeProjection {
		// Processing the custom projection pipeline and adding it to the pipe
		pipelineBuilder.getProjectedFieldsPipeline()
	}
	if builderOptions.IncludeSort {
		// Processing the custom sorting pipeline and adding it to the pipe
		pipelineBuilder.getSortingFieldsPipeline()
		pipelineBuilder.getSortingAfterFieldsPipeline()
	}
	if builderOptions.IncludeAddFields {
		// Processing the addfields pipeline
		pipelineBuilder.getAddFieldsPipeline()
	}

	pipelineBuilder.sortPipelines()
	for _, pipeline := range pipelineBuilder.pipelines {
		if (pipeline.pipelineType == Count) && builderOptions.IncludeCount {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if (pipeline.pipelineType == Limit) && builderOptions.IncludeLimit {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if (pipeline.pipelineType == Match ||
			pipeline.pipelineType == MatchField) && builderOptions.IncludeMatch {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if pipeline.pipelineType == Page {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if (pipeline.pipelineType == Project ||
			pipeline.pipelineType == ProjectField) && builderOptions.IncludeProjection {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if (pipeline.pipelineType == Skip) && builderOptions.IncludeSkip {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if (pipeline.pipelineType == Sort ||
			pipeline.pipelineType == SortAfter ||
			pipeline.pipelineType == SortField ||
			pipeline.pipelineType == SortAfterField) && builderOptions.IncludeSort {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if pipeline.pipelineType == Unwind && builderOptions.IncludeUnwind {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if pipeline.pipelineType == AddFields && builderOptions.IncludeAddFields {
			pipelines = append(pipelines, pipeline.primitive)
		}

		if pipeline.pipelineType == Custom {
			pipelines = append(pipelines, pipeline.primitive)
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) has(key pipelineType) (bool, int) {
	for index, pipeline := range pipelineBuilder.pipelines {
		if key == pipeline.pipelineType {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) getNextIndex() int {
	return len(pipelineBuilder.pipelines)
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

func (pipelineBuilder *PipelineBuilder) hasAddField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.addFields {
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

func (pipelineBuilder *PipelineBuilder) hasSortingEndField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.sortingEndFields {
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
		stringFilter = getOperationString(filter.field, filter.operation, filter.value)

		filterParser := NewFilterParser(stringFilter)
		parsedFilter, err := filterParser.Parse()

		if err == nil {
			fields[filter.field] = parsedFilter.(primitive.M)[filter.field]
		}
	}

	matchPipeline := Pipeline{
		pipelineType: MatchField,
		primitive: bson.D{
			{
				Key:   "$match",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(MatchField)
	if has {
		pipelineBuilder.pipelines[index].primitive = matchPipeline.primitive
		return true
	}

	return false
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
		pipelineType: ProjectField,
		primitive: bson.D{
			{
				Key:   "$project",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(ProjectField)
	if has {
		pipelineBuilder.pipelines[index].primitive = projectPipeline.primitive
		return true
	}

	return false
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
		pipelineType: SortField,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(SortField)
	if has {
		pipelineBuilder.pipelines[index].primitive = sortingPipeline.primitive
		return true
	}

	return false
}

func (pipelineBuilder *PipelineBuilder) getSortingAfterFieldsPipeline() bool {
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
		pipelineType: SortAfterField,
		primitive: bson.D{
			{
				Key:   "$sort",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(SortAfterField)
	if has {
		pipelineBuilder.pipelines[index].primitive = sortingPipeline.primitive
		return true
	}

	return false
}

func (pipelineBuilder *PipelineBuilder) sortPipelines() {
	sort.SliceStable(pipelineBuilder.pipelines, func(i, j int) bool {
		return pipelineBuilder.pipelines[i].executionOrder < pipelineBuilder.pipelines[j].executionOrder
	})
}

func (pipelineBuilder *PipelineBuilder) getAddFieldsPipeline() bool {
	if len(pipelineBuilder.addFields) == 0 {
		return false
	}

	fields := primitive.D{}

	for _, addField := range pipelineBuilder.addFields {
		primitiveField := bson.D{
			{
				Key:   addField.field,
				Value: addField.value,
			},
		}
		fields = append(fields, primitiveField...)
	}
	addPipeline := Pipeline{
		pipelineType: AddFields,
		primitive: bson.D{
			{
				Key:   "$addFields",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.has(AddFields)
	if has {
		pipelineBuilder.pipelines[index].primitive = addPipeline.primitive
		return true
	}

	return false
}
