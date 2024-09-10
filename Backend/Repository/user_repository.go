package repository

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) interfaces.UserRepositoryInterface {
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

// createUser creates a new user in the database
func (ur *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	usr := bson.M{
		"username":        user.Username,
		"password":        user.Password,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
		"bio":             user.Bio,
		"role":            "user",
	}

	// if no user in the system the first user is admin
	noUser, _ := ur.Collection.CountDocuments(context.Background(), bson.M{})
	log.Println("no user", noUser)
	if noUser == int64(0) {
		usr["role"] = "admin"
	}
	_, err := ur.Collection.InsertOne(context.Background(), usr)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// getUserByID retrieves a user from the database by their ID
func (ur *UserRepository) GetUserByID(id string) (*models.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user := &models.User{}
	usr_dto := map[string]interface{}{
		"_id":             primitive.ObjectID{},
		"username":        "",
		"password":        "",
		"email":           "",
		"profile_picture": "",
		"bio":             "",
		"role":            "",
	}
	err = ur.Collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(usr_dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user does not exist")
		}
		return nil, err
	}
	user.ID = usr_dto["_id"].(primitive.ObjectID).Hex()
	user.Username = usr_dto["username"].(string)
	user.Password = usr_dto["password"].(string)
	user.Email = usr_dto["email"].(string)
	user.ProfilePicture = usr_dto["profile_picture"].(string)
	user.Bio = usr_dto["bio"].(string)
	user.Role = usr_dto["role"].(string)

	return user, nil
}

// getUserByEmailOrUsername retrieves a user from the database by their email or username
func (ur *UserRepository) GetUserByEmailOrUsername(email string, username string) (*models.User, error) {
	user := map[string]interface{}{
		"_id":             primitive.ObjectID{},
		"username":        "",
		"password":        "",
		"email":           "",
		"profile_picture": "",
		"bio":             "",
		"role":            "",
	}

	err := ur.Collection.FindOne(context.Background(), bson.M{"$or": []bson.M{{"email": email}, {"username": username}}}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user does not exist")
		}
		return nil, err
	}
	// find the id of the user
	userID := user["_id"].(primitive.ObjectID).Hex()
	log.Println("user id", userID)
	return &models.User{
		ID:             userID,
		Username:       user["username"].(string),
		Password:       user["password"].(string),
		Email:          user["email"].(string),
		ProfilePicture: user["profile_picture"].(string),
		Bio:            user["bio"].(string),
		Role:           user["role"].(string),
	}, nil

}

// updateUser updates a user's profile in the database
func (ur *UserRepository) UpdateUser(user *models.User) (*models.User, error) {
	updateFeilds := bson.M{}

	if user.Username != "" {
		updateFeilds["username"] = user.Username
	}
	if user.Password != "" {
		updateFeilds["password"] = user.Password
	}
	if user.Email != "" {
		updateFeilds["email"] = user.Email
	}

	if user.ProfilePicture != "" {
		updateFeilds["profile_picture"] = user.ProfilePicture
	}
	if user.Bio != "" {
		updateFeilds["bio"] = user.Bio
	}

	userID, Ierr := primitive.ObjectIDFromHex(user.ID)
	if Ierr != nil {
		return nil, Ierr
	}

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": updateFeilds}
	_, err := ur.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	updated_user, _ := ur.GetUserByID(user.ID)
	return updated_user, nil
}

// deleteUser deletes a user from the database
func (ur *UserRepository) DeleteUser(id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = ur.Collection.DeleteOne(context.Background(), bson.M{"_id": userID})
	if err != nil {
		return err
	}
	return nil
}
