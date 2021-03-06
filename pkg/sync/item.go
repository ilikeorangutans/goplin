package sync

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var IDLineRegex = regexp.MustCompile("^id: [a-z0-9]{32}$")

func NewItems(items []*Item) *Items {
	byID := make(map[string]*Item)
	for _, item := range items {
		byID[item.ID] = item
	}

	return &Items{
		Items: items,
		byID:  byID,
	}
}

// Items is all the sync items in a sync directory
type Items struct {
	Items []*Item
	byID  map[string]*Item
}

func (i *Items) Len() int {
	return len(i.Items)
}

func (i *Items) FindByID(id string) (*Item, error) {
	item, ok := i.byID[id]
	if !ok {
		return nil, fmt.Errorf("no item with id %q", id)
	}
	return item, nil
}

// ReadItem reads a given file path and parses it into a Item.
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
			fmt.Fprintln(&bodyBuilder, line)
		} else {
			parts := strings.SplitN(line, ": ", 2)
			field := parts[0]
			value := parts[1]

			fieldSetters := map[string]func(string, interface{}) error{
				"altitude":                  asFloat(item.setField),
				"application_data":          item.setField,
				"author":                    item.setField,
				"body_diff":                 item.setField,
				"created_time":              asTime(item.setField),
				"encryption_applied":        asBool(item.setField),
				"encryption_blob_encrypted": item.setField,
				"encryption_cipher_text":    item.setField,
				"file_extension":            item.setField,
				"filename":                  item.setField,
				"id":                        item.setField,
				"is_conflict":               asBool(item.setField),
				"is_shared":                 asBool(item.setField),
				"is_todo":                   asBool(item.setField),
				"item_id":                   item.setField,
				"item_type":                 asInt(item.setField),
				"item_updated_time":         asUnixTimestamp(item.setField),
				"latitude":                  asFloat(item.setField),
				"longitude":                 asFloat(item.setField),
				"markup_language":           asBool(item.setField),
				"metadata_diff":             item.setField,
				"mime":                      item.setField,
				"note_id":                   item.setField,
				"order":                     asInt(item.setField),
				"parent_id":                 item.setField,
				"size":                      asInt(item.setField),
				"source":                    item.setField,
				"source_application":        item.setField,
				"source_url":                asURL(item.setField),
				"tag_id":                    item.setField,
				"title_diff":                item.setField,
				"todo_completed":            asBool(item.setField),
				"todo_due":                  asUnixTimestamp(item.setField),
				"type_":                     asInt(item.setField),
				"updated_time":              asTime(item.setField),
				"user_created_time":         asTime(item.setField),
				"user_updated_time":         asTime(item.setField),
			}

			setter, ok := fieldSetters[field]
			if !ok {
				return nil, fmt.Errorf("no setter for field %q", field)
			}

			err := setter(field, value)
			if err != nil {
				return nil, err
			}
		}
	}

	item.Body = bodyBuilder.String()
	return item, nil
}

type Item struct {
	Altitude                float64
	ApplicationData         string
	Author                  string
	Body                    string
	BodyDiff                string
	CreatedTime             time.Time
	EncryptionApplied       bool
	EncryptionBlobEncrypted string
	EncryptionCipherText    string
	FileExtension           string
	Filename                string
	ID                      string
	IsConflict              bool
	IsShared                bool
	IsTodo                  bool
	ItemID                  string
	ItemType                int
	ItemUpdatedTime         time.Time
	Latitude                float64
	Longitude               float64
	MarkupLanguage          bool
	MetadataDiff            string
	Mime                    string
	NoteID                  string
	Order                   int
	ParentID                string
	Size                    int
	Source                  string
	SourceApplication       string
	SourceURL               *url.URL
	TagID                   string
	TitleDiff               string
	TodoCompleted           bool
	TodoDue                 time.Time
	Type                    Type
	UpdatedTime             time.Time
	UserCreatedTime         time.Time
	UserUpdatedTime         time.Time
}

func (item *Item) setField(field string, value interface{}) error {
	switch field {
	case "id":
		item.ID = value.(string)
	case "user_updated_time":
		item.UserUpdatedTime = value.(time.Time)
	case "created_time":
		item.CreatedTime = value.(time.Time)
	case "mime":
		item.Mime = value.(string)
	case "filename":
		item.Filename = value.(string)
	case "file_extension":
		item.FileExtension = value.(string)
	case "encryption_cipher_text":
		item.EncryptionCipherText = value.(string)
	case "encryption_applied":
		item.EncryptionApplied = value.(bool)
	case "encryption_blob_encrypted":
		item.EncryptionCipherText = value.(string)
	case "size":
		item.Size = value.(int)
	case "is_shared":
		item.IsShared = value.(bool)
	case "type_":
		t := Type(value.(int))
		item.Type = t
	case "parent_id":
		item.ParentID = value.(string)
	case "item_type":
		item.ItemType = value.(int)
	case "item_id":
		item.ItemID = value.(string)
	case "item_updated_time":
		item.ItemUpdatedTime = value.(time.Time)
	case "updated_time":
		item.UpdatedTime = value.(time.Time)
	case "user_created_time":
		item.UserCreatedTime = value.(time.Time)
	case "title_diff":
		item.TitleDiff = value.(string)
	case "body_diff":
		item.BodyDiff = value.(string)
	case "metadata_diff":
		item.MetadataDiff = value.(string)
	case "is_conflict":
		item.IsConflict = value.(bool)
	case "is_todo":
		item.IsTodo = value.(bool)
	case "latitude":
		item.Latitude = value.(float64)
	case "longitude":
		item.Longitude = value.(float64)
	case "altitude":
		item.Altitude = value.(float64)
	case "author":
		item.Author = value.(string)
	case "source_url":
		item.SourceURL = value.(*url.URL)
	case "todo_due":
		item.TodoDue = value.(time.Time)
	case "todo_completed":
		item.TodoCompleted = value.(bool)
	case "source":
		item.Source = value.(string)
	case "source_application":
		item.SourceApplication = value.(string)
	case "application_data":
		item.ApplicationData = value.(string)
	case "order":
		item.Order = value.(int)
	case "markup_language":
		item.MarkupLanguage = value.(bool)
	case "note_id":
		item.NoteID = value.(string)
	case "tag_id":
		item.TagID = value.(string)
	default:
		return fmt.Errorf("don't know how to set field %q", field)
	}
	return nil
}
