package action

import (
	"time"
)

type FirebaseNotification struct {
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
	ChannelID   string `json:"android_channel_id,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Sound       string `json:"sound,omitempty"`
	Badge       string `json:"badge,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Color       string `json:"color,omitempty"`
	ClickAction string `json:"click_action,omitempty"`
}

type messages struct {
	To                       string                 `json:"to,omitempty"`
	RegistrationIDs          []string               `json:"registration_ids,omitempty"`
	Condition                string                 `json:"condition,omitempty"`
	CollapseKey              string                 `json:"collapse_key,omitempty"`
	Priority                 string                 `json:"priority,omitempty"`
	ContentAvailable         bool                   `json:"content_available,omitempty"`
	MutableContent           bool                   `json:"mutable_content,omitempty"`
	DelayWhileIdle           bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive               time.Duration          `json:"time_to_live,omitempty"`
	DeliveryReceiptRequested bool                   `json:"delivery_receipt_requested,omitempty"`
	DryRun                   bool                   `json:"dry_run,omitempty"`
	RestrictedPackageName    string                 `json:"restricted_package_name,omitempty"`
	Notification             *FirebaseNotification  `json:"notification,omitempty"`
	Data                     map[string]interface{} `json:"data,omitempty"`
}

type FirebaseMessages interface {
	apply(*messages)
}

type messageFunc func(*messages)

func (f messageFunc) apply(m *messages) {
	f(m)
}

func (c *FirebaseClient) WithToken(t string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.To = t
	})
}

func (c *FirebaseClient) WithRegistrationIDs(ids []string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.RegistrationIDs = ids
	})
}

func (c *FirebaseClient) WithCondition(cd string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.Condition = cd
	})
}

func (c *FirebaseClient) WithCollapseKey(k string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.CollapseKey = k
	})
}

func (c *FirebaseClient) WithPriority(p string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.Priority = p
	})
}

func (c *FirebaseClient) WithContentAvailable(b bool) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.ContentAvailable = b
	})
}

func (c *FirebaseClient) WithMutableContent(b bool) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.MutableContent = b
	})
}

func (c *FirebaseClient) WithDelay(b bool) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.DelayWhileIdle = b
	})
}

func (c *FirebaseClient) WithTimeToLive(t time.Duration) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.TimeToLive = t
	})
}

func (c *FirebaseClient) WithDeliveryReceiptRequested(b bool) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.DeliveryReceiptRequested = b
	})
}

func (c *FirebaseClient) WithDryRun(b bool) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.DryRun = b
	})
}

func (c *FirebaseClient) WithRestrictedPackageName(r string) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.RestrictedPackageName = r
	})
}

func (c *FirebaseClient) WithNotification(nf *FirebaseNotification) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.Notification = nf
	})
}

func (c *FirebaseClient) WithData(d map[string]interface{}) FirebaseMessages {
	return messageFunc(func(m *messages) {
		m.Data = d
	})
}
