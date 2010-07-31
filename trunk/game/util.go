// Copyright (C) 2010 Alex Gontmakher (gsasha@gmail.com)
// This code is licensed under GPLv3

package game

import (
    "crypto/rand"
    "encoding/base64"
)

func create_random_id(bytes int) string {
    random_bytes := make([]byte, bytes)
    random_id := make([]byte, base64.URLEncoding.EncodedLen(bytes))
    rand.Read(random_bytes)
    base64.URLEncoding.Encode(random_id, random_bytes)
    return string(random_id)
}
