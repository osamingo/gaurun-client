package gaurun

import "github.com/mercari/gaurun/gaurun"

const (
	// PlatformAndroid is enum for FCM/GCM.
	PlatformAndroid Platform = gaurun.PlatFormAndroid
	// PlatformIOS is enum for APNs.
	PlatformIOS Platform = gaurun.PlatFormIos
)

type (
	// A Payload has notifications.
	Payload struct {
		Notifications []*Notification `json:"notifications"`
	}
	// A Notification has gaurun notification data.
	Notification struct {
		// Common
		Tokens     []string `json:"token"`
		Platform   Platform `json:"platform"`
		Message    string   `json:"message"`
		Identifier string   `json:"identifier,omitempty"`
		// Android
		AndroidSetting
		// iOS
		IOSSetting
		// Metadata
		ID uint64 `json:"seq_id,omitempty"`
	}
	// An AndroidSetting has setting fields for FCM/GCM.
	AndroidSetting struct {
		CollapseKey    string `json:"collapse_key,omitempty"`
		DelayWhileIdle bool   `json:"delay_while_idle,omitempty"`
		TimeToLive     int    `json:"time_to_live,omitempty"`
	}
	// An IOSSetting has setting fields for APNs.
	IOSSetting struct {
		Title            string    `json:"title,omitempty"`
		Subtitle         string    `json:"subtitle,omitempty"`
		Badge            int       `json:"badge,omitempty"`
		Category         string    `json:"category,omitempty"`
		Sound            string    `json:"sound,omitempty"`
		ContentAvailable bool      `json:"content_available,omitempty"`
		MutableContent   bool      `json:"mutable_content,omitempty"`
		Expiry           int       `json:"expiry,omitempty"`
		Retry            int       `json:"retry,omitempty"`
		Extend           []*Extend `json:"extend,omitempty"`
	}
	// An Extend is alias gaurun.ExtendJSON.
	Extend gaurun.ExtendJSON
	// A Platform is alias gaurun platform enum.
	Platform int
)
