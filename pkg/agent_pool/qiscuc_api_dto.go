package agentpool

type QiscusAget struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	AuthenticationToken  string `json:"authentication_token"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	SdkEmail             string `json:"sdk_email"`
	SdkKey               string `json:"sdk_key"`
	IsAvailable          bool   `json:"is_available"`
	Type                 int    `json:"type"`
	AvatarURL            string `json:"avatar_url"`
	AppID                int    `json:"app_id"`
	IsVerified           bool   `json:"is_verified"`
	NotificationsRoomID  any    `json:"notifications_room_id"`
	BubbleColor          any    `json:"bubble_color"`
	QismoKey             string `json:"qismo_key"`
	DirectLoginToken     any    `json:"direct_login_token"`
	LastLogin            any    `json:"last_login"`
	ForceOffline         bool   `json:"force_offline"`
	DeletedAt            any    `json:"deleted_at"`
	IsTocAgree           bool   `json:"is_toc_agree"`
	TotpToken            any    `json:"totp_token"`
	IsReqOtpReset        any    `json:"is_req_otp_reset"`
	LastPasswordUpdate   string `json:"last_password_update"`
	LatestService        any    `json:"latest_service"`
	AssignedRules        []any  `json:"assigned_rules"`
	CurrentCustomerCount int    `json:"current_customer_count"`
	TotalResolved        string `json:"total_resolved"`
	TotalCustomers       string `json:"total_customers"`
	AssignedAgentRoles   []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		IsDefaultRole bool   `json:"is_default_role"`
	} `json:"assigned_agent_roles"`
	IsSupervisor bool `json:"is_supervisor"`
}

type QiscusAgentListDataApiResponse struct {
	Data []QiscusAget `json:"data"`
}

type QiscusListAgentsApiResponse struct {
	Data struct {
		Agents struct {
			CurrentPage  int          `json:"current_page"`
			Data         []QiscusAget `json:"data"`
			FirstPageURL string       `json:"first_page_url"`
			From         int          `json:"from"`
			LastPage     int          `json:"last_page"`
			LastPageURL  string       `json:"last_page_url"`
			NextPageURL  any          `json:"next_page_url"`
			Path         string       `json:"path"`
			PrevPageURL  any          `json:"prev_page_url"`
			To           int          `json:"to"`
			Total        int          `json:"total"`
		} `json:"agents"`
		Total       int `json:"total"`
		CurrentPage int `json:"current_page"`
	} `json:"data"`
}
