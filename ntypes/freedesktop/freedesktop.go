// +build !darwin
// +build !windows

package freedesktop

import (
	"fmt"

	"github.com/godbus/dbus"
)

type Notification struct {
	AppName    string
	ReplacesID uint
	AppIcon    string
	Summary    string
	Body       string
	Actions    []string
	// Hints         string
	ExpireTimeout int
}

func (n *Notification) Send() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("dbus connect: %s", err)
	}
	defer conn.Close()

	fdn := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	resp := fdn.Call(
		"org.freedesktop.Notifications.Notify", dbus.FlagNoAutoStart,
		n.AppName,
		uint32(n.ReplacesID),
		n.AppIcon,
		n.Summary,
		n.Body,
		n.Actions,
		map[string]dbus.Variant{},
		int32(n.ExpireTimeout),
	)

	if resp.Err != nil {
		return fmt.Errorf("notify: %s", resp.Err)
	}

	return nil
}

func (n *Notification) SetMessage(m string) {
	n.Body = m
}

func (n *Notification) Message() string {
	return n.Body
}
