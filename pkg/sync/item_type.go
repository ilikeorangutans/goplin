package sync

const (
	TypeNote               Type = 1
	TypeFolder             Type = 2
	TypeSetting            Type = 3
	TypeResource           Type = 4
	TypeTag                Type = 5
	TypeNoteTag            Type = 6
	TypeSearch             Type = 7
	TypeAlarm              Type = 8
	TypeMasterKey          Type = 9
	TypeItemChange         Type = 10
	TypeNoteResource       Type = 11
	TypeResourceLocalState Type = 12
	TypeRevision           Type = 13
	TypeMigration          Type = 14
)

// Type is the type of a sync item
type Type int

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
