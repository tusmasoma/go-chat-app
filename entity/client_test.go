package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewClient(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	userID := uuid.New().String()
	hub := &Hub{}

	patterns := []struct {
		name string
		arg  struct {
			id     string
			userID string
			hub    *Hub
		}
		want struct {
			client *Client
			err    error
		}
	}{
		{
			name: "Success: id is not empty",
			arg: struct {
				id     string
				userID string
				hub    *Hub
			}{
				id:     id,
				userID: userID,
				hub:    hub,
			},
			want: struct {
				client *Client
				err    error
			}{
				client: &Client{
					ID:       id,
					UserID:   userID,
					Hub:      hub,
					Channels: make(map[*Channel]bool),
				},
				err: nil,
			},
		},
		{
			name: "Success: id is empty",
			arg: struct {
				id     string
				userID string
				hub    *Hub
			}{
				id:     "",
				userID: userID,
				hub:    hub,
			},
			want: struct {
				client *Client
				err    error
			}{
				client: &Client{
					UserID:   userID,
					Hub:      hub,
					Channels: make(map[*Channel]bool),
				},
				err: nil,
			},
		},
		{
			name: "Fail: userID is empty",
			arg: struct {
				id     string
				userID string
				hub    *Hub
			}{
				id:     id,
				userID: "",
				hub:    hub,
			},
			want: struct {
				client *Client
				err    error
			}{
				client: nil,
				err:    errors.New("userID is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(tt.arg.id, tt.arg.userID, tt.arg.hub)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(client, tt.want.client, cmpopts.IgnoreFields(Client{}, "ID")); len(d) != 0 {
				t.Errorf("NewClient() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestClient_JoinChannel(t *testing.T) {
	t.Parallel()

	channle1 := &Channel{
		ID:      uuid.New().String(),
		Name:    "channel1",
		Private: false,
	}

	client := &Client{
		Channels: make(map[*Channel]bool),
	}

	patterns := []struct {
		name string
		arg  struct {
			channel *Channel
		}
		want struct {
			client *Client
		}
	}{
		{
			name: "Success: join channel",
			arg: struct {
				channel *Channel
			}{
				channel: channle1,
			},
			want: struct {
				client *Client
			}{
				client: &Client{
					Channels: map[*Channel]bool{
						channle1: true,
					},
				},
			},
		},
		{
			name: "Success: join channel twice",
			arg: struct {
				channel *Channel
			}{
				channel: channle1,
			},
			want: struct {
				client *Client
			}{
				client: &Client{
					Channels: map[*Channel]bool{
						channle1: true,
					},
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client.JoinChannel(tt.arg.channel)

			if d := cmp.Diff(client.Channels, tt.want.client.Channels); len(d) != 0 {
				t.Errorf("JoinChannel() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestClient_LeaveChannel(t *testing.T) {
	t.Parallel()

	channle1 := &Channel{
		ID:      uuid.New().String(),
		Name:    "channel1",
		Private: false,
	}
	channle2 := &Channel{
		ID:      uuid.New().String(),
		Name:    "channel2",
		Private: false,
	}

	patterns := []struct {
		name string
		arg  struct {
			channel *Channel
		}
		want struct {
			client *Client
		}
	}{
		{
			name: "Success: leave channel",
			arg: struct {
				channel *Channel
			}{
				channel: channle1,
			},
			want: struct {
				client *Client
			}{
				client: &Client{
					Channels: map[*Channel]bool{},
				},
			},
		},
		{
			name: "Success: leave channel not joined",
			arg: struct {
				channel *Channel
			}{
				channel: channle2,
			},
			want: struct {
				client *Client
			}{
				client: &Client{
					Channels: map[*Channel]bool{
						channle1: true,
					},
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &Client{
				Channels: map[*Channel]bool{
					channle1: true,
				},
			}

			client.LeaveChannel(tt.arg.channel)

			if d := cmp.Diff(client.Channels, tt.want.client.Channels); len(d) != 0 {
				t.Errorf("LeaveChannel() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
