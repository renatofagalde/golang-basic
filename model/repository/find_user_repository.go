package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang-basic/config/logger"
	"golang-basic/config/rest_err"
	"golang-basic/model"
	"golang-basic/model/repository/entity"
	"golang-basic/model/repository/entity/convert"
	"os"
)

func (ur *userRepository) FindUserByEmail(email string) (model.UserDomainInterface, *rest_err.RestErr) {
	logger.Info("init findUserByEmail user repository", zap.String("journey", "findUserByEmail"))

	collection := ur.databaseConnection.Collection(os.Getenv(MONGO_DB_COLLECTION))

	userEntity := &entity.UserEntity{}

	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(context.Background(), filter).Decode(userEntity)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMessage := fmt.Sprintf("Usuário não localizado com este email: %s", email)
			logger.Error(errorMessage, err, zap.String("journey", "findUserByEmail"))
			return nil, rest_err.NewNotFoundError(errorMessage)
		}
		errorMessage := "Erro ao pesquisar usuário"
		logger.Error(errorMessage, err, zap.String("journey", "findUserByEmail"))
		return nil, rest_err.NewInternalServerError(errorMessage)
	}

	logger.Info("init findUserByEmail user repository successfuly",
		zap.String("journey", "findUserByEmail"),
		zap.String("email", email),
		zap.String("userId", userEntity.ID.Hex()),
	)
	return convert.ConvertEntityToDomain(*userEntity), nil
}
