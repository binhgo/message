package config

import "os"

type OutboundAPICredential struct {
	Url  string
	Auth string
	//For SSOv2
	AppKey string
	AppSecret string
}

type config struct {
	SSO            *OutboundAPICredential
	SSOv2		   *OutboundAPICredential
	DefaultOrgCode string
	MainDBName     string
	QueueDBName    string
	LogDBName      string

	MainAuthDB  string
	QueueAuthDB string
	LogAuthDB   string

	PingInterval      int
	FirebaseServerKey string
}

// Config main config object
var Config config

func init() {
	env := os.Getenv("env")

	switch env {

	// config for staging
	case "stg":
		Config = config{
			DefaultOrgCode:    "ghnexpress",
			MainDBName:        "internal-tools_stg_message",
			QueueDBName:       "internal-tools_stg_message_queue",
			LogDBName:         "internal-tools_stg_message_log",
			MainAuthDB:        "admin",
			QueueAuthDB:       "admin",
			LogAuthDB:         "admin",
			PingInterval:      30,
			FirebaseServerKey: "AAAAoSmr0Ow:APA91bFJCq2Kzl4FYIX0mANrnHt4cJPR9JsgEDM9MTt8lb231fReotw4urXIEihG6MWiNc2SALIpVWQWukUNPO1s7f2duwzWHIXinuZzXM-oTNiLWGlo_y39iHHExVI106QQy1NuXi5C",
			SSO: &OutboundAPICredential{
				Url:  "https://dev-online-gateway.ghn.vn/sso/api/staff/get",
				Auth: "Basic b25saW5lOnl6MmdxeWdiSGJKV0JENXpra2pRSmdqU1o2MzRWYkFx",
			},
			SSOv2: &OutboundAPICredential{
				Url: "https://dev-online-gateway.ghn.vn/sso-v2/public-api/staff",
				AppKey: "6e3132ca-8833-4db0-8a7a-3447d7a09d63",
				AppSecret: "b81e3762-530f-4061-ba41-089c1887f327",
			},
		}
		return
	// config for uat
	case "uat":
		Config = config{
			DefaultOrgCode:    "ghnexpress",
			MainDBName:        "internal-tools_prd_message",
			QueueDBName:       "internal-tools_prd_message_queue",
			LogDBName:         "internal-tools_prd_message_log",
			MainAuthDB:        "internal-tools_prd_message",
			QueueAuthDB:       "internal-tools_prd_message_queue",
			LogAuthDB:         "internal-tools_prd_message_log",
			PingInterval:      120,
			FirebaseServerKey: "AAAAoSmr0Ow:APA91bFJCq2Kzl4FYIX0mANrnHt4cJPR9JsgEDM9MTt8lb231fReotw4urXIEihG6MWiNc2SALIpVWQWukUNPO1s7f2duwzWHIXinuZzXM-oTNiLWGlo_y39iHHExVI106QQy1NuXi5C",
		}
		return
	// config for production
	case "prd":
		Config = config{
			DefaultOrgCode: "ghnexpress",
			MainDBName:     "internal-tools_prd_message",
			QueueDBName:    "internal-tools_prd_message_queue",
			LogDBName:      "internal-tools_prd_message_log",
			MainAuthDB:     "internal-tools_prd_message",
			QueueAuthDB:    "internal-tools_prd_message_queue",
			LogAuthDB:      "internal-tools_prd_message_log",
			PingInterval:   120,
		}
		return
	}
}
