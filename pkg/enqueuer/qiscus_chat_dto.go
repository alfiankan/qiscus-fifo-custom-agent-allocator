package enqueuer

import "time"

type QiscusWebhookChatReqBody struct {
	AppID          string `json:"app_id"`
	AvatarURL      string `json:"avatar_url"`
	CandidateAgent struct {
		AvatarURL    any       `json:"avatar_url"`
		CreatedAt    time.Time `json:"created_at"`
		Email        string    `json:"email"`
		ForceOffline bool      `json:"force_offline"`
		ID           int       `json:"id"`
		IsAvailable  bool      `json:"is_available"`
		IsVerified   bool      `json:"is_verified"`
		LastLogin    any       `json:"last_login"`
		Name         string    `json:"name"`
		SdkEmail     string    `json:"sdk_email"`
		SdkKey       string    `json:"sdk_key"`
		Type         int       `json:"type"`
		TypeAsString string    `json:"type_as_string"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"candidate_agent"`
	Email         string `json:"email"`
	Extras        string `json:"extras"`
	IsNewSession  bool   `json:"is_new_session"`
	IsResolved    bool   `json:"is_resolved"`
	LatestService any    `json:"latest_service"`
	Name          string `json:"name"`
	RoomID        string `json:"room_id"`
	Source        string `json:"source"`
}
