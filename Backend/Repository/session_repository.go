package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	SessionCollection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		SessionCollection: db.Collection("sessions"),
	}
}

// SaveTokens saves the access and refresh tokens to the database
func (sr *SessionRepository) SaveTokens(userID string, accessToken string, refreshToken string) error {
	session := bson.M{}

	existingSessionerr := sr.SessionCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&session)
	if existingSessionerr != mongo.ErrNoDocuments {
		_,err := sr.SessionCollection.DeleteOne(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			return err
		}
	}

	_, err := sr.SessionCollection.InsertOne(context.Background(), bson.M{
		"user_id":       userID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if err != nil {
		return err
	}
	return nil
}


// DeleteTokens deletes the access and refresh tokens from the database
func (sr *SessionRepository) DeleteTokens(userID string) error {
	_, err := sr.SessionCollection.DeleteOne(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	return nil
}

// GetTokens retrieves the access and refresh tokens from the database
func (sr *SessionRepository) GetTokens(userID string) (string, string, error) {
	var session bson.M
	err := sr.SessionCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&session)
	if err != nil {
		return "", "", err
	}
	return session["access_token"].(string), session["refresh_token"].(string), nil
}