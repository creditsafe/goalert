package heartbeat

import (
	"time"

	"github.com/lib/pq"
	"github.com/target/goalert/validation/validate"
)

// A Monitor will generate an alert if it does not receive a heartbeat within the configured TimeoutMinutes.
type Monitor struct {
	ID        string        `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	ServiceID string        `json:"service_id,omitempty"`
	Timeout   time.Duration `json:"timeout,omitempty"`

	lastState     State
	lastHeartbeat time.Time
}

// LastState returns the last known state.
func (m Monitor) LastState() State { return m.lastState }

// LastHeartbeat returns the timestamp of the last successful heartbeat.
func (m Monitor) LastHeartbeat() time.Time { return m.lastHeartbeat }

// Normalize performs validation and returns a new copy.
func (m Monitor) Normalize() (*Monitor, error) {
	err := validate.Many(
		validate.UUID("ServiceID", m.ServiceID),
		validate.IDName("Name", m.Name),
		validate.Duration("Timeout", m.Timeout, 5*time.Minute, 9000*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	m.Timeout = m.Timeout.Truncate(time.Minute)

	return &m, nil
}

func (m *Monitor) scanFrom(scanFn func(...interface{}) error) error {
	var t pq.NullTime
	var timeoutVal int

	err := scanFn(&m.ID, &m.Name, &m.ServiceID, &timeoutVal, &m.lastState, &t)
	if err != nil {
		return err
	}
	// Postgres EPOCH is seconds
	m.Timeout = time.Second * time.Duration(timeoutVal)
	m.lastHeartbeat = t.Time
	return nil
}
