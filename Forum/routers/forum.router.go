package routers

import (
	"exemple_api/controllers"

	"github.com/gorilla/mux"
)

// RegisterForumRoutes enregistre les routes du forum
func RegisterForumRoutes(r *mux.Router, forumController *controllers.ForumControllers) {
	// CATEGORIES
	r.HandleFunc("/forum/categories", forumController.GetAllCategories).Methods("GET")
	r.HandleFunc("/forum/categories", forumController.CreateCategory).Methods("POST")
	r.HandleFunc("/forum/categories/{id}", forumController.GetCategoryByID).Methods("GET")

	// TOPICS
	r.HandleFunc("/forum/topics", forumController.GetAllTopics).Methods("GET")
	r.HandleFunc("/forum/topics", forumController.CreateTopic).Methods("POST")
	r.HandleFunc("/forum/topics/search", forumController.Search).Methods("GET")
	r.HandleFunc("/forum/topics/popular", forumController.GetPopularTopics).Methods("GET")
	r.HandleFunc("/forum/topics/{topicId}", forumController.GetTopicByID).Methods("GET")
	r.HandleFunc("/forum/topics/{topicId}", forumController.UpdateTopic).Methods("PUT")
	r.HandleFunc("/forum/topics/{topicId}", forumController.DeleteTopic).Methods("DELETE")
	r.HandleFunc("/forum/categories/{categoryId}/topics", forumController.GetTopicsByCategory).Methods("GET")

	// COMMENTS
	r.HandleFunc("/forum/topics/{topicId}/comments", forumController.GetCommentsByTopic).Methods("GET")
	r.HandleFunc("/forum/topics/{topicId}/comments", forumController.CreateComment).Methods("POST")
	r.HandleFunc("/forum/comments/{commentId}", forumController.UpdateComment).Methods("PUT")
	r.HandleFunc("/forum/comments/{commentId}", forumController.DeleteComment).Methods("DELETE")

	// LIKES
	r.HandleFunc("/forum/topics/{topicId}/like", forumController.LikeTopic).Methods("POST")
	r.HandleFunc("/forum/comments/{commentId}/like", forumController.LikeComment).Methods("POST")
}
