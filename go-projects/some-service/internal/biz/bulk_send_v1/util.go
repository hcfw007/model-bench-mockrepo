package bulk_send_v1

import "time"

func nowMillis() int64 { return time.Now().UnixMilli() }
