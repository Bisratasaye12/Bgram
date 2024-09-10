package repository

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VerificationURLRepository struct {
	UrlCollection *mongo.Collection
}

func NewVerificationURLRepository(db *mongo.Database) interfaces.VerificationURLRepositoryInterface {
	return &VerificationURLRepository{
		UrlCollection: db.Collection("verification_urls"),
	}
}

func (vr *VerificationURLRepository) SaveUrl(url *models.VerificationURL) (*models.VerificationURL, error) {
	url_id := url.UrlID
	log.Println(url)
	_, err := vr.UrlCollection.InsertOne(context.Background(), bson.M{
		"url_id": url_id,
		"url":    url.URL,
	})
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (vr *VerificationURLRepository) GetUrlByID(urlID string) (*models.VerificationURL, error) {
	url := &models.VerificationURL{}

	err := vr.UrlCollection.FindOne(context.Background(), bson.M{"url_id": urlID}).Decode(url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (vr *VerificationURLRepository) DeleteUrlByID(urlID string) error {
	_, err := vr.UrlCollection.DeleteOne(context.Background(), bson.M{"url_id": urlID})
	if err != nil {
		return err
	}
	return nil
}
