package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Basu008/Better-ESPN/helpers"
	"github.com/Basu008/Better-ESPN/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const teamsCollection string = "teams"

func CreateTeam(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var teamRequest helpers.Team
	var team model.Team
	decodingErr := json.NewDecoder(r.Body).Decode(&teamRequest)
	if decodingErr != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Issue with JSON Body", nil})
		return
	}
	team.Name = teamRequest.Name
	var playerIds = []primitive.ObjectID{}
	for _, playerId := range teamRequest.Players {
		idInBsonFormat, err := primitive.ObjectIDFromHex(playerId)
		if err != nil {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid Ids", nil})
			return
		}
		playerIds = append(playerIds, idInBsonFormat)
	}
	if len(playerIds) <= 1 {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Add at least 2 players to the team :)", nil})
		return
	}
	team.Players = playerIds
	_, err := db.Collection(teamsCollection).InsertOne(context.Background(), team)
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't add data to the server", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", true})
}

func GetTeamById(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idFromQuery := r.URL.Query().Get("id")
	if idFromQuery == "" {
		GetAllTeams(db, w, r)
		return
	}
	id, err := primitive.ObjectIDFromHex(idFromQuery)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Invalid Id", nil})
		return
	}
	filter := bson.D{
		{Key: "_id", Value: id},
	}
	var team model.Team
	mongoErr := db.Collection(teamsCollection).FindOne(context.TODO(), filter).Decode(&team)
	if mongoErr != nil {
		switch mongoErr {
		case mongo.ErrNoDocuments:
			CreateNewResponse(w, http.StatusNotFound, &Response{false, "No team with this id exists", nil})
		default:
			CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		}
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", team})
}

func GetAllTeams(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var teams = []model.Team{}
	curr, err := db.Collection(teamsCollection).Find(context.TODO(), bson.D{})
	if err != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Error Fetching data", nil})
		return
	}
	defer func() {
		err := curr.Close(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}()
	cursorErr := curr.All(context.TODO(), &teams)
	if cursorErr != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Couldn't fetch documents", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", teams})

}

func AddPlayerToTeam(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateTeamRequestBody helpers.UpdateTeamPlayers
	err := json.NewDecoder(r.Body).Decode(&updateTeamRequestBody)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Badly formed JSON body", nil})
		return
	}
	teamId, err := primitive.ObjectIDFromHex(updateTeamRequestBody.ID)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Incorrect Team Id.", nil})
		return
	}
	if len(updateTeamRequestBody.PlayerIds) < 1 {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Provide at least one player id", nil})
		return
	}
	var playerIds = []primitive.ObjectID{}
	for _, id := range updateTeamRequestBody.PlayerIds {
		playerId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Incorrect Player Id.", nil})
			return
		}
		playerIds = append(playerIds, playerId)
	}
	filter := bson.D{
		{Key: "_id", Value: teamId},
	}
	addPlayerQuery := bson.D{
		{Key: "$each", Value: playerIds},
	}
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "players", Value: addPlayerQuery}}},
	}
	_, updateErr := db.Collection(teamsCollection).UpdateOne(context.TODO(), filter, update)
	if updateErr != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Coudln't update data", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", true})
}

func RemovePlayerFromTeam(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateTeamRequestBody helpers.UpdateTeamPlayers
	err := json.NewDecoder(r.Body).Decode(&updateTeamRequestBody)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "JSON bana bhai shi se", nil})
		return
	}
	teamId, err := primitive.ObjectIDFromHex(updateTeamRequestBody.ID)
	if err != nil {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "Team id galat hai", nil})
		return
	}
	if len(updateTeamRequestBody.PlayerIds) < 1 {
		CreateNewResponse(w, http.StatusBadRequest, &Response{false, "ek player ki id toh de bhai mere", nil})
		return
	}
	var playerIds = []primitive.ObjectID{}
	for _, id := range updateTeamRequestBody.PlayerIds {
		playerId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			CreateNewResponse(w, http.StatusBadRequest, &Response{false, "PLayer id galat hai", nil})
			return
		}
		playerIds = append(playerIds, playerId)
	}
	filter := bson.D{{Key: "_id", Value: teamId}}
	subQuery := bson.D{{Key: "players", Value: bson.D{{Key: "$in", Value: playerIds}}}}
	updateQuery := bson.D{{Key: "$pull", Value: subQuery}}
	_, updaterErr := db.Collection(teamsCollection).UpdateOne(context.TODO(), filter, updateQuery)
	if updaterErr != nil {
		CreateNewResponse(w, http.StatusInternalServerError, &Response{false, "Nhi ho paaya update", nil})
		return
	}
	CreateNewResponse(w, http.StatusOK, &Response{true, "", true})
}
