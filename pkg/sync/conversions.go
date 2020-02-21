package sync

import (
	"net/url"
	"strconv"
	"time"
)

type Setter func(string, interface{}) error

func asURL(setter Setter) Setter {
	return func(field string, value interface{}) error {
		s := value.(string)
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		return setter(field, u)
	}
}

func asFloat(setter Setter) Setter {
	return func(field string, value interface{}) error {
		s := value.(string)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		return setter(field, f)
	}
}

func asInt(setter Setter) Setter {
	return func(field string, value interface{}) error {
		s := value.(string)
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		return setter(field, i)
	}
}

func asBool(setter Setter) Setter {
	return asInt(func(field string, value interface{}) error {
		i := value.(int)
		return setter(field, i == 1)
	})
}

func asTime(setter Setter) Setter {
	return func(field string, value interface{}) error {
		s := value.(string)
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return nil
		}
		return setter(field, t)
	}
}

func asUnixTimestamp(setter Setter) Setter {
	return func(field string, value interface{}) error {
		s := value.(string)
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		parsedTime := time.Unix(t, 0)
		return setter(field, parsedTime)
	}
}
