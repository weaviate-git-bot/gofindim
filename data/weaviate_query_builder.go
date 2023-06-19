package data

import (
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type QueryBuilder struct {
	*graphql.GetBuilder
	Client  *weaviate.Client
	Results []map[string]interface{}
}

func NewQueryBuilder(client *weaviate.Client) *QueryBuilder {
	builder := client.GraphQL().Get()

	qb := &QueryBuilder{GetBuilder: builder, Client: client}
	return qb
}

func (qb *QueryBuilder) NearImage(i *ImageFile, distance float32) *QueryBuilder {
	vector, err := i.ToVector()
	if err != nil {
		panic(err)
	}
	nearVector := qb.Client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithDistance(distance)
	qb.WithNearVector(nearVector)
	return qb
}

func (qb *QueryBuilder) NearText(text string, distance float32) *QueryBuilder {
	vector, err := VectorizeText(text)
	if err != nil {
		log.Fatal(err)
	}
	nearVector := qb.Client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithDistance(distance)
	qb.WithNearVector(nearVector)
	return qb
}

func (qb *QueryBuilder) NearVector(vector []float32, distance float32) *QueryBuilder {
	nearVector := qb.Client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithDistance(distance)
	qb.WithNearVector(nearVector)
	return qb
}

// func (qb *QueryBuilder) SelectFields(fields []string) *QueryBuilder {
// 	var queryFields []graphql.Field
// 	for _, fieldName := range fields {
// 		if fieldName == "id" {
// 			field := graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}}}
// 			queryFields = append(queryFields, field)
// 		} else {
// 			field := graphql.Field{Name: fieldName}
// 			queryFields = append(queryFields, field)
// 		}
// 	}
// 	qb.WithFields(queryFields...)
// 	return qb
// }

func (qb *QueryBuilder) SelectFields(fields []string) *QueryBuilder {
	var queryFields []graphql.Field
	additionalField := graphql.Field{Name: "_additional"}
	var additionalFields []graphql.Field
	for _, fieldName := range fields {
		var field graphql.Field
		if fieldName != "id" && fieldName != "distance" {
			println("field name: " + fieldName)
			field = graphql.Field{Name: fieldName}
			queryFields = append(queryFields, field)
		}
		if fieldName == "id" {
			additionalFields = append(additionalFields, graphql.Field{Name: "id"})
		} else if fieldName == "distance" {
			additionalFields = append(additionalFields, graphql.Field{Name: "distance"})
		}

	}
	if len(additionalFields) > 0 {
		additionalField.Fields = additionalFields
		queryFields = append(queryFields, additionalField)

	}
	qb.WithFields(queryFields...)
	return qb
}

// func (qb *QueryBuilder) GetResults() (*models.GraphQLResponse, error) {
// 	return qb.Do(context.Background())

// }
