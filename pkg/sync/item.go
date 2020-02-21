package sync

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TypeNote               = 1
	TypeFolder             = 2
	TypeSetting            = 3
	TypeResource           = 4
	TypeTag                = 5
	TypeNoteTag            = 6
	TypeSearch             = 7
	TypeAlarm              = 8
	TypeMasterKey          = 9
	TypeItemChange         = 10
	TypeNoteResource       = 11
	TypeResourceLocalState = 12
	TypeRevision           = 13
	TypeMigration          = 14
)

type Type uint

func (t Type) String() string {
	switch t {
	case TypeNote:
		return "Note"
	case TypeFolder:
		return "Folder"
	case TypeSetting:
		return "Setting"
	case TypeResource:
		return "Resource"
	case TypeTag:
		return "Tag"
	case TypeNoteTag:
		return "NoteTag"
	case TypeSearch:
		return "Search"
	case TypeAlarm:
		return "Alarm"
	case TypeMasterKey:
		return "MasterKey"
	case TypeItemChange:
		return "ItemChange"
	case TypeNoteResource:
		return "NoteResource"
	case TypeResourceLocalState:
		return "ResourceLocalState"
	case TypeRevision:
		return "Revision"
	case TypeMigration:
		return "Migration"

	default:
		panic("unknown type")
	}
}

var IDLineRegex = regexp.MustCompile("^id: [a-z0-9]{32}$")

func ReadItem(path string) (*Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if err != nil {
		return nil, err
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Find the last line that starts with id:
	lastIdPos := 0
	for i, line := range lines {
		if IDLineRegex.MatchString(line) {
			lastIdPos = i
		}
	}

	var bodyBuilder strings.Builder
	item := &Item{}
	for i, line := range lines {
		if i < lastIdPos {
			fmt.Fprint(&bodyBuilder, line)
		} else {
			parts := strings.SplitN(line, ": ", 2)
			field := parts[0]
			value := parts[1]

			switch field {
			case "id":
				item.ID = value
			case "parent_id":
				item.ParentID = value
			case "item_id":
				item.ItemID = value
			case "item_type":
				t, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				item.ItemType = t
			case "type_":
				t, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				item.Type = Type(t)
			case "item_updated_time":
				t, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, err
				}
				item.ItemUpdatedTime = time.Unix(t, 0)
			}
		}
	}

	//item.Body = bodyBuilder.String()
	return item, nil
}

type Item struct {
	ID              string
	ParentID        string
	ItemType        int
	ItemID          string
	Body            string
	Type            Type
	ItemUpdatedTime time.Time
}
