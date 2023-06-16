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

func (qb *QueryBuilder) NearImage(i *ImageFile, certainty float32) *QueryBuilder {
	vector, err := i.ToVector()
	if err != nil {
		log.Fatal(err)
	}
	nearVector := qb.Client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	qb.WithNearVector(nearVector)
	return qb
}

func (qb *QueryBuilder) NearText(text string, certainty float32) *QueryBuilder {
	vector, err := VectorizeText(text)
	if err != nil {
		log.Fatal(err)
	}
	nearVector := qb.Client.GraphQL().NearVectorArgBuilder().WithVector(vector).WithCertainty(certainty)
	qb.WithNearVector(nearVector)
	return qb
}

func (qb *QueryBuilder) SelectFields(fields []string) *QueryBuilder {
	var qlFields []graphql.Field
	for _, fieldName := range fields {
		field := graphql.Field{Name: fieldName}
		qlFields = append(qlFields, field)
	}
	qb.WithFields(qlFields...)
	return qb
}

// func (qb *QueryBuilder) GetResults() (*models.GraphQLResponse, error) {
// 	return qb.Do(context.Background())

// }
