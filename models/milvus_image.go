package models

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"gocv.io/x/gocv"
	"google.golang.org/appengine/log"
)

var (
	databaseName   = "Test"
	collectionName = "Images"
	// ImageSchema entity.Field
)

type ImageEntity struct {
	SHA256      string
	Descriptors []byte
	KeyPoints   []float32
	Path        string
	descDim     int
	kpDim       int
}

var ImageSchema = &entity.Schema{
	CollectionName: collectionName,
	Description:    "Test book search",
	Fields: []*entity.Field{
		{
			Name:        "image_sha256",
			Description: "Image sha256 hash",
			DataType:    entity.FieldTypeVarChar,
			PrimaryKey:  true,
			AutoID:      false,
			TypeParams: map[string]string{
				"max_length": "64",
			},
		},
		{
			Name:        "descriptors",
			Description: "Descriptors from ORB",
			DataType:    entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "32768",
			},
		},
		// {
		// 	Name:        "keypoints",
		// 	Description: "KeyPoints from ORB",
		// 	DataType:    entity.FieldTypeFloatVector,
		// 	TypeParams: map[string]string{
		// 		"dim": "7",
		// 	},
		// },
	},
	EnableDynamicField: true,
}

func NewImageEntity(id string, descriptors []byte, keypoints []float32, path string) *ImageEntity {
	imageEntity := &ImageEntity{SHA256: id, Descriptors: descriptors, KeyPoints: keypoints, Path: path}
	return imageEntity
}

func (i *ImageEntity) toImgMat() *gocv.Mat {
	imgMat := gocv.IMRead(i.Path, gocv.IMReadColor)
	return &imgMat
}

func (i *ImageEntity) toMat() *gocv.Mat {
	imgMat := i.toImgMat()
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(*imgMat, &gray, gocv.ColorBGRToGray)
	return imgMat
}

func (i *ImageEntity) toKeypointsDescriptors(orb *gocv.ORB) ([]gocv.KeyPoint, gocv.Mat) {
	imgMat := i.toMat()
	defer imgMat.Close()
	kp, desc := orb.DetectAndCompute(*imgMat, gocv.NewMat())
	return kp, desc
}

func NewImageEntityFromFile(path string, orb *gocv.ORB) (*ImageEntity, error) {
	imgEntity := &ImageEntity{Path: path, descDim: 32768, kpDim: 7}
	kp, desc := imgEntity.toKeypointsDescriptors(orb)
	kpFloats := KeyPointsToFloat32(&kp)
	imgEntity.KeyPoints = kpFloats
	// descVector := MatToBinaryVector(&desc)
	imgEntity.Descriptors = FormatDescriptors(&desc)
	fmt.Println(len(imgEntity.Descriptors))
	uuid, err := FileToUUID(path)
	if err != nil {
		log.Errorf(context.Background(), "Error getting UUID from file %s", path)
		return nil, err
	}
	imgEntity.SHA256 = uuid

	return imgEntity, nil
}

func (i *ImageEntity) ToColumns() *[]entity.Column {
	hash := entity.NewColumnVarChar("image_sha256", []string{i.SHA256})
	descHolder := make([][]byte, 1)
	descHolder[0] = i.Descriptors
	descriptors := entity.NewColumnBinaryVector("descriptors", i.descDim, descHolder)
	return &[]entity.Column{hash, descriptors}
}

func (i *ImageEntity) Insert(m client.Client) error {
	cols := *i.ToColumns()
	column, err := m.Insert(
		context.Background(),
		collectionName,
		"",
		cols[0],
		cols[1],
	)
	if err != nil {
		return err
	}
	fmt.Printf("Insert result: %v\n", column)
	return nil
}
