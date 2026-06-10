package main

import "time"

type Paste struct {
    Title      string     `json:"title"`
    Content    string     `json:"content"`
    Language   string     `json:"language"`
    Views      int        `json:"views"`
    CreatedAt  time.Time  `json:"created_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
    ShortID    string     `json:"short_id"`
    PrivateKey string 	  `json:"private_key,omitempty"`
    IsPublic   bool       `json:"is_public"`
}
