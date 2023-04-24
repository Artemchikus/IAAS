package models

import "time"

type KeyPair struct {
	CreatedAt   time.Time `json:"created_at"`
	Fingerprint string    `json:"fingerprint"`
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	IsDeleted   bool      `json:"is_deleted"`
	PrivateKey  string    `json:"private_key"`
	Type        string    `json:"type"`
	AccountID   string    `json:"user_id"`
	PublicKey   string    `json:"public_key"`
}

func NewKeyPair(PublicKey, Name, Type string) *KeyPair {
	return &KeyPair{
		PublicKey: PublicKey,
		Name:      Name,
		Type:      Type,
	}
}

// {
// 	"created_at": "2023-04-03T13:25:06.000000",
// 	"fingerprint": "af:60:d7:f5:12:ac:64:d9:aa:9d:b8:3f:9a:5a:26:24",
// 	"id": "mykey",
// 	"is_deleted": false,
// 	"name": "mykey",
// 	"private_key": null,
// 	"type": "ssh",
// 	"user_id": "63cec20a0cd44e689d889fc164d179b7"
//   }
