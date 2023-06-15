package data

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate/modules/multi2vec-clip/clients"
)

var clipOrigin = "http://127.0.0.1:9090"

func VectorizeImage(img *ImageFile) ([]float32, error) {
	vectorizer := clients.New(clipOrigin, logrus.New())
	res, err := vectorizer.Vectorize(context.Background(), []string{}, []string{img.Base64})
	if err != nil {
		return nil, err
	}
	return res.ImageVectors[0], nil
}

func VectorizeImages(images []*ImageFile) ([][]float32, error) {
	vectorizer := clients.New(clipOrigin, logrus.New())
	texts := make([]string, 0, len(images))
	base64s := make([]string, 0, len(images))
	for _, img := range images {
		base64s = append(base64s, img.Base64)
	}
	res, err := vectorizer.Vectorize(context.Background(), texts, base64s)
	if err != nil {
		return nil, err
	}
	return res.ImageVectors, nil
}

func VectorizeText(text string) ([]float32, error) {
	vectorizer := clients.New(clipOrigin, logrus.New())
	res, err := vectorizer.Vectorize(context.Background(), []string{text}, []string{})
	if err != nil {
		return nil, err
	}
	return res.TextVectors[0], nil
}

func VectorizeTexts(texts []string) ([][]float32, error) {
	vectorizer := clients.New(clipOrigin, logrus.New())
	res, err := vectorizer.Vectorize(context.Background(), texts, []string{})
	if err != nil {
		return nil, err
	}
	return res.TextVectors, nil
}

func VectorizeTextAndImage(text string, img *ImageFile) ([]float32, []float32, error) {
	vectorizer := clients.New(clipOrigin, logrus.New())
	res, err := vectorizer.Vectorize(context.Background(), []string{text}, []string{img.Base64})
	if err != nil {
		return nil, nil, err
	}
	return res.TextVectors[0], res.ImageVectors[0], nil
}
