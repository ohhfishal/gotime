package entry

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	EncodeDecoder EncodeDecoder
}

func (c Config) Encode(entry Entry) (string, error) {
	return c.EncodeDecoder.Encode(entry)
}

func (c Config) Decode(stream string) (*Entry, error) {
	return c.EncodeDecoder.Decode(stream)
}

type EncodeDecoder interface {
	Encode(Entry) (string, error)
	Decode(string) (*Entry, error)
}

func defaultConfig() Config {
	return Config{
		EncodeDecoder: defaultEncodeDecoder{
			TimeLayout: time.DateTime,
		},
	}
}

func DefaultConfig(config ...Config) Config {
	if len(config) == 0 {
		return defaultConfig()
	}

	cfg := config[0]
	if cfg.EncodeDecoder == nil {
		cfg.EncodeDecoder = defaultEncodeDecoder{
			TimeLayout: time.DateTime,
		}
	}
	return cfg
}

type defaultEncodeDecoder struct {
	TimeLayout string
}

func (ed defaultEncodeDecoder) Encode(entry Entry) (string, error) {
	return fmt.Sprintf(`%v,%s: %s`,
		entry.Time.Format(ed.TimeLayout),
		entry.Category,
		entry.Note,
	), nil
}

func (ed defaultEncodeDecoder) Decode(stream string) (*Entry, error) {
	var entry Entry
	var err error

	// Parse time
	timeLength := len(ed.TimeLayout)
	entry.Time, err = time.Parse(ed.TimeLayout, stream[:timeLength])
	if err != nil {
		return nil, fmt.Errorf("parsing time: %w", err)
	}

	// Parse Category
	stream = stream[timeLength:]
	if stream == `` {
		return nil, fmt.Errorf("unexpected eof")
	}

	if stream[0] != ',' {
		return nil, fmt.Errorf(`expected :",": found "%b"`, stream[0])
	}
	cutIndex := strings.Index(stream, ": ")
	if cutIndex == -1 {
		return nil, fmt.Errorf("unexpected eof")
	}
	entry.Category = stream[1:cutIndex]
	entry.Note = strings.TrimSpace(stream[cutIndex+1:])
	return &entry, nil
}
