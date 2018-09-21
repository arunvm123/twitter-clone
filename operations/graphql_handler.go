package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/arunvm/twitter-clone/pkg/auth"
	"github.com/arunvm/twitter-clone/pkg/mysql"
	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"user_name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"post_id": &graphql.Field{
			Type: graphql.Int,
		},
		"text": &graphql.Field{
			Type: graphql.String,
		},
		"user_name": &graphql.Field{
			Type: graphql.Int,
		},
		"time_stamp": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"getOwnPosts": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Returns all posts of the user",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userName, err := auth.ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Token not Valid: %v", err)
				}

				db := params.Info.RootValue.(map[string]interface{})["db"].(*mysql.MySQL)

				posts, err := db.GetAllPostsOfUser(userName.(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Could not retrieve posts: %v", err)
				}

				return posts, nil
			},
		},
		"getUsers": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "Returns list of users",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userName, err := auth.ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Token not Valid: %v", err)
				}

				db := params.Info.RootValue.(map[string]interface{})["db"].(*mysql.MySQL)

				users, err := db.GetAllUsers(userName.(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Could not retrieve posts: %v", err)
				}

				return users, nil
			},
		},
		"getPostFeed": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "Returns list of posts of people that the user follows",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userName, err := auth.ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Token not Valid: %v", err)
				}

				db := params.Info.RootValue.(map[string]interface{})["db"].(*mysql.MySQL)

				posts, err := db.GetPostsFeed(userName.(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Could not retrieve feed: %v", err)
				}

				return posts, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"addPost": &graphql.Field{
			Type:        postType,
			Description: "Inserts a new post into db",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				userName, err := auth.ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Token not Valid: %v", err)
				}

				db := params.Info.RootValue.(map[string]interface{})["db"].(*mysql.MySQL)
				post, err := db.AddPost(userName.(string), params.Args["text"].(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Could not insert post: %v", err)
				}

				return post, nil
			},
		},
		"followUser": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Follows the specified user",
			Args: graphql.FieldConfigArgument{
				"user_name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userName, err := auth.ValidateJWT(params.Context.Value("token").(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Token not Valid: %v", err)
				}

				db := params.Info.RootValue.(map[string]interface{})["db"].(*mysql.MySQL)
				err = db.FollowUser(userName.(string), params.Args["user_name"].(string))
				if err != nil {
					log.Println(err)
					return nil, fmt.Errorf("Could not insert post: %v", err)
				}

				return true, nil
			},
		},
	},
})

func executeQuery(query, token string, schem graphql.Schema, db *mysql.MySQL) *graphql.Result {
	rootObject := make(map[string]interface{})
	rootObject["db"] = db

	result := graphql.Do(graphql.Params{
		Schema:        schem,
		RequestString: query,
		Context:       context.WithValue(context.Background(), "token", token),
		RootObject:    rootObject,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result

}

func graphqlHandler(db *mysql.MySQL) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var schema, err = graphql.NewSchema(graphql.SchemaConfig{
				Query:    rootQuery,
				Mutation: rootMutation,
			})

			if err != nil {
				panic(err)
			}

			result := executeQuery(r.URL.Query().Get("query"), r.Header.Get("token"), schema, db)
			if len(result.Errors) > 0 {
				log.Println(result.Errors)
				respondWithError(w, http.StatusInternalServerError, "Couldnt process request")
			}
			json.NewEncoder(w).Encode(result.Data)
		},
	)
}
