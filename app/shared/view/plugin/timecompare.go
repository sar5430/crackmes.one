package plugin

import (
    "html/template"
    "time"
)

// TimeCompare compares 2 times, and returns true, if
// the difference between them are larger than the specified seconds
func TimeCompare() template.FuncMap {
    f := make(template.FuncMap)

    f["TIMECOMPARE"] = func(t1, t2 time.Time, seconds int64) bool {
        diff := t1.Unix() - t2.Unix()
        if diff < 0 {
            diff *= -1
        }
        return diff >= seconds
    }

    return f
}
