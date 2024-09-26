package entity

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestEntity_NewChannel(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()

	patterns := []struct {
		name string
		arg  struct {
			id      string
			name    string
			private bool
		}
		want struct {
			channel *Channel
			err     error
		}
	}{
		{
			name: "Success: id is not empty",
			arg: struct {
				id      string
				name    string
				private bool
			}{
				id:      id,
				name:    "channel",
				private: false,
			},
			want: struct {
				channel *Channel
				err     error
			}{
				channel: &Channel{
					ID:      id,
					Name:    "channel",
					Private: false,
					Clients: make(map[*Client]bool),
				},
				err: nil,
			},
		},
		{
			name: "Success: id is empty",
			arg: struct {
				id      string
				name    string
				private bool
			}{
				id:      "",
				name:    "channel",
				private: false,
			},
			want: struct {
				channel *Channel
				err     error
			}{
				channel: &Channel{
					Name:    "channel",
					Private: false,
					Clients: make(map[*Client]bool),
				},
				err: nil,
			},
		},
		{
			name: "Fail: name is empty",
			arg: struct {
				id      string
				name    string
				private bool
			}{
				id:      id,
				name:    "",
				private: false,
			},
			want: struct {
				channel *Channel
				err     error
			}{
				channel: nil,
				err:     errors.New("name is required"),
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channel, err := NewChannel(tt.arg.id, tt.arg.name, tt.arg.private)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("NewChannel() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("NewChannel() error = %v, wantErr %v", err, tt.want.err)
			}

			if d := cmp.Diff(channel, tt.want.channel, cmpopts.IgnoreFields(Channel{}, "ID")); len(d) != 0 {
				t.Errorf("NewChannel() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Channel_RegisterClientInChannel(t *testing.T) {
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
			channel *Channel
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
				channel *Channel
			}{
				channel: &Channel{
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
				channel *Channel
			}{
				channel: &Channel{
					Clients: make(map[*Client]bool),
				},
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			channel := &Channel{
				Clients: make(map[*Client]bool),
			}

			channel.RegisterClientInChannel(tt.arg.client)

			if d := cmp.Diff(channel.Clients, tt.want.channel.Clients); len(d) != 0 {
				t.Errorf("RegisterClientInChannel() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Channel_UnRegisterClientInChannel(t *testing.T) {
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
			channel *Channel
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
				channel *Channel
			}{
				channel: &Channel{
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
				channel *Channel
			}{
				channel: &Channel{
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

			channel := &Channel{
				Clients: map[*Client]bool{
					client1: true,
				},
			}

			channel.UnRegisterClientInChannel(tt.arg.client)

			if d := cmp.Diff(channel.Clients, tt.want.channel.Clients); len(d) != 0 {
				t.Errorf("UnRegisterClientInChannel() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Channel_FindClientByID(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	userID := uuid.New().String()
	client1, _ := NewClient(
		clientID,
		userID,
		nil,
	)

	channel := &Channel{
		Clients: map[*Client]bool{
			client1: true,
		},
	}

	patterns := []struct {
		name string
		arg  struct {
			id string
		}
		want struct {
			client *Client
		}
	}{
		{
			name: "Success",
			arg: struct {
				id string
			}{
				id: clientID,
			},
			want: struct {
				client *Client
			}{
				client: client1,
			},
		},
		{
			name: "Fail: client is not found",
			arg: struct {
				id string
			}{
				id: uuid.New().String(),
			},
			want: struct {
				client *Client
			}{
				client: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := channel.FindClientByID(tt.arg.id)

			if d := cmp.Diff(client, tt.want.client); len(d) != 0 {
				t.Errorf("FindClientByID() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestEntity_Channel_FindClientByUserID(t *testing.T) {
	t.Parallel()

	clientID := uuid.New().String()
	userID := uuid.New().String()
	client1, _ := NewClient(
		clientID,
		userID,
		nil,
	)

	channel := &Channel{
		Clients: map[*Client]bool{
			client1: true,
		},
	}

	patterns := []struct {
		name string
		arg  struct {
			userID string
		}
		want struct {
			client *Client
		}
	}{
		{
			name: "Success",
			arg: struct {
				userID string
			}{
				userID: userID,
			},
			want: struct {
				client *Client
			}{
				client: client1,
			},
		},
		{
			name: "Fail: channel is not found",
			arg: struct {
				userID string
			}{
				userID: uuid.New().String(),
			},
			want: struct {
				client *Client
			}{
				client: nil,
			},
		},
	}

	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := channel.FindClientByUserID(tt.arg.userID)

			if d := cmp.Diff(client, tt.want.client); len(d) != 0 {
				t.Errorf("FindClientByUserID() mismatch (-got +want):\n%s", d)
			}
		})
	}
}
