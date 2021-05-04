package config

import "os"

type OutboundAPICredential struct {
	Url  string
	Auth string
	// For SSOv2
	AppKey    string
	AppSecret string
}

type config struct {
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
			DefaultOrgCode:    "",
			MainDBName:        "internal-tools_stg_message",
			QueueDBName:       "internal-tools_stg_message_queue",
			LogDBName:         "internal-tools_stg_message_log",
			MainAuthDB:        "admin",
			QueueAuthDB:       "admin",
			LogAuthDB:         "admin",
			PingInterval:      30,
			FirebaseServerKey: "AAAAoSmr0Ow:APA91bFJCq2Kzl4FYIX0mANrnHt4cJPR9JsgEDM9MTt8lb231fReotw4urXIEihG6MWiNc2SALIpVWQWukUNPO1s7f2duwzWHIXinuZzXM-oTNiLWGlo_y39iHHExVI106QQy1NuXi5C",
			
		}
		return
	// config for uat
	case "uat":
		Config = config{
			DefaultOrgCode:    "",
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
			DefaultOrgCode: "",
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
