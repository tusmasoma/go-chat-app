package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewHub(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()

	patterns := []struct {
		name string
		arg  struct {
			id   string
			name string
		}
		want struct {
			hub *Hub
			err error
		}
	}{
		{
			name: "Success: id is not empty",
			arg: struct {
				id   string
				name string
			}{
				id:   id,
				name: "name",
			},
			want: struct {
				hub *Hub
				err error
			}{
				hub: &Hub{
					ID:       id,
					Name:     "name",
					Clients:  make(map[*Client]bool),
					Channels: make(map[*Channel]bool),
				},
				err: nil,
			},
		},
		{
			name: "Success: id is empty",
			arg: struct {
				id   string
				name string
			}{
				id:   "",
				name: "name",
			},
			want: struct {
				hub *Hub
				err error
			}{
				hub: &Hub{
					Name:     "name",
					Clients:  make(map[*Client]bool),
					Channels: make(map[*Channel]bool),
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				id   string
				name string
			}{
				id:   id,
				name: "",
			},
			want: struct {
				hub *Hub
				err error
			}{
				hub: nil,
				err: errors.New("name is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hub, err := NewHub(tt.arg.id, tt.arg.name)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewHub() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewHub() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(hub, tt.want.hub, cmpopts.IgnoreFields(Hub{}, "ID")); len(d) != 0 {
				t.Errorf("NewHub() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Hub_RegisterClient(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	userID := uuid.New().String()
	client1, _ := NewClient(
		clientID,
		userID,
		nil,
	)

	patterns := []struct {
		name string
		arg  struct {
			client *Client
		}
		want struct {
			hub *Hub
		}
	}{
		{
			name: "Success",
			arg: struct {
				client *Client
			}{
				client: client1,
			},
			want: struct {
				hub *Hub
			}{
				hub: &Hub{
					Clients: map[*Client]bool{
						client1: true,
					},
				},
			},
		},
		{
			name: "Fail: client is nil",
			arg: struct {
				client *Client
			}{
				client: nil,
			},
			want: struct {
				hub *Hub
			}{
				hub: &Hub{
					Clients: make(map[*Client]bool),
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hub := &Hub{
				Clients: make(map[*Client]bool),
			}

			hub.RegisterClient(tt.arg.client)

			if d := cmp.Diff(hub.Clients, tt.want.hub.Clients); len(d) != 0 {
				t.Errorf("RegisterClient() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Hub_UnRegisterClient(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	userID := uuid.New().String()
	client1, _ := NewClient(
		clientID,
		userID,
		nil,
	)

	patterns := []struct {
		name string
		arg  struct {
			client *Client
		}
		want struct {
			hub *Hub
		}
	}{
		{
			name: "Success",
			arg: struct {
				client *Client
			}{
				client: client1,
			},
			want: struct {
				hub *Hub
			}{
				hub: &Hub{
					Clients: make(map[*Client]bool),
				},
			},
		},
		{
			name: "Fail: client is nil",
			arg: struct {
				client *Client
			}{
				client: nil,
			},
			want: struct {
				hub *Hub
			}{
				hub: &Hub{
					Clients: map[*Client]bool{
						client1: true,
					},
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hub := &Hub{
				Clients: map[*Client]bool{
					client1: true,
				},
			}

			hub.UnRegisterClient(tt.arg.client)

			if d := cmp.Diff(hub.Clients, tt.want.hub.Clients); len(d) != 0 {
				t.Errorf("UnRegisterClient() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Hub_FindChannelByID(t *testing.T) {
	t.Parallel()

	channelID := uuid.New().String()
	channel1, _ := NewChannel(
		channelID,
		"name",
		false,
	)

	hub := &Hub{
		Channels: map[*Channel]bool{
			channel1: true,
		},
	}

	patterns := []struct {
		name string
		arg  struct {
			id string
		}
		want struct {
			channel *Channel
		}
	}{
		{
			name: "Success",
			arg: struct {
				id string
			}{
				id: channelID,
			},
			want: struct {
				channel *Channel
			}{
				channel: channel1,
			},
		},
		{
			name: "Fail: channel is not found",
			arg: struct {
				id string
			}{
				id: uuid.New().String(),
			},
			want: struct {
				channel *Channel
			}{
				channel: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channel := hub.FindChannelByID(tt.arg.id)

			if d := cmp.Diff(channel, tt.want.channel); len(d) != 0 {
				t.Errorf("FindChannelByID() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Hub_FindChannelByName(t *testing.T) {
	t.Parallel()

	channelID := uuid.New().String()
	channel1, _ := NewChannel(
		channelID,
		"channel1",
		false,
	)

	hub := &Hub{
		Channels: map[*Channel]bool{
			channel1: true,
		},
	}

	patterns := []struct {
		name string
		arg  struct {
			name string
		}
		want struct {
			channel *Channel
		}
	}{
		{
			name: "Success",
			arg: struct {
				name string
			}{
				name: "channel1",
			},
			want: struct {
				channel *Channel
			}{
				channel: channel1,
			},
		},
		{
			name: "Fail: channel is not found",
			arg: struct {
				name string
			}{
				name: "channel2",
			},
			want: struct {
				channel *Channel
			}{
				channel: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channel := hub.FindChannelByName(tt.arg.name)

			if d := cmp.Diff(channel, tt.want.channel); len(d) != 0 {
				t.Errorf("FindChannelByName() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
