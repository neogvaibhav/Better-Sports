package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Basu008/Better-ESPN/auth"
	"github.com/Basu008/Better-ESPN/helpers"
	"github.com/Basu008/Better-ESPN/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const playerCollection = "player"

func CreatePlayer(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	//Set the heder meaning the type of data that is being used ie. JSON
	w.Header().Set("Content-Type", "application/json")
	//Then we will fetch the JSON body that the user must've provided
	authToken := r.Header.Get("Authorization")
	_, authErr := auth.VerifyAuthToken(authToken)
	if authErr != nil {
		CreateNewResponse(w, http.StatusUnauthorized, &Response{false, "Auth token required", nil})
		return
	}
	var player model.Player
	var requestBody helpers.CreatePlayerRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		//This will let the user know about to misformed JSON body
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Issue with JSON Body", nil})
		return
	}
	if !requestBody.IsCreatePlayerRequestBodyValid() {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Inputs given are invalid", nil})
		return
	}
	//Now we set the data
	player.Name = requestBody.Name
	player.PlayerProfile = &model.PlayerProfile{
		Grade:    requestBody.Grade,
		Position: requestBody.Position,
		Stats: &model.Stats{
			Goals:   0,
			Assists: 0,
			Fouls:   0,
		},
	}
	player.CreatedAt = time.Now().UTC()

	//Now, if there is no error, we will send the player data to the DB
	result, err := db.Collection(playerCollection).InsertOne(context.Background(), player)

	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Player can't be create. Try Again!", nil})
		return
	}

	//Finally, we create a response body for success state

	player.ID = result.InsertedID.(primitive.ObjectID)
	CreateNewResponse(w, http.StatusCreated, &Response{true, "", player})
}

func GetAllPlayers(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//First we get the cursor
	var players []model.Player
	//To get the page number
	page := r.URL.Query().Get("page")

	//We create an option now to apply pagination
	//skip will be used to skip the documents
	skip, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Page number should be an integer", nil})
		return
	}
	if skip < 0 {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Page number should be positive", nil})
		return
	}
	skip = skip * 20
	var limit int64 = 20
	options := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	cursor, err := db.Collection(playerCollection).Find(context.TODO(), bson.D{}, &options)
	defer func() {
		err := cursor.Close(context.Background())
		if err != nil {
			log.Printf("Couldn't close cursor")
		}
	}()
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Coudln't find data", nil})
		return
	}
	cursorErr := cursor.All(context.Background(), &players)
	if cursorErr != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Coudln't find data", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", players})
}

func GetPlayerById(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//all the quries from the end point
	params := mux.Vars(r)
	//Then we get the id from the map of queries
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid player ID", nil})
		return
	}
	//then we check in the db for the existance of document
	var player model.Player
	//How to get a document based on ID
	mongoErr := db.Collection(playerCollection).FindOne(context.Background(), model.Player{ID: id}).Decode(&player)
	if mongoErr != nil {
		switch mongoErr {
		case mongo.ErrNoDocuments:
			CreateNewResponse(w, http.StatusNotFound, &Response{false, "No player with this id exists", nil})
		default:
			CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		}
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", player})

}

func UpdatePlayer(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//First we will check whether the body that we are recieving is valid or not
	var playerSchema helpers.Player
	err := json.NewDecoder(r.Body).Decode(&playerSchema)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "JSON body incorrect", nil})
		return
	}
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid player ID", nil})
		return
	}
	var fieldsToBeUpdated primitive.M = primitive.M{}
	if playerSchema.Name != "" {
		fieldsToBeUpdated["name"] = playerSchema.Name
	}
	if playerSchema.Grade != "" {
		if helpers.IsGradeValid(playerSchema.Grade) {
			fieldsToBeUpdated["player_profile.grade"] = playerSchema.Grade
		} else {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid player grade", nil})
			return
		}
	}
	if playerSchema.Position != "" {
		if helpers.IsPositionValid(playerSchema.Position) {
			fieldsToBeUpdated["player_profile.position"] = playerSchema.Position
		} else {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid player position", nil})
			return
		}
	}
	if playerSchema.Goals != nil {
		goals := *(playerSchema.Goals)
		if goals < 0 {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Goals can't be negative", nil})
			return
		}
		fieldsToBeUpdated["player_profile.stats.goals"] = goals
	}
	if playerSchema.Assists != nil {
		assists := *(playerSchema.Assists)
		if assists < 0 {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Assists can't be negative", nil})
			return
		}
		fieldsToBeUpdated["player_profile.stats.assists"] = assists
	}
	if playerSchema.Fouls != nil {
		fouls := *(playerSchema.Fouls)
		if fouls < 0 {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Fouls can't be negative", nil})
			return
		}
		fieldsToBeUpdated["player_profile.stats.fouls"] = fouls
	}
	filter := bson.M{"_id": id}
	options := bson.M{
		"$set": fieldsToBeUpdated,
	}

	result, err := db.Collection(playerCollection).UpdateOne(context.Background(), filter, options)
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "There was some issue", nil})
		return
	}

	if result.MatchedCount == 1 {
		CreateNewResponse(w, http.StatusOK, &Response{true, "", true})
		return
	}

	CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Can't be updated!", nil})

}

func DeletePlayer(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Here we will only delete one player
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid player ID", nil})
		return
	}
	_, deletionError := db.Collection(playerCollection).DeleteOne(context.Background(), model.Player{ID: id})
	if deletionError != nil {
		switch deletionError {
		case mongo.ErrNoDocuments:
			CreateNewResponse(w, http.StatusNotFound, &Response{false, "No player with this id exists", nil})
		default:
			CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		}
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", true})

}

func GetTopScorers(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var players = []model.Player{}
	sortOption := options.Find().SetSort(bson.D{
		{Key: "player_profile.stats.goals", Value: -1},
		{Key: "name", Value: 1},
	}).SetLimit(5)
	cursor, err := db.Collection(playerCollection).Find(context.Background(), bson.D{}, sortOption)
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		return
	}
	if err := cursor.All(context.TODO(), &players); err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", players})

}

func GetPlayersFromPositions(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var players = []model.Player{}
	position := r.URL.Query().Get("position")
	if !helpers.IsPositionValid(position) {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid Position", nil})
		return
	}
	filter := bson.D{{Key: "player_profile.position", Value: position}}
	cur, err := db.Collection(playerCollection).Find(context.TODO(), filter)
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Can't fetch documents", nil})
		return
	}
	cursorError := cur.All(context.TODO(), &players)
	if cursorError != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Ill formed data", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", players})
}
